package models

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"userName"`
}

// TableName Database Table Name of this model
// func (e *Example) TableName() string {
// 	return "examples"
// }
