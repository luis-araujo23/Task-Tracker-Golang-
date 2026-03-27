# Task Tracker (Go)

Aplicacion de linea de comandos para gestionar tareas con almacenamiento local en JSON.

## Estructura

- `main.go`: entrada CLI y parseo de comandos posicionales
- `tracker/task.go`: modelos (`Task` y `User`) y estados
- `tracker/store.go`: persistencia JSON en `tasks.json` y `users.json`
- `tracker/service.go`: logica de negocio y validaciones

## Ejecutar

```bash
go run .
```

## Comandos

- `create-user "nombre"`
- `list-users`
- `add <user_id> "descripcion"`
- `update <id> "nueva descripcion"`
- `delete <id>`
- `mark-in-progress <id>`
- `mark-done <id>`
- `list`
- `list done`
- `list todo`
- `list in-progress`

## Flujo recomendado

1. Crear usuario: `create-user "Juan Perez"`
2. Ver usuarios y su ID: `list-users`
3. Crear tarea asociada a un usuario: `add 1 "Comprar pan"`
4. Listar tareas mostrando creador: `list`
