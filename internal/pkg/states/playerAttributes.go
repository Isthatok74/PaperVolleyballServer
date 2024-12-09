package states

// represents the constant variables of a player that is connected to the game
type PlayerAttributes struct {
	DisplayName string  `json:"DisplayName"`
	Strength    float32 `json:"Strength"`
	Speed       float32 `json:"Speed"`
	Jump        float32 `json:"Jump"`
	Size        float32 `json:"Size"`
	Tier        int     `json:"Tier"`
}
