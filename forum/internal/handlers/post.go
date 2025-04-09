package handlers

import (
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// HomeHandler affiche la page d'accueil avec la liste des posts récents
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier si l'URL est exactement "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Récupérer l'utilisateur actuel s'il est connecté
	currentUser := middleware.GetUserFromContext(r)
	var currentUserID int
	if currentUser != nil {
		currentUserID = currentUser.ID
	}

	// Récupérer les paramètres de filtrage et de pagination
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	perPage := 10 // Nombre de posts par page

	categoryID := 0
	if categoryStr := r.URL.Query().Get("category"); categoryStr != "" {
		if catID, err := strconv.Atoi(categoryStr); err == nil && catID > 0 {
			categoryID = catID
		}
	}

	userID := 0
	if userStr := r.URL.Query().Get("user"); userStr != "" {
		if uID, err := strconv.Atoi(userStr); err == nil && uID > 0 {
			userID = uID
		}
	}

	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" {
		sortBy = "date_desc" // Tri par défaut
	}

	// Récupérer les posts en fonction des filtres
	posts, total, err := models.GetPosts(page, perPage, categoryID, userID, sortBy, currentUserID)
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		log.Printf("Error fetching posts: %v", err)
		return
	}

	// Récupérer toutes les catégories pour le filtre
	categories, err := models.GetAllCategories()
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		log.Printf("Error fetching categories: %v", err)
		return
	}

	// Calculer le nombre total de pages
	totalPages := (total + perPage - 1) / perPage // Arrondi au supérieur

	// Préparer les données pour le template
	data := map[string]interface{}{
		"Posts":       posts,
		"Categories":  categories,
		"CurrentUser": currentUser,
		"CurrentPage": page,
		"TotalPages":  totalPages,
		"CategoryID":  categoryID,
		"UserID":      userID,
		"SortBy":      sortBy,
	}

	// Charger et exécuter le template
	tmpl, err := utils.ParseTemplate("templates/base.html", "templates/home.html")
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

// CreatePostHandler gère la création d'un nouveau post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que l'utilisateur est connecté
	currentUser := middleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Si c'est une requête GET, afficher le formulaire de création
	if r.Method == http.MethodGet {
		// Récupérer toutes les catégories pour le formulaire
		categories, err := models.GetAllCategories()
		if err != nil {
			http.Error(w, "Error fetching categories", http.StatusInternalServerError)
			log.Printf("Error fetching categories: %v", err)
			return
		}

		// Préparer les données pour le template
		data := map[string]interface{}{
			"Categories":  categories,
			"CurrentUser": currentUser,
		}

		// Charger et exécuter le template
		tmpl, err := utils.ParseTemplate("templates/base.html", "templates/create_post.html")
		if err != nil {
			http.Error(w, "Error loading templates", http.StatusInternalServerError)
			log.Printf("Error parsing template: %v", err)
			return
		}

		err = tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
		}
		return
	}

	// Si c'est une requête POST, traiter le formulaire
	if r.Method == http.MethodPost {
		// Analyser le formulaire
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Récupérer les données du formulaire
		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryIDs := r.Form["categories"] // Peut contenir plusieurs valeurs

		// Convertir les IDs de catégories en entiers
		var categoryIDsInt []int
		for _, idStr := range categoryIDs {
			id, err := strconv.Atoi(idStr)
			if err == nil && id > 0 {
				categoryIDsInt = append(categoryIDsInt, id)
			}
		}

		// Validation de base
		if title == "" || content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		// Créer le post
		postID, err := models.CreatePost(title, content, currentUser.ID, categoryIDsInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error creating post: %v", err)
			return
		}

		// Rediriger vers la page du post nouvellement créé
		http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
		return
	}

	// Si la méthode n'est ni GET ni POST
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// ViewPostHandler affiche un post spécifique avec ses commentaires
func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID du post de l'URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.NotFound(w, r)
		return
	}

	postID, err := strconv.Atoi(parts[2])
	if err != nil || postID <= 0 {
		http.NotFound(w, r)
		return
	}

	// Récupérer l'utilisateur actuel s'il est connecté
	currentUser := middleware.GetUserFromContext(r)
	var currentUserID int
	if currentUser != nil {
		currentUserID = currentUser.ID
	}

	// Récupérer le post
	post, err := models.GetPostByID(postID, currentUserID)
	if err != nil {
		if err.Error() == "post not found" {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Error fetching post", http.StatusInternalServerError)
			log.Printf("Error fetching post: %v", err)
		}
		return
	}

	// Récupérer les commentaires du post
	comments, err := models.GetCommentsByPostID(postID, currentUserID)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		log.Printf("Error fetching comments: %v", err)
		return
	}

	// Préparer les données pour le template
	data := map[string]interface{}{
		"Post":        post,
		"Comments":    comments,
		"CurrentUser": currentUser,
		"CanEdit":     currentUser != nil && (currentUser.ID == post.UserID || currentUser.Role == "admin"),
	}

	// Charger et exécuter le template
	tmpl, err := utils.ParseTemplate("templates/base.html", "templates/view_post.html")
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

