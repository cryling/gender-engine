package domain

type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return e.Name + " could not be found or is not clearly assignable to a single gender"
}
