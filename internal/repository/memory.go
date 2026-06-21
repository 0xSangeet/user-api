package repository

import (
	"context"
	"sync"

	"github.com/0xSangeet/user-api/internal/domain"
)

type MemoryRepo struct {
	users map[string]domain.User
	mu    sync.RWMutex
}

func NewMemoryRepository() *MemoryRepo {
	return &MemoryRepo{
		users: make(map[string]domain.User),
	}
}

func (m *MemoryRepo) Create(ctx context.Context, u *domain.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[u.ID]; ok {
		return domain.ErrUserAlreadyExists
	}

	m.users[u.ID] = *u

	return nil
}

func (m *MemoryRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[id]

	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return &user, nil
}

func (m *MemoryRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]domain.User, 0, len(m.users))

	for _, u := range m.users {
		users = append(users, u)
	}

	return users, nil
}

func (m *MemoryRepo) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[id]; !ok {
		return domain.ErrUserNotFound
	}

	delete(m.users, id)
	return nil
}

