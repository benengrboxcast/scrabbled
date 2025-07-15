package game

import (
	"testing"
)

// TestBoardInitialization tests board creation and initial state
func TestBoardInitialization(t *testing.T) {
	board := NewBoard()

	// Test center position
	if board.Center.Row != 7 || board.Center.Col != 7 {
		t.Errorf("Center position should be (7,7), got (%d,%d)", board.Center.Row, board.Center.Col)
	}

	// Test center position string representation
	if board.Center.String() != "H8" {
		t.Errorf("Center position string should be 'H8', got '%s'", board.Center.String())
	}

	// Test that board is initially empty
	if !board.IsFirstMove() {
		t.Errorf("New board should be empty for first move")
	}

	// Test all squares are initially unoccupied
	occupiedCount := 0
	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			if board.Grid[row][col].Occupied {
				occupiedCount++
			}
		}
	}

	if occupiedCount != 0 {
		t.Errorf("New board should have no occupied squares, found %d", occupiedCount)
	}

	// Test board validation
	if err := board.ValidateBoard(); err != nil {
		t.Errorf("New board should be valid: %v", err)
	}
}

// TestSquareStateManagement tests individual square operations
func TestSquareStateManagement(t *testing.T) {
	board := NewBoard()
	pos := Position{Row: 7, Col: 7} // Center
	tile := Tile{Letter: 'A', Points: 1, IsBlank: false}

	// Test initial state
	if !board.IsEmpty(pos) {
		t.Errorf("Square should initially be empty")
	}

	if board.HasTileAt(pos) {
		t.Errorf("Square should not have tile initially")
	}

	if board.GetTile(pos) != nil {
		t.Errorf("GetTile should return nil for empty square")
	}

	// Test tile placement
	err := board.PlaceTile(tile, pos)
	if err != nil {
		t.Errorf("Should be able to place tile: %v", err)
	}

	// Test state after placement
	if board.IsEmpty(pos) {
		t.Errorf("Square should not be empty after placing tile")
	}

	if !board.HasTileAt(pos) {
		t.Errorf("Square should have tile after placement")
	}

	retrievedTile := board.GetTile(pos)
	if retrievedTile == nil {
		t.Errorf("GetTile should return the placed tile")
	}

	if retrievedTile.Letter != tile.Letter || retrievedTile.Points != tile.Points {
		t.Errorf("Retrieved tile should match placed tile")
	}

	// Test double placement (should fail)
	err = board.PlaceTile(tile, pos)
	if err == nil {
		t.Errorf("Should not be able to place tile on occupied square")
	}

	// Test tile removal
	removedTile, err := board.RemoveTile(pos)
	if err != nil {
		t.Errorf("Should be able to remove tile: %v", err)
	}

	if removedTile.Letter != tile.Letter {
		t.Errorf("Removed tile should match original tile")
	}

	// Test state after removal
	if !board.IsEmpty(pos) {
		t.Errorf("Square should be empty after removing tile")
	}

	// Test removing from empty square (should fail)
	_, err = board.RemoveTile(pos)
	if err == nil {
		t.Errorf("Should not be able to remove tile from empty square")
	}
}

// TestPositionValidation tests position validation and boundaries
func TestPositionValidation(t *testing.T) {
	board := NewBoard()

	// Test valid positions
	validPositions := []Position{
		{Row: 0, Col: 0},   // A1
		{Row: 7, Col: 7},   // H8 (center)
		{Row: 14, Col: 14}, // O15
		{Row: 0, Col: 14},  // O1
		{Row: 14, Col: 0},  // A15
	}

	for _, pos := range validPositions {
		if !board.IsValidPosition(pos) {
			t.Errorf("Position %s should be valid", pos.String())
		}

		if !pos.IsValid() {
			t.Errorf("Position %s should be valid", pos.String())
		}
	}

	// Test invalid positions
	invalidPositions := []Position{
		{Row: -1, Col: 0},  // Before first row
		{Row: 0, Col: -1},  // Before first column
		{Row: 15, Col: 0},  // After last row
		{Row: 0, Col: 15},  // After last column
		{Row: -1, Col: -1}, // Both invalid
		{Row: 20, Col: 20}, // Way out of bounds
	}

	for _, pos := range invalidPositions {
		if board.IsValidPosition(pos) {
			t.Errorf("Position (%d,%d) should be invalid", pos.Row, pos.Col)
		}

		if pos.IsValid() {
			t.Errorf("Position (%d,%d) should be invalid", pos.Row, pos.Col)
		}
	}

	// Test operations on invalid positions
	invalidPos := Position{Row: -1, Col: -1}

	// Should handle invalid positions gracefully
	if board.GetTile(invalidPos) != nil {
		t.Errorf("GetTile on invalid position should return nil")
	}

	if board.GetSquare(invalidPos) != nil {
		t.Errorf("GetSquare on invalid position should return nil")
	}

	if board.IsEmpty(invalidPos) {
		t.Errorf("IsEmpty on invalid position should return false")
	}

	if board.HasTileAt(invalidPos) {
		t.Errorf("HasTileAt on invalid position should return false")
	}
}

