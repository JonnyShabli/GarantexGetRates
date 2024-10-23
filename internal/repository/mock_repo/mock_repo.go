package mock_repo

import (
	"context"
	"sync"

	"github.com/JonnyShabli/GarantexGetRates/internal/models"
)

// Fake repository сохраняет данные в мапу вместо БД
type MockRepo struct {
	data map[uint]models.RatesToDB
	mu   sync.Mutex
	Idx  uint
}

func NewMockRepo() *MockRepo {
	return &MockRepo{
		data: make(map[uint]models.RatesToDB),
		mu:   sync.Mutex{},
		Idx:  0,
	}
}

func (m *MockRepo) InsertRates(ctx context.Context, db models.RatesToDB) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Idx++
	m.data[m.Idx] = db
	return nil
}

func GetRates(storage *MockRepo, id uint) models.RatesToDB {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	return storage.data[id]
}
