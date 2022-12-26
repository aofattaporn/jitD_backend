package models

type Book struct {
	Auther     string `json:"auther"`
	Name       string `json:"name"`
	Characters []string 
}
