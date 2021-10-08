package http

import (
	log2 "git.dustess.com/mk-base/log"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
	"net/http"

	"github.com/bxcodec/go-clean-arch/domain"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewArticleHandler)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticleHandler  represent the httphandler for article
type ArticleHandler struct {
	log      *log2.LoggerTrace
	AUsecase domain.ArticleUsecase
}

// NewArticleHandler will initialize the articles/ resources endpoint
func NewArticleHandler(e *gin.Engine, us domain.ArticleUsecase, log *log2.LoggerTrace) error {
	handler := &ArticleHandler{
		AUsecase: us,
		log:      log,
	}
	e.GET("/articles", handler.FetchArticle)
	e.POST("/articles", handler.Store)
	e.GET("/articles/:id", handler.GetByID)
	e.DELETE("/articles/:id", handler.Delete)
	return nil
}

// FetchArticle will fetch the article based on given params
func (a *ArticleHandler) FetchArticle(ctx *gin.Context) {
	listAr, _, err := a.AUsecase.Fetch(ctx, "cursor", 1)
	if err != nil {
		ctx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, listAr)
}

// GetByID will get article by given id
func (a *ArticleHandler) GetByID(ctx *gin.Context) {

	art, err := a.AUsecase.GetByID(ctx, 1)
	if err != nil {
		ctx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, art)
}

func isRequestValid(m *domain.Article) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store will store the article by given request body
func (a *ArticleHandler) Store(ctx *gin.Context) {
	var article domain.Article
	err := ctx.Bind(&article)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var ok bool
	if ok, err = isRequestValid(&article); !ok {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err = a.AUsecase.Store(ctx, &article)
	if err != nil {
		ctx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, article)
}

// Delete will delete article by given param
func (a *ArticleHandler) Delete(ctx *gin.Context) {
	err := a.AUsecase.Delete(ctx, 1)
	if err != nil {
		ctx.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
