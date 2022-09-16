package handlers

import "gopkg.in/go-playground/validator.v9"

var (
	v = validator.New()
)

// productValidator validate struct
type productValidator struct {
	validator *validator.Validate
}

func (p *productValidator) Validate(i interface{}) error {
	return p.validator.Struct(i)
}

// userValidator validate struct
type userValidator struct {
	validator *validator.Validate
}

func (u *userValidator) Validate(i interface{}) error {
	return u.validator.Struct(i)
}
