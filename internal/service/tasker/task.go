package tasker

import (
	"context"
	"errors"
	"fmt"
	"hard/internal/domain/task"
	"hard/pkg/store"
	"time"
)

func (s *Service) ListTasks(ctx context.Context) (res []task.Response, err error) {
	data, err := s.taskRepository.List(ctx)
	if err != nil {
		fmt.Printf("failed to select: %v", err)
		return
	}

	res = task.ParseFromEntities(data)

	return
}

func (s *Service) CreateTask(ctx context.Context, req task.Request) (res task.Response, err error) {
	var completedAt time.Time
	if req.CompletedAt != nil {
		completedAt, err = time.Parse("2006-01-02", *req.CompletedAt)
		if err != nil {
			return res, fmt.Errorf("failed to parse: %v", err)
		}
	}
	data := task.Entity{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
		ProjectID:   req.ProjectID,
		CompletedAt: &completedAt,
	}
	data.ID, err = s.taskRepository.Add(ctx, data)
	if err != nil {
		fmt.Printf("failed to create: %v\n", err)
		return
	}

	res = task.ParseFromEntity(data)

	return
}

func (s *Service) GetTask(ctx context.Context, id string) (res task.Response, err error) {
	data, err := s.taskRepository.Get(ctx, id)
	if err != nil {
		fmt.Printf("failed to get by id: %v", err)
		return
	}

	res = task.ParseFromEntity(data)

	return
}

func (s *Service) UpdateTask(ctx context.Context, id string, req task.Request) (err error) {
	var completedAt time.Time
	if req.CompletedAt != nil {
		completedAt, err = time.Parse("2006-01-02", *req.CompletedAt)
		if err != nil {
			return fmt.Errorf("failed to parse: %v", err)
		}
	}

	data := task.Entity{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
		ProjectID:   req.ProjectID,
		CompletedAt: &completedAt,
	}
	err = s.taskRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to update by id: %v\n", err)
		return
	}

	return
}

func (s *Service) DeleteTask(ctx context.Context, id string) (err error) {
	err = s.taskRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to delete by id: %v\n", err)
		return
	}

	return
}

func (s *Service) SearchTasks(ctx context.Context, req task.Request) (res []task.Response, err error) {
	data := task.Entity{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		AssigneeID:  req.AssigneeID,
		ProjectID:   req.ProjectID,
	}

	data2, err := s.taskRepository.Search(ctx, data)
	if err != nil {
		fmt.Printf("failed to search tasks: %v\n", err)
		return
	}

	res = task.ParseFromEntities(data2)

	return
}
