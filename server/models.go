package main

import (
	"time"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"unique;notNull"`
	PublicKey []byte `gorm:"unique;notNull"`
	Admin     bool
	Songs     []Song `gorm:"many2many:stats;" json:"-"`
}

type Song struct {
	ID       string `gorm:"primaryKey"`
	Name     string `gorm:"unique;notNull"`
	Artist   string
	Duration time.Duration
}

type Stat struct {
	UserID string `gorm:"primaryKey"`
	SongID string `gorm:"primaryKey"`
	Count  uint64
	Boost  uint64
}

type StatQuery struct {
	ID    string
	Name  string
	Count uint64
	Boost uint64
}
