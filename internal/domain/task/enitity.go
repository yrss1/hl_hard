package task

type Entity struct {
	ID          string  `db:"id"`
	Title       *string `db:"title"`
	Description *string `db:"description"`
	Priority    *string `db:"priority"`
	Status      *string `db:"status"`
	AssigneeID  *string `db:"assignee_id"`
	ProjectID   *string `db:"project_id"`
	CompletedAt *string `db:"completed_at"`
}
