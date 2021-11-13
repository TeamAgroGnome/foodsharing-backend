package domain

type Permission uint64

const (
	Admin Permission = 1 << iota

	CreateUser
	ReadUser
	EditUser

	CreateAct
	ReadAct
	EditAct

	AddCity
	ReadCity
	EditCity

	AddCompany
	ReadCompany
	EditCompany

	CreateGroup
	ReadGroup
	EditGroup
)

func (p Permission) IsAdmin() bool {
	return p&Admin != 0
}

func (p Permission) canCreateUser() bool {
	return p&CreateUser != 0
}

func (p Permission) canReadUser() bool {
	return p&ReadUser != 0
}

func (p Permission) canEditUser() bool {
	return p&EditUser != 0
}

func (p Permission) canCreateAct() bool {
	return p&CreateAct != 0
}

func (p Permission) canReadAct() bool {
	return p&ReadAct != 0
}

func (p Permission) canEditAct() bool {
	return p&EditAct != 0
}

func (p Permission) canAddCity() bool {
	return p&AddCity != 0
}

func (p Permission) canReadCity() bool {
	return p&ReadCity != 0
}

func (p Permission) canEditCity() bool {
	return p&EditCity != 0
}

func (p Permission) canAddCompany() bool {
	return p&AddCompany != 0
}

func (p Permission) canReadCompany() bool {
	return p&ReadCompany != 0
}

func (p Permission) canEditCompany() bool {
	return p&EditCompany != 0
}

func (p Permission) canCreateGroup() bool {
	return p&CreateGroup != 0
}

func (p Permission) canReadGroup() bool {
	return p&ReadGroup != 0
}

func (p Permission) canEditGroup() bool {
	return p&EditGroup != 0
}
