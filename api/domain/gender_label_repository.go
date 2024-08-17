package domain

type GenderLabelRepository interface {
	FindByNameAndCountry(name string, country string) (*GenderCountryLabel, error)
	FindByName(name string) (*GenderLabel, error)
}
