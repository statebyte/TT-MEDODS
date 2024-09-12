package authservice

import (
	"backend/src/config"
	"backend/src/pkg/jwt"
	"backend/src/pkg/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const PORT string = "8080"

var (
	router *mux.Router
	DB     AuthDB
	jwtGen jwt.JWTGenerator
)

func InitRoutes() {
	router = mux.NewRouter()
	router.HandleFunc("/auth/token/{user_id}", IssueTokens).Methods("GET")
	router.HandleFunc("/auth/token/refresh", RefreshTokens).Methods("PATCH")

	router.Use(utils.LoggingMiddleware)
}

func InitJWT() {
	jwtGen.Init(config.Env.Secret)
	log.Println("Module JWT inited")
}

func InitDB() {
	err := DB.DBInstance.Connect(config.Env.Database)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected")
}

func Run() {
	InitRoutes()
	InitJWT()
	InitDB()

	log.Println("Service running on port", PORT)

	http.ListenAndServe(":"+PORT, router)
}
