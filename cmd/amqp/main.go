package main

import (
	"log"

	"github.com/ishola-faazele/taskflow/internal/emailservice"
	"github.com/ishola-faazele/taskflow/internal/utils/amqp"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("WARNING: .env FILE_NOT_FOUND")
	}
	conn := amqp.InitAMQP()
	defer conn.Close()
	emailservice.RegisterRoutes(conn)
}
