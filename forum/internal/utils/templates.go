package utils

import (
	"html/template"
	"path/filepath"
)

// ParseTemplate charge les templates et y ajoute les fonctions personnalisées
func ParseTemplate(baseTemplate string, templates ...string) (*template.Template, error) {
	// Créer un nouveau template avec les fonctions personnalisées
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"subtract": func(a, b int) int {
			return a - b
		},
		"sequence": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
	}
	
	// Ajouter le template de base
	allTemplates := append([]string{baseTemplate}, templates...)
	
	// Créer et parser le template
	return template.New(filepath.Base(baseTemplate)).Funcs(funcMap).ParseFiles(allTemplates...)
}