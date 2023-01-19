package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
    Insert(string, string, string) error 
    Exists(int) (bool, error) 
    Authenticate(string, string) (int, error) 
    Get(int) (*User, error)
}


type User struct {
    ID int
    Username string
    Email string
    Hashed_password []byte 
    Created time.Time 
}

type UserModel struct {
    DB *sql.DB
}

func (m *UserModel) Insert(username, email, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }

    query := `INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`

    _, err = m.DB.Exec(query, username, email, string(hashedPassword))
    if err != nil {
        var mySQLError *mysql.MySQLError
        if errors.As(err, &mySQLError) {
            if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
                return ErrDuplicateEmail
            }
        }
        return err
    }

    return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
    var id int
    var hashedPassword []byte

    query := `
        SELECT id, hashed_password
        FROM users
        WHERE email = ?
    `
    err := m.DB.QueryRow(query, email).Scan(&id, &hashedPassword)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return 0, ErrInvalidCredentials
        }
        return 0, err
    }
    
    err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    
    if err != nil {
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return 0, ErrInvalidCredentials
        }
        return 0, err 
    }

    return id, nil 
}

func (m *UserModel) Exists(id int) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`
    err := m.DB.QueryRow(query, id).Scan(&exists)
    return exists, err
}

func (m *UserModel) Get(id int) (*User, error) {
    user := &User{} 

    query := `SELECT name, email, created
    FROM users
    WHERE id = ?`
    
    err := m.DB.QueryRow(query, id).Scan(&user.Username, &user.Email, &user.Created)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNoRecord
        }
        return nil, err
    }

    return user, nil
}

