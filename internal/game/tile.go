package game

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// Tile represents a single letter tile in Scrabble
type Tile struct {
	Letter  rune `json:"letter"`   // The letter on the tile ('A', 'B', etc.) or 0 for blank
	Points  int  `json:"points"`   // Point value of the tile
	IsBlank bool `json:"is_blank"` // True if this is a blank tile
}

// String returns a string representation of the tile
func (t Tile) String() string {
	if t.IsBlank {
		return "BLANK"
	}
	return string(t.Letter)
}

// TileBag manages the collection of tiles that can be drawn from
type TileBag struct {
	tiles []Tile
	mu    sync.Mutex
}

// Standard Scrabble tile distribution (100 tiles total)
var standardTileDistribution = map[rune]struct {
	quantity int
	points   int
}{
	'A': {9, 1}, 'B': {2, 3}, 'C': {2, 3}, 'D': {4, 2}, 'E': {12, 1},
	'F': {2, 4}, 'G': {3, 2}, 'H': {2, 4}, 'I': {9, 1}, 'J': {1, 8},
	'K': {1, 5}, 'L': {4, 1}, 'M': {2, 3}, 'N': {6, 1}, 'O': {8, 1},
	'P': {2, 3}, 'Q': {1, 10}, 'R': {6, 1}, 'S': {4, 1}, 'T': {6, 1},
	'U': {4, 1}, 'V': {2, 4}, 'W': {2, 4}, 'X': {1, 8}, 'Y': {2, 4},
	'Z': {1, 10},
}

// Number of blank tiles
const blankTileCount = 2

// NewTileBag creates a new tile bag with the standard Scrabble distribution
func NewTileBag() *TileBag {
	bag := &TileBag{
		tiles: make([]Tile, 0, 100), // Pre-allocate for 100 tiles
	}

	// Add letter tiles according to standard distribution
	for letter, dist := range standardTileDistribution {
		for i := 0; i < dist.quantity; i++ {
			bag.tiles = append(bag.tiles, Tile{
				Letter:  letter,
				Points:  dist.points,
				IsBlank: false,
			})
		}
	}

	// Add blank tiles
	for i := 0; i < blankTileCount; i++ {
		bag.tiles = append(bag.tiles, Tile{
			Letter:  0, // 0 represents blank
			Points:  0,
			IsBlank: true,
		})
	}

	// Shuffle the tiles
	bag.shuffle()

	return bag
}

// shuffle randomizes the order of tiles in the bag
func (tb *TileBag) shuffle() {
	rand.Seed(time.Now().UnixNano())
	for i := len(tb.tiles) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		tb.tiles[i], tb.tiles[j] = tb.tiles[j], tb.tiles[i]
	}
}

// DrawTiles removes and returns up to 'count' tiles from the bag
// Returns fewer tiles if the bag doesn't have enough tiles
func (tb *TileBag) DrawTiles(count int) []Tile {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if count <= 0 {
		return []Tile{}
	}

	// Can't draw more tiles than available
	available := len(tb.tiles)
	if count > available {
		count = available
	}

	// Draw tiles from the end of the slice
	drawn := make([]Tile, count)
	copy(drawn, tb.tiles[available-count:])

	// Remove drawn tiles from the bag
	tb.tiles = tb.tiles[:available-count]

	return drawn
}

// ReturnTiles adds tiles back to the bag and shuffles
// This is used when players exchange tiles
func (tb *TileBag) ReturnTiles(tiles []Tile) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// Add tiles back to the bag
	tb.tiles = append(tb.tiles, tiles...)

	// Shuffle to ensure randomness
	tb.shuffle()
}

// RemainingCount returns the number of tiles left in the bag
func (tb *TileBag) RemainingCount() int {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return len(tb.tiles)
}

// IsEmpty returns true if the tile bag is empty
func (tb *TileBag) IsEmpty() bool {
	return tb.RemainingCount() == 0
}

// GetTileValue returns the point value for a given letter
// Returns 0 for blank tiles or invalid letters
func GetTileValue(letter rune) int {
	if dist, exists := standardTileDistribution[letter]; exists {
		return dist.points
	}
	return 0 // Blank or invalid letter
}

// GetTileQuantity returns the standard quantity for a given letter
// Returns 0 for invalid letters
func GetTileQuantity(letter rune) int {
	if dist, exists := standardTileDistribution[letter]; exists {
		return dist.quantity
	}
	if letter == 0 { // Blank tile
		return blankTileCount
	}
	return 0
}

// ValidateTileDistribution verifies that the standard distribution totals 100 tiles
func ValidateTileDistribution() error {
	total := 0

	// Count letter tiles
	for _, dist := range standardTileDistribution {
		total += dist.quantity
	}

	// Add blank tiles
	total += blankTileCount

	if total != 100 {
		return errors.New("tile distribution does not total 100 tiles")
	}

	return nil
}

// GetAllTileInfo returns information about all tiles for testing/debugging
func GetAllTileInfo() map[rune]struct {
	Quantity int
	Points   int
} {
	result := make(map[rune]struct {
		Quantity int
		Points   int
	})

	// Copy letter tiles
	for letter, dist := range standardTileDistribution {
		result[letter] = struct {
			Quantity int
			Points   int
		}{dist.quantity, dist.points}
	}

	// Add blank tiles (using rune 0 to represent blank)
	result[0] = struct {
		Quantity int
		Points   int
	}{blankTileCount, 0}

	return result
}
