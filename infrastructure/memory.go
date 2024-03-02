package infrastructure

import "github.com/cryling/gender-engine/domain"

type MemoryHandler struct {
	storage map[string]string
}

func NewMemoryHandler() MemoryHandler {
	return MemoryHandler{storage: make(map[string]string)}
}

func (handler MemoryHandler) FindByName(name string) (*domain.GenderLabel, error) {
	result, err := handler.storage[name]
	if err {
		return &domain.GenderLabel{}, nil
	}

	return &domain.GenderLabel{Name: name, Gender: result}, nil
}
