package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	store := Create()
	route := mux.NewRouter()

	route.HandleFunc("/store/{key}", StoreGet(store)).Methods("GET")
	route.HandleFunc("/store/{key}/{val}", StorePut(store)).Methods("PUT")

	log.Printf("Server started on port 8080")

	err := http.ListenAndServe(":8080", route)
	log.Fatal(err)
}

func StoreGet(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		err := store.View(func(tx *Transaction) error {
			if val, ok := store.data[key]; !ok {
				return errors.New("entry not found")
			} else {
				fmt.Fprint(w, key+": "+val)
				return nil
			}
		})
		if err != nil {
			http.Error(w, "entry not found", http.StatusNotFound)
		}
	}
}

func StorePut(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		val := vars["val"]

		store.Update(func(tx *Transaction) error {
			tx.Set(key, val)
			return nil
		})
	}
}
