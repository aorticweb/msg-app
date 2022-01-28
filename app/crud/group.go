package crud

import (
	"errors"

	"gorm.io/gorm"
)

type Group struct {
	ID        int64  `gorm:"column:id;type:bigserial;primary_key" json:"-"`
	Groupname string `gorm:"column:groupname;type:varchar(240);unique" json:"groupname"`
}

func (g *Group) TableName() string {
	return "public.group"
}

type UserGroup struct {
	ID      int64 `gorm:"column:id;type:bigserial;primary_key" json:"-"`
	GroupID int64 `gorm:"column:group_id;integer"`
	Group   Group `gorm:"foreignKey:group_id"`
	UserID  int64 `gorm:"column:user_id;integer"`
	User    User  `gorm:"foreignKey:user_id"`
}

func (u *UserGroup) TableName() string {
	return "public.user_group"
}

func FindGroup(db *gorm.DB, groupname string) (*Group, bool, error) {
	var group Group
	err := db.Where("groupname = ?", groupname).First(&group).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &group, true, nil
}

func GroupExists(db *gorm.DB, groupname string) (bool, error) {
	_, exist, err := FindGroup(db, groupname)
	return exist, err
}

func CreateGroup(db *gorm.DB, groupname string, users []User) (*Group, error) {
	group := Group{Groupname: groupname}
	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&group)
		if result.Error != nil {
			return result.Error
		}
		var userGroups []UserGroup
		for _, user := range users {
			userGroups = append(userGroups, UserGroup{Group: group, User: user})
		}
		result = tx.Create(&userGroups)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &group, nil
}
