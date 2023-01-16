package models

type User struct {
	UserID  string `json:"UserID"`
	PetName string `json:"petName"`
	Point   int    `json:"point"`
	Config  struct {
		Noti bool `json:"noti"`
	} `json:"config"`
	Posts    []string `json:"Post"`
	Comments []string `json:"comments"`
	Likes    []string `json:"likes"`
}
