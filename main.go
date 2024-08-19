package main

import (
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Todo struct {
	Id      int
	Message string
}

func main() {

	data := map[string][]Todo{
		"Todos": {
			{Id: 1, Message: "Read newspaper"},
		},
	}

	tmpl := template.Must(template.ParseFiles("index.html"))

	todosHandler := func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data)
	}

	addTodoHandler := func(w http.ResponseWriter, r *http.Request) {
		message := r.PostFormValue("message")
		todo := Todo{Id: len(data["Todos"]) + 1, Message: message}
		data["Todos"] = append(data["Todos"], todo)
		tmpl.ExecuteTemplate(w, "todo-list-element", todo)
	}

	deleteTodoHandler := func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PostFormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		for i, todo := range data["Todos"] {
			if todo.Id == id {
				data["Todos"] = append(data["Todos"][:i], data["Todos"][i+1:]...)
				break
			}
		}

		w.WriteHeader(http.StatusOK)
	}

	http.HandleFunc("/", todosHandler)
	http.HandleFunc("/add-todo", addTodoHandler)
	http.HandleFunc("/delete-todo", deleteTodoHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
