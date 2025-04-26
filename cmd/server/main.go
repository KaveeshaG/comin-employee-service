// cmd/server/main.go
package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Axontik/comin-employee-service/internal/handler"
	"github.com/Axontik/comin-employee-service/internal/middleware"
	"github.com/Axontik/comin-employee-service/internal/repository"
	"github.com/Axontik/comin-employee-service/internal/service"
	"github.com/Axontik/comin-employee-service/pkg/auth"
	"github.com/Axontik/comin-employee-service/pkg/organization"
)

type Application struct {
	db              *gorm.DB
	employeeHandler *handler.EmployeeHandler
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	app := &Application{}

	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	app.db = db

	// Initialize dependencies
	app.initializeDependencies()

	// Setup router
	router := setupRouter(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB() (*gorm.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgresql://comin_owner:Ye5rfjcIB7FX@ep-flat-shadow-a8onelva.eastus2.azure.neon.tech/comin?sslmode=require"
	}

	// Run migrations
	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Printf("Warning: Failed to initialize migrations: %v", err)
	} else {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Printf("Warning: Failed to run migrations: %v", err)
		}
	}

	config := &gorm.Config{
		// Add any GORM configurations here
	}

	return gorm.Open(postgres.Open(dbURL), config)
}

func (app *Application) initializeDependencies() {
	// Initialize repositories
	employeeRepo := repository.NewEmployeeRepository(app.db)

	// Initialize services
	employeeService := service.NewEmployeeService(employeeRepo)

	// Initialize handlers
	app.employeeHandler = handler.NewEmployeeHandler(employeeService)
}

func setupRouter(app *Application) *gin.Engine {
	// Initialize auth client
	authClient := auth.NewAuthClient("https://comin.kaveeshagimhana.com/api/v1/auth")
	if authClient == nil {
		authClient = auth.NewAuthClient("https://comin.kaveeshagimhana.com/api/v1/auth")
	}

	orgClient := organization.NewOrganizationClient("https://comin.kaveeshagimhana.com/api/v1")
	if orgClient == nil {
		orgClient = organization.NewOrganizationClient("https://comin.kaveeshagimhana.com/api/v1")
	}

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Employee routes
		employees := api.Group("/employees")
		employees.Use(organization.ValidateOrganizationAccess(authClient, orgClient))
		{
			employees.POST("/", app.employeeHandler.Create)
			employees.GET("/:id", app.employeeHandler.GetByID)
			employees.PUT("/:id", app.employeeHandler.Update)
			employees.DELETE("/:id", app.employeeHandler.Delete)
			employees.PUT("/:id/status", app.employeeHandler.UpdateStatus)
			employees.PUT("/:id/manager", app.employeeHandler.AssignManager)
		}

		// Organization specific routes
		org := api.Group("employees/organizations/:organization_id")
		org.Use(organization.ValidateOrganizationAccess(authClient, orgClient))
		{
			org.GET("/employees", app.employeeHandler.ListByOrganization)
		}

		// Department specific routes
		dept := api.Group("/departments/:department_id")
		dept.Use(organization.ValidateOrganizationAccess(authClient, orgClient))
		{
			dept.GET("/employees", app.employeeHandler.ListByDepartment)
		}

		// Manager specific routes
		manager := api.Group("/managers/:manager_id")
		manager.Use(organization.ValidateOrganizationAccess(authClient, orgClient))
		{
			manager.GET("/employees", app.employeeHandler.ListByManager)
		}

		// User specific routes
		users := api.Group("/users")
		users.Use(organization.ValidateOrganizationAccess(authClient, orgClient))
		{
			users.GET("/:user_id/employee", app.employeeHandler.GetByUserID)
		}
	}

	return router
}
