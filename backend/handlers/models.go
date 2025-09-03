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
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Difficulty    string    `json:"difficulty"`
	TimeLimit     int32     `json:"time_limit"`
	MemoryLimitMb int32     `json:"memory_limit_mb"`
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
	Status    string         `json:"status"`
	RuntimeMs sql.NullInt32  `json:"runtime_ms"`
	MemoryKb  sql.NullInt32  `json:"memory_kb"`
	Output    sql.NullString `json:"output"`
}

func DatabaseSubmissiontoSubmission(dbSubmission database.Submission) Submission {
	return Submission{
		ID:        dbSubmission.ID,
		CreatedAt: dbSubmission.CreatedAt,
		UserID:    dbSubmission.UserID,
		ProblemID: dbSubmission.ProblemID,
		Code:      dbSubmission.Code,
		Language:  dbSubmission.Language,
		Status:    dbSubmission.Status,
		RuntimeMs: dbSubmission.RuntimeMs,
		MemoryKb:  dbSubmission.MemoryKb,
		Output:    dbSubmission.Output,
	}
}

func DatabaseSubmissionsToSubmissions(dbSubmissions []database.Submission) []Submission {
	submissions := []Submission{}
	for _, dbSubmission := range dbSubmissions {
		submissions = append(submissions, DatabaseSubmissiontoSubmission(dbSubmission))
	}
	return submissions
}

type ProblemTest struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ProblemID uuid.UUID `json:"problem_id"`
	TestCode  string    `json:"test_code"`
}

func DatabaseProblemTestToProblemTest(dbProblemTest database.ProblemTest) ProblemTest {
	return ProblemTest{
		ID:        dbProblemTest.ID,
		CreatedAt: dbProblemTest.CreatedAt,
		UpdatedAt: dbProblemTest.UpdatedAt,
		ProblemID: dbProblemTest.ProblemID,
		TestCode:  dbProblemTest.TestCode,
	}
}

func DatabaseProblemTestsToProblemTests(dbProblemTests []database.ProblemTest) []ProblemTest {
	problemtests := []ProblemTest{}
	for _, dbProblemTest := range dbProblemTests {
		problemtests = append(problemtests, DatabaseProblemTestToProblemTest(dbProblemTest))
	}
	return problemtests
}

type Lobby struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	InviteCode string    `json:"invite_code"`
	MaxUsers   int32     `json:"max_users"`
	Status     string    `json:"status"`
	ReadyState bool      `json:"ready_state"`
}

func DatabaseLobbyToLobby(dbLobby database.Lobby) Lobby {
	return Lobby{
		ID:         dbLobby.ID,
		InviteCode: dbLobby.InviteCode,
		MaxUsers:   dbLobby.MaxUsers,
		Status:     dbLobby.Status,
		ReadyState: dbLobby.ReadyState,
	}
}

type LobbyUser struct {
	LobbyID  uuid.UUID `json:"lobby_id"`
	UserID   uuid.UUID `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
	Role     string    `json:"role"`
}

func DatabaseLobbyUserToLobbyUser(dbLobbyUser database.LobbyUser) LobbyUser {
	return LobbyUser{
		LobbyID:  dbLobbyUser.LobbyID,
		UserID:   dbLobbyUser.UserID,
		JoinedAt: dbLobbyUser.JoinedAt,
		Role:     dbLobbyUser.Role,
	}
}

func DatabaseLobbyUsersToLobbyUsers(dbLobbyUsers []database.LobbyUser) []LobbyUser {
	LobbyUsers := []LobbyUser{}
	for _, dbLobbyUser := range dbLobbyUsers {
		LobbyUsers = append(LobbyUsers, DatabaseLobbyUserToLobbyUser(dbLobbyUser))
	}
	return LobbyUsers
}
