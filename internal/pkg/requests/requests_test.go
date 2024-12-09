package requests

import (
	"strconv"
	"testing"

	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/structures"
)

func TestSerializePingRequest(t *testing.T) {
	rq := PingRequest{
		PingTime: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq PingRequest) string { return rq.PingTime })
}

func TestSerializeAdmitRequest(t *testing.T) {
	rq := AdmissionRequest{
		ClientPlayerID: 42,
		ServerPlayerID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq AdmissionRequest) string { return strconv.Itoa(rq.ClientPlayerID) + rq.ServerPlayerID })
}

func TestSerializeAddPlayerGameRequest(t *testing.T) {
	rq := AddPlayerGameRequest{
		GameID:         "xyzguid",
		ServerPlayerID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq AddPlayerGameRequest) string { return rq.GameID + rq.ServerPlayerID })
}

func TestSerializeAddPlayerLobbyRequest(t *testing.T) {
	rq := AddPlayerLobbyRequest{
		RoomCode:       "QBPX",
		ServerPlayerID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq AddPlayerLobbyRequest) string { return rq.RoomCode + rq.ServerPlayerID })
}

func TestSerializeCreateGameRequest(t *testing.T) {
	rq := CreateGameRequest{
		GameID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq CreateGameRequest) string { return rq.GameID })
}

func TestSerializeCreateLobbyRequest(t *testing.T) {
	rq := CreateLobbyRequest{
		ErrMsg:   "",
		RoomCode: "JXPQ",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq CreateLobbyRequest) string { return rq.ErrMsg + rq.RoomCode })
}
