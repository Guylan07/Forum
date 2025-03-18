package main

import ( //Les importation de packages
	"fmt" // Sert à afficher du texte et des réponses HTTP
	"net/http" //Fournit des fonctionnalités nécessaires pour créer un serveur web.

)

func main () {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //Associe une route HTTP (ça /) à une fonction qui gère les 
	// requêtes des utilisateurs (en gros quand un utioisateur accède à http://localhost:8081/, la fonction associée est éxécutée)
		fmt.Fprint(w, "Bienvenue sur mon forum !") //
	})
	
	port := "8081" //Port d'écoute définit 8081
	fmt.Printf("Lancement du serveur sur le port %s\n", port)
	err := http.ListenAndServe(":"+port, nil) //Lancement du serveur
	if err != nil { //Gestion des erreurs
		fmt.Printf("Pas de démarrage le serveur bug : %s\n", err)
	}
}