package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Student struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Class string `json:"class"`
}

var students = []Student{
	{ID: "1", Name: "Іван Перший", Class: "10A"},
	{ID: "2", Name: "Другий Тестовий", Class: "10A"},
}

var teachers = map[string]bool{
	"teacher123": true,
}

var mutex sync.RWMutex

func main() {
	http.HandleFunc("/student/", handleStudent)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не підтримується", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/student/"):]

	userID := r.Header.Get("UserID")
	if !teachers[userID] {
		http.Error(w, "Доступ заборонено", http.StatusForbidden)
		return
	}

	mutex.RLock()
	defer mutex.RUnlock()

	student, found := findStudentByID(id)
	if !found {
		http.Error(w, "Учень не знайдений", http.StatusNotFound)
		return
	}

	jsonStudent, err := json.Marshal(student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonStudent)
}

func findStudentByID(id string) (Student, bool) {
	for _, student := range students {
		if student.ID == id {
			return student, true
		}
	}
	return Student{}, false
}
