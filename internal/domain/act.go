package domain

import "time"

type Act struct {
	Object
	UserID         ID
	DonorCompanyID ID
}

type ActContent struct {
	Object
	ActID          ID
	Number         int
	Name           string
	Count          int
	Price          int
	ExpirationDate time.Time
	Comment        string
}
