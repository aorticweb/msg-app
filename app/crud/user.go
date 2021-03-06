package crud

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	ID       int64  `gorm:"column:id;type:bigserial;primary_key" json:"-"`
	Username string `gorm:"column:username;type:varchar(240);unique" json:"username" validate:"required"`
}

func (c *User) TableName() string {
	return "public.user"
}

func FindUsers(db *gorm.DB, usernames []string) ([]User, error) {
	var users []User
	err := db.Where("username in ?", usernames).Find(&users).Error
	if err != nil {
		return []User{}, err
	}
	return users, nil
}

func FindUser(db *gorm.DB, username string) (*User, bool, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &user, true, nil
}

func UserExist(db *gorm.DB, username string) (bool, error) {
	var users []User
	err := db.Where("username = ?", username).Limit(1).Find(&users).Error
	if err != nil {
		return false, err
	}
	return len(users) == 1, nil
}

func CreateUser(db *gorm.DB, user User) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&user)
		return result.Error
	})
}
