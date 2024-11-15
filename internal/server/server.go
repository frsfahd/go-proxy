package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/dotenv-org/godotenvvault/autoload"
	"github.com/frsfahd/go-proxy/config"
	"github.com/frsfahd/go-proxy/internal/database"
)

type Server struct {
	port   int
	target string
	db     database.Service
}

func NewServer(port int, target string, configs config.Config) *http.Server {
	// port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:   port,
		target: target,
		db:     database.New(configs),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	slog.Info(fmt.Sprintf("🛜 database status: %v", NewServer.db.Health()["redis_status"]))
	slog.Info(fmt.Sprintf("✅ server up at port %v", port))

	return server
}
