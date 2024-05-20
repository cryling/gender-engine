package domain

type GenderLabelRepository interface {
	FindByNameAndCountry(name string, country string) (*GenderLabel, error)
}
