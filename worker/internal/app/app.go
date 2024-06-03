package app

import (
	"log"
	"merge-api/shared/config"
	"merge-api/shared/pkg/database"
	"merge-api/shared/pkg/rabbitmq"
)

func Run() {
	cfg, err := config.NewConfig(nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = database.NewDatabaseConnection(cfg)
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

	//repositories := repo.NewRepositories(&db)
	//serviceDeps := service.Dependencies{
	//	Repositories: *repositories,
	//	RabbitMQ:     rmq,
	//}
	//_ := service.NewServices(serviceDeps)
}
