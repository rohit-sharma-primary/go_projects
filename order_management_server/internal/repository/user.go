package repository

import (
	"sync"

	"server/internal/model"
)

type UserRepository struct {
	mu     sync.RWMutex
	nextID int
	users  map[int]model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		nextID: 1,
		users:  make(map[int]model.User),
	}
}

func (r *UserRepository) Create(user model.User) model.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.nextID++
	r.users[user.ID] = user

	return user
}

func (r *UserRepository) GetByID(id int) (model.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	return user, ok
}

func (r *UserRepository) List() []model.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users
}
