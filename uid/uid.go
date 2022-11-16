package uid

import (
	"encoding/json"
	"time"
)

type Uid struct {
	ID         int64      `json:"id"`
	Group      string     `json:"group"`
	Version    int32      `json:"version"`
	CurrentId  int64      `json:"current_id"`
	UpdateTime *time.Time `json:"update_time"`
}

func (u *Uid) String() string {
	bytes, _ := json.Marshal(u)
	return string(bytes)
}
