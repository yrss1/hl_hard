package project

import "time"

type Entity struct {
	ID          string     `db:"id"`
	Title       *string    `db:"title"`
	Description *string    `db:"description"`
	StartDate   *time.Time `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
	ManagerID   *string    `db:"manager_id"`
}
