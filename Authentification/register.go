package authentification

import (
	"fmt"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		fmt.Printf("Enregistrement réussi : %s\n,", username)
		fmt.Fprintln(w, "Utilisateur enregistré avec succès")
	} else {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}
