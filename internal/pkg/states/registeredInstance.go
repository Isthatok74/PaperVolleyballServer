package states

import (
	"sync"

	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/util"
)

// Base struct for instances containing live players
type RegisteredInstance struct {
	ExpirableInstance
	Players sync.Map `json:"Players"` // key: string; value: dummy flag (boolean)
	HostID  string   `json:"HostID"`  // the id of the hosting player
}

// create a clone of the stored instance
func (r *RegisteredInstance) Clone() *RegisteredInstance {
	retVal := &RegisteredInstance{
		Players: *util.CopySyncMap(&r.Players),
	}
	retVal.ExpirableInstance.LastUpdate = r.LastUpdate
	return retVal
}
