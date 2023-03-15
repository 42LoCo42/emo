package main

import (
	"time"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string
	PublicKey []byte `gorm:"unique;notNull"`
	Admin     bool
	Songs     []*Song `gorm:"many2many:user_songs;" json:"-"`
}

type Song struct {
	File     string `gorm:"primaryKey"`
	Name     string `gorm:"unique;notNull"`
	Artist   string
	Duration time.Duration
	Users    []*User `gorm:"many2many:user_songs;" json:"-"`
}
