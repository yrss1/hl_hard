package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hard/pkg/helpers"
	"hard/pkg/store"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"hard/internal/domain/task"
	"hard/internal/service/tasker"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) List(ctx context.Context) (dest []task.Entity, err error) {
	args := m.Called(ctx)
	return args.Get(0).([]task.Entity), args.Error(1)
}

func (m *MockTaskRepository) Add(ctx context.Context, data task.Entity) (id string, err error) {
	args := m.Called(ctx, data)
	return args.String(0), args.Error(1)
}

func (m *MockTaskRepository) Get(ctx context.Context, id string) (dest task.Entity, err error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return task.Entity{}, store.ErrorNotFound
	}

	return args.Get(0).(task.Entity), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, id string, dest task.Entity) (err error) {
	args := m.Called(ctx, id, dest)

	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) (err error) {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTaskRepository) Search(ctx context.Context, data task.Entity) (dest []task.Entity, err error) {
	args := m.Called(ctx, data)
	return args.Get(0).([]task.Entity), args.Error(1)
}

func TestList(t *testing.T) {
	mockTasks := []task.Entity{
		{
			ID:          "1",
			Title:       helpers.GetStringPtr("Task 1"),
			Description: helpers.GetStringPtr("Description of Task 1"),
			Priority:    helpers.GetStringPtr("High"),
			Status:      helpers.GetStringPtr("Active"),
			AssigneeID:  helpers.GetStringPtr("1"),
			ProjectID:   helpers.GetStringPtr("1"),
		},
		{
			ID:          "2",
			Title:       helpers.GetStringPtr("Task 2"),
			Description: helpers.GetStringPtr("Description of Task 2"),
			Priority:    helpers.GetStringPtr("Medium"),
			Status:      helpers.GetStringPtr("Pending"),
			AssigneeID:  helpers.GetStringPtr("2"),
			ProjectID:   helpers.GetStringPtr("2"),
		},
	}

	tests := []struct {
		name           string
		mockRepoOutput []task.Entity
		mockRepoError  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful List",
			mockRepoOutput: mockTasks,
			mockRepoError:  nil,
			expectedStatus: http.StatusOK,
			expectedBody: `{"data":[{"id":"1","title":"Task 1","description":"Description of Task 1",
				"priority":"High","status":"Active","assignee_id":"1","project_id":"1","completed_at":""},
				{"id":"2","title":"Task 2","description":"Description of Task 2",
				"priority":"Medium","status":"Pending","assignee_id":"2","project_id":"2","completed_at":""}],
				"success":true}`,
		},
		{
			name:           "Empty List",
			mockRepoOutput: []task.Entity{},
			mockRepoError:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":[],"success":true}`,
		},
		{
			name:           "Internal Server Error",
			mockRepoOutput: nil,
			mockRepoError:  errors.New("repository error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"repository error","success":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			mockRepo.On("List", mock.Anything).Return(tt.mockRepoOutput, tt.mockRepoError)

			taskService, _ := tasker.New(tasker.WithTaskRepository(mockRepo))
			taskHandler := NewTaskHandler(taskService)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/tasks", taskHandler.list)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/tasks", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      string
		inputData      task.Entity
		mockRepoOutput string
		mockRepoError  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Successful Add",
			inputBody: `{"title":"test tester","description":"A new task description","priority":"High","status":"Active","assignee_id":"1","project_id":"2"}","description":"A new task description","priority":"High","status":"Active","assignee_id":"1","project_id":"2"}`,
			inputData: task.Entity{
				Title:       helpers.GetStringPtr("test tester"),
				Description: helpers.GetStringPtr("A new task description"),
				Priority:    helpers.GetStringPtr("High"),
				Status:      helpers.GetStringPtr("Active"),
				AssigneeID:  helpers.GetStringPtr("1"),
				ProjectID:   helpers.GetStringPtr("2"),
			},
			mockRepoOutput: "new-task-id",
			mockRepoError:  nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"data":{"id":"new-task-id","title":"test tester","description":"A new task description","priority":"High","status":"Active","assignee_id":"1","project_id":"2","completed_at":""},"success":true}`,
		},
		{
			name:           "Bad Request: Missing Title",
			inputBody:      `{"description":"A new task description","priority":"High","status":"Active","assignee_id":"1","project_id":"2"}`,
			mockRepoOutput: "new-task-id",
			mockRepoError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"data":{"id":"","title":null,"description":"A new task description","priority":"High","status":"Active","assignee_id":"1","project_id":"2","completed_at":null},"message":"title: cannot be blank","success":false}`,
		},
		{
			name:           "Invalid JSON Payload",
			inputBody:      `{"title":"test tester","description":"A new task description","priority":"High","status":"Active","assignee_id":"1","project_id":"2",}`,
			inputData:      task.Entity{},
			mockRepoOutput: "",
			mockRepoError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"data":{"id":"","title":null,"description":null,"priority":null,"status":null,"assignee_id":null,"project_id":null,"completed_at":null},"message":"invalid character '}' looking for beginning of object key string","success":false}`,
		},
		{
			name:           "Internal Server Error",
			inputBody:      `{"title":"test tester","description":"A new task description","priority":"High","status":"Active","assignee_id":"1","project_id":"2"}","description":"A new task description","priority":"High","status":"Active","assignee_id":"1","project_id":"2"}`,
			mockRepoOutput: "new-task-id",
			inputData: task.Entity{
				Title:       helpers.GetStringPtr("test tester"),
				Description: helpers.GetStringPtr("A new task description"),
				Priority:    helpers.GetStringPtr("High"),
				Status:      helpers.GetStringPtr("Active"),
				AssigneeID:  helpers.GetStringPtr("1"),
				ProjectID:   helpers.GetStringPtr("2"),
			},
			mockRepoError:  errors.New("repository error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"repository error","success":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)

			mockRepo.On("Add", mock.Anything, tt.inputData).Return(tt.mockRepoOutput, tt.mockRepoError)

			taskService, _ := tasker.New(tasker.WithTaskRepository(mockRepo))

			taskHandler := NewTaskHandler(taskService)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/tasks", taskHandler.add)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/tasks",
				bytes.NewBufferString(tt.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestGet(t *testing.T) {
	mockTask := task.Entity{
		ID:          "1",
		Title:       helpers.GetStringPtr("Design Homepage"),
		Description: helpers.GetStringPtr("Create a responsive homepage design"),
		Priority:    helpers.GetStringPtr("High"),
		Status:      helpers.GetStringPtr("Active"),
		AssigneeID:  helpers.GetStringPtr("4"),
		ProjectID:   helpers.GetStringPtr("1"),
	}

	tests := []struct {
		name           string
		mockRepoOutput task.Entity
		mockRepoError  error
		taskID         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful Get",
			mockRepoOutput: mockTask,
			mockRepoError:  nil,
			taskID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody: `{"data":{"id":"1","title":"Design Homepage","description":"Create a responsive homepage design",
		     "priority":"High","status":"Active","assignee_id":"4","project_id":"1","completed_at":""},"success":true}`,
		},
		{
			name:           "Task Not Found",
			mockRepoOutput: mockTask,
			mockRepoError:  store.ErrorNotFound,
			taskID:         "99",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"error not found","success":false}`,
		},
		{
			name:           "Internal Server Error",
			mockRepoOutput: mockTask,
			mockRepoError:  errors.New("repository error"),
			taskID:         "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"repository error","success":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			mockRepo.On("Get", mock.Anything, tt.taskID).Return(tt.mockRepoOutput, tt.mockRepoError)

			taskService, _ := tasker.New(tasker.WithTaskRepository(mockRepo))
			taskHandler := NewTaskHandler(taskService)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/tasks/:id", taskHandler.get)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/tasks/"+tt.taskID, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      string
		mockRepoError  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful Update",
			inputBody:      `{"title":"Updated Task","description":"This is an updated task","priority":"Medium","status":"InProgress","assignee_id":"2","project_id":"3"}`,
			mockRepoError:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":"ok","success":true}`,
		},
		{
			name:           "Invalid JSON Payload",
			inputBody:      `{"title":"Updated Task","description":"This is an updated task","priority":"Medium","status":"InProgress","assignee_id":"2","project_id":"3",}`,
			mockRepoError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"data":{"id":"","title":null,"description":null,"priority":null,"status":null,"assignee_id":null,"project_id":null,"completed_at":null},"message":"invalid character '}' looking for beginning of object key string","success":false}`,
		},
		{
			name:           "Task Not Found",
			inputBody:      `{"title":"Updated Task","description":"This is an updated task","priority":"Medium","status":"InProgress","assignee_id":"2","project_id":"3"}`,
			mockRepoError:  store.ErrorNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"error not found","success":false}`,
		},
		{
			name:           "Internal Server Error",
			inputBody:      `{"title":"Updated Task","description":"This is an updated task","priority":"Medium","status":"InProgress","assignee_id":"2","project_id":"3"}`,
			mockRepoError:  errors.New("repository error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"repository error","success":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			mockRepo.On("Update", mock.Anything, "mock-task-id", mock.AnythingOfType("task.Entity")).Return(tt.mockRepoError)

			taskService, _ := tasker.New(tasker.WithTaskRepository(mockRepo))
			taskHandler := NewTaskHandler(taskService)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.PUT("/tasks/:id", taskHandler.update)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/tasks/mock-task-id", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name           string
		mockRepoError  error
		taskID         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Successful Delete",
			mockRepoError:  nil,
			taskID:         "1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":"1","success":true}`,
		},
		{
			name:           "Task Not Found",
			mockRepoError:  store.ErrorNotFound,
			taskID:         "99",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"error not found","success":false}`,
		},
		{
			name:           "Internal Server Error",
			mockRepoError:  errors.New("repository error"),
			taskID:         "1",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"repository error","success":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			mockRepo.On("Delete", mock.Anything, tt.taskID).Return(tt.mockRepoError)

			taskService, _ := tasker.New(tasker.WithTaskRepository(mockRepo))
			taskHandler := NewTaskHandler(taskService)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.DELETE("/tasks/:id", taskHandler.delete)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/tasks/"+tt.taskID, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		mockRepoOutput []task.Entity
		mockRepoError  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful Search",
			queryParams: map[string]string{
				"title":       "test",
				"priority":    "High",
				"status":      "Active",
				"assignee_id": "1",
				"project_id":  "2",
			},
			mockRepoOutput: []task.Entity{
				{
					ID:          "1",
					Title:       helpers.GetStringPtr("test task"),
					Description: helpers.GetStringPtr("A test task description"),
					Priority:    helpers.GetStringPtr("High"),
					Status:      helpers.GetStringPtr("Active"),
					AssigneeID:  helpers.GetStringPtr("1"),
					ProjectID:   helpers.GetStringPtr("2"),
				},
				// Add more mock data as needed
			},
			mockRepoError:  nil,
			expectedStatus: http.StatusOK,
			expectedBody: `{"data":[{"id":"1","title":"test task","description":"A test task description",
			"priority":"High","status":"Active","assignee_id":"1","project_id":"2","completed_at":""}],"success":true}`,
		},
		{
			name:           "Missing Query Parameters",
			queryParams:    map[string]string{},
			mockRepoOutput: nil,
			mockRepoError:  nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"query parameters required","success":false}`,
		},
		{
			name:           "Empty Result Search",
			queryParams:    map[string]string{"title": "nonexistent"},
			mockRepoOutput: []task.Entity{},
			mockRepoError:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":[],"success":true}`,
		},
		{
			name: "Internal Server Error",
			queryParams: map[string]string{
				"title": "test",
			},
			mockRepoOutput: nil,
			mockRepoError:  errors.New("repository error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"repository error","success":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			mockRepo.On("Search", mock.Anything, mock.MatchedBy(func(data task.Entity) bool {
				return reflect.DeepEqual(data, task.Entity{
					Title:      helpers.GetStringPtr(tt.queryParams["title"]),
					Priority:   helpers.GetStringPtr(tt.queryParams["priority"]),
					Status:     helpers.GetStringPtr(tt.queryParams["status"]),
					AssigneeID: helpers.GetStringPtr(tt.queryParams["assignee_id"]),
					ProjectID:  helpers.GetStringPtr(tt.queryParams["project_id"]),
				})
			})).Return(tt.mockRepoOutput, tt.mockRepoError)

			taskService, _ := tasker.New(tasker.WithTaskRepository(mockRepo))
			taskHandler := NewTaskHandler(taskService)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/tasks/search", taskHandler.search)

			url := "/tasks/search?" + encodeQueryParams(tt.queryParams)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", url, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func encodeQueryParams(params map[string]string) string {
	var encodedParams []string
	for key, value := range params {
		encodedParams = append(encodedParams, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(encodedParams, "&")
}
