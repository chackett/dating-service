package repository

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

type Repository struct {
	logger *slog.Logger
	db     *gorm.DB
}

func New(user string, pass string, host string, port int, dbName string) (*Repository, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True", user, pass, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to DB: %w", err)
	}

	result := &Repository{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		db:     db,
	}

	return result, nil
}

func (r *Repository) CreateUser(ctx context.Context, newUser *User) (*User, error) {
	res := r.db.WithContext(ctx).Create(newUser)
	if res.Error != nil {
		return nil, fmt.Errorf("create user: %w", res.Error)
	}
	return newUser, nil
}

func (r *Repository) GetUser(ctx context.Context, emailAddress string) (User, error) {
	u := User{}
	res := r.db.WithContext(ctx).Where("email = ?", emailAddress).First(&u)
	if res.Error != nil {
		return User{}, fmt.Errorf("retrieve user by email: %w", res.Error)
	}

	return u, nil
}

func (r *Repository) CreateUserAuthSession(ctx context.Context, session Session) error {
	res := r.db.WithContext(ctx).Create(&session)
	if res.Error != nil {
		return fmt.Errorf("create user auth session: %w", res.Error)
	}
	return nil
}

func (r *Repository) GetUnratedUsers(ctx context.Context, userID int) ([]User, error) {
	var unratedUsers []User

	subquery := r.db.WithContext(ctx).Table("swipes").Select("candidate_id").Where("user_id = ?", userID)

	res := r.db.WithContext(ctx).Where("id NOT IN (?) AND id != ?", subquery, userID).Find(&unratedUsers)
	if res.Error != nil {
		return nil, fmt.Errorf("error retrieving unrated users: %w", res.Error)
	}

	return unratedUsers, nil
}

func (r *Repository) SubmitSwipe(ctx context.Context, input Swipe) error {
	res := r.db.WithContext(ctx).Create(input)
	if res.Error != nil {
		return fmt.Errorf("submit swipe to db: %w", res.Error)
	}

	return nil
}

func (r *Repository) IsUserMatch(ctx context.Context, userID int, candidateID int) (bool, error) {
	var count int64
	// Check if there is a mutual like between userID1 and userID2
	err := r.db.WithContext(ctx).Table("swipes").
		Where("(user_id = ? AND candidate_id = ? AND likes = ?) OR (user_id = ? AND candidate_id = ? AND likes = ?)",
			userID, candidateID, true, candidateID, userID, true).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("query for mutual matches: %w", err)
	}

	return count == 2, nil
}

func (r *Repository) GetUserFromAuthToken(ctx context.Context, token string) (*User, error) {
	user := &User{}
	res := r.db.WithContext(ctx).Joins("JOIN sessions ON sessions.user_id = users.id").
		Where("sessions.token = ?", token).
		First(user)
	if res.Error != nil {
		return nil, fmt.Errorf("user not found for auth token: %w", res.Error)
	}
	return user, nil
}
