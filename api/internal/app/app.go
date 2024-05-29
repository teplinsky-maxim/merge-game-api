package app

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"merge-api/config"
	v1 "merge-api/internal/contoller/http/v1"
	"merge-api/internal/repo"
	"merge-api/internal/service"
	"merge-api/pkg/database"
	"merge-api/pkg/rabbitmq"
)

func Run() {
	cfg, err := config.NewConfig(nil)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = database.RunMigrations(cfg, "public")
	if err != nil {
		log.Fatal(err)
	}

	rmq, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = rabbitmq.Initialize(&rmq)
	if err != nil {
		log.Fatal(err)
	}

	repositories := repo.NewRepositories(&db)
	serviceDeps := service.Dependencies{
		Repositories: *repositories,
		RabbitMQ:     rmq,
	}
	services := service.NewServices(serviceDeps)

	app := fiber.New()
	v1.NewRouter(app, services)

	err = app.Listen("0.0.0.0:3000")
	if err != nil {
		log.Fatal(err)
	}
}
