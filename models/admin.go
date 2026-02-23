package models

import "time"

type Admin struct {
			ID                int
			Email             string
			EncryptedPassword string
			CreatedAt         time.Time
			UpdatedAt         time.Time
}