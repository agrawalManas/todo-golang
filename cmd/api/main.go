/*Important commands used
- migrate create -ext sql -dir migrations  -seq create_todos_table => for creating migration files for our DB
- migrate -path migrations -database "DATABASE_URL" up => for running an upgrade migration
- migrate -path migrations -database "DATABASE_URL" down $count => for running a downgrade migration with a specific count if needed
- psql -U manasagarwal-d postgres => for running postgres in terminal
- brew services start postgresql => for launching postgres services (installed via homebrew)
*/

package main

import (
	"log"
	"todo_with_auth/internals/config"
	"todo_with_auth/internals/database"
	"todo_with_auth/internals/handlers"
	"todo_with_auth/internals/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var cfg *config.Config
	var err error

	cfg, err = config.LoadConfig()

	if err != nil {
		log.Fatal("Cannot get the config: ", err)
	}

	//------------DB-Setup----------
	var pool *pgxpool.Pool
	pool, err = database.ConnectDB(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	defer pool.Close()

	//-------------Api-Setup-----------
	router := gin.Default()
	router.SetTrustedProxies(nil)

	//--------Route-Handling--------
	protectedAuth := router.Group("/app/auth")
	{
		protectedAuth.POST("/register-user", handlers.CreateUserHandler(pool))
		protectedAuth.POST("/signin", handlers.SigninUserHandler(pool, cfg))
	}

	protected := router.Group("/app")
	protected.Use(middlewares.AuthMiddleware(cfg))
	{
		protected.GET("/get-todos", handlers.GetTodoHandler(pool))
		protected.GET("/get-todos/:id", handlers.GetTodoByIdHandler(pool))
		protected.POST("/create-todo", handlers.CreateTodoHandler(pool))
		protected.PUT("/update-todo/:id", handlers.UpdateTodoHandler(pool))
		protected.DELETE("/delete-todo/:id", handlers.DeleteTodoHandler(pool))
	}

	router.Run("localhost:" + cfg.Port)
}
