package model

type LoginData struct {
	Login    string
	Password string
}

type RegData struct {
	Login           string
	Password        string
	ConfirmPassword string
}

type CredsData struct {
	Name     string
	Login    string
	Password string
}

type BankData struct {
	Name   string
	Number string
	Exp    string
	CVV    string
}

type FileData struct {
	Name string
	Path string
}
