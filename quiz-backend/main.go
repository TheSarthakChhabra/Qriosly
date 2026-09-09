package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"quiz-backend/handler"
	"quiz-backend/model"
	"quiz-backend/repository"
	"quiz-backend/service"
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
	mux := http.NewServeMux()
	quizRepo := repository.NewQuizRepository(pool)
	quizService := service.NewQuizService(quizRepo)
	quizHandler := handler.NewQuizHandler(quizService)

	mux.HandleFunc("GET /quizzes", quizHandler.GetAllQuizzes)
	mux.HandleFunc("GET /quizzes/{id}", quizHandler.GetQuizByID)
	// mux.HandleFunc("POST /quizzes", quizHandler.CreateQuiz)

	questionRepo := repository.NewQuestionRepository(pool)
	questionService := service.NewQuestionService(questionRepo, quizRepo)
	questionHandler := handler.NewQuestionHandler(questionService)

	// mux.HandleFunc("POST /quizzes/{id}/questions", questionHandler.CreateQuestion)
	mux.HandleFunc("GET /quizzes/{id}/questions", questionHandler.GetQuestionsByQuizID)
	mux.HandleFunc("GET /questions/{id}", questionHandler.GetQuestionByID)

	optionRepo := repository.NewOptionRepository(pool)
	optionService := service.NewOptionService(optionRepo, questionRepo)
	optionHandler := handler.NewOptionHandler(optionService)

	// mux.HandleFunc("POST /questions/{id}/options", optionHandler.CreateOption)
	mux.HandleFunc("GET /questions/{id}/options", optionHandler.GetOptionsByQuestionID)

	answerRepo := repository.NewAnswerRepository(pool)
	attemptRepo := repository.NewAttemptRepository(pool)
	answerService := service.NewAnswerService(answerRepo, attemptRepo, questionRepo, optionRepo)
	answerHandler := handler.NewAnswerHandler(answerService)
	attemptService := service.NewAttemptService(attemptRepo, quizRepo, answerRepo, optionRepo)
	attemptHandler := handler.NewAttemptHandler(attemptService)

	mux.HandleFunc("POST /attempts/{attemptID}/submit", attemptHandler.SubmitAttempt)
	// mux.HandleFunc("POST /quizzes/{id}/attempts", attemptHandler.StartAttempt)
	mux.HandleFunc("POST /attempts/{attemptID}/answers", answerHandler.SubmitAnswer)
	mux.HandleFunc("PUT /attempts/{attemptID}/answers/{questionID}", answerHandler.UpdateAnswer)

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.Handle("POST /quizzes/{id}/attempts", handler.AuthMiddleware(http.HandlerFunc(attemptHandler.StartAttempt)))

	mux.Handle("POST /quizzes",
		handler.AuthMiddleware(
			handler.RequireRole(model.RoleTeacher, model.RoleAdmin)(
				http.HandlerFunc(quizHandler.CreateQuiz),
			),
		),
	)

	mux.Handle("POST /quizzes/{id}/questions",
		handler.AuthMiddleware(
			handler.RequireRole(model.RoleTeacher, model.RoleAdmin)(
				http.HandlerFunc(questionHandler.CreateQuestion),
			),
		),
	)

	mux.Handle("POST /questions/{id}/options",
		handler.AuthMiddleware(
			handler.RequireRole(model.RoleTeacher, model.RoleAdmin)(
				http.HandlerFunc(optionHandler.CreateOption),
			),
		),
	)
	mux.Handle("GET /attempts/{attemptID}",
		handler.AuthMiddleware(http.HandlerFunc(attemptHandler.GetAttemptByID)),
	)

	mux.Handle("GET /my-attempts",
		handler.AuthMiddleware(http.HandlerFunc(attemptHandler.GetMyAttempts)),
	)

	mux.Handle("GET /quizzes/{quizID}/attempts",
		handler.AuthMiddleware(handler.RequireRole(model.RoleTeacher, model.RoleAdmin)(
			http.HandlerFunc(attemptHandler.GetAttemptsForQuiz),
		)),
	)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
