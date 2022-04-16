package vision

import (
	"log"
	"time"

	router "github.com/fasthttp/router"
	"github.com/patrickmn/go-cache"

	session_http "github.com/perlinleo/vision/internal/session/delivery/http"
	session_redis "github.com/perlinleo/vision/internal/session/repository/redis"
	session_usecase "github.com/perlinleo/vision/internal/session/usecase"

	user_http "github.com/perlinleo/vision/internal/user/delivery/http"
	user_psql "github.com/perlinleo/vision/internal/user/repository/psql"
	user_usecase "github.com/perlinleo/vision/internal/user/usecase"

	status_http "github.com/perlinleo/vision/internal/status/delivery/http"
	status_psql "github.com/perlinleo/vision/internal/status/repository/psql"
	status_usecase "github.com/perlinleo/vision/internal/status/usecase"

	"github.com/valyala/fasthttp"
	// forum_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/forum/delivery"
	// thread_http "github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread/delivery"
)

func Start() error {
	config := NewConfig()

	_, err := NewServer(config)
	if err != nil {
		return err
	}

	PSQLConnPool, err := NewPostgreSQLDataBase(config.App.DatabaseURL)
	if err != nil {
		return err
	}
	log.Printf("PSQL Database connection success on %s", config.App.DatabaseURL)

	RedisClient, err := NewRedisDataBase(config.RedisDB.addr, config.RedisDB.password, config.RedisDB.db)
	if err != nil {
		return err
	}

	log.Printf("Redis Database connection success on %s", config.App.DatabaseURL)

	router := router.New()

	// usERR
	userCache := cache.New(5*time.Minute, 10*time.Minute)
	userRepository := user_psql.NewUserPSQLRepository(PSQLConnPool, userCache)
	userUsecase := user_usecase.NewUserUsecase(userRepository)
	user_http.NewUserHandler(router, userUsecase)

	//auth
	authCache := cache.New(5*time.Minute, 10*time.Minute)
	authRepository := session_redis.NewSessionRedisRepository(&RedisClient, authCache)
	authUsecase := session_usecase.NewSessionUsecase(authRepository, userRepository)
	session_http.NewSessionHandler(router, authUsecase)

	// status
	statusCache := cache.New(5*time.Minute, 10*time.Minute)
	statusRepository := status_psql.NewStatusPSQLRepository(PSQLConnPool, statusCache)
	statusUsecase := status_usecase.NewStatusUsecase(statusRepository)
	status_http.NewStatusHandler(router, statusUsecase)

	log.Printf("STARTING SERVICE ON PORT %s\n", config.App.Port)

	err = fasthttp.ListenAndServe(config.App.Port, router.Handler)
	if err != nil {
		return err
	}

	return nil
}
