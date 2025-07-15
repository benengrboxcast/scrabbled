package game

import (
	"errors"
	"fmt"
	"strings"
)

// PremiumType represents the type of premium square
type PremiumType int

const (
	Normal            PremiumType = iota
	DoubleLetterScore             // Light blue - multiplies letter by 2
	TripleLetterScore             // Dark blue - multiplies letter by 3
	DoubleWordScore               // Pink/Light red - multiplies word by 2
	TripleWordScore               // Red - multiplies word by 3
)

// String returns a string representation of the premium type
func (pt PremiumType) String() string {
	switch pt {
	case Normal:
		return "NORMAL"
	case DoubleLetterScore:
		return "DLS"
	case TripleLetterScore:
		return "TLS"
	case DoubleWordScore:
		return "DWS"
	case TripleWordScore:
		return "TWS"
	default:
		return "UNKNOWN"
	}
}

// Position represents a coordinate on the board
type Position struct {
	Row int `json:"row"` // 0-based row (0-14)
	Col int `json:"col"` // 0-based column (0-14)
}

// String returns a string representation of the position (e.g., "H8")
func (p Position) String() string {
	if !p.IsValid() {
		return "INVALID"
	}
	// Convert to 1-based and use letter notation (A-O for columns, 1-15 for rows)
	col := string(rune('A' + p.Col))
	row := p.Row + 1
	return fmt.Sprintf("%s%d", col, row)
}

// IsValid checks if the position is within the board boundaries
func (p Position) IsValid() bool {
	return p.Row >= 0 && p.Row < 15 && p.Col >= 0 && p.Col < 15
}

// NewPositionFromString creates a Position from string notation (e.g., "H8")
func NewPositionFromString(s string) (Position, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if len(s) < 2 || len(s) > 3 {
		return Position{}, fmt.Errorf("invalid position format: %s", s)
	}

	// Parse column (A-O)
	col := int(s[0] - 'A')
	if col < 0 || col > 14 {
		return Position{}, fmt.Errorf("invalid column: %c", s[0])
	}

	// Parse row (1-15)
	var row int
	if len(s) == 2 {
		row = int(s[1] - '0')
	} else {
		// Handle two-digit rows (10-15)
		if s[1] != '1' {
			return Position{}, fmt.Errorf("invalid row: %s", s[1:])
		}
		row = 10 + int(s[2]-'0')
	}

	if row < 1 || row > 15 {
		return Position{}, fmt.Errorf("invalid row: %d", row)
	}

	return Position{Row: row - 1, Col: col}, nil
}

// Square represents a single square on the board
type Square struct {
	Tile     *Tile       `json:"tile"`     // Tile placed on this square (nil if empty)
	Premium  PremiumType `json:"premium"`  // Premium type for this square
	Occupied bool        `json:"occupied"` // True if a tile is placed here
}

// IsEmpty returns true if the square has no tile
func (s *Square) IsEmpty() bool {
	return s.Tile == nil || !s.Occupied
}

// Board represents the 15x15 Scrabble game board
type Board struct {
	Grid   [15][15]Square `json:"grid"`   // 15x15 grid of squares
	Center Position       `json:"center"` // Center position (H8)
}

// NewBoard creates a new Scrabble board with premium squares initialized
func NewBoard() *Board {
	board := &Board{
		Center: Position{Row: 7, Col: 7}, // H8 (0-based: row 7, col 7)
	}

	// Initialize all squares as normal
	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			board.Grid[row][col] = Square{
				Tile:     nil,
				Premium:  Normal,
				Occupied: false,
			}
		}
	}

	// Set premium squares according to official Scrabble board layout
	board.initializePremiumSquares()

	return board
}

