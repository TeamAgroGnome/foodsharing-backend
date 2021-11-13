package domain

type Group struct {
	Object
	Name        string
	Permissions Permission
}
