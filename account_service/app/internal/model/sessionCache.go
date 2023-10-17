package model

import "time"

type SessionCache struct {
	SessionID    string    `json:"-"`
	AccountID    string    `json:"account_id"`
	ClientIP     string    `json:"client_ip"`
	SessionInfo  string    `json:"session_info"` // like device or browser
	LastActivity time.Time `json:"last_activity"`
}
