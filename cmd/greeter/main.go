package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/leosunmo/zapchi"
	"github.com/locurateam/gitops-greeter/assets"
	"go.uber.org/zap"
)

const (
	redisHostEnvVar   = "REDIS_HOST"
	serverPortEnvVar  = "SERVER_PORT"
	environmentEnvVar = "ENVIRONMENT"
	redisDB           = 0
	redisPeopleKey    = "people"
)

func mustHaveEnvVariable(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Errorf("could not get the required env variable: %s", key))
	}
	return value
}

func main() {
	redisHost := mustHaveEnvVariable(redisHostEnvVar)
	environment := mustHaveEnvVariable(environmentEnvVar)
	serverPort := mustHaveEnvVariable(serverPortEnvVar)

	rc := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		DB:       redisDB,
		Password: "",
	})

	indexTmpl := template.Must(template.New("index").Parse(assets.IndexTemplate))

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	if environment == "production" {
		l, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		logger = l
	}

	logger.Info("greeter started")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(zapchi.Logger(logger, "router"))
	r.Use(middleware.Recoverer)

	r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		people := rc.Incr(r.Context(), redisPeopleKey)
		if people.Err() != nil {
			logger.Error("error during redis key retriving", zap.Error(err))
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("server error"))
			return
		}
		indexTmpl.Execute(rw, struct {
			People int64
		}{People: people.Val()})
	})

	http.ListenAndServe(":"+serverPort, r)
}
