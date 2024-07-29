package game

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/nsfisis/iosdc-2024-albatross/backend/api"
	"github.com/nsfisis/iosdc-2024-albatross/backend/db"
)

type playerClientState int

const (
	playerClientStateWaitingEntries playerClientState = iota
	playerClientStateEntried
	playerClientStateReady
)

type gameHub struct {
	ctx               context.Context
	game              *game
	q                 *db.Queries
	players           map[*playerClient]playerClientState
	registerPlayer    chan *playerClient
	unregisterPlayer  chan *playerClient
	playerC2SMessages chan *playerMessageC2SWithClient
	watchers          map[*watcherClient]bool
	registerWatcher   chan *watcherClient
	unregisterWatcher chan *watcherClient
}

func newGameHub(ctx context.Context, game *game, q *db.Queries) *gameHub {
	return &gameHub{
		ctx:               ctx,
		game:              game,
		q:                 q,
		players:           make(map[*playerClient]playerClientState),
		registerPlayer:    make(chan *playerClient),
		unregisterPlayer:  make(chan *playerClient),
		playerC2SMessages: make(chan *playerMessageC2SWithClient),
		watchers:          make(map[*watcherClient]bool),
		registerWatcher:   make(chan *watcherClient),
		unregisterWatcher: make(chan *watcherClient),
	}
}

