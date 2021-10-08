package main

import (
	"database/sql"
	"fmt"
	log2 "git.dustess.com/mk-base/log"
	"github.com/bxcodec/go-clean-arch/internal/article/delivery/http"
	"github.com/bxcodec/go-clean-arch/internal/article/repository/mysql"
	"github.com/bxcodec/go-clean-arch/internal/article/usecase"
	mysql2 "github.com/bxcodec/go-clean-arch/internal/author/repository/mysql"
	"log"

	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
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
	// r.Use(middleware.CORS())

	logger := log2.StartNameTrace("main") // 可以替换log
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	err = InitApp2(r, dbConn, logger, timeoutContext)
	if err != nil {
		panic(err)
	}
	logger.Fatal(r.Run(viper.GetString("server.address")))
}

func InitApp2(engine *gin.Engine, db *sql.DB, loggerTrace *log2.LoggerTrace, duration time.Duration) error {
	articleRepository := mysql.NewMysqlArticleRepository(db, loggerTrace)
	authorRepository := mysql2.NewMysqlAuthorRepository(db, loggerTrace)
	articleUsecase := usecase.NewArticleUsecase(articleRepository, authorRepository, duration, loggerTrace)
	error2 := http.NewArticleHandler(engine, articleUsecase, loggerTrace)
	return error2
}
