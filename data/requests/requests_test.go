package requests

import (
	"pv-server/data/structures"
	"testing"
)

func TestSerializePingRequest(t *testing.T) {
	rq := PingRequest{
		PingTime: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq PingRequest) string { return rq.PingTime })
}

func TestSerializePlayerRequest(t *testing.T) {
	rq := PlayerRequest{
		PlayerID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq PlayerRequest) string { return rq.PlayerID })
}

func TestSerializeCreateRequest(t *testing.T) {
	rq := CreateRequest{
		GameID: "anyString",
	}
	structures.CompareSerializeDeserialize(t, rq, func(rq CreateRequest) string { return rq.GameID })
}
