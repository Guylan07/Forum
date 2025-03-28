package main

import (
	"fmt"
	"net/http"
	"authentification" // Assurez-vous que ce package est bien importé et accessible
)

func main() {
	// Définir les routes AVANT le lancement du serveur
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Bienvenue sur mon forum !")
	})

	// Enregistrer les handlers pour /register et /login
	http.HandleFunc("/register", authentification.RegisterHandler)
	http.HandleFunc("/login", authentification.LoginHandler)

	port := "8081"
	fmt.Printf("Lancement du serveur sur le port %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Pas de démarrage le serveur bug : %s\n", err)
	}
}
