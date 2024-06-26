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

// DateService is the core component of this project, sitting between the HTTP layer and DB repository.
// Business logic is to be performed here.
type DateService struct {
	logger *slog.Logger
	repo   *repository.Repository
}

// New returns a new instance of DateService
func New(repo *repository.Repository) (*DateService, error) {
	result := &DateService{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		repo:   repo,
	}

	return result, nil
}

// CreateUser persists a new user into the DB.
// Note that passwords are not persisted "as is" but rather hashed using a PBKDF.
// If successful, the created user is returned with its unique identifer (`ID`) populated and the password removed.
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

// SetUserPreferences stores an updated set of preferences for a user. Note that the underlying DB operation is "upsert"
// so existing preferences will be overridden.
func (s *DateService) SetUserPreferences(ctx context.Context, prefs repository.UserPreferences) error {
	err := s.repo.UpsertUserPreferences(ctx, prefs)
	if err != nil {
		return fmt.Errorf("unable to upsert user preferences: %w", err)
	}
	return nil
}

// Login is used to create an authenticated session for a user, so subsequent authenticated calls can be made. Here a username
// and password is provided, and if valid, a session token is returned.
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

// Discover returns a collection of profiles that have been ranked and matched against the logged-in user. The intention
// is that these are presented to the user and subsequently "swiped", "yes" or "no" by the user.
// The returned results are ranked in decreasing order and some sensitive information has been removed for privacy reasons.
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

	userPrefs, err := s.repo.GetUserPreferences(ctx, sessionUserID)
	if err != nil {
		return rankingservice.RankedResultSet{}, fmt.Errorf("get user preferences from repo: %w", err)
	}
	rankedMatches := rankingservice.NewRankedResultSet()

	for _, cand := range candidateMatches {
		canPrefs, err := s.repo.GetUserPreferences(ctx, cand.ID)
		if err != nil {
			return rankingservice.RankedResultSet{}, fmt.Errorf("get user preferences from repo: %w", err)
		}

		score, err := currentUser.RankCandidate(cand, userPrefs, canPrefs)
		if err != nil {
			s.logger.Error("error ranking user (%d) with candidate (%d): %w", currentUser.ID, cand.ID, err)
		}

		if score == -1 {
			// Don't add candidate to results
			continue
		}

		candidateDistance := currentUser.DistanceFromUser(cand)
		cand.Age = cand.CalculateAge()
		cand.MaskPrivateFields()

		rankedMatch := rankingservice.RankedMatch{
			User:           cand,
			Ranking:        score,
			DistanceFromMe: candidateDistance,
		}

		rankedMatches.AddMatch(rankedMatch)
	}

	return rankedMatches, nil
}

// Swipe enables a user to specify if they like a discovered profile or not.
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

// AuthenticateUserToken verifies the tokens created during calls to Login. If the token is valid, the linked users ID is
// returned.
func (s *DateService) AuthenticateUserToken(ctx context.Context, token string) (int, error) {
	u, err := s.repo.GetUserFromAuthToken(ctx, token)
	if err != nil {
		return 0, fmt.Errorf("get user from auth token: %w", err)
	}
	return u.ID, nil
}
