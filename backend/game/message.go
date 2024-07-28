package game

import (
	"encoding/json"
	"fmt"
)

type MessageWithClient struct {
	Client  *GameClient
	Message *Message
}

type Message struct {
	Type string      `json:"type"`
	Data MessageData `json:"data"`
}

type MessageData interface{}

type MessageDataConnect struct {
}

type MessageDataPrepare struct {
	Problem string `json:"problem"`
}

type MessageDataReady struct {
}

type MessageDataStart struct {
	StartTime string `json:"startTime"`
}

type MessageDataCode struct {
	Code string `json:"code"`
}

type MessageDataScore struct {
	Score int `json:"score"`
}

type MessageDataFinish struct {
	YourScore     *int `json:"yourScore"`
	OpponentScore *int `json:"opponentScore"`
}

type MessageDataWatch struct {
	Problem string `json:"problem"`
	ScoreA  *int   `json:"scoreA"`
	CodeA   string `json:"codeA"`
	ScoreB  *int   `json:"scoreB"`
	CodeB   string `json:"codeB"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if err := json.Unmarshal(raw["type"], &m.Type); err != nil {
		return err
	}

	var err error
	switch m.Type {
	case "connect":
		var data MessageDataConnect
		err = json.Unmarshal(raw["data"], &data)
		m.Data = data
	case "prepare":
		var data MessageDataPrepare
		err = json.Unmarshal(raw["data"], &data)
		m.Data = data
	case "ready":
		var data MessageDataReady
		err = json.Unmarshal(raw["data"], &data)
		m.Data = data
	case "start":
		var data MessageDataStart
		err = json.Unmarshal(raw["data"], &data)
		m.Data = data
	case "code":
		var data MessageDataCode
		err = json.Unmarshal(raw["data"], &data)
		m.Data = data
	case "score":
		var data MessageDataScore
		err = json.Unmarshal(raw["data"], &data)
		m.Data = data
	case "finish":
		var data MessageDataFinish
		err = json.Unmarshal(raw["data"], &data)
		m.Data = data
	case "watch":
		var data MessageDataWatch
		err = json.Unmarshal(raw["data"], &data)
		m.Data = data
	default:
		err = fmt.Errorf("unknown message type: %s", m.Type)
	}

	return err
}
