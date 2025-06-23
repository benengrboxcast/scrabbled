# Scrabble Game Implementation Plan (Go)

## Project Overview
Develop a multiplayer Scrabble game using Go with client-server architecture. The server manages game state and validates moves, while clients provide the user interface for gameplay.

## Architecture

### High-Level Design
```
┌─────────────┐    WebSocket/TCP    ┌─────────────┐
│   Client 1  │◄──────────────────►│             │
├─────────────┤                    │             │
│   Client 2  │◄──────────────────►│   Server    │
├─────────────┤                    │             │
│   Client 3  │◄──────────────────►│             │
├─────────────┤                    │             │
│   Client 4  │◄──────────────────►│             │
└─────────────┘                    └─────────────┘
                                         │
                                         ▼
                                   ┌─────────────┐
                                   │ Dictionary  │
                                   │  (words.txt)│
                                   └─────────────┘
```

### Technology Stack
- **Language**: Go 1.21+
- **Communication**: WebSocket (gorilla/websocket) or gRPC
- **Dictionary**: Text file (words.txt)
- **Client UI**: Terminal-based (initially) or Web-based
- **Concurrency**: Goroutines and channels
- **Persistence**: SQLite (development) / PostgreSQL (production)
- **Session Management**: Redis for player sessions
- **Cleanup**: Background goroutines for game expiration

## Project Structure

```
scrabbled/
├── cmd/
│   ├── server/
│   │   └── main.go           # Server entry point
│   └── client/
│       └── main.go           # Client entry point
├── internal/
│   ├── game/
│   │   ├── board.go          # Board logic and validation
│   │   ├── tile.go           # Tile management
│   │   ├── player.go         # Player state
│   │   ├── game.go           # Game state and rules
│   │   ├── scoring.go        # Scoring calculations
│   │   └── persistence.go    # Game state serialization
│   ├── dictionary/
│   │   └── dictionary.go     # Word validation
│   ├── storage/
│   │   ├── database.go       # Database abstraction layer
│   │   ├── game_store.go     # Game persistence operations
│   │   └── session_store.go  # Player session management
│   ├── server/
│   │   ├── server.go         # HTTP/WebSocket server
│   │   ├── handlers.go       # Game event handlers
│   │   ├── room.go           # Game room management
│   │   ├── session.go        # Player session handling
│   │   └── cleanup.go        # Game expiration cleanup
│   └── client/
│       ├── client.go         # Client connection logic
│       ├── ui.go             # User interface
│       ├── renderer.go       # Board rendering
│       └── reconnection.go   # Reconnection logic
├── pkg/
│   └── protocol/
│       └── messages.go       # Client-server message types
├── data/
│   └── words.txt            # Dictionary file
├── docs/
│   ├── rules_of_scrabble.md
│   └── implementation_plan.md
├── go.mod
└── README.md
```

## Core Components

### 1. Game State Management (`internal/game/`)

#### Board (`board.go`)
```go
type Board struct {
    Grid     [15][15]Square
    Center   Position
}

type Square struct {
    Tile      *Tile
    Premium   PremiumType
    Occupied  bool
}

type Position struct {
    Row, Col int
}

type PremiumType int
const (
    Normal PremiumType = iota
    DoubleLetterScore
    TripleLetterScore
    DoubleWordScore
    TripleWordScore
)
```

#### Tile Management (`tile.go`)
```go
type Tile struct {
    Letter rune
    Points int
    IsBlank bool
}

type TileBag struct {
    tiles []Tile
    mu    sync.Mutex
}

// Methods: DrawTiles(), ReturnTiles(), RemainingCount()
```

#### Player (`player.go`)
```go
type Player struct {
    ID       string
    Name     string
    Rack     []Tile
    Score    int
    IsActive bool
}
```

#### Game Logic (`game.go`)
```go
type Game struct {
    ID          string
    Board       *Board
    Players     []*Player
    TileBag     *TileBag
    Dictionary  *Dictionary
    CurrentTurn int
    GameState   GameState
    History     []Move
    CreatedAt   time.Time
    LastActivity time.Time
    ExpiresAt   time.Time
    mu          sync.RWMutex
}

type GameState int
const (
    WaitingForPlayers GameState = iota
    InProgress
    Finished
)

type Move struct {
    PlayerID   string
    Tiles      []PlacedTile
    Score      int
    CreatedWords []string
    Timestamp  time.Time
}

type PlacedTile struct {
    Tile     Tile
    Position Position
}
```

### 2. Dictionary Service (`internal/dictionary/`)

