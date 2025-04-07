package main
import (
	"forum/internal/database"
	"forum/internal/handlers"
	"forum/internal/middleware"
	"log"
	"net/http"
)

func main() {
	// On initialise la base de données au démarrage du programme
	// C'est comme ouvrir le grand livre des utilisateurs avant d'accueillir les visiteurs
	err := database.InitDB("./forum.db")
	if err != nil {
		// Si on n'arrive pas à initialiser la base de données, on arrête tout
		// C'est comme fermer le bâtiment si le registre des membres est inaccessible
		log.Fatalf("Error initializing database: %v", err)
	}
	// On s'assure que la base de données sera fermée proprement à la fin du programme
	// C'est comme prévoir de ranger et fermer le registre quand la journée sera terminée
	defer database.CloseDB()
	
	// On crée un nouveau routeur pour gérer les différentes adresses du site
	// C'est comme installer un panneau d'orientation qui indique où aller selon ce qu'on cherche
	mux := http.NewServeMux()
	
	// On configure le serveur pour qu'il puisse fournir des fichiers statiques (images, CSS, JavaScript)
	// C'est comme désigner un espace où les visiteurs peuvent prendre librement des brochures et documents
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// On applique le middleware d'authentification à toutes les routes
	// C'est comme placer un portier discret à l'entrée qui vérifie les identités sans bloquer le passage
	authMiddleware := middleware.AuthMiddleware(mux)
	
	// On configure les routes pour l'authentification
	// C'est comme indiquer où se trouvent les guichets d'inscription, de connexion et de déconnexion
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
	
	// Exemple de routes protégées (qui seront implémentées plus tard)
	// C'est comme créer une zone réservée aux membres
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		// Cette page affiche simplement un message indiquant qu'elle nécessite une authentification
		w.Write([]byte("Profile page - requires authentication"))
	})
	// On ajoute un gardien spécial pour cette zone qui vérifie que l'utilisateur est bien connecté
	mux.Handle("/profile", middleware.RequireAuthMiddleware(protectedMux))
	
	// Page d'accueil du site
	// C'est comme configurer ce qu'on voit en entrant dans le bâtiment
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// On vérifie que l'URL est exactement "/" (la racine)
		// C'est comme s'assurer que le visiteur est bien à l'entrée principale
		if r.URL.Path != "/" {
			// Si l'URL ne correspond à aucune page connue, on affiche une erreur 404
			// C'est comme dire "Cette pièce n'existe pas dans notre bâtiment"
			http.NotFound(w, r)
			return
		}
		// Affiche un message simple sur la page d'accueil
		w.Write([]byte("Home page"))
	})
	
	// On démarre le serveur sur le port 8080
	// C'est comme ouvrir officiellement les portes du bâtiment et commencer à accueillir les visiteurs
	log.Println("Starting server on :8085...")
	// Si le serveur rencontre une erreur fatale, on enregistre l'erreur
	log.Fatal(http.ListenAndServe(":8085", authMiddleware))
}