# Hexagonal Architecture Migration — Status & Plan

## ✅ What's Already Done

| Layer | Files | Status |
|---|---|---|
| **Ports (interfaces)** | `internal/ports/player.go`, `track.go`, `searcher.go`, `platforme.go` | ✅ Defined |
| **Domain services** | `internal/domain/player.go` (PlayerService), `internal/domain/tracks.go` (TrackService) | ✅ Scaffolded, depends only on ports |
| **Infra directory structure** | `internal/infra/mpv/`, `internal/infra/ytdlp/` | ✅ Created |
| **Partial adapter** | `internal/infra/ytdlp/ytdlp.go` | ⚠️ Has `StreamURL()` method — port expects `Resolve()` |
| **Empty directories** | `internal/ui/`, `internal/domain/entities/`, `cmd/` | ✅ Reserved |

---

## ❌ Remaining Tasks (Priority Order)

### 1. Fix `internal/infra/ytdlp/` adapter
- Rename `StreamURL()` → `Resolve()` to match `ports.StreamURLResolver` interface

### 2. Implement `internal/infra/mpv/` adapter
- Implement `ports.Player` interface by extracting logic from `player/player.go` (mpv spawn, stdin parsing, progress regex) and `player/ipcCommand.go` (named pipe IPC)
- Currently 3 stub files (`client.go`, `event.go`, `player.go`) — all just `package mpv`

### 3. Implement search platform adapters
- Move `player.SearchYoutube()` into a proper adapter implementing `ports.PlatformeInterface`
- Implement adapters for Spotify and Deezer (listed in sidebar, not implemented anywhere)

### 4. Implement `ports.Tracks`
- Like, playlist CRUD — currently no implementation exists anywhere

### 5. Refactor TUI to use domain services
- `tui/tui.go`, `tui/trackItem.go`, `tui/footerPanel.go` depend directly on `player.Player`
  and its Bubble Tea message types (`SearchCompleteMsg`, `PlayerProgressMsg`, etc.)
- Replace with `domain.PlayerService` + `domain.TrackService`
- Extract Bubble Tea messages from `player/` into the TUI layer

### 6. Clean up stubs & dead code
- `ipc/client.go`, `pkg/ipc/pipeUnix.go`, `pkg/ipc/pipeWindows.go` — all stubs
- `player/player.go` — entire legacy module to delete once migration is complete

### 7. Fix naming/typos
- `PlatformeInterface` → `PlatformInterface`
- `Traks` → `Tracks` in `domain.TrackService`
- `plateformModel` → `platformModel` in TUI
- `PlatfomDeleget` → `PlatformDelegate`
- `Platfome` → `Platform` in ports

### 8. Wire up `main.go`
- Inject concrete adapters into domain services and pass them to TUI

### 9. (Optional) Domain entities
- `internal/domain/entities/` is empty — could hold `Track` as a rich domain model

### 10. Tests
- Move `player/player_test.go` to the new adapter layer
- Add unit tests for domain services with mocked ports
