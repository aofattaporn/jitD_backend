package models

type User struct {
	TokenID string `json:"TokenId"`
	PetName string `json:"petName"`
	Point   int    `json:"point"`
	Config  struct {
		Noti bool `json:"noti"`
	} `json:"config"`
	Posts    []string `json:"Post"`
	Comments []string `json:"comments"`
	Likes    []string `json:"likes"`
}