```go
type Dictionary struct {
    words map[string]bool
    mu    sync.RWMutex
}

func NewDictionary(filename string) (*Dictionary, error)
func (d *Dictionary) IsValidWord(word string) bool
func (d *Dictionary) LoadFromFile(filename string) error
```

### 3. Server (`internal/server/`)

#### Server Structure
```go
type Server struct {
    games      map[string]*Game
    clients    map[string]*Client
    dictionary *Dictionary
    gameStore  storage.GameStore
    sessionStore storage.SessionStore
    cleanup    *CleanupService
    mu         sync.RWMutex
}

type Client struct {
    ID       string
    Name     string
    Conn     *websocket.Conn
    GameID   string
    PlayerID string
    SessionID string
    Send     chan []byte
}

type CleanupService struct {
    gameStore storage.GameStore
    ticker    *time.Ticker
    done      chan bool
}
```

### 4. Storage Layer (`internal/storage/`)

#### Game Persistence
```go
type GameStore interface {
    SaveGame(game *Game) error
    LoadGame(gameID string) (*Game, error)
    DeleteGame(gameID string) error
    GetExpiredGames() ([]string, error)
    UpdateLastActivity(gameID string) error
}

type SessionStore interface {
    SaveSession(session *PlayerSession) error
    GetSession(sessionID string) (*PlayerSession, error)
    DeleteSession(sessionID string) error
    GetPlayerSessions(playerID string) ([]*PlayerSession, error)
}

type PlayerSession struct {
    ID          string
    PlayerID    string
    GameID      string
    PlayerName  string
    CreatedAt   time.Time
    LastSeen    time.Time
    ExpiresAt   time.Time
}
```

#### Message Protocol (`pkg/protocol/`)
```go
type MessageType string
const (
    JoinGame     MessageType = "join_game"
    RejoinGame   MessageType = "rejoin_game"
    PlaceTiles   MessageType = "place_tiles"
    ExchangeTiles MessageType = "exchange_tiles"
    PassTurn     MessageType = "pass_turn"
    GameUpdate   MessageType = "game_update"
    PlayerReconnected MessageType = "player_reconnected"
    PlayerDisconnected MessageType = "player_disconnected"
    Error        MessageType = "error"
    Challenge    MessageType = "challenge"
)

type Message struct {
    Type      MessageType `json:"type"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
}

type JoinGameRequest struct {
    PlayerName string `json:"player_name"`
    GameID     string `json:"game_id,omitempty"`
}

type PlaceTilesRequest struct {
    Tiles []PlacedTile `json:"tiles"`
}

type GameUpdateResponse struct {
    Game    GameSnapshot `json:"game"`
    Players []Player     `json:"players"`
}
```

## Implementation Phases

### Phase 1: Core Game Engine (Weeks 1-2)
- [ ] Implement basic data structures (Board, Tile, Player)
- [ ] Create tile bag with proper distribution
- [ ] Implement board validation logic
- [ ] Basic scoring calculations
- [ ] Unit tests for core logic

**Deliverables:**
- Working game engine with move validation
- Comprehensive test suite
- Dictionary loading functionality

### Phase 2: Server Infrastructure (Weeks 3-4)
- [ ] WebSocket server setup
- [ ] Game room management
- [ ] Client connection handling
- [ ] Message routing and validation
- [ ] Concurrent game support

**Deliverables:**
- Multi-client server capable of managing multiple games
- WebSocket communication protocol
- Game state synchronization

### Phase 2.5: Persistence Layer (Weeks 4-5)
- [ ] Database schema design and implementation
- [ ] Game state serialization/deserialization
- [ ] Storage layer interfaces and implementations
- [ ] Session management system
- [ ] Game expiration and cleanup service

**Deliverables:**
- Persistent game state across server restarts
- Player session management
- Automatic cleanup of expired games

### Phase 3: Basic Client (Weeks 6-7)
- [ ] Terminal-based client interface
- [ ] Board rendering in text format
- [ ] User input handling
- [ ] Real-time game updates
- [ ] Move validation feedback
- [ ] Client-side session persistence
- [ ] Reconnection logic and UI

**Deliverables:**
- Functional terminal client
- Complete game playable end-to-end
- Client-server integration tested
- Players can reconnect to games after disconnection

### Phase 4: Advanced Features (Weeks 7-8)
- [ ] Word challenge system
- [ ] Tile exchange functionality
- [ ] Game replay system
- [ ] Spectator mode
- [ ] Enhanced error handling

### Phase 5: Polish & Optimization (Weeks 9-10)
- [ ] Performance optimization
- [ ] Comprehensive logging
- [ ] Configuration management
- [ ] Documentation completion
- [ ] Deployment preparation

## Key Implementation Details

### Dictionary Integration
```go
// Load dictionary on server startup
dict, err := dictionary.NewDictionary("data/words.txt")
if err != nil {
    log.Fatal("Failed to load dictionary:", err)
}

