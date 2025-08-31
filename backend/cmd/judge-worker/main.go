package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/LaurelEdison/clashcoder/backend/handlers/auth"
	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/LaurelEdison/clashcoder/backend/utils"
	"go.uber.org/zap"
)

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
