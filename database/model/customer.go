package model

import "time"

type Customer struct {
	ID             uint      `gorm:"column:id;type;int unsigned;not null;primaryKey;autoIncrement" json:"id"`
	IdentityNumber string    `gorm:"column:identity_number;type:varchar(256);unique;default:null" json:"identity_number"`
	FullName       string    `gorm:"column:full_name;type:text;default:null" json:"full_name"`
	Status         string    `gorm:"column:status;type:varchar(20);default:null" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;not null;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;not null;autoCreateTime;autoUpdateTime" json:"updated_at"`
}

func (c *Customer) TableName() string {
	return "customers"
}
