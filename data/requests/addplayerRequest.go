package requests

// represents a game ball, along with all of its state variables
type AddPlayerRequest struct {
	ClientPlayerID int    `json:"ClientPlayerID"`
	ServerPlayerID string `json:"ServerPlayerID"`
}
