package domain

import "strings"

type GenderFinder struct {
	Repo    GenderLabelRepository
	Name    string
	Country string
}

func NewGenderFinder(repo GenderLabelRepository, name string, country string) GenderFinder {
	return GenderFinder{Repo: repo, Name: name, Country: country}
}

func (finder GenderFinder) Find() (*GenderLabel, error) {
	return finder.Repo.FindByNameAndCountry(strings.ToLower(finder.Name), finder.Country)
}
