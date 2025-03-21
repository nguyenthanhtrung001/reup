package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	Phone          string             `bson:"phone"`
	Password       string             `bson:"password"`
	FullName       string             `bson:"fullname"`
	Verified       bool               `bson:"verified"`
	GroupID        primitive.ObjectID `bson:"group_id,omitempty"`
	GroupRole      string             `bson:"group_role,omitempty"`
	GroupName      string             `bson:"group_name,omitempty"`
	BlackListToken *[]string          `bson:"black_list_token,omitempty"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	CreatedAt      time.Time          `bson:"created_at"`
	DeletedAt      *time.Time         `bson:"deleted_at,omitempty"`
	RefererID      primitive.ObjectID `bson:"referer_id,omitempty"`
	Credit         float64            `bson:"credit"`
	CreditUsed     float64            `bson:"credit_used"`
	Discount       map[int]float64    `bson:"discount"`
	LastOrderAt    *time.Time         `bson:"last_order_at,omitempty"`
	ApiKey         string             `bson:"api_key,omitempty"`
	Priority       bool               `bson:"priority"`
}

var ListGroupRole = map[string]string{
	"admin":    "admin",
	"reseller": "reseller",
	"user":     "user",
}

func (u User) IsAdmin() bool {
	return u.GroupRole == "admin"
}
func (u User) IsUser() bool {
	return u.GroupRole == "user"
}

func (u User) IsReseller() bool {
	return u.GroupRole == "reseller"
}

func (u User) IsSuperAdmin() bool {
	return u.GroupRole == "superadmin"
}

type DefaultGroup struct {
	GroupID   string
	GroupRole string
	GroupName string
}

func (u User) GetDefaultGroup() DefaultGroup {
	return DefaultGroup{
		GroupID:   "60f445259d83d3bdd84dfc99",
		GroupRole: "user",
		GroupName: "User",
	}
}
