package models

type Compound struct {
    Id        string
    Type      string
    Class     string
    Name      string
    Mass      string
    Formula   string
    Fragments []Fragment
}

type Fragment struct {
    Id      string
    Name    string
    Mass    string
    Formula string
}
