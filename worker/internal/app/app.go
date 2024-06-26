package app

import (
	"log"
	"merge-api/shared/config"
	"merge-api/shared/pkg/database"
	"merge-api/shared/pkg/rabbitmq"
	"merge-api/worker/internal/repo"
	"merge-api/worker/internal/service"
	"merge-api/worker/pkg/redis"
	"merge-api/worker/pkg/task"
	"merge-api/worker/pkg/task/executors"
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

	rmq, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = rabbitmq.Initialize(&rmq)
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := redis.NewRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repositories := repo.NewRepositories(&db, &redisClient)
	dependencies := service.Dependencies{
		Repositories: *repositories,
		Redis:        redisClient,
	}
	services := service.NewServices(dependencies)

	newNewBoardTaskExecutor := executors.NewNewBoardTaskExecutor(services.Board)
	newMoveItemTaskExecutor := executors.NewMoveItemTaskExecutor(services.Board)
	newMergeItemsTaskExecutor := executors.NewMergeItemsTaskExecutor(services.Board, services.Collection)
	newClickItemTaskExecutor := executors.NewClickItemTaskExecutor(services.Board, services.Collection)

	taskExecutorsManager := task.NewTaskExecutorsManager([]task.Executor{
		newNewBoardTaskExecutor,
		newMoveItemTaskExecutor,
		newMergeItemsTaskExecutor,
		newClickItemTaskExecutor,
	})

	err = task.StartPullTasks(&rmq, taskExecutorsManager, services.Task)
	if err != nil {
		log.Fatal(err)
	}
}
