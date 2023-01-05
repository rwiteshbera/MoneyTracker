package models

type User struct {
	UserId      string `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber uint64 `json:"phoneNumber"`
	Password    string `json:"password"`
}
