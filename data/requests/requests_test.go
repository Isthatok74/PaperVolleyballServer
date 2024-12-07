package requests

import (
	"strconv"
	"testing"

	"github.com/Isthatok74/PaperVolleyballServer/data/structures"
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

func TestSerializeCreateRequest(t *testing.T) {
	rq := CreateGameRequest{
		GameID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq CreateGameRequest) string { return rq.GameID })
}