// initializePremiumSquares sets up all premium squares according to official Scrabble rules
func (b *Board) initializePremiumSquares() {
	// Triple Word Score (TWS) - Red squares
	twsPositions := []string{"A1", "A8", "A15", "H1", "H15", "O1", "O8", "O15"}
	for _, posStr := range twsPositions {
		pos, _ := NewPositionFromString(posStr)
		b.Grid[pos.Row][pos.Col].Premium = TripleWordScore
	}

	// Double Word Score (DWS) - Pink/Light red squares (including center star)
	dwsPositions := []string{"B2", "C3", "D4", "E5", "H8", "K5", "L4", "M3", "N2",
		"B14", "C13", "D12", "E11", "K11", "L12", "M13", "N14"}
	for _, posStr := range dwsPositions {
		pos, _ := NewPositionFromString(posStr)
		b.Grid[pos.Row][pos.Col].Premium = DoubleWordScore
	}

	// Triple Letter Score (TLS) - Dark blue squares
	tlsPositions := []string{"B6", "B10", "F2", "F6", "F10", "F14", "J2", "J6", "J10", "J14", "N6", "N10"}
	for _, posStr := range tlsPositions {
		pos, _ := NewPositionFromString(posStr)
		b.Grid[pos.Row][pos.Col].Premium = TripleLetterScore
	}

	// Double Letter Score (DLS) - Light blue squares
	dlsPositions := []string{"A4", "A12", "C7", "C9", "D1", "D8", "D15", "G3", "G7", "G9", "G13",
		"H4", "H12", "I3", "I7", "I9", "I13", "L1", "L8", "L15", "M7", "M9", "O4", "O12"}
	for _, posStr := range dlsPositions {
		pos, _ := NewPositionFromString(posStr)
		b.Grid[pos.Row][pos.Col].Premium = DoubleLetterScore
	}
}

// IsValidPosition checks if a position is within the board boundaries
func (b *Board) IsValidPosition(pos Position) bool {
	return pos.IsValid()
}

// PlaceTile places a tile at the specified position
func (b *Board) PlaceTile(tile Tile, pos Position) error {
	if !b.IsValidPosition(pos) {
		return fmt.Errorf("invalid position: %s", pos.String())
	}

	square := &b.Grid[pos.Row][pos.Col]
	if square.Occupied {
		return fmt.Errorf("position %s is already occupied", pos.String())
	}

	// Place the tile
	square.Tile = &tile
	square.Occupied = true

	return nil
}

// GetTile returns the tile at the specified position (nil if empty)
func (b *Board) GetTile(pos Position) *Tile {
	if !b.IsValidPosition(pos) {
		return nil
	}

	square := &b.Grid[pos.Row][pos.Col]
	if !square.Occupied {
		return nil
	}

	return square.Tile
}

// GetSquare returns the square at the specified position
func (b *Board) GetSquare(pos Position) *Square {
	if !b.IsValidPosition(pos) {
		return nil
	}

	return &b.Grid[pos.Row][pos.Col]
}

// RemoveTile removes a tile from the specified position
func (b *Board) RemoveTile(pos Position) (*Tile, error) {
	if !b.IsValidPosition(pos) {
		return nil, fmt.Errorf("invalid position: %s", pos.String())
	}

	square := &b.Grid[pos.Row][pos.Col]
	if !square.Occupied {
		return nil, fmt.Errorf("no tile at position %s", pos.String())
	}

	tile := square.Tile
	square.Tile = nil
	square.Occupied = false

	return tile, nil
}

// IsEmpty returns true if the square at the given position is empty
func (b *Board) IsEmpty(pos Position) bool {
	if !b.IsValidPosition(pos) {
		return false
	}

	return b.Grid[pos.Row][pos.Col].IsEmpty()
}

// GetAdjacentPositions returns all valid adjacent positions (up, down, left, right)
func (b *Board) GetAdjacentPositions(pos Position) []Position {
	if !b.IsValidPosition(pos) {
		return nil
	}

	adjacent := []Position{}

	// Check all four directions: up, down, left, right
	directions := []Position{
		{Row: -1, Col: 0}, // Up
		{Row: 1, Col: 0},  // Down
		{Row: 0, Col: -1}, // Left
		{Row: 0, Col: 1},  // Right
	}

	for _, dir := range directions {
		newPos := Position{
			Row: pos.Row + dir.Row,
			Col: pos.Col + dir.Col,
		}

		if b.IsValidPosition(newPos) {
			adjacent = append(adjacent, newPos)
		}
	}

	return adjacent
}

// GetPremiumType returns the premium type for the given position
func (b *Board) GetPremiumType(pos Position) PremiumType {
	if !b.IsValidPosition(pos) {
		return Normal
	}

	return b.Grid[pos.Row][pos.Col].Premium
}

