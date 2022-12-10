package models

import (
	"time"
)

type Dong struct {
	OwnerID     int64
	Size        int
	GeneratedAt time.Time
}
