package model

import "time"

type SessionCache struct {
	SessionID    string    `json:"-"`
	AccountID    string    `json:"account_id"`
	MachineID    string    `json:"machine_id"`
	ClientIP     string    `json:"client_ip"`
	LastActivity time.Time `json:"last_activity"`
}
