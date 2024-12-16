package states

// represents the constant variables of a player that is connected to the game
type PlayerAttributes struct {

	// display name
	DisplayName string `json:"DisplayName"`

	// player stat levels
	Strength float32 `json:"Strength"`
	Speed    float32 `json:"Speed"`
	Jump     float32 `json:"Jump"`
	Size     float32 `json:"Size"`
	Tier     int     `json:"Tier"`

	// visual accessory codes
	Emblem     string `json:"Emblem"`
	Hair       string `json:"Hair"`
	LowerFace  string `json:"LowerFace"`
	Expression string `json:"Expression"`
}
