package game

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// Player represents a Scrabble player with their game state
type Player struct {
	ID       string       `json:"id"`        // Unique player identifier
	Name     string       `json:"name"`      // Player's display name
	Rack     []Tile       `json:"rack"`      // Player's tile rack (max 7 tiles)
	Score    int          `json:"score"`     // Current game score
	IsActive bool         `json:"is_active"` // Whether player is active in game
	mu       sync.RWMutex `json:"-"`         // Mutex for thread-safe operations
}

const (
	// MaxRackSize is the maximum number of tiles a player can have
	MaxRackSize = 7
)

// NewPlayer creates a new player with the given ID and name
func NewPlayer(id, name string) *Player {
	if id == "" {
		id = generatePlayerID()
	}

	return &Player{
		ID:       id,
		Name:     name,
		Rack:     make([]Tile, 0, MaxRackSize),
		Score:    0,
		IsActive: true,
	}
}

// generatePlayerID creates a simple player ID (in real implementation, use UUID)
func generatePlayerID() string {
	// Simple implementation for now - in production use proper UUID
	return fmt.Sprintf("player_%d", len("temp"))
}

// GetID returns the player's ID (thread-safe)
func (p *Player) GetID() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ID
}

// GetName returns the player's name (thread-safe)
func (p *Player) GetName() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Name
}

// SetName updates the player's name (thread-safe)
func (p *Player) SetName(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Name = name
}

// GetScore returns the player's current score (thread-safe)
func (p *Player) GetScore() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Score
}

// AddScore adds points to the player's score (thread-safe)
func (p *Player) AddScore(points int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Score += points
}

// SetScore sets the player's score to a specific value (thread-safe)
func (p *Player) SetScore(score int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Score = score
}

// IsPlayerActive returns whether the player is active (thread-safe)
func (p *Player) IsPlayerActive() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.IsActive
}

// SetActive sets the player's active status (thread-safe)
func (p *Player) SetActive(active bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.IsActive = active
}

// GetRackSize returns the current number of tiles in the player's rack
func (p *Player) GetRackSize() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.Rack)
}

// GetRack returns a copy of the player's current rack (thread-safe)
func (p *Player) GetRack() []Tile {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rack := make([]Tile, len(p.Rack))
	copy(rack, p.Rack)
	return rack
}

