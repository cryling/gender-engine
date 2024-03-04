package infrastructure

import (
	"log"
	"strings"

	"github.com/cryling/gender-engine/domain"
)

type MemoryHandler struct {
	storage map[string]string
}

func NewMemoryHandler(data map[string]string) MemoryHandler {
	return MemoryHandler{storage: data}
}

func (handler MemoryHandler) FindByName(name string) (*domain.GenderLabel, error) {
	result, ok := handler.storage[strings.ToLower(name)]
	if !ok {
		log.Printf("Name not found")
		return &domain.GenderLabel{}, &domain.NotFoundError{Name: name}
	}

	return &domain.GenderLabel{Name: name, Gender: result}, nil
}
