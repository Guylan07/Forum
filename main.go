package main

import (
	"fmt"
	"net/http"
)

func main () {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Bienvenue sur mon forum !")
	})
	
	port := "8081"
	fmt.Printf("Lancement du serveur sur le port 8081", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Pas de démarrage le serveur bug", err)
	}
}