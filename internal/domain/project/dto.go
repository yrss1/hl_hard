package project

import (
	"errors"
	"time"
)

type Request struct {
	ID          string  `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	ManagerID   *string `json:"manager_id"`
}

func (s *Request) Validate() error {

	if s.Title == nil {
		return errors.New("title: cannot be blank")
	}

	if s.StartDate == nil {
		return errors.New("start_date: cannot be blank")
	}

	if s.ManagerID == nil {
		return errors.New("manager_id: cannot be blank")
	}

	return nil
}

func (s *Request) IsEmpty() bool {
	return s.Title == nil && s.Description == nil &&
		s.StartDate == nil &&
		s.EndDate == nil &&
		s.ManagerID == nil
}

func TimeParse(data *string) (res time.Time, err error) {
	res, err = time.Parse("2006-01-02", *data)
	if err != nil {
		return res, err
	}

	return
}

type Response struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	ManagerID   string    `json:"manager_id"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:          data.ID,
		Title:       *data.Title,
		Description: *data.Description,
		ManagerID:   *data.ManagerID,
	}
	if data.StartDate != nil {
		res.StartDate = *data.StartDate
	}
	if data.EndDate != nil {
		res.EndDate = *data.EndDate
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
