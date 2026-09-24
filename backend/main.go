package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"quiz-backend/config"
	"quiz-backend/handler"
	"quiz-backend/model"
	"quiz-backend/repository"
	"quiz-backend/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func maxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	if err := godotenv.Load(); err != nil {
		slog.Info("no .env file found, relying on real environment variables")
	}
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration: %v", "error", err)
		os.Exit(1)
	}
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("unable to create connection pool: %v", "error", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		slog.Error("unable to reach the database: %v", "error", err)
		os.Exit(1)
	}
	slog.Info("Connected to PostgreSQL")
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
	attemptService := service.NewAttemptService(attemptRepo, quizRepo, answerRepo, optionRepo, pool)
	attemptHandler := handler.NewAttemptHandler(attemptService)

	mux.HandleFunc("POST /attempts/{attemptID}/submit", attemptHandler.SubmitAttempt)
	// mux.HandleFunc("POST /quizzes/{id}/attempts", attemptHandler.StartAttempt)
	mux.HandleFunc("POST /attempts/{attemptID}/answers", answerHandler.SubmitAnswer)
	mux.HandleFunc("PUT /attempts/{attemptID}/answers/{questionID}", answerHandler.UpdateAnswer)

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	healthHandler := handler.NewHealthHandler(pool)
	mux.HandleFunc("GET /health", healthHandler.Health)
	mux.HandleFunc("GET /ready", healthHandler.Ready)

	mux.Handle("POST /login", handler.RateLimitMiddleware(http.HandlerFunc(authHandler.Login)))
	mux.Handle("POST /register", handler.RateLimitMiddleware(http.HandlerFunc(authHandler.Register)))
	mux.Handle("POST /quizzes/{id}/attempts", handler.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(attemptHandler.StartAttempt)))

	mux.Handle("POST /quizzes",
		handler.AuthMiddleware(cfg.JWTSecret)(
			handler.RequireRole(model.RoleTeacher, model.RoleAdmin)(
				http.HandlerFunc(quizHandler.CreateQuiz),
			),
		),
	)

	mux.Handle("POST /quizzes/{id}/questions",
		handler.AuthMiddleware(cfg.JWTSecret)(
			handler.RequireRole(model.RoleTeacher, model.RoleAdmin)(
				http.HandlerFunc(questionHandler.CreateQuestion),
			),
		),
	)

	mux.Handle("POST /questions/{id}/options",
		handler.AuthMiddleware(cfg.JWTSecret)(
			handler.RequireRole(model.RoleTeacher, model.RoleAdmin)(
				http.HandlerFunc(optionHandler.CreateOption),
			),
		),
	)
	mux.Handle("GET /attempts/{attemptID}",
		handler.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(attemptHandler.GetAttemptByID)),
	)

	mux.Handle("GET /my-attempts",
		handler.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(attemptHandler.GetMyAttempts)),
	)

	mux.Handle("GET /quizzes/{quizID}/attempts",
		handler.AuthMiddleware(cfg.JWTSecret)(handler.RequireRole(model.RoleTeacher, model.RoleAdmin)(
			http.HandlerFunc(attemptHandler.GetAttemptsForQuiz),
		)),
	)

	mux.HandleFunc("GET /docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.yaml")
	})
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/docs/openapi.yaml"),
	))
	corsHandler := handler.CORSMiddleware(cfg.AllowedOrigin)(mux)
	loggedMux := handler.LoggingMiddleware(handler.RequestIDMiddleware(maxBodySize(1 << 20)(corsHandler)))
	slog.Info("starting server", "port", cfg.ServerPort, "environment", cfg.Environment)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, loggedMux))
}
