package main

import (
    "net/http"
    "log"
    "math/rand"
    "time"
)

func main() {
    rand.Seed(time.Now().UnixNano())
    Setup()
    http.HandleFunc("/", ServePage)
    err := http.ListenAndServe(":300", nil)
    if err != nil {
        log.Fatal(err)
    }
}
