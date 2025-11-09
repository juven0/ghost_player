package player

import (
	"testing"
)

func TestSearchYoutube(t *testing.T) {
	// Skip si on veut éviter les tests qui nécessitent internet
	if testing.Short() {
		t.Skip("Skipping test requiring internet connection")
	}

	// Test de recherche basique
	query := "lofi hip hop"
	maxResults := 5

	results, err := SearchYoutube(query, maxResults)
	if err != nil {
		t.Fatalf("SearchYoutube() error = %v", err)
	}

	if len(results) == 0 {
		t.Fatal("SearchYoutube() returned no results")
	}

	// Vérifier qu'on n'a pas plus de résultats que demandé
	if len(results) > maxResults {
		t.Errorf("SearchYoutube() returned %d results, want max %d", len(results), maxResults)
	}

	// Vérifier que chaque résultat a les champs requis
	for i, video := range results {
		if video.ID == "" {
			t.Errorf("Result %d: ID is empty", i)
		}
		if video.Title == "" {
			t.Errorf("Result %d: Title is empty", i)
		}
		if video.Duration == 0 {
			t.Errorf("Result %d: Duration is empty", i)
		}

		// Afficher le résultat (optionnel)
		t.Logf("Result %d: %s (ID: %s)", i+1, video.Title, video.ID)
	}
}
