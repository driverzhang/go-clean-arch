package http

import (
	"git.dustess.com/mk-base/log"
	"github.com/bxcodec/go-clean-arch/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewAuthorHandler)

type AuthorHandler struct {
	log       *log.LoggerTrace
	AtUsecase domain.AuthorUsecase
}

func NewAuthorHandler(e *gin.Engine, us domain.AuthorUsecase, log *log.LoggerTrace) error {
	handler := &AuthorHandler{
		log:       log,
		AtUsecase: us,
	}
	e.GET("/author", handler.GetAuthorById)
	return nil
}

// GetAuthorById
func (a *AuthorHandler) GetAuthorById(ctx *gin.Context) {
	data, err := a.AtUsecase.GetByID(ctx, 1)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{"msg": "getById error"})
		return
	}

	ctx.JSON(http.StatusOK, data)
}
