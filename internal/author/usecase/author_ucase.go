package usecase

import (
	"context"
	"git.dustess.com/mk-base/log"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/google/wire"
	"time"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewAuthorUsecase)

type authoreUsecase struct {
	authorRepo     domain.AuthorRepository
	contextTimeout time.Duration
	log            *log.LoggerTrace
}

// NewAuthorUsecase will create new an authoreUsecase object representation of domain.AuthorUsecase interface
func NewAuthorUsecase(ar domain.AuthorRepository, timeout time.Duration, log *log.LoggerTrace) domain.AuthorUsecase {
	return &authoreUsecase{
		authorRepo:     ar,
		contextTimeout: timeout,
		log:            log,
	}
}

func (a *authoreUsecase) GetByID(c context.Context, id int64) (author domain.Author, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	a.log.Info("GetByID test")
	return a.authorRepo.GetByID(ctx, id)
}
