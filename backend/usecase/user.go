package usecase

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UserUsecase struct {
	db *sql.DB
}

func NewUserUsecase(db *sql.DB) *UserUsecase {
	return &UserUsecase{
		db: db,
	}
}

func (uc *UserUsecase) Register(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = uc.db.Exec(
		"INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		name, email, string(hashedPassword),
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}
