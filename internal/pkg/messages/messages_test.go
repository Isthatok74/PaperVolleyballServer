package messages

import (
	"strconv"
	"testing"

	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/structures"
)

func TestSerializePingRequest(t *testing.T) {
	rq := PingMessage{
		PingTime: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq PingMessage) string { return rq.PingTime })
}

func TestSerializeAdmitRequest(t *testing.T) {
	rq := AdmissionMessage{
		ClientPlayerID: 42,
		ServerPlayerID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq AdmissionMessage) string { return strconv.Itoa(rq.ClientPlayerID) + rq.ServerPlayerID })
}

func TestSerializeAddPlayerGameRequest(t *testing.T) {
	rq := AddPlayerGameMessage{
		GameID:         "xyzguid",
		ServerPlayerID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq AddPlayerGameMessage) string { return rq.GameID + rq.ServerPlayerID })
}

func TestSerializeAddPlayerLobbyRequest(t *testing.T) {
	rq := AddPlayerLobbyMessage{
		RoomCode:       "QBPX",
		ServerPlayerID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq AddPlayerLobbyMessage) string { return rq.RoomCode + rq.ServerPlayerID })
}

func TestSerializeCreateGameRequest(t *testing.T) {
	rq := CreateGameMessage{
		GameID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq CreateGameMessage) string { return rq.GameID })
}

func TestSerializeCreateLobbyRequest(t *testing.T) {
	rq := CreateLobbyMessage{
		ErrMsg:   "",
		RoomCode: "JXPQ",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq CreateLobbyMessage) string { return rq.ErrMsg + rq.RoomCode })
}
