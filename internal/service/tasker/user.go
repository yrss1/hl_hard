package tasker

import (
	"context"
	"errors"
	"fmt"
	"hard/internal/domain/task"
	"hard/internal/domain/user"
	"hard/pkg/store"
)

func (s *Service) ListUsers(ctx context.Context) (res []user.Response, err error) {
	data, err := s.userRepository.List(ctx)
	if err != nil {
		fmt.Printf("failed to select: %v", err)
		return
	}

	res = user.ParseFromEntities(data)

	return
}

func (s *Service) CreateUser(ctx context.Context, req user.Request) (res user.Response, err error) {
	data := user.Entity{
		FullName: req.FullName,
		Email:    req.Email,
		Role:     req.Role,
	}

	data.ID, err = s.userRepository.Add(ctx, data)
	if err != nil {
		fmt.Printf("failed to create: %v\n", err)
		return
	}

	res = user.ParseFromEntity(data)

	return
}

func (s *Service) GetUser(ctx context.Context, id string) (res user.Response, err error) {
	data, err := s.userRepository.Get(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to get by id: %v", err)
		return
	}

	res = user.ParseFromEntity(data)

	return
}

func (s *Service) UpdateUser(ctx context.Context, id string, req user.Request) (err error) {
	data := user.Entity{
		FullName: req.FullName,
		Email:    req.Email, // тут было типа &req.Email
		Role:     req.Role,
	}
	err = s.userRepository.Update(ctx, id, data)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to update by id: %v\n", err)
		return
	}

	return
}

func (s *Service) DeleteUser(ctx context.Context, id string) (err error) {
	err = s.userRepository.Delete(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		//fmt.Printf("failed to delete by id: %v\n", err)
		return
	}

	return
}

func (s *Service) SearchUser(ctx context.Context, name string, email string) (res []user.Response, err error) {
	data, err := s.userRepository.Search(ctx, name, email)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to search user: %v\n", err)
		return
	}

	res = user.ParseFromEntities(data)

	return
}

func (s *Service) GetTasksByUser(ctx context.Context, id string) (res []task.Response, err error) {
	data, err := s.userRepository.ListTasks(ctx, id)
	if err != nil && !errors.Is(err, store.ErrorNotFound) {
		fmt.Printf("failed to search tasks by user: %v\n", err)
		return
	}

	res = task.ParseFromEntities(data)

	return
}
