package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/thienhi/fusionstart/internal/configs"
	"github.com/thienhi/fusionstart/internal/databases"
	"github.com/thienhi/fusionstart/internal/middleware"
	"github.com/thienhi/fusionstart/internal/rabbitmq"
	router "github.com/thienhi/fusionstart/internal/routes"
	"github.com/thienhi/fusionstart/internal/workers"
)

func main() {
	envPath := ".env"
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("No .env file found: %s", err)
	}
	var config = configs.Load()

	// databases
	db, err := databases.ConnectDatabase(config)
	if err != nil {
		log.Fatalf("Failed to connect database: %s", err)
	}
	defer databases.CloseConnectionDatabase(db)

	// rabbitmq
	rbmq, err := rabbitmq.ConnectRabbitMQ(config)
	if err != nil {
		log.Fatalf("Failed to connect rabbitmq: %s", err)
	}
	defer rabbitmq.CloseConnectionRabbitMQ(rbmq)

	// setup rabbitmq queues
	rabbitmq.SetupRabbitMQ(rbmq)

	// routers
	routers := gin.Default()
	routers.Use(middleware.CORSMiddleware())
	routers.Use(gin.Logger())
	routers.Use(gin.Recovery())

	// channel rabbitmq
	channel, err := rbmq.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer channel.Close()
	router.SetupRouter(routers, db, channel)

	// workers
	bookingChannel, err := rbmq.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	workers.BookingCancelTicketConsumer(bookingChannel, db)
	defer bookingChannel.Close()

	paymentChannel, err := rbmq.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	workers.PaymentTicketConsumer(paymentChannel, db)
	defer paymentChannel.Close()

	routers.Run(":8000")
}