// TestPositionStringConversion tests position to/from string conversion
func TestPositionStringConversion(t *testing.T) {
	testCases := []struct {
		pos      Position
		expected string
	}{
		{Position{Row: 0, Col: 0}, "A1"},
		{Position{Row: 7, Col: 7}, "H8"},
		{Position{Row: 14, Col: 14}, "O15"},
		{Position{Row: 0, Col: 14}, "O1"},
		{Position{Row: 14, Col: 0}, "A15"},
		{Position{Row: 9, Col: 4}, "E10"},
	}

	for _, tc := range testCases {
		// Test Position.String()
		result := tc.pos.String()
		if result != tc.expected {
			t.Errorf("Position %v String() = %s, want %s", tc.pos, result, tc.expected)
		}

		// Test NewPositionFromString()
		parsedPos, err := NewPositionFromString(tc.expected)
		if err != nil {
			t.Errorf("NewPositionFromString(%s) error: %v", tc.expected, err)
		}

		if parsedPos.Row != tc.pos.Row || parsedPos.Col != tc.pos.Col {
			t.Errorf("NewPositionFromString(%s) = %v, want %v", tc.expected, parsedPos, tc.pos)
		}
	}

	// Test invalid string positions
	invalidStrings := []string{
		"",
		"A",
		"A0",
		"A16",
		"P1",
		"Z1",
		"A1X",
		"11",
		"AA",
	}

	for _, invalid := range invalidStrings {
		_, err := NewPositionFromString(invalid)
		if err == nil {
			t.Errorf("NewPositionFromString(%s) should return error", invalid)
		}
	}
}

// TestPremiumSquarePlacement tests that all 61 premium squares are placed correctly
func TestPremiumSquarePlacement(t *testing.T) {
	board := NewBoard()

	// Count premium squares
	counts := board.CountPremiumSquares()

	expectedCounts := map[PremiumType]int{
		Normal:            164, // 225 total - 61 premium squares
		DoubleLetterScore: 24,
		TripleLetterScore: 12,
		DoubleWordScore:   17, // Including center star
		TripleWordScore:   8,
	}

	for premiumType, expected := range expectedCounts {
		if actual := counts[premiumType]; actual != expected {
			t.Errorf("Premium square count for %s: expected %d, got %d",
				premiumType.String(), expected, actual)
		}
	}

	// Test specific premium square positions

	// Triple Word Score positions
	twsPositions := []string{"A1", "A8", "A15", "H1", "H15", "O1", "O8", "O15"}
	for _, posStr := range twsPositions {
		pos, _ := NewPositionFromString(posStr)
		if board.GetPremiumType(pos) != TripleWordScore {
			t.Errorf("Position %s should be Triple Word Score", posStr)
		}
	}

	// Double Word Score positions (including center)
	dwsPositions := []string{"B2", "C3", "D4", "E5", "H8", "K5", "L4", "M3", "N2",
		"B14", "C13", "D12", "E11", "K11", "L12", "M13", "N14"}
	for _, posStr := range dwsPositions {
		pos, _ := NewPositionFromString(posStr)
		if board.GetPremiumType(pos) != DoubleWordScore {
			t.Errorf("Position %s should be Double Word Score", posStr)
		}
	}

	// Triple Letter Score positions
	tlsPositions := []string{"B6", "B10", "F2", "F6", "F10", "F14", "J2", "J6", "J10", "J14", "N6", "N10"}
	for _, posStr := range tlsPositions {
		pos, _ := NewPositionFromString(posStr)
		if board.GetPremiumType(pos) != TripleLetterScore {
			t.Errorf("Position %s should be Triple Letter Score", posStr)
		}
	}

	// Double Letter Score positions
	dlsPositions := []string{"A4", "A12", "C7", "C9", "D1", "D8", "D15", "G3", "G7", "G9", "G13",
		"H4", "H12", "I3", "I7", "I9", "I13", "L1", "L8", "L15", "M7", "M9", "O4", "O12"}
	for _, posStr := range dlsPositions {
		pos, _ := NewPositionFromString(posStr)
		if board.GetPremiumType(pos) != DoubleLetterScore {
			t.Errorf("Position %s should be Double Letter Score", posStr)
		}
	}

	// Test center square specifically
	center := Position{Row: 7, Col: 7}
	if board.GetPremiumType(center) != DoubleWordScore {
		t.Errorf("Center square should be Double Word Score")
	}
}

