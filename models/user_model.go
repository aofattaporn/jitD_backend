package models

type User struct {
	TokenID  string `json:"TokenId"`
	PetName  string `json:"petName"`
	PetName0 string `json:"PetName"`
	Point    int    `json:"point"`
	Config   struct {
		Noti bool `json:"noti"`
	} `json:"config"`
	Post     []string `json:"Post"`
	Comments []string `json:"comments"`
	Likes    []string `json:"likes"`
}

// TableName Database Table Name of this model
// func (e *Example) TableName() string {
// 	return "examples"
// }
