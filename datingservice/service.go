package datingservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/chackett/dating-service/pkg/security"
	"github.com/chackett/dating-service/rankingservice"
	"github.com/chackett/dating-service/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"time"
)

var ErrDuplicateSwipe = errors.New("already swiped this user")

const ctxKeySessionUserID = "session_user_id"

type DateService struct {
	logger *slog.Logger
	repo   *repository.Repository
}

func New(repo *repository.Repository) (*DateService, error) {
	result := &DateService{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		repo:   repo,
	}

	return result, nil
}

func (s *DateService) CreateUser(ctx context.Context, user repository.User) (*repository.User, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("unable to hash password: %w", err)
	}

	user.Password = string(h)
	createdUser, err := s.repo.CreateUser(ctx, &user)
	if err != nil {
		return nil, errors.New("")
	}
	// Clear password as soon as is appropriate.
	createdUser.Password = ""
	createdUser.Age = createdUser.CalculateAge()
	createdUser.DateOfBirth = nil
	return createdUser, nil
}

func (s *DateService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("get user password hash: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("compare hash and password: %w", err)
	}

	sessionTokenSize := 32
	st, err := security.CreateSecureSessionToken(sessionTokenSize)
	if err != nil {
		return "", fmt.Errorf("create session token: %w", err)
	}

	now := time.Now()
	userSession := repository.Session{
		UserID:    user.ID,
		Token:     st,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Hour * 24),
	}

	err = s.repo.CreateUserAuthSession(ctx, userSession)
	if err != nil {
		return "", fmt.Errorf("create user auth session: %w", err)
	}

	return userSession.Token, nil
}

func (s *DateService) Discover(ctx context.Context, userID int) (rankingservice.RankedResultSet, error) {
	sessionUserID, ok := ctx.Value(ctxKeySessionUserID).(int)
	if !ok {
		return rankingservice.RankedResultSet{}, errors.New("cannot find user id in context")
	}
	currentUser, err := s.repo.GetUserByID(ctx, sessionUserID)
	if err != nil {
		return rankingservice.RankedResultSet{}, fmt.Errorf("find user by d in repo: %w", err)
	}

	candidateMatches, err := s.repo.GetUnratedUsers(ctx, userID)
	if err != nil {
		return rankingservice.RankedResultSet{}, fmt.Errorf("discover candidateMatches in repo: %w", err)
	}

	rankedMatches := rankingservice.NewRankedResultSet()

	for _, c := range candidateMatches {
		score, err := currentUser.RankCandidate(c)
		if err != nil {
			s.logger.Error("error ranking user (%d) with candidate (%d): %w", currentUser.ID, c.ID, err)
		}

		rankedMatch := rankingservice.RankedMatch{
			User:    c,
			Ranking: score,
		}

		rankedMatches.AddMatch(rankedMatch)
	}

	return rankedMatches, nil
}

func (s *DateService) Swipe(ctx context.Context, swipeMessage repository.Swipe) (bool, error) {
	err := s.repo.SubmitSwipe(ctx, swipeMessage)
	if err != nil {
		if errors.As(gorm.ErrDuplicatedKey, &err) {
			return false, ErrDuplicateSwipe
		}
		return false, fmt.Errorf("submit swipe to repo: %w", err)
	}

	match, err := s.repo.IsUserMatch(ctx, swipeMessage.UserID, swipeMessage.CandidateID)
	if err != nil {
		return false, fmt.Errorf("check for user match: %w", err)
	}

	return match, nil
}

func (s *DateService) AuthenticateUserToken(ctx context.Context, token string) (int, error) {
	u, err := s.repo.GetUserFromAuthToken(ctx, token)
	if err != nil {
		return 0, fmt.Errorf("get user from auth token: %w", err)
	}
	return u.ID, nil
}
