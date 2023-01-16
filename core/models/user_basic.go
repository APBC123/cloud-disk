package models

import "time"

type UserBasic struct {
	Id          int
	Identity    string
	Name        string
	Password    string
	Email       string
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
	DeletedAt   time.Time `xorm:"deleted"`
	NowVolume   int64     `xorm:"now_volume"`
	TotalVolume int64     `xorm:"total_volume"`
}

func (table UserBasic) TableName() string {
	return "user_basic"

}
