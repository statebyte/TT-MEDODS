package main

import (
	"backend/src/config"
	"backend/src/internal/authservice"
)

func init() {
	config.Env.Init()
}

func main() {
	authservice.Run()
}
