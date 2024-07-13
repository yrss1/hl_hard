package user

import "time"

type Entity struct {
	ID               string     `db:"id"`
	FullName         *string    `db:"full_name"`
	Email            *string    `db:"email"`
	RegistrationDate *time.Time `db:"registration_date"`
	Role             *string    `db:"role"`
}
