package game

import (
	"encoding/json"
	"fmt"

	"github.com/nsfisis/iosdc-2024-albatross/backend/api"
)

const (
	playerMessageTypeS2CPrepare = "player:s2c:prepare"
	playerMessageTypeS2CStart   = "player:s2c:start"
	playerMessageTypeC2SEntry   = "player:c2s:entry"
	playerMessageTypeC2SReady   = "player:c2s:ready"
)

type playerMessageC2SWithClient struct {
	client  *playerClient
	message playerMessageC2S
}

type playerMessage = api.GamePlayerMessage

type playerMessageS2C = interface{}
type playerMessageS2CPrepare = api.GamePlayerMessageS2CPrepare
type playerMessageS2CPreparePayload = api.GamePlayerMessageS2CPreparePayload
type playerMessageS2CStart = api.GamePlayerMessageS2CStart
type playerMessageS2CStartPayload = api.GamePlayerMessageS2CStartPayload

type playerMessageC2S = interface{}
type playerMessageC2SEntry = api.GamePlayerMessageC2SEntry
type playerMessageC2SReady = api.GamePlayerMessageC2SReady

func asPlayerMessageC2S(raw map[string]json.RawMessage) (playerMessageC2S, error) {
	var typ string
	if err := json.Unmarshal(raw["type"], &typ); err != nil {
		return nil, err
	}

	switch typ {
	case playerMessageTypeC2SEntry:
		return &playerMessageC2SEntry{
			Type: playerMessageTypeC2SEntry,
		}, nil
	case playerMessageTypeC2SReady:
		return &playerMessageC2SReady{
			Type: playerMessageTypeC2SReady,
		}, nil
	default:
		return nil, fmt.Errorf("unknown message type: %s", typ)
	}
}

type watcherMessageS2C = interface{}
