// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"database/sql"
	"git.dustess.com/mk-base/log"
	articleHttpDelivery "github.com/bxcodec/go-clean-arch/internal/article/delivery/http"
	articleRepo "github.com/bxcodec/go-clean-arch/internal/article/repository/mysql"
	articleUcase "github.com/bxcodec/go-clean-arch/internal/article/usecase"
	authorRepo "github.com/bxcodec/go-clean-arch/internal/author/repository/mysql"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"time"
)

// InitApp init kratos application.
func InitApp(*gin.Engine, *sql.DB, *log.LoggerTrace, time.Duration) error {
	panic(wire.Build(articleHttpDelivery.ProviderSet, articleUcase.ProviderSet, authorRepo.ProviderSet, articleRepo.ProviderSet))
}
