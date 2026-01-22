package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Food", Description: "Description for Food"},
	{ID: 2, Name: "Drink", Description: "Description for Drink"},
	{ID: 3, Name: "Clothes", Description: "Description for Clothes"},
}

func main() {

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			var newCategory Category
			err := json.NewDecoder(r.Body).Decode(&newCategory)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			newCategory.ID = len(categories) + 1
			categories = append(categories, newCategory)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newCategory)
		}
	})

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid Categories ID", http.StatusBadRequest)
				return
			}

			// Cari id category
			for _, category := range categories {
				if category.ID == id {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(category)
					return
				}
			}
			http.Error(w, "Category not found", http.StatusNotFound)
		} else if r.Method == "PUT" {
			idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid Category ID", http.StatusBadRequest)
				return
			}

			var updateCategory Category
			err = json.NewDecoder(r.Body).Decode(&updateCategory)
			if err != nil {
				http.Error(w, "Invalid Category ID", http.StatusBadRequest)
				return
			}

			for i := range categories {
				if categories[i].ID == id {
					updateCategory.ID = id
					categories[i] = updateCategory

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(updateCategory)
					return
				}
			}
			http.Error(w, "Category Not Found", http.StatusNotFound)
		} else if r.Method == "DELETE" {
			idStr := strings.TrimPrefix(r.URL.Path, "/categories/")

			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid Category ID", http.StatusBadRequest)
				return
			}

			for i, category := range categories {
				if category.ID == id {
					categories = append(categories[:i], categories[i+1:]...)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(map[string]string{
						"message": "success deleted",
					})

					return
				}
			}

			http.Error(w, "Category not found", http.StatusNotFound)
		}

	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "API Running",
		})
	})

	fmt.Println("Server started at :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("error starting server:", err)
	}

}