// EditPostHandler gère la modification d'un post existant
func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID du post de l'URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[2] != "edit" {
		http.NotFound(w, r)
		return
	}

	postID, err := strconv.Atoi(parts[3])
	if err != nil || postID <= 0 {
		http.NotFound(w, r)
		return
	}

	// Vérifier que l'utilisateur est connecté
	currentUser := middleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer le post
	post, err := models.GetPostByID(postID, currentUser.ID)
	if err != nil {
		if err.Error() == "post not found" {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Error fetching post", http.StatusInternalServerError)
			log.Printf("Error fetching post: %v", err)
		}
		return
	}

	// Vérifier que l'utilisateur est autorisé à modifier ce post
	if currentUser.ID != post.UserID && currentUser.Role != "admin" {
		http.Error(w, "You don't have permission to edit this post", http.StatusForbidden)
		return
	}

	// Si c'est une requête GET, afficher le formulaire d'édition
	if r.Method == http.MethodGet {
		// Récupérer toutes les catégories pour le formulaire
		categories, err := models.GetAllCategories()
		if err != nil {
			http.Error(w, "Error fetching categories", http.StatusInternalServerError)
			log.Printf("Error fetching categories: %v", err)
			return
		}

		// Préparer un tableau des IDs des catégories du post pour le formulaire
		var selectedCategoryIDs []int
		for _, cat := range post.Categories {
			selectedCategoryIDs = append(selectedCategoryIDs, cat.ID)
		}

		// Préparer les données pour le template
		data := map[string]interface{}{
			"Post":               post,
			"Categories":         categories,
			"SelectedCategories": selectedCategoryIDs,
			"CurrentUser":        currentUser,
		}

		// Charger et exécuter le template
		tmpl, err := utils.ParseTemplate("templates/base.html", "templates/edit_post.html")
		if err != nil {
			http.Error(w, "Error loading templates", http.StatusInternalServerError)
			log.Printf("Error parsing template: %v", err)
			return
		}

		err = tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
		}
		return
	}

	// Si c'est une requête POST, traiter le formulaire
	if r.Method == http.MethodPost {
		// Analyser le formulaire
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Récupérer les données du formulaire
		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryIDs := r.Form["categories"] // Peut contenir plusieurs valeurs

		// Convertir les IDs de catégories en entiers
		var categoryIDsInt []int
		for _, idStr := range categoryIDs {
			id, err := strconv.Atoi(idStr)
			if err == nil && id > 0 {
				categoryIDsInt = append(categoryIDsInt, id)
			}
		}

		// Validation de base
		if title == "" || content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		// Mettre à jour le post
		err = models.UpdatePost(postID, currentUser.ID, title, content, categoryIDsInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error updating post: %v", err)
			return
		}

		// Rediriger vers la page du post mis à jour
		http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
		return
	}

	// Si la méthode n'est ni GET ni POST
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// DeletePostHandler gère la suppression d'un post
func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID du post de l'URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[2] != "delete" {
		http.NotFound(w, r)
		return
	}

	postID, err := strconv.Atoi(parts[3])
	if err != nil || postID <= 0 {
		http.NotFound(w, r)
		return
	}

	// Vérifier que l'utilisateur est connecté
	currentUser := middleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Supprimer le post
	isAdmin := currentUser.Role == "admin"
	err = models.DeletePost(postID, currentUser.ID, isAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error deleting post: %v", err)
		return
	}

	// Rediriger vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ReactToPostHandler gère les réactions (like/dislike) à un post
func ReactToPostHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier que la méthode est POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Vérifier que l'utilisateur est connecté
	currentUser := middleware.GetUserFromContext(r)
	if currentUser == nil {
		http.Error(w, "You must be logged in to react to posts", http.StatusUnauthorized)
		return
	}

	// Analyser le formulaire
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Récupérer les données du formulaire
	postIDStr := r.FormValue("post_id")
	reactionType := r.FormValue("reaction_type") // "like", "dislike" ou "" (pour supprimer)

	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Enregistrer la réaction
	err = models.ReactToPost(postID, currentUser.ID, reactionType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error reacting to post: %v", err)
		return
	}

	// Rediriger vers la page du post
	http.Redirect(w, r, "/post/"+postIDStr, http.StatusSeeOther)
}