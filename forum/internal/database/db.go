package database
import (
	"database/sql"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
)
// Cette variable DB est accessible dans tout le programme pour interagir avec la base de données
// C'est comme un canal de communication ouvert avec notre fichier de base de données
var DB *sql.DB

// La fonction InitDB prépare notre base de données pour qu'on puisse l'utiliser
// Elle prend le chemin du fichier où notre base de données sera stockée
func InitDB(filepath string) error {
	// On vérifie d'abord si le fichier de base de données existe déjà sur l'ordinateur
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		// Si le fichier n'existe pas, on le crée comme on créerait un nouveau cahier vide
		file, err := os.Create(filepath)
		if err != nil {
			// Si on n'arrive pas à créer le fichier, on signale l'erreur
			return err
		}
		file.Close()
	}
	
	// On ouvre la connexion avec notre base de données SQLite
	// C'est comme ouvrir le cahier pour pouvoir écrire dedans
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		// Si on n'arrive pas à ouvrir la base de données, on signale l'erreur
		return err
	}
	DB = db
	
	// On crée les tableaux dans notre base de données si ils n'existent pas encore
	// C'est comme préparer les pages de notre cahier avec des colonnes et des titres
	err = createTables()
	if err != nil {
		// Si on n'arrive pas à créer les tableaux, on signale l'erreur
		return err
	}
	
	// On note dans le journal que tout s'est bien passé
	log.Println("Database initialized successfully")
	return nil
}

// La fonction createTables crée les structures nécessaires dans notre base de données
// C'est comme préparer différentes sections dans notre cahier
func createTables() error {
	// On crée un tableau pour stocker les informations des utilisateurs
	// Chaque ligne du tableau représentera un utilisateur avec ses informations
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := DB.Exec(usersTable)
	if err != nil {
		// Si on n'arrive pas à créer le tableau des utilisateurs, on signale l'erreur
		return err
	}
	
	// On crée un tableau pour stocker les sessions de connexion
	// Une session représente une période pendant laquelle un utilisateur est connecté
	sessionsTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		user_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	_, err = DB.Exec(sessionsTable)
	if err != nil {
		// Si on n'arrive pas à créer le tableau des sessions, on signale l'erreur
		return err
	}
	
	// Tout s'est bien passé, on ne signale aucune erreur
	return nil
}

// La fonction CloseDB ferme proprement la connexion avec la base de données
// C'est comme fermer notre cahier quand on a fini de l'utiliser
func CloseDB() error {
	if DB != nil {
		// Si la base de données est ouverte, on la ferme
		return DB.Close()
	}
	// Si la base de données n'est pas ouverte, il n'y a rien à faire
	return nil
}