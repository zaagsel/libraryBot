package main

import "time"

type User struct {
    ID        int64  `gorm:"primaryKey"`
    Nickname  string `gorm:"not null"`
    UserName  string
    Verify    bool      `gorm:"default:false"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Genre struct {
    ID    int64  `gorm:"primaryKey"`
    Genre string `gorm:"not null"`
}

type Status struct {
    ID     int64  `gorm:"primaryKey"`
    Status string `gorm:"not null"`
}

type Book struct {
    ID        int64  `gorm:"primaryKey"`
    UserID    int64  `gorm:"not null"` 
    Book      string `gorm:"not null"`
    Author    string `gorm:"not null"`
    GenreID   int64  `gorm:"not null"` 
    StatusID  int64  `gorm:"not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}