package models

import "errors"

var (
    ErrNoRecord = errors.New("models: no matching record found")
    ErrInvalidCredentials = errors.New("models: credentials are invalid")
    ErrDuplicateEmail = errors.New("models: mail already exists")
)


