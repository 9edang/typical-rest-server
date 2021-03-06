package infra

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-rest-server/pkg/typrest"
	"go.uber.org/dig"
)

type (
	// HealthCheck for profiler
	HealthCheck struct {
		dig.In
		PG    *sql.DB `name:"pg"`
		MySQL *sql.DB `name:"mysql"`
		Redis *redis.Client
	}
)

// Handle echo function
func (h *HealthCheck) Handle(ec echo.Context) error {
	healthy, detail := typrest.HealthStatus(typrest.HealthMap{
		"postgres": h.PG.Ping(),
		"mysql":    h.MySQL.Ping(),
		"redis":    h.Redis.Ping().Err(),
	})

	status := http.StatusOK
	if !healthy {
		status = http.StatusServiceUnavailable
	}

	return ec.JSON(status, map[string]interface{}{
		"name":   fmt.Sprintf("%s (%s)", typgo.ProjectName, typgo.ProjectVersion),
		"status": detail,
	})
}
