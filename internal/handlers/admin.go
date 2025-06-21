package handlers

import (
    "fmt"
    "net/http"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Добро пожаловать, админ!")
}