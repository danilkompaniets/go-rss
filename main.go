package main

import (
	"database/sql"
	"fmt"
	"github.com/danilkompaniets/go-rss/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func startDbMigrations(dbConnection *sql.DB) error {
	goose.SetBaseFS(nil)
	goose.SetDialect("postgres")
	goose.SetLogger(log.New(os.Stdout, "[migrate] ", log.LstdFlags))
	time.Sleep(5 * time.Second)

	err := goose.Up(dbConnection, "sql/schema")

	if err != nil {
		return err
	}

	return nil
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	router := chi.NewRouter()

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("DB_HOST")

	if dbHost == "" {
		log.Fatal("environment variable not set")
	}

	dbUrl := fmt.Sprintf("postgres://%v:%v@%v:5432/%v?sslmode=disable", dbUser, dbPassword, dbHost, dbName)

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	err = startDbMigrations(conn)
	if err != nil {
		fmt.Println(err)
	}

	apiConfig := apiConfig{
		DB: database.New(conn),
	}

	go startScraping(
		apiConfig.DB, 10, time.Minute,
	)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handleReadiness)
	v1Router.Get("/error", handlerErr)

	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handlerGetUser))
	v1Router.Post("/users", apiConfig.handlerCreateUser)

	v1Router.Get("/posts", apiConfig.middlewareAuth(apiConfig.handlerGetPostsForUser))

	v1Router.Post("/feed", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feed", apiConfig.handlerGetFeeds)

	v1Router.Post("/feed-follows", apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollows))
	v1Router.Get("/feed-follows", apiConfig.middlewareAuth(apiConfig.handlerGetFeedFollows))
	v1Router.Delete("/feed-follows/{feedFollowId}", apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollows))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Listening on port %s", port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
