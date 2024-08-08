package game

import (
	"encoding/json"
	"fmt"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/api"
)

const (
	playerMessageTypeS2CPrepare    = "player:s2c:prepare"
	playerMessageTypeS2CStart      = "player:s2c:start"
	playerMessageTypeS2CExecResult = "player:s2c:execresult"
	playerMessageTypeC2SEntry      = "player:c2s:entry"
	playerMessageTypeC2SReady      = "player:c2s:ready"
	playerMessageTypeC2SCode       = "player:c2s:code"
	playerMessageTypeC2SSubmit     = "player:c2s:submit"
)

type playerMessageC2SWithClient struct {
	client  *playerClient
	message playerMessageC2S
}

type playerMessageS2C = interface{}
type playerMessageS2CPrepare = api.GamePlayerMessageS2CPrepare
type playerMessageS2CPreparePayload = api.GamePlayerMessageS2CPreparePayload
type playerMessageS2CStart = api.GamePlayerMessageS2CStart
type playerMessageS2CStartPayload = api.GamePlayerMessageS2CStartPayload
type playerMessageS2CExecResult = api.GamePlayerMessageS2CExecResult
type playerMessageS2CExecResultPayload = api.GamePlayerMessageS2CExecResultPayload

type playerMessageC2S = interface{}
type playerMessageC2SEntry = api.GamePlayerMessageC2SEntry
type playerMessageC2SReady = api.GamePlayerMessageC2SReady
type playerMessageC2SCode = api.GamePlayerMessageC2SCode
type playerMessageC2SCodePayload = api.GamePlayerMessageC2SCodePayload
type playerMessageC2SSubmit = api.GamePlayerMessageC2SSubmit
type playerMessageC2SSubmitPayload = api.GamePlayerMessageC2SSubmitPayload

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
	case playerMessageTypeC2SCode:
		var payload playerMessageC2SCodePayload
		if err := json.Unmarshal(raw["data"], &payload); err != nil {
			return nil, err
		}
		return &playerMessageC2SCode{
			Type: playerMessageTypeC2SCode,
			Data: payload,
		}, nil
	case playerMessageTypeC2SSubmit:
		var payload playerMessageC2SSubmitPayload
		if err := json.Unmarshal(raw["data"], &payload); err != nil {
			return nil, err
		}
		return &playerMessageC2SSubmit{
			Type: playerMessageTypeC2SSubmit,
			Data: payload,
		}, nil
	default:
		return nil, fmt.Errorf("unknown message type: %s", typ)
	}
}

const (
	watcherMessageTypeS2CStart      = "watcher:s2c:start"
	watcherMessageTypeS2CExecResult = "watcher:s2c:execresult"
	watcherMessageTypeS2CCode       = "watcher:s2c:code"
)

type watcherMessageS2C = interface{}
type watcherMessageS2CStart = api.GameWatcherMessageS2CStart
type watcherMessageS2CStartPayload = api.GameWatcherMessageS2CStartPayload
type watcherMessageS2CCode = api.GameWatcherMessageS2CCode
type watcherMessageS2CCodePayload = api.GameWatcherMessageS2CCodePayload
