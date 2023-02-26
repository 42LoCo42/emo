package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"unique;notNull"`
	PublicKey []byte `gorm:"unique;notNull"`
	Admin     bool
	Songs     []*Song `gorm:"many2many:user_songs;"`
}

type Song struct {
	gorm.Model
	Path     string `gorm:"unique;notNull"`
	Name     string `gorm:"notNull"`
	Artist   string
	Duration time.Duration
	Users    []*User `gorm:"many2many:user_songs;"`
}
