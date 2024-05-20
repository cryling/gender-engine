package domain

type GenderLabelRepository interface {
	FindByName(name string) (*GenderLabel, error)
}
