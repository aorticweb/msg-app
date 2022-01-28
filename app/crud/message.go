package crud

import "time"

type Message struct {
	ID        int64     `gorm:"column:id;type:bigserial;primary_key" json:"id"`
	RE        *Message  `gorm:"foreignKey:id"`
	Sender    User      `gorm:"foreignKey:id"`
	Recipient Group     `gorm:"foreignKey:id"`
	Subject   string    `gorm:"column:subject;type:text;" json:"subject"`
	Body      string    `gorm:"column:body;type:text;" json:"body"`
	SentAt    time.Time `gorm:"column:sent_at;type:timestamp with time zone;" json:"sentAt"`
}
