package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int
	Name    string
	Content string
}

type allTask []task

var tasks = allTask{
	{
		ID:      1,
		Name:    " La primera Tarea",
		Content: "Nunca Rendirse",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Inserte una tarea validad")
	}
	json.Unmarshal(reqBody, &newTask)

	//se crea el id de la tarea y se agrega la nueva tarea
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTaskOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID invalido, ingrea solamente numeros")
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(task)
		}

	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	delID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID invalido, ingrea solamente numeros")
		return
	}

	for i, task := range tasks {
		if task.ID == delID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "la tarea con ID %v eliminada", delID)
		}
	}

}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(w, "ID invalido, ingrea solamente numeros")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "ingrese datos validos")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = taskID
			tasks = append(tasks, updatedTask)

			fmt.Fprintf(w, "Tarea con el ID %v actualizada", taskID)
		}
	}
}

func indexRouter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "___bienvenido a mi API____")
}

// urls

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTaskOne).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/", indexRouter)
	log.Fatal(http.ListenAndServe(":3000", router))
}
