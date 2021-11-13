package domain

import "time"

type DonorCompany struct {
	Object
	Name           string
	CityID         ID
	ContractDate   time.Time
	ContractNumber int
}
