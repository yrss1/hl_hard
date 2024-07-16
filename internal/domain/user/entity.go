package user

type Entity struct {
	ID       string  `db:"id"`
	FullName *string `db:"full_name"`
	Email    *string `db:"email"`
	Role     *string `db:"role"`
}
