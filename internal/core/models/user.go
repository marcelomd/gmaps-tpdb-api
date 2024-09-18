package models

type User struct {
    Id           string
    Name         string
    Email        string
    Role         Role
    PasswordHash []byte
}

type NewUser struct {
    Name     string
    Email    string
    Role     Role
    Password string
}
