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

func TestSerializePlayerRequest(t *testing.T) {
	rq := AddPlayerRequest{
		ClientPlayerID: 42,
		ServerPlayerID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq AddPlayerRequest) string { return strconv.Itoa(rq.ClientPlayerID) + rq.ServerPlayerID })
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