// TestTilePlacement tests tile placement with various conditions
func TestTilePlacement(t *testing.T) {
	board := NewBoard()

	// Test placing on valid positions
	validPlacements := []struct {
		tile Tile
		pos  Position
	}{
		{Tile{Letter: 'A', Points: 1}, Position{Row: 7, Col: 7}},   // Center
		{Tile{Letter: 'B', Points: 3}, Position{Row: 0, Col: 0}},   // Corner
		{Tile{Letter: 'C', Points: 3}, Position{Row: 14, Col: 14}}, // Opposite corner
	}

	for _, placement := range validPlacements {
		err := board.PlaceTile(placement.tile, placement.pos)
		if err != nil {
			t.Errorf("Should be able to place tile %c at %s: %v",
				placement.tile.Letter, placement.pos.String(), err)
		}

		// Verify tile is placed correctly
		retrievedTile := board.GetTile(placement.pos)
		if retrievedTile == nil || retrievedTile.Letter != placement.tile.Letter {
			t.Errorf("Tile not placed correctly at %s", placement.pos.String())
		}
	}

	// Test placing on invalid positions
	invalidPos := Position{Row: -1, Col: -1}
	err := board.PlaceTile(Tile{Letter: 'X', Points: 8}, invalidPos)
	if err == nil {
		t.Errorf("Should not be able to place tile on invalid position")
	}

	// Test placing on occupied square
	occupiedPos := Position{Row: 7, Col: 7} // Already has tile from above
	err = board.PlaceTile(Tile{Letter: 'Y', Points: 4}, occupiedPos)
	if err == nil {
		t.Errorf("Should not be able to place tile on occupied square")
	}
}

// TestTileRetrieval tests getting tiles from various positions
func TestTileRetrieval(t *testing.T) {
	board := NewBoard()

	// Place some tiles
	testTiles := []struct {
		tile Tile
		pos  Position
	}{
		{Tile{Letter: 'T', Points: 1}, Position{Row: 5, Col: 5}},
		{Tile{Letter: 'E', Points: 1}, Position{Row: 5, Col: 6}},
		{Tile{Letter: 'S', Points: 1}, Position{Row: 5, Col: 7}},
		{Tile{Letter: 'T', Points: 1}, Position{Row: 5, Col: 8}},
	}

	for _, tt := range testTiles {
		board.PlaceTile(tt.tile, tt.pos)
	}

	// Test retrieving placed tiles
	for _, tt := range testTiles {
		tile := board.GetTile(tt.pos)
		if tile == nil {
			t.Errorf("Should retrieve tile at %s", tt.pos.String())
		}

		if tile.Letter != tt.tile.Letter {
			t.Errorf("Retrieved tile letter %c != expected %c at %s",
				tile.Letter, tt.tile.Letter, tt.pos.String())
		}
	}

	// Test retrieving from empty positions
	emptyPositions := []Position{
		{Row: 0, Col: 0},
		{Row: 10, Col: 10},
		{Row: 14, Col: 14},
	}

	for _, pos := range emptyPositions {
		tile := board.GetTile(pos)
		if tile != nil {
			t.Errorf("Should not retrieve tile from empty position %s", pos.String())
		}
	}

	// Test retrieving from invalid positions
	invalidPos := Position{Row: -1, Col: -1}
	tile := board.GetTile(invalidPos)
	if tile != nil {
		t.Errorf("Should not retrieve tile from invalid position")
	}
}

