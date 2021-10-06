package main

import (
	"database/sql"
	"fmt"
	log2 "git.dustess.com/mk-base/log"
	"log"

	"net/url"
	"time"

	"git.dustess.com/mk-base/gin-ext/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"

	_articleHttpDelivery "github.com/bxcodec/go-clean-arch/internal/article/delivery/http"
	_articleRepo "github.com/bxcodec/go-clean-arch/internal/article/repository/mysql"
	_articleUcase "github.com/bxcodec/go-clean-arch/internal/article/usecase"
	_authorRepo "github.com/bxcodec/go-clean-arch/internal/author/repository/mysql"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	r := gin.New()
	// todo: 各种中间件
	r.Use(middleware.CORS())

	// todo: 依赖注入wire
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)
	logger := log2.StartNameTrace("")
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	// 这里可以有多组 多个领域实体
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext, logger)
	_articleHttpDelivery.NewArticleHandler(r, au, logger)

	log.Fatal(r.Run(viper.GetString("server.address")))
}
