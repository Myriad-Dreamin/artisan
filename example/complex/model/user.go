package model

import (
	"time"
)

type User struct {
	ID        uint      `dorm:"id" gorm:"column:id;primary_key;not_null"`
	CreatedAt time.Time `dorm:"created_at" gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null" json:"created_at"`
	UpdatedAt time.Time `dorm:"updated_at" gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null;" json:"updated_at"`
	LastLogin time.Time `dorm:"last_login" gorm:"column:last_login;default:CURRENT_TIMESTAMP;not null;" json:"last_login"`

	NickName string `dorm:"nick_name" gorm:"column:nick_name;unique;not_null"`
	Name     string `dorm:"name" gorm:"column:name;not_null"`
	Password string `dorm:"password" gorm:"column:password;not_null"`
	Phone    string `dorm:"phone" gorm:"column:phone;unique;not_null"`
}
