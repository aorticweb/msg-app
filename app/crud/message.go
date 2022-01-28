package crud

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID          int64     `gorm:"column:id;type:bigserial;primary_key"`
	REID        *int64    `gorm:"column:re_id;integer"`
	RE          *Message  `gorm:"foreignKey:re_id"`
	SenderID    *int64    `gorm:"column:sender_id;integer"`
	Sender      *User     `gorm:"foreignKey:sender_id"`
	RecipientID *int64    `gorm:"column:recipient_id;integer"`
	Recipient   *User     `gorm:"foreignKey:recipient_id"`
	GroupID     *int64    `gorm:"column:group_id;integer"`
	Group       *Group    `gorm:"foreignKey:group_id"`
	Subject     string    `gorm:"column:subject;type:text;" json:"subject"`
	Body        string    `gorm:"column:body;type:text;" json:"body"`
	SentAt      time.Time `gorm:"column:sent_at;type:timestamp with time zone;" json:"sentAt"`
}

func (m *Message) TableName() string {
	return "public.message"
}
func CreateMessage(db *gorm.DB, message *Message) (*Message, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&message)
		return result.Error
	})
	return message, err
}

func GetMessage(db *gorm.DB, messageID int64) (*Message, bool, error) {
	var msg Message
	err := db.Where("id = ?", messageID).First(&msg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}
