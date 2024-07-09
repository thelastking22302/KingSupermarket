package usermodels

type Roles int

const (
	ADMIN Roles = iota
	MEMBER
)

func (r Roles) String() string {
	return []string{"ADMIN", "MEMBER"}[r]
}
