package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var ErrCorruptedJSON = errors.New("el archivo JSON esta corrupto o mal formado")

type Store struct {
	TasksFilePath string
	UsersFilePath string
}

// Crea un store asociado al archivo JSON de tareas.
func NewStore(tasksFilePath string, usersFilePath string) *Store {
	return &Store{
		TasksFilePath: tasksFilePath,
		UsersFilePath: usersFilePath,
	}
}

// Carga tareas desde JSON; crea el archivo si no existe.
func (s *Store) LoadTasks() ([]Task, error) {
	if _, err := os.Stat(s.TasksFilePath); errors.Is(err, os.ErrNotExist) {
		if err := s.SaveTasks([]Task{}); err != nil {
			return nil, fmt.Errorf("no se pudo crear el archivo JSON: %w", err)
		}
		return []Task{}, nil
	}

	data, err := os.ReadFile(s.TasksFilePath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo JSON: %w", err)
	}

	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorruptedJSON, err)
	}

	return tasks, nil
}

// Guarda todas las tareas en el archivo JSON.
func (s *Store) SaveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("no se pudo serializar el JSON: %w", err)
	}

	if err := os.WriteFile(s.TasksFilePath, data, 0o644); err != nil {
		return fmt.Errorf("no se pudo escribir el archivo JSON: %w", err)
	}

	return nil
}

// Carga usuarios desde JSON; crea el archivo si no existe.
func (s *Store) LoadUsers() ([]User, error) {
	if _, err := os.Stat(s.UsersFilePath); errors.Is(err, os.ErrNotExist) {
		if err := s.SaveUsers([]User{}); err != nil {
			return nil, fmt.Errorf("no se pudo crear el archivo JSON de usuarios: %w", err)
		}
		return []User{}, nil
	}

	data, err := os.ReadFile(s.UsersFilePath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo JSON de usuarios: %w", err)
	}

	if len(data) == 0 {
		return []User{}, nil
	}

	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorruptedJSON, err)
	}

	return users, nil
}

// Guarda todos los usuarios en el archivo JSON.
func (s *Store) SaveUsers(users []User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("no se pudo serializar el JSON de usuarios: %w", err)
	}

	if err := os.WriteFile(s.UsersFilePath, data, 0o644); err != nil {
		return fmt.Errorf("no se pudo escribir el archivo JSON de usuarios: %w", err)
	}

	return nil
}
