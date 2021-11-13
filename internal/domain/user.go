package domain

import "time"

type User struct {
	Object
	Surname     string
	Name        string
	Patronymic  string
	DateOfBirth time.Time
	PhoneNumber string
	Email       string
	CityID      ID
}
