package service

import (
	"merge-api/worker/internal/repo"
	"merge-api/worker/pkg/redis"
)

type Dependencies struct {
	Repositories repo.Repositories
	Redis        redis.Redis
}
