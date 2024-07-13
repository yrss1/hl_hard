package tasker

import (
	"hard/internal/domain/project"
	"hard/internal/domain/task"
	"hard/internal/domain/user"
)

type Configuration func(s *Service) error

type Service struct {
	userRepository    user.Repository
	taskRepository    task.Repository
	projectRepository project.Repository
}

func New(configs ...Configuration) (s *Service, err error) {
	s = &Service{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithUserRepository(userRepository user.Repository) Configuration {
	return func(s *Service) error {
		s.userRepository = userRepository
		return nil
	}
}

func WithTaskRepository(taskRepository task.Repository) Configuration {
	return func(s *Service) error {
		s.taskRepository = taskRepository
		return nil
	}
}

func WithProjectRepository(projectRepository project.Repository) Configuration {
	return func(s *Service) error {
		s.projectRepository = projectRepository
		return nil
	}
}
