package requests

// represents a game ball, along with all of its state variables
type CreateRequest struct {
	GameID string `json:"GameID"`
}
