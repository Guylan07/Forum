// On crée un tableau pour stocker les catégories de posts
// C'est comme créer une section spéciale dans notre cahier pour lister tous les thèmes possibles des discussions
categoriesTable := `
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT
);`
_, err = DB.Exec(categoriesTable)
if err != nil {
    // Si on n'arrive pas à créer le tableau des catégories, on signale l'erreur
    return err
}

// On crée un tableau pour stocker les posts
// C'est comme préparer des pages où chaque utilisateur pourra écrire ses messages principaux
postsTable := `
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);`
_, err = DB.Exec(postsTable)
if err != nil {
    // Si on n'arrive pas à créer le tableau des posts, on signale l'erreur
    return err
}

// On crée un tableau pour associer des catégories aux posts
// C'est comme créer une liste qui indique dans quel rayon de bibliothèque ranger chaque livre
postCategoriesTable := `
CREATE TABLE IF NOT EXISTS post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);`
_, err = DB.Exec(postCategoriesTable)
if err != nil {
    // Si on n'arrive pas à créer le tableau des associations, on signale l'erreur
    return err
}

// On crée un tableau pour stocker les commentaires
// C'est comme préparer des petites fiches où les gens peuvent répondre à chaque message principal
commentsTable := `
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);`
_, err = DB.Exec(commentsTable)
if err != nil {
    // Si on n'arrive pas à créer le tableau des commentaires, on signale l'erreur
    return err
}

// On crée un tableau pour stocker les likes/dislikes des posts
// C'est comme préparer un système où chaque personne peut mettre un pouce en haut ou en bas sur un message
postReactionsTable := `
CREATE TABLE IF NOT EXISTS post_reactions (
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    reaction_type TEXT NOT NULL,  -- 'like' ou 'dislike'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);`
_, err = DB.Exec(postReactionsTable)
if err != nil {
    // Si on n'arrive pas à créer le tableau des réactions aux posts, on signale l'erreur
    return err
}

// On crée un tableau pour stocker les likes/dislikes des commentaires
// C'est comme le système précédent, mais pour les petites fiches de réponse
commentReactionsTable := `
CREATE TABLE IF NOT EXISTS comment_reactions (
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    reaction_type TEXT NOT NULL,  -- 'like' ou 'dislike'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, comment_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);`
_, err = DB.Exec(commentReactionsTable)
if err != nil {
    // Si on n'arrive pas à créer le tableau des réactions aux commentaires, on signale l'erreur
    return err
}

// On ajoute quelques catégories par défaut si elles n'existent pas déjà
// C'est comme préparer les rayons principaux de notre bibliothèque avant d'ouvrir
defaultCategories := []string{"Général", "Technologie", "Sport", "Musique", "Cinéma", "Jeux vidéo", "Science", "Art", "Politique", "Autre"}
for _, category := range defaultCategories {
    // Pour chaque catégorie de notre liste, on l'ajoute si elle n'existe pas déjà
    _, err = DB.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", category)
    if err != nil {
        // Si on n'arrive pas à ajouter une catégorie, on signale l'erreur
        return err
    }
}