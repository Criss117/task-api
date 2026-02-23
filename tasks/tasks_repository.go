package tasks

import "strings"

type Filters struct {
	Select string
	Name   string
}

type TasksRepository struct {
	Tasks []*Task
}

func NewTasksRepository() *TasksRepository {
	learnGoTask := NewTask("Learn Go")
	learnDockerTask := NewTask("Learn Docker")
	learnKubernetesTask := NewTask("Learn Kubernetes")

	return &TasksRepository{
		Tasks: []*Task{
			learnGoTask,
			learnDockerTask,
			learnKubernetesTask,
		},
	}
}

func (r *TasksRepository) AddTask(task *Task) {
	r.Tasks = append(r.Tasks, task)
}

func (r *TasksRepository) GetTask(id string) *Task {
	for _, task := range r.Tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func (r *TasksRepository) GetAllTasks(filters Filters) []*Task {
	filteredTasks := make([]*Task, len(r.Tasks))
	copy(filteredTasks, r.Tasks)

	filteredTasks = filterTaskByName(filteredTasks, filters.Name)
	filteredTasks = filterTaskBySelect(filteredTasks, filters.Select)

	if len(filteredTasks) == 0 {
		return []*Task{}
	}

	return filteredTasks
}

func (r *TasksRepository) UpdateTask(task *Task) {
	for i, t := range r.Tasks {
		if t.ID == task.ID {
			r.Tasks[i] = task
		}
	}
}

func (r *TasksRepository) DeleteTask(id string) {
	var newTasks []*Task
	for _, task := range r.Tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)
		}
	}
	r.Tasks = newTasks
}

func filterTaskByName(tasks []*Task, name string) []*Task {

	if name == "" {
		return tasks
	}
	var newTasks []*Task

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Name), strings.ToLower(name)) {
			newTasks = append(newTasks, task)
		}
	}

	return newTasks
}

func filterTaskBySelect(tasks []*Task, selectFilter string) []*Task {
	filter := "all"

	switch selectFilter {
	case "completed":
		filter = "completed"
	case "uncompleted":
		filter = "uncompleted"
	case "all":
		filter = "all"
	default:
		filter = "all"
	}

	if filter == "all" {
		return tasks
	}

	var newTasks []*Task

	for _, task := range tasks {
		if task.Completed == true && filter == "completed" {
			newTasks = append(newTasks, task)
		}

		if task.Completed == false && filter == "uncompleted" {
			newTasks = append(newTasks, task)
		}
	}

	return newTasks
}
