package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"task-tracker/tracker"
)

const jsonFileName = "tasks.json"
const usersFileName = "users.json"

// Inicializa el store/servicio y enruta el comando de la CLI.
func main() {
	store := tracker.NewStore(jsonFileName, usersFileName)
	service := tracker.NewService(store)

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	if err := execute(service, command, args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Ejecuta el comando recibido y delega en la capa de servicio.
func execute(service *tracker.Service, command string, args []string) error {
	switch command {
	case "create-user":
		if len(args) < 1 {
			return errors.New("uso: create-user \"nombre\"")
		}
		name := strings.Join(args, " ")
		u, err := service.CreateUser(name)
		if err != nil {
			return err
		}
		fmt.Printf("Usuario creado correctamente (id: %d, nombre: %s).\n", u.ID, u.Name)
		return nil

	case "list-users":
		if len(args) != 0 {
			return errors.New("uso: list-users")
		}
		users, err := service.ListUsers()
		if err != nil {
			return err
		}
		printUsers(users)
		return nil

	case "add":
		if len(args) < 2 {
			return errors.New("uso: add <user_id> \"descripcion\"")
		}
		userID, err := tracker.ParseID(args[0])
		if err != nil {
			return err
		}
		description := strings.Join(args[1:], " ")
		t, err := service.AddTask(description, userID)
		if err != nil {
			return err
		}
		fmt.Printf("Tarea agregada correctamente (id: %d, creada por: %s).\n", t.ID, t.CreatedBy)
		return nil

	case "update":
		if len(args) < 2 {
			return errors.New("uso: update <id> \"nueva descripcion\"")
		}
		id, err := tracker.ParseID(args[0])
		if err != nil {
			return err
		}
		description := strings.Join(args[1:], " ")
		t, err := service.UpdateTask(id, description)
		if err != nil {
			return err
		}
		fmt.Printf("Tarea %d actualizada correctamente.\n", t.ID)
		return nil

	case "delete":
		if len(args) != 1 {
			return errors.New("uso: delete <id>")
		}
		id, err := tracker.ParseID(args[0])
		if err != nil {
			return err
		}
		if err := service.DeleteTask(id); err != nil {
			return err
		}
		fmt.Printf("Tarea %d eliminada correctamente.\n", id)
		return nil

	case "mark-in-progress":
		if len(args) != 1 {
			return errors.New("uso: mark-in-progress <id>")
		}
		id, err := tracker.ParseID(args[0])
		if err != nil {
			return err
		}
		t, err := service.MarkTask(id, tracker.StatusInProgress)
		if err != nil {
			return err
		}
		fmt.Printf("Tarea %d marcada como in-progress.\n", t.ID)
		return nil

	case "mark-done":
		if len(args) != 1 {
			return errors.New("uso: mark-done <id>")
		}
		id, err := tracker.ParseID(args[0])
		if err != nil {
			return err
		}
		t, err := service.MarkTask(id, tracker.StatusDone)
		if err != nil {
			return err
		}
		fmt.Printf("Tarea %d marcada como done.\n", t.ID)
		return nil

	case "list":
		if len(args) > 1 {
			return errors.New("uso: list [todo|in-progress|done]")
		}

		filter := ""
		if len(args) == 1 {
			filter = args[0]
		}

		tasks, err := service.ListTasks(filter)
		if err != nil {
			return err
		}
		printTasks(tasks)
		return nil

	case "help", "--help", "-h":
		printHelp()
		return nil

	default:
		printHelp()
		return fmt.Errorf("comando no reconocido: %s", command)
	}
}

// Muestra las tareas en formato de tabla simple.
func printTasks(tasks []tracker.Task) {
	if len(tasks) == 0 {
		fmt.Println("No hay tareas para mostrar.")
		return
	}

	fmt.Println("ID | STATUS      | CREATED BY     | DESCRIPTION            | CREATED AT                | UPDATED AT")
	fmt.Println("---+-------------+----------------+------------------------+---------------------------+---------------------------")
	for _, t := range tasks {
		creator := t.CreatedBy
		if strings.TrimSpace(creator) == "" {
			creator = "N/A"
		}
		fmt.Printf("%-2d | %-11s | %-14s | %-22s | %-25s | %-25s\n", t.ID, t.Status, truncate(creator, 14), truncate(t.Description, 22), t.CreatedAt, t.UpdatedAt)
	}
}

// Muestra los usuarios en formato de tabla simple.
func printUsers(users []tracker.User) {
	if len(users) == 0 {
		fmt.Println("No hay usuarios para mostrar.")
		return
	}

	fmt.Println("ID | NAME                 | CREATED AT")
	fmt.Println("---+----------------------+---------------------------")
	for _, u := range users {
		fmt.Printf("%-2d | %-20s | %-25s\n", u.ID, truncate(u.Name, 20), u.CreatedAt)
	}
}

// Recorta texto largo para mantener columnas legibles.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// Imprime la ayuda con los comandos disponibles.
func printHelp() {
	fmt.Println("Task Tracker - CLI")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  create-user \"nombre\"")
	fmt.Println("  list-users")
	fmt.Println("  add <user_id> \"descripcion\"")
	fmt.Println("  update <id> \"nueva descripcion\"")
	fmt.Println("  delete <id>")
	fmt.Println("  mark-in-progress <id>")
	fmt.Println("  mark-done <id>")
	fmt.Println("  list")
	fmt.Println("  list done")
	fmt.Println("  list todo")
	fmt.Println("  list in-progress")
}
