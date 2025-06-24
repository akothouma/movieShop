package main

import (
    "log"
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("./frontend/static"))
    http.Handle("/", fs)

    log.Println("Starting frontend server on :3000")
    err := http.ListenAndServe(":3000", nil)
    if err != nil {
        log.Fatalf("Frontend server error: %v", err)
    }
}