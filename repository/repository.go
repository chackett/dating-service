package repository

import (
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, port, dbName)
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

func (r *Repository) CloseConnection() error {
	return nil
}

func (r *Repository) CreateUser(newUser *User) (*User, error) {
	res := r.db.Create(newUser)
	if res.Error != nil {
		return nil, fmt.Errorf("create user: %w", res.Error)
	}
	return newUser, nil
}

func (r *Repository) GetUser(emailAddress string) (User, error) {
	u := User{}
	res := r.db.Where("email = ?", emailAddress).First(&u)
	if res.Error != nil {
		return User{}, fmt.Errorf("retrieve user by email: %w", res.Error)
	}

	return u, nil
}

func (r *Repository) CreateUserAuthSession(session Session) error {
	res := r.db.Create(&session)
	if res.Error != nil {
		return fmt.Errorf("create user auth session: %w", res.Error)
	}
	return nil
}

func (r *Repository) GetUnratedUsers(userID int) ([]User, error) {
	var unratedUsers []User

	subquery := r.db.Table("swipes").Select("candidate_id").Where("user_id = ?", userID)

	res := r.db.Where("id NOT IN (?) AND id != ?", subquery, userID).Find(&unratedUsers)
	if res.Error != nil {
		return nil, fmt.Errorf("error retrieving unrated users: %w", res.Error)
	}

	return unratedUsers, nil
}

func (r *Repository) SubmitSwipe(input Swipe) error {
	res := r.db.Create(input)
	if res.Error != nil {
		return fmt.Errorf("submit swipe to db: %w", res.Error)
	}

	return nil
}

func (r *Repository) GetUserFromAuthToken(token string) (*User, error) {
	user := &User{}
	res := r.db.Joins("JOIN sessions ON sessions.user_id = users.id").
		Where("sessions.token = ?", token).
		First(user)
	if res.Error != nil {
		return nil, fmt.Errorf("user not found for auth token: %w", res.Error)
	}
	return user, nil
}
