package model

type LoginData struct {
	Login    string
	Password string
}

type RegData struct {
	Login    string
	Password string
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

type MetaListData struct {
	UUID string
	Type string
	Name string
	File string
	Date string
}
