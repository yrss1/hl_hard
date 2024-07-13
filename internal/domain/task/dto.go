package task

import (
	"errors"
	"time"
)

type Request struct {
	ID          string  `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority"`
	Status      *string `json:"status"`
	AssigneeID  *string `json:"assignee_id"`
	ProjectID   *string `json:"project_id"`
	CompletedAt *string `json:"completed_at"`
}

func (s *Request) Validate() error {
	if s.Title == nil {
		return errors.New("title: cannot be blank")
	}

	if s.Description == nil {
		return errors.New("description: cannot be blank")
	}

	if s.Priority == nil {
		return errors.New("priority: cannot be blank")
	}

	if s.Status == nil {
		return errors.New("status: cannot be blank")
	}

	if s.AssigneeID == nil {
		return errors.New("assignee_id: cannot be blank")
	}

	if s.ProjectID == nil {
		return errors.New("project_id: cannot be blank")
	}

	return nil
}

func IsEmpty(data Request) bool {
	return data.Title == nil &&
		data.Priority == nil &&
		data.Status == nil &&
		data.AssigneeID == nil &&
		data.ProjectID == nil
}

type Response struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	AssigneeID  string    `json:"assignee_id"`
	ProjectID   string    `json:"project_id"`
	CompletedAt time.Time `json:"completed_at"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:          data.ID,
		Title:       *data.Title,
		Description: *data.Description,
		Priority:    *data.Priority,
		AssigneeID:  *data.AssigneeID,
		ProjectID:   *data.ProjectID,
		Status:      *data.Status,
	}
	if data.CompletedAt != nil {
		res.CompletedAt = *data.CompletedAt
	}
	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
