# Tâches Restantes — Migration Architecture Hexagonale

**État actuel** : ~30% structurellement complété, 0% fonctionnellement opérationnel

---

## Tâche 1 — Corriger l'adaptateur ytdlp

**Fichier** : `internal/infra/ytdlp/ytdlp.go`

**Problème** : La méthode `StreamURL()` ne respecte pas le contrat de `ports.StreamURLResolver`.

**Actions** :
- Renommer `StreamURL()` → `Resolve()` pour matcher l'interface `ports.StreamURLResolver`
- Vérifier la signature complète : `Resolve(ctx context.Context, url string) (string, error)`

---

## Tâche 2 — Implémenter l'adaptateur mpv

**Fichiers à implémenter** :
- `internal/infra/mpv/client.go` — Client IPC (named pipe Windows)
- `internal/infra/mpv/event.go` — Mapping des événements mpv vers `ports.PlayerEvent`
- `internal/infra/mpv/player.go` — Implémentation de `ports.Player`

**Logique à extraire depuis** :
- `player/player.go` : spawning mpv, parsing stdout, regex progression
- `player/ipcCommand.go` : IPC named pipe via go-winio, `sendCommand()`, `GetPercentPos()`

**Port à implémenter** : `ports.Player`
- `Connect(ctx) error`
- `Play(url string) error`
- `Pause() error`
- `Resume() error`
- `Stop() error`
- `SetVolume(level float64) error`
- `GetVolume() (float64, error)`
- `Seek(mode ports.SeekMode) error`
- `SetMute(mute bool) error`
- `GetMute() (bool, error)`
- `Close() error`

---

## Tâche 3 — Implémenter les adaptateurs de plateformes

**Fichier source** : `player/player.go` → fonction `SearchYoutube()`

**Actions** :
- Créer `internal/infra/ytdlp/platform.go` (ou autre nom) implémentant `ports.PlatformeInterface`
- Déplacer la logique de recherche YouTube vers cet adaptateur
- **Bonus** : Créer des stubs pour Spotify et Deezer (affichés dans la sidebar)

**Port à implémenter** : `ports.PlatformeInterface`
- `Search(ctx context.Context, query string) ([]ports.Track, error)`
- `StreamUrlFormat() string`

---

## Tâche 4 — Implémenter `ports.Tracks`

**Statut** : Aucune implém. n'existe nulle part

**Actions** :
- Créer `internal/infra/memory/tracks.go` (ou `internal/infra/sqlite/tracks.go`)
- Implémenter le like et la gestion des playlists

**Port à implémenter** : `ports.Tracks`
- `Like(track ports.Track) error`
- `Playlist(track ports.Track) error`
- `NewPlaylist(name string) error`
- `DeletePlaylist(name string) error`

---

## Tâche 5 — Refactorer le TUI pour utiliser les domain services

**Fichiers concernés** :
- `tui/tui.go` — dépend directement de `player.Player`
- `tui/trackItem.go` — utilise `player.SearchYTCmd`, `player.TrackItem`
- `tui/footerPanel.go` — utilise `player.PlayerProgressMsg`
- `tui/trackItemDelegate.go` — utilise `player.Player`

**Actions** :
- Remplacer les dépendances `player.Player` par `domain.PlayerService` + `domain.TrackService`
- Extraire les messages BubbleTea (`SearchCompleteMsg`, `PlayerProgressMsg`, etc.) de `player/` vers la couche TUI
- Refactorer `internal/ui/` pour recevoir les services injectés

---

## Tâche 6 — Nettoyer les stubs et code mort

**Fichiers à supprimer** :
- `ipc/client.go` — stub vide
- `pkg/ipc/pipeUnix.go` — stub vide
- `pkg/ipc/pipeWindows.go` — stub vide

**Dossier à supprimer une fois migration terminée** :
- `player/` — tout le module legacy

---

## Tâche 7 — Corriger les typos de noms

| Actuel | Corrigé | Fichier |
|---|---|---|
| `PlatformeInterface` | `PlatformInterface` | `internal/ports/platforme.go` |
| `Platfome` | `Platform` | `internal/ports/platforme.go` |
| `Traks` | `Tracks` | `internal/domain/tracks.go` |
| `SteamURL` | `StreamURL` | `internal/ports/track.go` |
| `Pauased` | `Paused` | `internal/ports/player.go` |
| `plateformModel` | `platformModel` | `tui/tracklist.go` |
| `PlatfomDeleget` | `PlatformDelegate` | `tui/palteformDelegete.go` |

---

## Tâche 8 — câbler `main.go`

**Fichier** : `main.go`

**Actions** :
- Instancier les adaptateurs concrets (mpv, ytdlp, etc.)
- Injecter dans les domain services (`NewPlayerService(...)`, `NewTrack(...)`)
- Passer les services au TUI
- Remplacer `tui.NewModel()` par un constructeur injectant les dépendances

---

## Tâche 9 — (Optionnel) Entités de domaine

**Dossier** : `internal/domain/entities/` (actuellement vide)

**Actions** :
- Déplacer `Track` depuis `internal/ports/track.go` vers `internal/domain/entities/track.go`
- Enrichir avec des méthodes de validation ou de comportement métier

---

## Tâche 10 — Tests

**Fichier existant** : `player/player_test.go` (1 test d'intégration)

**Actions** :
- Déplacer `player/player_test.go` vers `internal/infra/ytdlp/` (test d'adaptateur)
- Créer des mocks pour chaque port interface
- Ajouter des tests unitaires pour `domain.PlayerService` et `domain.TrackService`
- Ajouter des tests pour les adaptateurs mpv et ytdlp

---

## Ordre de Priorité Recommandé

1. **Tâche 1** — Fix ytdlp (rapide, débloque les tests)
2. **Tâche 2** — Impl mpv (cœur fonctionnel)
3. **Tâche 8** — câbler main.go (rend le tout utilisable)
4. **Tâche 5** — Refactor TUI (découple du legacy)
5. **Tâche 7** — Corriger typos (propreté)
6. **Tâche 6** — Nettoyer stubs (réduction dette)
7. **Tâche 4** — Impl Tracks (fonctionnalité)
8. **Tâche 3** — Adaptateurs plateformes (extensions)
9. **Tâche 10** — Tests (fiabilité)
10. **Tâche 9** — Entités domaine (optionnel)
