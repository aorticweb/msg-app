package crud

import (
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

func FindGroup(db *gorm.DB, groupname string) (*Group, error) {
	var groups []Group
	err := db.Where("groupname = ?", groupname).Limit(1).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, nil
	}
	return &groups[0], nil
}

func GroupExists(db *gorm.DB, groupname string) (bool, error) {
	group, err := FindGroup(db, groupname)
	return group != nil, err
}

func CreateGroup(db *gorm.DB, groupname string, users []User) error {
	return db.Transaction(func(tx *gorm.DB) error {
		group := Group{Groupname: groupname}
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
}
