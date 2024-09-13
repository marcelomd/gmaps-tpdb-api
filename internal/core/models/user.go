package models

type User struct {
	Id           string
	Role         string
	Name         string
	Email        string
	PasswordHash []byte
}


type NewUser struct {
	Role     string
	Name     string
	Email    string
	Password string
}