// TestAdjacentPositions tests adjacency logic for edges, corners, and center
func TestAdjacentPositions(t *testing.T) {
	board := NewBoard()

	testCases := []struct {
		pos              Position
		expectedCount    int
		description      string
		expectedAdjacent []Position
	}{
		{
			pos:           Position{Row: 7, Col: 7}, // Center H8
			expectedCount: 4,
			description:   "center",
			expectedAdjacent: []Position{
				{Row: 6, Col: 7}, {Row: 8, Col: 7}, // Up, Down
				{Row: 7, Col: 6}, {Row: 7, Col: 8}, // Left, Right
			},
		},
		{
			pos:           Position{Row: 0, Col: 0}, // Corner A1
			expectedCount: 2,
			description:   "corner A1",
			expectedAdjacent: []Position{
				{Row: 1, Col: 0}, {Row: 0, Col: 1}, // Down, Right
			},
		},
		{
			pos:           Position{Row: 14, Col: 14}, // Corner O15
			expectedCount: 2,
			description:   "corner O15",
			expectedAdjacent: []Position{
				{Row: 13, Col: 14}, {Row: 14, Col: 13}, // Up, Left
			},
		},
		{
			pos:           Position{Row: 0, Col: 7}, // Top edge H1
			expectedCount: 3,
			description:   "top edge",
			expectedAdjacent: []Position{
				{Row: 1, Col: 7},                   // Down
				{Row: 0, Col: 6}, {Row: 0, Col: 8}, // Left, Right
			},
		},
		{
			pos:           Position{Row: 14, Col: 7}, // Bottom edge H15
			expectedCount: 3,
			description:   "bottom edge",
			expectedAdjacent: []Position{
				{Row: 13, Col: 7},                    // Up
				{Row: 14, Col: 6}, {Row: 14, Col: 8}, // Left, Right
			},
		},
		{
			pos:           Position{Row: 7, Col: 0}, // Left edge A8
			expectedCount: 3,
			description:   "left edge",
			expectedAdjacent: []Position{
				{Row: 6, Col: 0}, {Row: 8, Col: 0}, // Up, Down
				{Row: 7, Col: 1}, // Right
			},
		},
		{
			pos:           Position{Row: 7, Col: 14}, // Right edge O8
			expectedCount: 3,
			description:   "right edge",
			expectedAdjacent: []Position{
				{Row: 6, Col: 14}, {Row: 8, Col: 14}, // Up, Down
				{Row: 7, Col: 13}, // Left
			},
		},
	}

	for _, tc := range testCases {
		adjacent := board.GetAdjacentPositions(tc.pos)

		if len(adjacent) != tc.expectedCount {
			t.Errorf("%s position %s: expected %d adjacent, got %d",
				tc.description, tc.pos.String(), tc.expectedCount, len(adjacent))
		}

		// Check that all expected positions are present
		for _, expectedPos := range tc.expectedAdjacent {
			found := false
			for _, actualPos := range adjacent {
				if actualPos.Row == expectedPos.Row && actualPos.Col == expectedPos.Col {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s position %s: missing expected adjacent position %s",
					tc.description, tc.pos.String(), expectedPos.String())
			}
		}

		// Verify all returned positions are valid
		for _, pos := range adjacent {
			if !board.IsValidPosition(pos) {
				t.Errorf("Adjacent position %s should be valid", pos.String())
			}
		}
	}

	// Test invalid position
	invalidPos := Position{Row: -1, Col: -1}
	adjacent := board.GetAdjacentPositions(invalidPos)
	if adjacent != nil {
		t.Errorf("GetAdjacentPositions on invalid position should return nil")
	}
}

// TestBoardStateValidation tests comprehensive board validation
func TestBoardStateValidation(t *testing.T) {
	// Test valid board
	board := NewBoard()
	if err := board.ValidateBoard(); err != nil {
		t.Errorf("New board should be valid: %v", err)
	}

	// Test board with modified center
	board.Center = Position{Row: 0, Col: 0}
	if err := board.ValidateBoard(); err == nil {
		t.Errorf("Board with wrong center should be invalid")
	}

	// Reset center
	board.Center = Position{Row: 7, Col: 7}

	// Test board with wrong premium square at center
	board.Grid[7][7].Premium = Normal
	if err := board.ValidateBoard(); err == nil {
		t.Errorf("Board with wrong center premium should be invalid")
	}

	// Reset center premium
	board.Grid[7][7].Premium = DoubleWordScore

	// Test board with wrong premium square counts
	board.Grid[0][0].Premium = Normal // Change A1 from TWS to Normal
	if err := board.ValidateBoard(); err == nil {
		t.Errorf("Board with wrong premium counts should be invalid")
	}
}

