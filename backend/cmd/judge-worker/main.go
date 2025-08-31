package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/handlers/auth"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/LaurelEdison/clashcoder/backend/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func workerLoop(zapLogger *zap.Logger, queries *database.Queries) {
	for {
		ctx := context.Background()

		submission, err := queries.SelectPendingSubmission(ctx)
		if err != nil {
			if err == sql.ErrNoRows {
				time.Sleep(2 * time.Second)
				continue
			}
			zapLogger.Error("Error fetching submission", zap.Error(err))
			time.Sleep(2 * time.Second)
			continue
		}

		err = queries.UpdateSubmissionStatus(ctx, database.UpdateSubmissionStatusParams{
			ID:     submission.ID,
			Status: "running",
		})
		if err != nil {
			zapLogger.Error("Error updating submission to running", zap.Error(err))
			continue
		}

		problem, err := queries.GetProblemByID(ctx, submission.ProblemID)
		if err != nil {
			zapLogger.Error("Error fetching problem", zap.Error(err))
			continue
		}
		tests, err := queries.GetProblemTestsByProblemID(ctx, submission.ProblemID)
		if err != nil {
			zapLogger.Error("Error fetching problem", zap.Error(err))
			continue
		}
		result, err := RunJudge(zapLogger, submission, problem, tests)
		if err != nil {
			zapLogger.Error("Error running judge", zap.Error(err))
			continue
		}
		status := "wrong_answer"
		if ctx.Err() == context.DeadlineExceeded {
			status = "time_limit_exceed"
		} else if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 137 {
				status = "memory_limit_exceed"
			}
		}
		if result.Passed {
			status = "accepted"
		} else if strings.Contains(result.Output, "panic") || strings.Contains(result.Output, "exit status 2") {
			status = "runtime_error"
		} else if strings.Contains(result.ErrorMsg, "time limit exceeded") {
			status = "time_limit_exceed"
		} else if strings.Contains(result.ErrorMsg, "memory limit exceeded") {
			status = "memory_limit_exceed"
		}
		zapLogger.Info("Judge results", zap.Bool("passed", result.Passed),
			zap.String("result", result.Output), zap.String("error", result.ErrorMsg))

		err = queries.UpdateSubmissionResult(ctx, database.UpdateSubmissionResultParams{
			ID:     submission.ID,
			Output: sql.NullString{String: result.Output},
			Status: status,
		})
		if err != nil {
			zapLogger.Error("Error updating submission to running", zap.Error(err))
			continue
		}

	}
}

type JudgeResult struct {
	Passed   bool
	Output   string
	ErrorMsg string
}

func RunJudge(zapLogger *zap.Logger, submission database.Submission, problem database.Problem, tests []database.ProblemTest) (JudgeResult, error) {
	tempDir, err := os.MkdirTemp("", "submission-*")
	if err != nil {
		zapLogger.Error("Could not create temporary directory", zap.Error(err))
	}

	if err := os.WriteFile(filepath.Join(tempDir, "solution.go"), []byte(submission.Code), 0644); err != nil {
		zapLogger.Error("Could not write submission to tempdir", zap.Error(err))
		err = os.RemoveAll(tempDir)
		if err != nil {
			zapLogger.Error("Could not delete temporary directory", zap.Error(err))
		}
		return JudgeResult{}, err
	}

	for _, test := range tests {
		if err := os.WriteFile(filepath.Join(tempDir, "solution_test.go"),
			[]byte(test.TestCode), 0644); err != nil {
			zapLogger.Error("Could not write testfile to tempdir", zap.Error(err))
			err = os.RemoveAll(tempDir)
			if err != nil {
				zapLogger.Error("Could not delete temporary directory", zap.Error(err))
			}
			return JudgeResult{}, err
		}
	}
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"),
		[]byte("module submission\n\ngo 1.22\n"), 0644); err != nil {
		zapLogger.Error("Could not write go.mod", zap.Error(err))
		_ = os.RemoveAll(tempDir)
		return JudgeResult{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(problem.TimeLimit)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"--cpus=1.0", "--memory=512m",
		"-v", fmt.Sprintf("%s:/app", tempDir),
		"-w", "/app", "golang:1.22",
		"go", "test", "-v", ".")

	output, err := cmd.CombinedOutput()
	result := JudgeResult{
		Passed:   err == nil,
		Output:   string(output),
		ErrorMsg: "",
	}
	if err != nil {
		result.ErrorMsg = err.Error()
		zapLogger.Error("Error found while running judge", zap.Error(err))
	}

	err = os.RemoveAll(tempDir)
	if err != nil {
		zapLogger.Error("Could not delete temporary directory", zap.Error(err))
	}

	return result, nil
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load .env : %v", err)
	}
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Error creating zap logger: %v", err))
	}
	dbURL := utils.GetDBUrl(zapLogger)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET is not set")
	}
	auth.InitJWT(secret)

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		zapLogger.Error("Failed to connect to database", zap.Error(err))
	}

	if err := conn.Ping(); err != nil {
		zapLogger.Error("Failed to connect to database", zap.Error(err))
	}

	queries := database.New(conn)

	zapLogger.Info("Judge worker started")

}