// ValidateBoard performs comprehensive board state validation
func (b *Board) ValidateBoard() error {
	// Check that center position is correct
	if b.Center.Row != 7 || b.Center.Col != 7 {
		return errors.New("center position must be H8 (row 7, col 7)")
	}

	// Verify premium square counts
	premiumCounts := make(map[PremiumType]int)
	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			premiumCounts[b.Grid[row][col].Premium]++
		}
	}

	// Expected premium square counts according to official Scrabble rules
	expectedCounts := map[PremiumType]int{
		Normal:            164, // 225 total - 61 premium squares
		DoubleLetterScore: 24,
		TripleLetterScore: 12,
		DoubleWordScore:   17, // Including center star
		TripleWordScore:   8,
	}

	for premiumType, expected := range expectedCounts {
		if actual := premiumCounts[premiumType]; actual != expected {
			return fmt.Errorf("premium square count mismatch for %s: expected %d, got %d",
				premiumType.String(), expected, actual)
		}
	}

	// Verify center square is Double Word Score
	if b.Grid[b.Center.Row][b.Center.Col].Premium != DoubleWordScore {
		return errors.New("center square must be a Double Word Score")
	}

	return nil
}

// IsFirstMove checks if this is the first move of the game (board is empty)
func (b *Board) IsFirstMove() bool {
	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			if b.Grid[row][col].Occupied {
				return false
			}
		}
	}
	return true
}

// HasTileAt returns true if there is a tile at the specified position
func (b *Board) HasTileAt(pos Position) bool {
	if !b.IsValidPosition(pos) {
		return false
	}

	return b.Grid[pos.Row][pos.Col].Occupied
}

// GetOccupiedPositions returns all positions that have tiles
func (b *Board) GetOccupiedPositions() []Position {
	positions := []Position{}

	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			if b.Grid[row][col].Occupied {
				positions = append(positions, Position{Row: row, Col: col})
			}
		}
	}

	return positions
}

// CountPremiumSquares returns the count of each premium square type
func (b *Board) CountPremiumSquares() map[PremiumType]int {
	counts := make(map[PremiumType]int)

	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			counts[b.Grid[row][col].Premium]++
		}
	}

	return counts
}

// GetRow returns all squares in the specified row
func (b *Board) GetRow(row int) []Square {
	if row < 0 || row >= 15 {
		return nil
	}

	squares := make([]Square, 15)
	copy(squares, b.Grid[row][:])
	return squares
}

// GetColumn returns all squares in the specified column
func (b *Board) GetColumn(col int) []Square {
	if col < 0 || col >= 15 {
		return nil
	}

	squares := make([]Square, 15)
	for row := 0; row < 15; row++ {
		squares[row] = b.Grid[row][col]
	}
	return squares
}

// String returns a string representation of the board for debugging
func (b *Board) String() string {
	var sb strings.Builder

	// Header with column letters
	sb.WriteString("   ")
	for col := 0; col < 15; col++ {
		sb.WriteString(fmt.Sprintf(" %c ", 'A'+col))
	}
	sb.WriteString("\n")

	// Board rows
	for row := 0; row < 15; row++ {
		sb.WriteString(fmt.Sprintf("%2d ", row+1))
		for col := 0; col < 15; col++ {
			square := &b.Grid[row][col]
			if square.Occupied && square.Tile != nil {
				sb.WriteString(fmt.Sprintf(" %s ", square.Tile.String()))
			} else {
				// Show premium square type
				switch square.Premium {
				case TripleWordScore:
					sb.WriteString(" * ")
				case DoubleWordScore:
					sb.WriteString(" + ")
				case TripleLetterScore:
					sb.WriteString(" ^ ")
				case DoubleLetterScore:
					sb.WriteString(" - ")
				default:
					sb.WriteString(" . ")
				}
			}
		}
		sb.WriteString(fmt.Sprintf(" %d\n", row+1))
	}

	// Footer with column letters
	sb.WriteString("   ")
	for col := 0; col < 15; col++ {
		sb.WriteString(fmt.Sprintf(" %c ", 'A'+col))
	}
	sb.WriteString("\n")

	return sb.String()
}