// TestRowColumnOperations tests getting rows and columns
func TestRowColumnOperations(t *testing.T) {
	board := NewBoard()

	// Test valid row
	row := board.GetRow(7) // Middle row
	if len(row) != 15 {
		t.Errorf("Row should have 15 squares, got %d", len(row))
	}

	// Test valid column
	col := board.GetColumn(7) // Middle column
	if len(col) != 15 {
		t.Errorf("Column should have 15 squares, got %d", len(col))
	}

	// Test invalid row
	invalidRow := board.GetRow(-1)
	if invalidRow != nil {
		t.Errorf("Invalid row should return nil")
	}

	invalidRow = board.GetRow(15)
	if invalidRow != nil {
		t.Errorf("Invalid row should return nil")
	}

	// Test invalid column
	invalidCol := board.GetColumn(-1)
	if invalidCol != nil {
		t.Errorf("Invalid column should return nil")
	}

	invalidCol = board.GetColumn(15)
	if invalidCol != nil {
		t.Errorf("Invalid column should return nil")
	}
}

// TestOccupiedPositions tests tracking of occupied squares
func TestOccupiedPositions(t *testing.T) {
	board := NewBoard()

	// Initially no occupied positions
	occupied := board.GetOccupiedPositions()
	if len(occupied) != 0 {
		t.Errorf("New board should have no occupied positions, got %d", len(occupied))
	}

	// Place some tiles
	positions := []Position{
		{Row: 7, Col: 7}, // H8
		{Row: 7, Col: 8}, // I8
		{Row: 7, Col: 9}, // J8
	}

	for i, pos := range positions {
		tile := Tile{Letter: rune('A' + i), Points: 1}
		board.PlaceTile(tile, pos)
	}

	// Check occupied positions
	occupied = board.GetOccupiedPositions()
	if len(occupied) != len(positions) {
		t.Errorf("Should have %d occupied positions, got %d", len(positions), len(occupied))
	}

	// Verify all placed positions are in occupied list
	for _, placedPos := range positions {
		found := false
		for _, occupiedPos := range occupied {
			if occupiedPos.Row == placedPos.Row && occupiedPos.Col == placedPos.Col {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Placed position %s not found in occupied positions", placedPos.String())
		}
	}

	// Test first move detection
	if board.IsFirstMove() {
		t.Errorf("Board with tiles should not be first move")
	}

	// Remove a tile and verify occupied positions update
	board.RemoveTile(positions[0])
	occupied = board.GetOccupiedPositions()
	if len(occupied) != len(positions)-1 {
		t.Errorf("Should have %d occupied positions after removal, got %d",
			len(positions)-1, len(occupied))
	}
}

// TestPremiumTypeString tests premium type string representations
func TestPremiumTypeString(t *testing.T) {
	testCases := []struct {
		premium  PremiumType
		expected string
	}{
		{Normal, "NORMAL"},
		{DoubleLetterScore, "DLS"},
		{TripleLetterScore, "TLS"},
		{DoubleWordScore, "DWS"},
		{TripleWordScore, "TWS"},
	}

	for _, tc := range testCases {
		result := tc.premium.String()
		if result != tc.expected {
			t.Errorf("PremiumType %d String() = %s, want %s", tc.premium, result, tc.expected)
		}
	}

	// Test unknown premium type
	unknown := PremiumType(999)
	if unknown.String() != "UNKNOWN" {
		t.Errorf("Unknown premium type should return 'UNKNOWN'")
	}
}

// TestBoardString tests board string representation
func TestBoardString(t *testing.T) {
	board := NewBoard()

	// Test that string representation can be generated without error
	boardStr := board.String()
	if len(boardStr) == 0 {
		t.Errorf("Board string representation should not be empty")
	}

	// Place a tile and verify it appears in string
	tile := Tile{Letter: 'A', Points: 1}
	pos := Position{Row: 7, Col: 7}
	board.PlaceTile(tile, pos)

	boardStr = board.String()
	// The string should contain the tile representation
	// Note: This is a basic test; more detailed string format testing could be added
	if len(boardStr) == 0 {
		t.Errorf("Board string with tile should not be empty")
	}
}
