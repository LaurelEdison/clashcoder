package handlers

import (
	"database/sql"
	"time"

	"github.com/LaurelEdison/clashcoder/backend/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID    `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	LastLoginAt sql.NullTime `json:"last_login_at"`
	Role        string       `json:"role"`
}

func DatabaseUserToUser(dbUser database.User) User {
	return User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Name:        dbUser.Name,
		Email:       dbUser.Email,
		LastLoginAt: dbUser.LastLoginAt,
		Role:        dbUser.Role,
	}
}

type Problem struct {
	ID            uuid.UUID      `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Difficulty    sql.NullString `json:"difficulty"`
	TimeLimit     int32          `json:"time_limit"`
	MemoryLimitMb int32          `json:"memory_limit_mb"`
}

func DatabaseProblemToProblem(dbProblem database.Problem) Problem {
	return Problem{
		ID:            dbProblem.ID,
		CreatedAt:     dbProblem.CreatedAt,
		UpdatedAt:     dbProblem.UpdatedAt,
		Title:         dbProblem.Title,
		Description:   dbProblem.Description,
		Difficulty:    dbProblem.Difficulty,
		TimeLimit:     dbProblem.TimeLimit,
		MemoryLimitMb: dbProblem.MemoryLimitMb,
	}
}

func DatabaseProblemsToProblems(dbProblems []database.Problem) []Problem {
	problems := []Problem{}
	for _, dbProblem := range dbProblems {
		problems = append(problems, DatabaseProblemToProblem(dbProblem))
	}
	return problems
}

type Submission struct {
	ID        uuid.UUID      `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UserID    uuid.UUID      `json:"user_id"`
	ProblemID uuid.UUID      `json:"problem_id"`
	Code      string         `json:"code"`
	Language  string         `json:"language"`
	Status    sql.NullString `json:"status"`
	RuntimeMs sql.NullInt32  `json:"runtime_ms"`
	MemoryKb  sql.NullInt32  `json:"memory_kb"`
	Output    sql.NullString `json:"output"`
}