// Validate words during move processing
func (g *Game) validateMove(move Move) error {
    words := g.getFormedWords(move)
    for _, word := range words {
        if !g.Dictionary.IsValidWord(word) {
            return fmt.Errorf("invalid word: %s", word)
        }
    }
    return nil
}
```

### Game Persistence & Session Management
```go
// Save game state after each move
func (s *Server) handleMove(client *Client, move Move) error {
    game := s.games[client.GameID]
    if err := game.ApplyMove(move); err != nil {
        return err
    }
    
    // Update last activity and save to database
    game.LastActivity = time.Now()
    if err := s.gameStore.SaveGame(game); err != nil {
        log.Error("Failed to save game:", err)
    }
    
    s.broadcastGameUpdate(game)
    return nil
}

// Handle player reconnection
func (s *Server) handleRejoin(client *Client, sessionID string) error {
    session, err := s.sessionStore.GetSession(sessionID)
    if err != nil {
        return fmt.Errorf("invalid session: %w", err)
    }
    
    game, err := s.gameStore.LoadGame(session.GameID)
    if err != nil {
        return fmt.Errorf("game not found: %w", err)
    }
    
    // Reconnect player to game
    client.GameID = game.ID
    client.PlayerID = session.PlayerID
    s.games[game.ID] = game
    
    // Notify other players
    s.broadcastPlayerReconnected(game, session.PlayerID)
    return nil
}

// Background cleanup of expired games
func (c *CleanupService) Start() {
    c.ticker = time.NewTicker(1 * time.Hour)
    go func() {
        for {
            select {
            case <-c.ticker.C:
                c.cleanupExpiredGames()
            case <-c.done:
                c.ticker.Stop()
                return
            }
        }
    }()
}

func (c *CleanupService) cleanupExpiredGames() {
    expiredGames, err := c.gameStore.GetExpiredGames()
    if err != nil {
        log.Error("Failed to get expired games:", err)
        return
    }
    
    for _, gameID := range expiredGames {
        if err := c.gameStore.DeleteGame(gameID); err != nil {
            log.Error("Failed to delete expired game:", gameID, err)
        }
    }
}
```

### Concurrency Strategy
- Use `sync.RWMutex` for game state protection
- Dedicated goroutines for each client connection
- Channel-based communication for game events
- Graceful shutdown handling

### Error Handling
- Comprehensive error types for different scenarios
- Client notification for invalid moves
- Server resilience against client disconnections
- Transaction-like move validation (all-or-nothing)

### Performance Considerations
- Efficient word lookup using hash maps
- Minimal JSON marshaling/unmarshaling
- Connection pooling for high concurrent loads
- Caching for frequently accessed game data

## Testing Strategy

### Unit Tests
- Game logic validation
- Scoring calculations
- Dictionary operations
- Board state management

### Integration Tests
- Client-server communication
- Multi-player scenarios
- Game flow from start to finish
- Error condition handling

### Load Testing
- Multiple concurrent games
- High-frequency move submissions
- Client connection/disconnection stress testing

## Configuration Management

```go
type ServerConfig struct {
    Port           int    `json:"port"`
    DictionaryFile string `json:"dictionary_file"`
    MaxGames       int    `json:"max_games"`
    MaxPlayersPerGame int `json:"max_players_per_game"`
    LogLevel       string `json:"log_level"`
    Database       DatabaseConfig `json:"database"`
    GameExpiration time.Duration  `json:"game_expiration"`
    CleanupInterval time.Duration `json:"cleanup_interval"`
}

type DatabaseConfig struct {
    Type     string `json:"type"`     // "sqlite" or "postgres"
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Name     string `json:"name"`
    User     string `json:"user"`
    Password string `json:"password"`
    SSLMode  string `json:"ssl_mode"`
}
```

## Deployment Considerations

### Docker Setup
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/data ./data
CMD ["./server"]
```

### Production Readiness
- Health check endpoints
- Metrics collection (Prometheus)
- Structured logging
- Graceful shutdown
- Load balancing considerations

## Future Enhancements

### Web Client
- React/Vue.js web interface
- Drag-and-drop tile placement
- Real-time animations
- Mobile responsiveness

### Advanced Features
- Tournament mode
- AI opponents
- Custom dictionaries
- Game statistics and analytics
- Replay system with visualization

### Scalability
- Database persistence
- Horizontal scaling with Redis
- Microservices architecture
- CDN for static assets

---

This implementation plan provides a solid foundation for building a production-ready Scrabble game in Go with room for future enhancements and scaling. 