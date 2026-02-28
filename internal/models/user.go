package models

import "time"

type User struct {
	Id        string    `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	RoleId    int       `json:"role_id"`
	Role      UserRole  `json:"role" gorm:"foreignkey:RoleId"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRole struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserGetRequest struct {
	Id string `uri:"id" validate:"required"`
}
