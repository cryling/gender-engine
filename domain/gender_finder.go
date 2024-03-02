package domain

type GenderFinder struct {
	Repo GenderLabelRepository
	Name string
}

func NewGenderFinder(repo GenderLabelRepository, name string) GenderFinder {
	return GenderFinder{Repo: repo, Name: name}
}

func (finder GenderFinder) Find() (*GenderLabel, error) {
	return finder.Repo.FindByName(finder.Name)
}