// AddTilesToRack adds tiles to the player's rack
func (p *Player) AddTilesToRack(tiles []Tile) error {
	if len(tiles) == 0 {
		return nil // No tiles to add is not an error
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if adding tiles would exceed rack capacity
	if len(p.Rack)+len(tiles) > MaxRackSize {
		return fmt.Errorf("cannot add %d tiles: would exceed rack capacity of %d (current: %d)",
			len(tiles), MaxRackSize, len(p.Rack))
	}

	// Add tiles to rack
	p.Rack = append(p.Rack, tiles...)

	return nil
}

// RemoveTilesFromRack removes tiles from the player's rack at specified indices
func (p *Player) RemoveTilesFromRack(indices []int) ([]Tile, error) {
	if len(indices) == 0 {
		return []Tile{}, nil // No indices specified is not an error
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Validate indices
	for _, index := range indices {
		if index < 0 || index >= len(p.Rack) {
			return nil, fmt.Errorf("invalid rack index: %d (rack size: %d)", index, len(p.Rack))
		}
	}

	// Check for duplicate indices
	indexSet := make(map[int]bool)
	for _, index := range indices {
		if indexSet[index] {
			return nil, fmt.Errorf("duplicate index: %d", index)
		}
		indexSet[index] = true
	}

	// Sort indices in descending order to remove from back to front
	sortedIndices := make([]int, len(indices))
	copy(sortedIndices, indices)
	sort.Sort(sort.Reverse(sort.IntSlice(sortedIndices)))

	// Extract tiles to return
	removedTiles := make([]Tile, len(indices))
	for i, index := range indices {
		removedTiles[i] = p.Rack[index]
	}

	// Remove tiles from rack (back to front to avoid index shifting)
	for _, index := range sortedIndices {
		p.Rack = append(p.Rack[:index], p.Rack[index+1:]...)
	}

	return removedTiles, nil
}

// RemoveTilesByValue removes specific tiles from the rack by their values
func (p *Player) RemoveTilesByValue(tiles []Tile) error {
	if len(tiles) == 0 {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Create a copy of the rack to work with
	rackCopy := make([]Tile, len(p.Rack))
	copy(rackCopy, p.Rack)

	// For each tile to remove, find and remove it from the copy
	for _, tileToRemove := range tiles {
		found := false
		for i, rackTile := range rackCopy {
			if tilesEqual(rackTile, tileToRemove) {
				// Remove this tile
				rackCopy = append(rackCopy[:i], rackCopy[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tile not found in rack: %v", tileToRemove)
		}
	}

	// Update the actual rack
	p.Rack = rackCopy
	return nil
}

// tilesEqual compares two tiles for equality
func tilesEqual(a, b Tile) bool {
	return a.Letter == b.Letter && a.Points == b.Points && a.IsBlank == b.IsBlank
}

// HasTile checks if the player has a specific tile in their rack
func (p *Player) HasTile(tile Tile) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, rackTile := range p.Rack {
		if tilesEqual(rackTile, tile) {
			return true
		}
	}
	return false
}

// HasTiles checks if the player has all specified tiles in their rack
func (p *Player) HasTiles(tiles []Tile) bool {
	if len(tiles) == 0 {
		return true
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	// Create a copy of the rack to work with
	rackCopy := make([]Tile, len(p.Rack))
	copy(rackCopy, p.Rack)

	// For each required tile, check if it exists in the rack copy
	for _, requiredTile := range tiles {
		found := false
		for i, rackTile := range rackCopy {
			if tilesEqual(rackTile, requiredTile) {
				// Remove this tile from the copy so it can't be used again
				rackCopy = append(rackCopy[:i], rackCopy[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// GetTileCount returns the count of a specific letter in the player's rack
func (p *Player) GetTileCount(letter rune) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	count := 0
	for _, tile := range p.Rack {
		if tile.Letter == letter {
			count++
		}
	}
	return count
}

// GetBlankCount returns the number of blank tiles in the player's rack
func (p *Player) GetBlankCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	count := 0
	for _, tile := range p.Rack {
		if tile.IsBlank {
			count++
		}
	}
	return count
}

// IsRackFull returns true if the player's rack is at maximum capacity
func (p *Player) IsRackFull() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.Rack) >= MaxRackSize
}

// IsRackEmpty returns true if the player has no tiles
func (p *Player) IsRackEmpty() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.Rack) == 0
}

// GetRackValue returns the total point value of tiles in the rack
func (p *Player) GetRackValue() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	total := 0
	for _, tile := range p.Rack {
		total += tile.Points
	}
	return total
}

// ClearRack removes all tiles from the player's rack and returns them
func (p *Player) ClearRack() []Tile {
	p.mu.Lock()
	defer p.mu.Unlock()

	tiles := make([]Tile, len(p.Rack))
	copy(tiles, p.Rack)
	p.Rack = p.Rack[:0] // Clear the rack but keep capacity

	return tiles
}

// ValidatePlayer performs comprehensive player state validation
func (p *Player) ValidatePlayer() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Check ID
	if p.ID == "" {
		return errors.New("player ID cannot be empty")
	}

	// Check name
	if p.Name == "" {
		return errors.New("player name cannot be empty")
	}

	// Check rack size
	if len(p.Rack) > MaxRackSize {
		return fmt.Errorf("rack size %d exceeds maximum %d", len(p.Rack), MaxRackSize)
	}

	// Check for negative score (might be valid in some variants, but flag for review)
	if p.Score < 0 {
		return fmt.Errorf("player has negative score: %d", p.Score)
	}

	// Validate each tile in rack
	for i, tile := range p.Rack {
		if tile.Points < 0 {
			return fmt.Errorf("tile at rack index %d has negative points: %d", i, tile.Points)
		}

		// Blank tiles should have 0 points and Letter should be 0
		if tile.IsBlank && (tile.Points != 0 || tile.Letter != 0) {
			return fmt.Errorf("blank tile at rack index %d has invalid properties: Letter=%c, Points=%d",
				i, tile.Letter, tile.Points)
		}

		// Non-blank tiles should have a valid letter
		if !tile.IsBlank && tile.Letter == 0 {
			return fmt.Errorf("non-blank tile at rack index %d has no letter", i)
		}
	}

	return nil
}

// String returns a string representation of the player for debugging
func (p *Player) String() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rackStr := "["
	for i, tile := range p.Rack {
		if i > 0 {
			rackStr += ", "
		}
		if tile.IsBlank {
			rackStr += "_"
		} else {
			rackStr += string(tile.Letter)
		}
	}
	rackStr += "]"

	activeStr := "active"
	if !p.IsActive {
		activeStr = "inactive"
	}

	return fmt.Sprintf("Player{ID: %s, Name: %s, Score: %d, Rack: %s (%d/7), Status: %s}",
		p.ID, p.Name, p.Score, rackStr, len(p.Rack), activeStr)
}
