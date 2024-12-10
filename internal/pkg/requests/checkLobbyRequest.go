package requests

// a request that checks whether a lobby with the given room code is open
type CheckLobbyRequest struct {
	Exists   bool   `json:"Exists"`
	RoomCode string `json:"RoomCode"`
}
