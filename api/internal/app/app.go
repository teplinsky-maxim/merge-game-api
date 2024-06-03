package app

import (
	"github.com/gofiber/fiber/v2"
	"log"
	v1 "merge-api/api/internal/contoller/http/v1"
	"merge-api/api/internal/repo"
	"merge-api/api/internal/service"
	"merge-api/shared/config"
	"merge-api/shared/pkg/database"
	"merge-api/shared/pkg/rabbitmq"
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
