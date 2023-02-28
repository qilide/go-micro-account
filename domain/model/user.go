package model

import (
	"time"
)

type User struct {
	//主键
	ID         int64     `gorm:"primary_key;not_null;auto_increment" json:"id"`
	UserID     int64     `gorm:"unique_index;not_null" json:"user_id"`
	UserName   string    `gorm:"unique_index;not_null" json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	PassWord   string    `json:"password"`
	Permission int64     `json:"permission"`
	CreateDate time.Time `json:"create_date"`
	UpdateDate time.Time `json:"update_date"`
	IsActive   int64     `json:"is_active"`
	Email      string    `json:"email"`
}
