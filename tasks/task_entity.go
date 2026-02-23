package tasks

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type CreateTaskDto struct {
	Name string `json:"name"`
}

func (c *CreateTaskDto) Validate() map[string][]string {
	var errorsMap map[string][]string
	var nameErrors []string

	if c.Name == "" {
		nameErrors = append(nameErrors, "name is required")
	}

	if len(c.Name) < 3 {
		nameErrors = append(nameErrors, "name must be at least 3 characters")
	}

	if len(nameErrors) > 0 {
		errorsMap = map[string][]string{
			"name": nameErrors,
		}

		return errorsMap
	}

	return nil
}

type UpdateTaskNameDto struct {
	Name string `json:"name"`
}

func (u *UpdateTaskNameDto) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

type Task struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Completed bool       `json:"completed"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func NewTask(name string) *Task {
	return &Task{
		ID:        uuid.NewString(),
		Name:      name,
		CreatedAt: time.Now(),
		Completed: false,
	}
}

func (t *Task) ToogleTaskCompleted() {
	now := time.Now()

	t.Completed = !t.Completed
	t.UpdatedAt = &now
}

func (t *Task) UpdateTaskName(name string) {
	now := time.Now()

	t.Name = name
	t.UpdatedAt = &now
}
