package models

type Book struct {
	// Auther     string `json:"auther"`
	Name      string   `json:"name"`
	Character []string `json:"character"`
	Seller    []string `json:"seller"`
}
