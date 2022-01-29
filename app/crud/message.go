package crud

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID          int64     `gorm:"column:id;type:bigserial;primary_key"`
	REID        *int64    `gorm:"column:re_id;integer"`
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
	err := db.Preload("Sender").Preload("Recipient").Preload("Group").Where("id = ?", messageID).First(&msg).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &msg, true, nil
}

func GetMessageReplies(db *gorm.DB, messageID int64) ([]Message, error) {
	var msgs []Message
	query := db.Preload("Sender").Preload("Recipient").Preload("Group")
	err := query.Where("re_id = ?", messageID).Order("sent_at DESC").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func GetUserMailbox(db *gorm.DB, userID int64) ([]Message, error) {
	// Query
	// Get Messages where RecipientID = <UserID>
	// Get Messages Where GroupID IN (Get UserGroup where UserID = <UserID>)
	groupIDs, err := FindGroupsByUserID(db, userID)
	if err != nil {
		return nil, err
	}
	var msgs []Message
	query := db.Preload("Sender").Preload("Recipient").Preload("Group")
	err = query.Where("(recipient_id = ? or group_id in ?)", userID, groupIDs).Order("sent_at desc").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil

}
