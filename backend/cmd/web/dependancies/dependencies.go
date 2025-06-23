package dependencies

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	DB    *sql.DB
	Cache *redis.Client
}
