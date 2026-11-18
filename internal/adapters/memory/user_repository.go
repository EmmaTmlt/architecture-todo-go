package memory

import (
	"context"
	"sync"

	"github.com/EmmaTmlt/architecture-todo-go/internal/domain"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[domain.Username]domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: make(map[domain.Username]domain.User)}
}

func (r *UserRepository) Save(ctx context.Context, user domain.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.users[user.Username];
	if exists {
		return domain.ErrUsernameTaken
	}

	r.users[user.Username] = user
	return nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username domain.Username) (domain.User, error) {
	if err := ctx.Err(); err != nil {
		return domain.User{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[username]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}
