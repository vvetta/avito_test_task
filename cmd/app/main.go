package main

import (
	"log"
	"net/http"
	"os"

	apphttp "avito_test_task/internal/adapters/http"
	"avito_test_task/internal/adapters/http/openapi"
	"avito_test_task/internal/adapters/repository/pr_repository"
	"avito_test_task/internal/adapters/repository/team_repository"
	"avito_test_task/internal/adapters/repository/user_repository"
	"avito_test_task/internal/usecase"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load(".env")

	dsn := getDSNFromEnv()

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}


	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}


	teamRepo := teamrepository.NewTeamRepository(db)
	userRepo := userrepository.NewUserRepository(db)
	prRepo := prrepository.NewPrRepository(db)


	teamService := usecase.NewTeamService(teamRepo, userRepo)
	userService := usecase.NewUserService(userRepo, prRepo)
	prService := usecase.NewPRService(prRepo, userRepo, teamRepo)


	handler := apphttp.NewHandler(userService, teamService, prService)


	mux := http.NewServeMux()
	openapi.HandlerFromMux(handler, mux)


	log.Printf("starting HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}

func getDSNFromEnv() string {
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresPASS := os.Getenv("POSTGRES_PASSWORD")
	postgresUSER := os.Getenv("POSTGRES_USER")
	postgresHOST := os.Getenv("POSTGRES_HOST")
	postgresPORT := os.Getenv("POSTGRES_P")

	dsn := "host="+postgresHOST+" port="+postgresPORT+" user="+postgresUSER+" password="+postgresPASS+" dbname="+postgresDB+" sslmode=disable"
	return dsn
}
