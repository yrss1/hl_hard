package tasker

import (
	"context"
	"errors"
	"fmt"
	"hard/internal/domain/project"
	"hard/internal/domain/task"
	"hard/pkg/store"
)

func (s *Service) ListProjects(ctx context.Context) (res []project.Response, err error) {
	data, err := s.projectRepository.List(ctx)
	if err != nil {
		fmt.Printf("failed to select: %v", err)
	}

	res = project.ParseFromEntities(data)

	return
}

func (s *Service) CreateProject(ctx context.Context, req project.Request) (res project.Response, err error) {
	data := project.Entity{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		ManagerID:   req.ManagerID,
	}

	data.ID, err = s.projectRepository.Add(ctx, data)
	if err != nil {
		fmt.Printf("failed to create: %v\n", err)
		return
	}

	res = project.ParseFromEntity(data)

	return
}

func (s *Service) GetProject(ctx context.Context, id string) (res project.Response, err error) {
	data, err := s.projectRepository.Get(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to get by id: %v", err)
		return
	}

	res = project.ParseFromEntity(data)

	return
}

func (s *Service) UpdateProject(ctx context.Context, id string, req project.Request) (err error) {
	data := project.Entity{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		ManagerID:   req.ManagerID,
	}

	err = s.projectRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to update by id: %v\n", err)
		return
	}

	return
}

func (s *Service) DeleteProject(ctx context.Context, id string) (err error) {
	err = s.projectRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to delete by id %v\n", err)
		return
	}

	return
}

func (s *Service) SearchProjects(ctx context.Context, req project.Request) (res []project.Response, err error) {
	searchData := project.Entity{
		Description: req.Description,
		ManagerID:   req.ManagerID,
	}

	data, err := s.projectRepository.Search(ctx, searchData)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to search projects: %v\n", err)
		return
	}

	res = project.ParseFromEntities(data)

	return
}
func (s *Service) GetTasksByProject(ctx context.Context, id string) (res []task.Response, err error) {
	data, err := s.projectRepository.ListTasks(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to search tasks by project: %v\n", err)
		return
	}

	res = task.ParseFromEntities(data)

	return
}
