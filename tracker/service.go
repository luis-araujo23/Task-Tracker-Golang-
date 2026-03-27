package tracker

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ErrTaskNotFound = errors.New("tarea no encontrada")

type Service struct {
	store *Store
}

// Crea el servicio de negocio usando el store de persistencia.
func NewService(store *Store) *Service {
	return &Service{store: store}
}

// Crea un nuevo usuario validando que su nombre no esté vacío.
func (s *Service) CreateUser(name string) (User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return User{}, errors.New("el nombre del usuario no puede estar vacio")
	}

	users, err := s.store.LoadUsers()
	if err != nil {
		return User{}, err
	}

	newUser := User{
		ID:        nextUserID(users),
		Name:      name,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	users = append(users, newUser)
	if err := s.store.SaveUsers(users); err != nil {
		return User{}, err
	}

	return newUser, nil
}

// Lista usuarios ordenados por ID ascendente.
func (s *Service) ListUsers() ([]User, error) {
	users, err := s.store.LoadUsers()
	if err != nil {
		return nil, err
	}

	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })
	return users, nil
}

// Agrega una nueva tarea con estado inicial "todo" y usuario creador.
func (s *Service) AddTask(description string, createdByID int) (Task, error) {
	description = strings.TrimSpace(description)
	if description == "" {
		return Task{}, errors.New("la descripcion no puede estar vacia")
	}
	if createdByID <= 0 {
		return Task{}, errors.New("el id del usuario debe ser un numero entero positivo")
	}

	tasks, err := s.store.LoadTasks()
	if err != nil {
		return Task{}, err
	}

	users, err := s.store.LoadUsers()
	if err != nil {
		return Task{}, err
	}

	creator, ok := findUserByID(users, createdByID)
	if !ok {
		return Task{}, fmt.Errorf("usuario no encontrado: id %d", createdByID)
	}

	now := time.Now().Format(time.RFC3339)
	newTask := Task{
		ID:          nextID(tasks),
		Description: description,
		Status:      StatusTodo,
		CreatedByID: creator.ID,
		CreatedBy:   creator.Name,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tasks = append(tasks, newTask)
	if err := s.store.SaveTasks(tasks); err != nil {
		return Task{}, err
	}

	return newTask, nil
}

// Actualiza la descripción de una tarea existente por ID.
func (s *Service) UpdateTask(id int, newDescription string) (Task, error) {
	newDescription = strings.TrimSpace(newDescription)
	if newDescription == "" {
		return Task{}, errors.New("la descripcion no puede estar vacia")
	}

	tasks, err := s.store.LoadTasks()
	if err != nil {
		return Task{}, err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Description = newDescription
			tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			if err := s.store.SaveTasks(tasks); err != nil {
				return Task{}, err
			}
			return tasks[i], nil
		}
	}

	return Task{}, fmt.Errorf("%w: id %d", ErrTaskNotFound, id)
}

// Elimina una tarea por ID.
func (s *Service) DeleteTask(id int) error {
	tasks, err := s.store.LoadTasks()
	if err != nil {
		return err
	}

	idx := -1
	for i := range tasks {
		if tasks[i].ID == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fmt.Errorf("%w: id %d", ErrTaskNotFound, id)
	}

	tasks = append(tasks[:idx], tasks[idx+1:]...)
	return s.store.SaveTasks(tasks)
}

// Cambia el estado de una tarea (todo, in-progress, done).
func (s *Service) MarkTask(id int, status string) (Task, error) {
	if !IsValidStatus(status) {
		return Task{}, errors.New("estado invalido")
	}

	tasks, err := s.store.LoadTasks()
	if err != nil {
		return Task{}, err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			if err := s.store.SaveTasks(tasks); err != nil {
				return Task{}, err
			}
			return tasks[i], nil
		}
	}

	return Task{}, fmt.Errorf("%w: id %d", ErrTaskNotFound, id)
}

// Lista tareas, con filtro opcional por estado.
func (s *Service) ListTasks(filterStatus string) ([]Task, error) {
	if filterStatus != "" && !IsValidStatus(filterStatus) {
		return nil, errors.New("filtro de estado invalido; use: todo, in-progress, done")
	}

	tasks, err := s.store.LoadTasks()
	if err != nil {
		return nil, err
	}

	if filterStatus == "" {
		sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
		return tasks, nil
	}

	filtered := make([]Task, 0, len(tasks))
	for _, t := range tasks {
		if t.Status == filterStatus {
			filtered = append(filtered, t)
		}
	}

	sort.Slice(filtered, func(i, j int) bool { return filtered[i].ID < filtered[j].ID })
	return filtered, nil
}

// Convierte y valida el ID recibido por CLI.
func ParseID(idArg string) (int, error) {
	id, err := strconv.Atoi(idArg)
	if err != nil || id <= 0 {
		return 0, errors.New("el id debe ser un numero entero positivo")
	}
	return id, nil
}

// Calcula el siguiente ID incremental disponible.
func nextID(tasks []Task) int {
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}
	return maxID + 1
}

// Calcula el siguiente ID incremental disponible para usuarios.
func nextUserID(users []User) int {
	maxID := 0
	for _, u := range users {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	return maxID + 1
}

// Busca un usuario por ID para asociarlo como creador de tarea.
func findUserByID(users []User, id int) (User, bool) {
	for _, u := range users {
		if u.ID == id {
			return u, true
		}
	}
	return User{}, false
}
