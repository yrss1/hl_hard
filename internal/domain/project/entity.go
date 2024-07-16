package project

type Entity struct {
	ID          string  `db:"id"`
	Title       *string `db:"title"`
	Description *string `db:"description"`
	StartDate   *string `db:"start_date"`
	EndDate     *string `db:"end_date"`
	ManagerID   *string `db:"manager_id"`
}
