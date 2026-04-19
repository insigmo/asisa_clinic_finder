package db

type User struct {
	ID           int64
	Username     string
	Name         string
	Lastname     string
	IsBot        bool
	City         string
	LanguageCode string
}
