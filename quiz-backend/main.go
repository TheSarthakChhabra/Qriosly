package main

import (
	"context"
	"log"
	"net/http"

	"quiz-backend/handler"
	"quiz-backend/repository"
	"quiz-backend/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	dbURL := "postgres://quizuser:quizpass@localhost:5432/quizdb"
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to reach the database: %v", err)
	}
	log.Println("Connected to PostgreSQL")
	quizRepo := repository.NewQuizRepository(pool)
	quizService := service.NewQuizService(quizRepo)
	quizHandler := handler.NewQuizHandler(quizService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /quizzes", quizHandler.GetAllQuizzes)
	mux.HandleFunc("GET /quizzes/{id}", quizHandler.GetQuizByID)
	mux.HandleFunc("POST /quizzes", quizHandler.CreateQuiz)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
