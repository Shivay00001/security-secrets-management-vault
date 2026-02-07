package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "sync"
    "github.com/gorilla/mux"
)

type Secret struct {
    ID    string `json:"id"`
    Value string `json:"value"`
}

var store = struct {
    sync.RWMutex
    data map[string]string
}{data: make(map[string]string)}

func CreateSecret(w http.ResponseWriter, r *http.Request) {
    var secret Secret
    if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    store.Lock()
    store.data[secret.ID] = secret.Value
    store.Unlock()
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"status": "stored"})
}

func GetSecret(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    store.RLock()
    val, ok := store.data[id]
    store.RUnlock()
    if !ok {
        http.Error(w, "Secret not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(w).Encode(Secret{ID: id, Value: val})
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/secrets", CreateSecret).Methods("POST")
    r.HandleFunc("/secrets/{id}", GetSecret).Methods("GET")
    fmt.Println("Vault listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}