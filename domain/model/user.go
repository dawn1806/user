package model

type User struct {
	// 用户ID
	ID int64 `gorm:"primary_key;not_null;auto_increment"`
	// 用户名
	UserName string `gorm:"not_null；unique_index"`
	// 添加需要的字段
	FirstName string
	// 密码
	HashPassword string
}
