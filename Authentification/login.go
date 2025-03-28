package authentification

import (
    "fmt"
    "net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        password := r.FormValue("password")

        fmt.Printf("Tentative de connexion : %s\n", username)
        fmt.Fprintln(w, "Connexion réussie")
    } else {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
    }
}
