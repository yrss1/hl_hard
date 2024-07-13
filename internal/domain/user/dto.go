package user

import (
	"errors"
)

type Request struct {
	ID       string  `json:"id"`
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Role     *string `json:"role"`
}

func (s *Request) Validate() error {
	if s.FullName == nil {
		return errors.New("full_name: cannot be blank")
	}

	if s.Email == nil {
		return errors.New("email: cannot be blank")
	}

	if s.Role == nil {
		return errors.New("role: cannot be blank")
	}

	return nil
}

type Response struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:       data.ID,
		FullName: *data.FullName,
		Email:    *data.Email,
		Role:     *data.Role,
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
