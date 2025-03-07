package main

import (
    "fmt"
    "log"
    "net/http"
    "text/template"
)

func main() {
    mux := http.NewServeMux()

    // 1) Route pour la page d'accueil
    mux.HandleFunc("/", HomeHandler)

    // 2) Servir les fichiers statiques (CSS, JS, images)
    fs := http.FileServer(http.Dir("./web/static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    // 3) Lancement du serveur sur le port 8080
    fmt.Println("Server running on http://127.0.0.1:8080")
    err := http.ListenAndServe(":8080", mux)
    if err != nil {
        log.Fatal("ListenAndServe error:", err)
    }
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Ici, on parse un template HTML situé dans web/templates/index.html
    tmpl := template.Must(template.ParseFiles(
        "web/templates/index.html",
        // vous pouvez ajouter d'autres fichiers comme "web/templates/layout.html"
    ))

    // On peut passer des données au template, ici on envoie juste une string
    data := "Bienvenue sur mon forum !"
    err := tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
