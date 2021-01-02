package main

import (
	"encoding/json" // json implementa la codifacion y decodificacionde JSON
	"fmt" // fmt implementa formateo de funciones analogas a print y scant de C
	"io/ioutil" // ioutil implementa algunas funciones de E / S
	"log" // log impleta paquetes de registro para formateo de salida en este caso utilize FATAL
	"net/http" // http propociona implementaciones de servidor y cliente HTTP ( GET, POST, DELETE, PUSH)
	"strconv" // strconv implementa conversiones hacia y desde representacion de cadenas 

	"github.com/gorilla/mux" // gorilla/mux implementa un enrutador de solicitudes y despachador para hacer coincidir las solicitudes
)

// Typos de datos, especificando 
type task struct {
	ID      int    `json:"ID"` // int = entero
	Name    string `json:"Name"` // string = cadena de caracteres
	Content string `json:"Content"` // string = cadena de caracteres
}

type allTasks []task

// Typos de datos constantes
var tasks = allTasks{
	{
		ID:      1,
		Name:    "Tarea inicial",
		Content: "Crear api go",
	},

}

func indexRoute(w http.ResponseWriter, r *http.Request) { // Ruta principal
	fmt.Fprintf(w, "¡Bienvenido a mi API GO!")
}

func createTask(w http.ResponseWriter, r *http.Request) { // Ruta de crear tareas
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insertar datos de tarea válidos")
	}

	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

func getTasks(w http.ResponseWriter, r *http.Request) { // Ruta de mirar tareas
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getOneTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) { // Ruta de actualizar tareas
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(w, "identificación invalida")
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Ingrese datos válidos")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)

			updatedTask.ID = t.ID
			tasks = append(tasks, updatedTask)


			fmt.Fprintf(w, "La tarea con ID% v se ha actualizado correctamente", taskID)
		}
	}

}

func deleteTask(w http.ResponseWriter, r *http.Request) { // Ruta de eliminar tareeas
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "ID de usuario invalido")
		return
	}

	for i, t := range tasks {
		if t.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "La tarea con ID %v se ha eliminado correctamente", taskID)
		}
	}
}

func main() { // Func main principal a donde se ejecuta todo el programa
	router := mux.NewRouter().StrictSlash(true) // Creando las rutas 

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", createTask).Methods("POST") // POST = enviar
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getOneTask).Methods("GET") // GET = obtener
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE") // DELETE = eliminar
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT") // PUT = poner

	log.Fatal(http.ListenAndServe(":3000", router)) // Si el servidor no funciona, automaticamente se abrira el la ruta 3000
}
