package storage

import (
	d "to-do-list/internal/storage/postgresql"
	task "to-do-list/models/task"
)

type TaskRepository struct {
	storage *d.Database
}

func New(storage *d.Database) *TaskRepository {
	return &TaskRepository{
		storage: storage,
	}
}


func (r *TaskRepository) GetTaskByID(id int) (*task.Task, error) {
	return r.storage.QueryTaskByID(id)
}

func (r *TaskRepository) GetAllTasks() ([]task.Task, error) {
	return r.storage.QueryTasksAll()
}

func (r *TaskRepository) CreateTask(task task.Task) (*task.Task, error) {
	return r.storage.InsertTaskRecord(&task)
}

func (r *TaskRepository) UpdateTask(task task.Task) error {
	return r.storage.UpdateTaskRecord(task)
}

func (r *TaskRepository) DeleteTaskByID(id int) error {
	return r.storage.DeleteTaskRecordByID(id)
}
