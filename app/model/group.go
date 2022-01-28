package model

type GroupPost struct {
	Groupname string   `json:"groupname" validate:"required"`
	Usernames []string `json:"usernames" validate:"required"`
}
