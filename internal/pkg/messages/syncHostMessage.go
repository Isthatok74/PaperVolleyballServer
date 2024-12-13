package messages

// a message that contains information about who the host is in the current lobby
type SyncHostMessage struct {
	HostID string `json:"HostID"`
}
