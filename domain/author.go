package domain

import "context"

// Author ...
type Author struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type AuthorUsecase interface {
	GetByID(ctx context.Context, id int64) (author Author, err error)
}

// AuthorRepository represent the author's repository contract
type AuthorRepository interface {
	GetByID(ctx context.Context, id int64) (Author, error)
}
