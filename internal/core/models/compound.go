package models

type Compound struct {
	Id        string
	Type      string
	Name      string
	Mass      string
	Fragments []Fragment
}

type Fragment struct {
	Id   string
	Name string
	Mass string
}