func (hub *gameHub) run() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case player := <-hub.registerPlayer:
			hub.players[player] = playerClientStateWaitingEntries
		case player := <-hub.unregisterPlayer:
			if _, ok := hub.players[player]; ok {
				hub.closePlayerClient(player)
			}
		case watcher := <-hub.registerWatcher:
			hub.watchers[watcher] = true
		case watcher := <-hub.unregisterWatcher:
			if _, ok := hub.watchers[watcher]; ok {
				hub.closeWatcherClient(watcher)
			}
		case message := <-hub.playerC2SMessages:
			switch msg := message.message.(type) {
			case *playerMessageC2SEntry:
				log.Printf("entry: %v", message.message)
				// TODO: assert state is waiting_entries
				hub.players[message.client] = playerClientStateEntried
				entriedPlayerCount := 0
				for _, state := range hub.players {
					if playerClientStateEntried <= state {
						entriedPlayerCount++
					}
				}
				if entriedPlayerCount == 2 {
					for player := range hub.players {
						player.s2cMessages <- &playerMessageS2CPrepare{
							Type: playerMessageTypeS2CPrepare,
							Data: playerMessageS2CPreparePayload{
								Problem: api.Problem{
									ProblemId:   1,
									Title:       "the answer",
									Description: "print 42",
								},
							},
						}
					}
					err := hub.q.UpdateGameState(hub.ctx, db.UpdateGameStateParams{
						GameID: int32(hub.game.gameID),
						State:  string(gameStatePrepare),
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.state = gameStatePrepare
				}
			case *playerMessageC2SReady:
				log.Printf("ready: %v", message.message)
				// TODO: assert state is prepare
				hub.players[message.client] = playerClientStateReady
				readyPlayerCount := 0
				for _, state := range hub.players {
					if playerClientStateReady <= state {
						readyPlayerCount++
					}
				}
				if readyPlayerCount == 2 {
					startAt := time.Now().Add(11 * time.Second).UTC()
					for player := range hub.players {
						player.s2cMessages <- &playerMessageS2CStart{
							Type: playerMessageTypeS2CStart,
							Data: playerMessageS2CStartPayload{
								StartAt: int(startAt.Unix()),
							},
						}
					}
					for watcher := range hub.watchers {
						watcher.s2cMessages <- &watcherMessageS2CStart{
							Type: watcherMessageTypeS2CStart,
							Data: watcherMessageS2CStartPayload{
								StartAt: int(startAt.Unix()),
							},
						}
					}
					err := hub.q.UpdateGameStartedAt(hub.ctx, db.UpdateGameStartedAtParams{
						GameID: int32(hub.game.gameID),
						StartedAt: pgtype.Timestamp{
							Time:             startAt,
							InfinityModifier: pgtype.Finite,
							Valid:            true,
						},
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.startedAt = &startAt
					err = hub.q.UpdateGameState(hub.ctx, db.UpdateGameStateParams{
						GameID: int32(hub.game.gameID),
						State:  string(gameStateStarting),
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.state = gameStateStarting
				}
			case *playerMessageC2SCode:
				// TODO: assert game state is gaming
				log.Printf("code: %v", message.message)
				code := msg.Data.Code
				score := len(code)
				message.client.s2cMessages <- &playerMessageS2CExecResult{
					Type: playerMessageTypeS2CExecResult,
					Data: playerMessageS2CExecResultPayload{
						Score:  &score,
						Status: api.GamePlayerMessageS2CExecResultPayloadStatusSuccess,
					},
				}
				for watcher := range hub.watchers {
					watcher.s2cMessages <- &watcherMessageS2CCode{
						Type: watcherMessageTypeS2CCode,
						Data: watcherMessageS2CCodePayload{
							PlayerId: message.client.playerID,
							Code:     code,
						},
					}
					watcher.s2cMessages <- &watcherMessageS2CExecResult{
						Type: watcherMessageTypeS2CExecResult,
						Data: watcherMessageS2CExecResultPayload{
							PlayerId: message.client.playerID,
							Score:    &score,
							Stdout:   "",
							Stderr:   "",
						},
					}
				}
			default:
				log.Fatalf("unexpected message type: %T", message.message)
			}
		case <-ticker.C:
			if hub.game.state == gameStateStarting {
				if time.Now().After(*hub.game.startedAt) {
					err := hub.q.UpdateGameState(hub.ctx, db.UpdateGameStateParams{
						GameID: int32(hub.game.gameID),
						State:  string(gameStateGaming),
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.state = gameStateGaming
				}
			} else if hub.game.state == gameStateGaming {
				if time.Now().After(hub.game.startedAt.Add(time.Duration(hub.game.durationSeconds) * time.Second)) {
					err := hub.q.UpdateGameState(hub.ctx, db.UpdateGameStateParams{
						GameID: int32(hub.game.gameID),
						State:  string(gameStateFinished),
					})
					if err != nil {
						log.Fatalf("failed to set game state: %v", err)
					}
					hub.game.state = gameStateFinished
				}
				hub.close()
				return
			}
		}
	}
}

func (hub *gameHub) close() {
	for client := range hub.players {
		hub.closePlayerClient(client)
	}
	close(hub.registerPlayer)
	close(hub.unregisterPlayer)
	close(hub.playerC2SMessages)
	for watcher := range hub.watchers {
		hub.closeWatcherClient(watcher)
	}
	close(hub.registerWatcher)
	close(hub.unregisterWatcher)
}

func (hub *gameHub) closePlayerClient(player *playerClient) {
	delete(hub.players, player)
	close(player.s2cMessages)
}

func (hub *gameHub) closeWatcherClient(watcher *watcherClient) {
	delete(hub.watchers, watcher)
	close(watcher.s2cMessages)
}

type GameHubs struct {
	hubs map[int]*gameHub
	q    *db.Queries
}

func NewGameHubs(q *db.Queries) *GameHubs {
	return &GameHubs{
		hubs: make(map[int]*gameHub),
		q:    q,
	}
}

func (hubs *GameHubs) Close() {
	for _, hub := range hubs.hubs {
		hub.close()
	}
}

func (hubs *GameHubs) getHub(gameID int) *gameHub {
	return hubs.hubs[gameID]
}

func (hubs *GameHubs) RestoreFromDB(ctx context.Context) error {
	games, err := hubs.q.ListGames(ctx)
	if err != nil {
		return err
	}
	for _, row := range games {
		var startedAt *time.Time
		if row.StartedAt.Valid {
			startedAt = &row.StartedAt.Time
		}
		var problem_ *problem
		if row.ProblemID != nil {
			if row.Title == nil || row.Description == nil {
				panic("inconsistent data")
			}
			problem_ = &problem{
				problemID:   int(*row.ProblemID),
				title:       *row.Title,
				description: *row.Description,
			}
		}
		hubs.hubs[int(row.GameID)] = newGameHub(ctx, &game{
			gameID:          int(row.GameID),
			durationSeconds: int(row.DurationSeconds),
			state:           gameState(row.State),
			displayName:     row.DisplayName,
			startedAt:       startedAt,
			problem:         problem_,
		}, hubs.q)
	}
	return nil
}

func (hubs *GameHubs) Run() {
	for _, hub := range hubs.hubs {
		go hub.run()
	}
}

func (hubs *GameHubs) SockHandler() *sockHandler {
	return newSockHandler(hubs)
}
