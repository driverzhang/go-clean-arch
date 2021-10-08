package mysql

import (
	"context"
	"database/sql"
	"git.dustess.com/mk-base/log"
	"github.com/google/wire"
	
	"github.com/bxcodec/go-clean-arch/domain"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewMysqlAuthorRepository)

type mysqlAuthorRepo struct {
	DB *sql.DB
	log *log.LoggerTrace
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewMysqlAuthorRepository(db *sql.DB, log *log.LoggerTrace) domain.AuthorRepository {
	return &mysqlAuthorRepo{
		DB: db,
		log: log,
	}
}

func (m *mysqlAuthorRepo) getOne(ctx context.Context, query string, args ...interface{}) (res domain.Author, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Author{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res = domain.Author{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return
}

func (m *mysqlAuthorRepo) GetByID(ctx context.Context, id int64) (domain.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return m.getOne(ctx, query, id)
}
