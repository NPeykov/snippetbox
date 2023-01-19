package mocks

import (
	"time"

	"github.com/NPeykov/snippetbox/internal/models"
)

type UserModel struct{}

var mockUser = &models.User{
    ID: 1, 
    Username: "pkv", 
    Email: "gg@gg.com", 
    Hashed_password: []byte("$2a$12$AKCfdUHY.pHWtuiBi0FXl.97yMqVO8tm82qvpJCL/gG7Pqnwab/2K"),
    Created: time.Now(),
}

func (m *UserModel) Insert(username, email, password string) error {
    switch email {
        case "gg@gg.com": return models.ErrDuplicateEmail
        default: return nil
    }
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
    if email == "gg@gg.com" && password == "12345678" {
        return 1, nil
    }

    return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
    switch id {
    case 1:
        return true, nil
    default:
        return false, nil 
    }
}

func (m *UserModel) Get(id int) (*models.User, error) {
    switch id {
    case 1:
        return mockUser, nil
    default:
        return nil, models.ErrNoRecord
    }
}

