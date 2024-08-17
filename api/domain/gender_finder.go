package domain

import "strings"

type GenderFinder struct {
	Repo    GenderLabelRepository
	Name    string
	Country string
}

func NewGenderFinder(repo GenderLabelRepository, name string, country string) *GenderFinder {
	return &GenderFinder{Repo: repo, Name: name, Country: country}
}

func (finder GenderFinder) FindByNameAndCountry() (*GenderCountryLabel, error) {
	return finder.Repo.FindByNameAndCountry(
		strings.ToLower(finder.Name),
		strings.ToUpper(finder.Country),
	)
}

func (finder GenderFinder) FindByName() (*GenderLabel, error) {
	return finder.Repo.FindByName(strings.ToLower(finder.Name))
}
