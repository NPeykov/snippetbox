package models

import (
	"database/sql"
	"time"
)

type User struct {
    ID int
    username string
    email string
    hashed_password []byte 
    created time.Time 
}

type UserModel struct {
    DB *sql.DB
}

func (m *UserModel) Insert(username, email, password string) error {
    return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
    return 1, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
    return false, nil
}

