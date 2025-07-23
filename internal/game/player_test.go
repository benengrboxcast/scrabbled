package game

import (
	"sync"
	"testing"
)

// TestPlayerCreationAndInitialization tests player creation with various parameters
func TestPlayerCreationAndInitialization(t *testing.T) {
	// Test creating player with ID and name
	player := NewPlayer("player1", "Alice")

	if player.GetID() != "player1" {
		t.Errorf("Expected ID 'player1', got '%s'", player.GetID())
	}

	if player.GetName() != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", player.GetName())
	}

	if player.GetScore() != 0 {
		t.Errorf("Expected initial score 0, got %d", player.GetScore())
	}

	if !player.IsPlayerActive() {
		t.Errorf("Expected player to be active initially")
	}

	if player.GetRackSize() != 0 {
		t.Errorf("Expected empty rack initially, got size %d", player.GetRackSize())
	}

	if !player.IsRackEmpty() {
		t.Errorf("Expected rack to be empty initially")
	}

	if player.IsRackFull() {
		t.Errorf("Expected rack not to be full initially")
	}

	// Test player validation
	if err := player.ValidatePlayer(); err != nil {
		t.Errorf("New player should be valid: %v", err)
	}

	// Test creating player with empty ID (should generate one)
	player2 := NewPlayer("", "Bob")
	if player2.GetID() == "" {
		t.Errorf("Player with empty ID should have ID generated")
	}

	if player2.GetName() != "Bob" {
		t.Errorf("Expected name 'Bob', got '%s'", player2.GetName())
	}
}

// TestPlayerScoreManagement tests score operations
func TestPlayerScoreManagement(t *testing.T) {
	player := NewPlayer("test", "TestPlayer")

	// Test initial score
	if player.GetScore() != 0 {
		t.Errorf("Expected initial score 0, got %d", player.GetScore())
	}

	// Test adding score
	player.AddScore(10)
	if player.GetScore() != 10 {
		t.Errorf("Expected score 10 after adding 10, got %d", player.GetScore())
	}

	// Test adding more score
	player.AddScore(25)
	if player.GetScore() != 35 {
		t.Errorf("Expected score 35 after adding 25, got %d", player.GetScore())
	}

	// Test setting score directly
	player.SetScore(100)
	if player.GetScore() != 100 {
		t.Errorf("Expected score 100 after setting, got %d", player.GetScore())
	}

	// Test adding negative score
	player.AddScore(-20)
	if player.GetScore() != 80 {
		t.Errorf("Expected score 80 after subtracting 20, got %d", player.GetScore())
	}
}

// TestPlayerActiveStatus tests active status management
func TestPlayerActiveStatus(t *testing.T) {
	player := NewPlayer("test", "TestPlayer")

	// Test initial active status
	if !player.IsPlayerActive() {
		t.Errorf("Expected player to be active initially")
	}

	// Test setting inactive
	player.SetActive(false)
	if player.IsPlayerActive() {
		t.Errorf("Expected player to be inactive after setting")
	}

	// Test setting active again
	player.SetActive(true)
	if !player.IsPlayerActive() {
		t.Errorf("Expected player to be active after setting")
	}
}

// TestPlayerNameManagement tests name operations
func TestPlayerNameManagement(t *testing.T) {
	player := NewPlayer("test", "OriginalName")

	if player.GetName() != "OriginalName" {
		t.Errorf("Expected name 'OriginalName', got '%s'", player.GetName())
	}

	// Test changing name
	player.SetName("NewName")
	if player.GetName() != "NewName" {
		t.Errorf("Expected name 'NewName', got '%s'", player.GetName())
	}
}

// TestAddTilesToRack tests adding tiles with various conditions
func TestAddTilesToRack(t *testing.T) {
	player := NewPlayer("test", "TestPlayer")

	// Test adding tiles to empty rack
	tiles := []Tile{
		{Letter: 'A', Points: 1, IsBlank: false},
		{Letter: 'B', Points: 3, IsBlank: false},
		{Letter: 'C', Points: 3, IsBlank: false},
	}

	err := player.AddTilesToRack(tiles)
	if err != nil {
		t.Errorf("Should be able to add tiles to empty rack: %v", err)
	}

	if player.GetRackSize() != 3 {
		t.Errorf("Expected rack size 3, got %d", player.GetRackSize())
	}

	// Test adding more tiles (up to limit)
	moreTiles := []Tile{
		{Letter: 'D', Points: 2, IsBlank: false},
		{Letter: 'E', Points: 1, IsBlank: false},
		{Letter: 'F', Points: 4, IsBlank: false},
		{Letter: 'G', Points: 2, IsBlank: false},
	}

	err = player.AddTilesToRack(moreTiles)
	if err != nil {
		t.Errorf("Should be able to add tiles up to limit: %v", err)
	}

	if player.GetRackSize() != 7 {
		t.Errorf("Expected rack size 7, got %d", player.GetRackSize())
	}

	if !player.IsRackFull() {
		t.Errorf("Expected rack to be full")
	}

	// Test adding tiles beyond limit (should fail)
	extraTile := []Tile{{Letter: 'H', Points: 4, IsBlank: false}}
	err = player.AddTilesToRack(extraTile)
	if err == nil {
		t.Errorf("Should not be able to add tiles beyond rack limit")
	}

	// Test adding empty slice (should not error)
	err = player.AddTilesToRack([]Tile{})
	if err != nil {
		t.Errorf("Adding empty slice should not error: %v", err)
	}

	// Test rack contents
	rack := player.GetRack()
	if len(rack) != 7 {
		t.Errorf("GetRack should return 7 tiles, got %d", len(rack))
	}

	// Verify rack is a copy (modifying returned slice shouldn't affect player)
	rack[0] = Tile{Letter: 'Z', Points: 10, IsBlank: false}
	originalRack := player.GetRack()
	if originalRack[0].Letter == 'Z' {
		t.Errorf("GetRack should return a copy, not the original rack")
	}
}

// TestRemoveTilesFromRack tests tile removal with invalid indices and empty rack
func TestRemoveTilesFromRack(t *testing.T) {
	player := NewPlayer("test", "TestPlayer")

	// Test removing from empty rack
	_, err := player.RemoveTilesFromRack([]int{0})
	if err == nil {
		t.Errorf("Should not be able to remove from empty rack")
	}

	// Add some tiles first
	tiles := []Tile{
		{Letter: 'A', Points: 1, IsBlank: false},
		{Letter: 'B', Points: 3, IsBlank: false},
		{Letter: 'C', Points: 3, IsBlank: false},
		{Letter: 'D', Points: 2, IsBlank: false},
		{Letter: 'E', Points: 1, IsBlank: false},
	}
	player.AddTilesToRack(tiles)

	// Test removing valid indices
	removedTiles, err := player.RemoveTilesFromRack([]int{1, 3}) // Remove B and D
	if err != nil {
		t.Errorf("Should be able to remove valid indices: %v", err)
	}

	if len(removedTiles) != 2 {
		t.Errorf("Expected 2 removed tiles, got %d", len(removedTiles))
	}

	// Check removed tiles (note: order matches original indices, not sorted removal order)
	if removedTiles[0].Letter != 'B' || removedTiles[1].Letter != 'D' {
		t.Errorf("Removed tiles don't match expected. Got: %c, %c", removedTiles[0].Letter, removedTiles[1].Letter)
	}

	if player.GetRackSize() != 3 {
		t.Errorf("Expected rack size 3 after removal, got %d", player.GetRackSize())
	}

	// Test removing invalid indices
	_, err = player.RemoveTilesFromRack([]int{10})
	if err == nil {
		t.Errorf("Should not be able to remove invalid index")
	}

	_, err = player.RemoveTilesFromRack([]int{-1})
	if err == nil {
		t.Errorf("Should not be able to remove negative index")
	}

	// Test removing duplicate indices
	_, err = player.RemoveTilesFromRack([]int{0, 0})
	if err == nil {
		t.Errorf("Should not be able to remove duplicate indices")
	}

	// Test removing empty slice (should not error)
	removedTiles, err = player.RemoveTilesFromRack([]int{})
	if err != nil {
		t.Errorf("Removing empty slice should not error: %v", err)
	}
	if len(removedTiles) != 0 {
		t.Errorf("Expected empty slice when removing no indices")
	}

	// Test removing all remaining tiles
	remainingIndices := []int{0, 1, 2} // All remaining tiles
	removedTiles, err = player.RemoveTilesFromRack(remainingIndices)
	if err != nil {
		t.Errorf("Should be able to remove all remaining tiles: %v", err)
	}

	if len(removedTiles) != 3 {
		t.Errorf("Expected 3 removed tiles, got %d", len(removedTiles))
	}

	if !player.IsRackEmpty() {
		t.Errorf("Expected rack to be empty after removing all tiles")
	}
}

// TestRemoveTilesByValue tests removing tiles by their values
func TestRemoveTilesByValue(t *testing.T) {
	player := NewPlayer("test", "TestPlayer")

	// Add some tiles
	tiles := []Tile{
		{Letter: 'A', Points: 1, IsBlank: false},
		{Letter: 'B', Points: 3, IsBlank: false},
		{Letter: 'A', Points: 1, IsBlank: false}, // Duplicate A
		{Letter: 'C', Points: 3, IsBlank: false},
		{Letter: 0, Points: 0, IsBlank: true}, // Blank tile
	}
	player.AddTilesToRack(tiles)

	// Test removing existing tiles
	tilesToRemove := []Tile{
		{Letter: 'B', Points: 3, IsBlank: false},
		{Letter: 'A', Points: 1, IsBlank: false}, // Should remove only one A
	}

	err := player.RemoveTilesByValue(tilesToRemove)
	if err != nil {
		t.Errorf("Should be able to remove existing tiles: %v", err)
	}

	if player.GetRackSize() != 3 {
		t.Errorf("Expected rack size 3 after removal, got %d", player.GetRackSize())
	}

	// Test that one A is still there
	if player.GetTileCount('A') != 1 {
		t.Errorf("Expected 1 A tile remaining, got %d", player.GetTileCount('A'))
	}

	// Test removing non-existing tile
	nonExistentTile := []Tile{{Letter: 'Z', Points: 10, IsBlank: false}}
	err = player.RemoveTilesByValue(nonExistentTile)
	if err == nil {
		t.Errorf("Should not be able to remove non-existent tile")
	}

	// Test removing blank tile
	blankTile := []Tile{{Letter: 0, Points: 0, IsBlank: true}}
	err = player.RemoveTilesByValue(blankTile)
	if err != nil {
		t.Errorf("Should be able to remove blank tile: %v", err)
	}

	if player.GetBlankCount() != 0 {
		t.Errorf("Expected 0 blank tiles after removal, got %d", player.GetBlankCount())
	}

	// Test removing empty slice
	err = player.RemoveTilesByValue([]Tile{})
	if err != nil {
		t.Errorf("Removing empty slice should not error: %v", err)
	}
}

// TestRackSizeCalculation tests rack size methods
func TestRackSizeCalculation(t *testing.T) {
	player := NewPlayer("test", "TestPlayer")

	// Test empty rack
	if player.GetRackSize() != 0 {
		t.Errorf("Expected rack size 0 for empty rack, got %d", player.GetRackSize())
	}

	if !player.IsRackEmpty() {
		t.Errorf("Expected empty rack to return true for IsRackEmpty")
	}

	if player.IsRackFull() {
		t.Errorf("Expected empty rack to return false for IsRackFull")
	}

	// Add tiles one by one and test size
	for i := 1; i <= MaxRackSize; i++ {
		tile := Tile{Letter: rune('A' + i - 1), Points: 1, IsBlank: false}
		player.AddTilesToRack([]Tile{tile})

		if player.GetRackSize() != i {
			t.Errorf("Expected rack size %d, got %d", i, player.GetRackSize())
		}

		if player.IsRackEmpty() {
			t.Errorf("Expected non-empty rack to return false for IsRackEmpty")
		}

		if i == MaxRackSize {
			if !player.IsRackFull() {
				t.Errorf("Expected full rack to return true for IsRackFull")
			}
		} else {
			if player.IsRackFull() {
				t.Errorf("Expected non-full rack to return false for IsRackFull")
			}
		}
	}
}

// TestHasTileMethods tests tile existence checking
func TestHasTileMethods(t *testing.T) {
	player := NewPlayer("test", "TestPlayer")

	// Add some tiles
	tiles := []Tile{
		{Letter: 'A', Points: 1, IsBlank: false},
		{Letter: 'B', Points: 3, IsBlank: false},
		{Letter: 'A', Points: 1, IsBlank: false}, // Duplicate A
		{Letter: 0, Points: 0, IsBlank: true},    // Blank tile
	}
	player.AddTilesToRack(tiles)

	// Test HasTile for existing tiles
	aTitle := Tile{Letter: 'A', Points: 1, IsBlank: false}
	if !player.HasTile(aTitle) {
		t.Errorf("Player should have A tile")
	}

	bTile := Tile{Letter: 'B', Points: 3, IsBlank: false}
	if !player.HasTile(bTile) {
		t.Errorf("Player should have B tile")
	}

	blankTile := Tile{Letter: 0, Points: 0, IsBlank: true}
	if !player.HasTile(blankTile) {
		t.Errorf("Player should have blank tile")
	}

	// Test HasTile for non-existing tile
	zTile := Tile{Letter: 'Z', Points: 10, IsBlank: false}
	if player.HasTile(zTile) {
		t.Errorf("Player should not have Z tile")
	}

	// Test HasTiles for multiple tiles
	existingTiles := []Tile{aTitle, bTile}
	if !player.HasTiles(existingTiles) {
		t.Errorf("Player should have A and B tiles")
	}

	// Test HasTiles for multiple A tiles (should succeed with 2 A's)
	twoATiles := []Tile{aTitle, aTitle}
	if !player.HasTiles(twoATiles) {
		t.Errorf("Player should have two A tiles")
	}

	// Test HasTiles for too many of same tile
	threeATiles := []Tile{aTitle, aTitle, aTitle}
	if player.HasTiles(threeATiles) {
		t.Errorf("Player should not have three A tiles")
	}

	// Test HasTiles with non-existing tile
	mixedTiles := []Tile{aTitle, zTile}
	if player.HasTiles(mixedTiles) {
		t.Errorf("Player should not have A and Z tiles")
	}

	// Test HasTiles with empty slice
	if !player.HasTiles([]Tile{}) {
		t.Errorf("Player should have zero tiles (empty slice)")
	}
}

// TestTileCountMethods tests tile counting functionality
func TestTileCountMethods(t *testing.T) {
	player := NewPlayer("test", "TestPlayer")

	// Add some tiles
	tiles := []Tile{
		{Letter: 'A', Points: 1, IsBlank: false},
		{Letter: 'A', Points: 1, IsBlank: false},
		{Letter: 'A', Points: 1, IsBlank: false}, // 3 A's
		{Letter: 'B', Points: 3, IsBlank: false}, // 1 B
		{Letter: 0, Points: 0, IsBlank: true},    // 1 Blank
		{Letter: 0, Points: 0, IsBlank: true},    // 2 Blanks
	}
	player.AddTilesToRack(tiles)

	// Test GetTileCount
	if player.GetTileCount('A') != 3 {
		t.Errorf("Expected 3 A tiles, got %d", player.GetTileCount('A'))
	}

	if player.GetTileCount('B') != 1 {
		t.Errorf("Expected 1 B tile, got %d", player.GetTileCount('B'))
	}

	if player.GetTileCount('Z') != 0 {
		t.Errorf("Expected 0 Z tiles, got %d", player.GetTileCount('Z'))
	}

	// Test GetBlankCount
	if player.GetBlankCount() != 2 {
		t.Errorf("Expected 2 blank tiles, got %d", player.GetBlankCount())
	}

	// Test GetRackValue
	expectedValue := 3*1 + 1*3 + 2*0 // 3 A's + 1 B + 2 blanks
	if player.GetRackValue() != expectedValue {
		t.Errorf("Expected rack value %d, got %d", expectedValue, player.GetRackValue())
	}
}

// TestClearRack tests clearing the entire rack
func TestClearRack(t *testing.T) {
	player := NewPlayer("test", "TestPlayer")

	// Add some tiles
	tiles := []Tile{
		{Letter: 'A', Points: 1, IsBlank: false},
		{Letter: 'B', Points: 3, IsBlank: false},
		{Letter: 'C', Points: 3, IsBlank: false},
	}
	player.AddTilesToRack(tiles)

	// Test ClearRack
	clearedTiles := player.ClearRack()

	if len(clearedTiles) != 3 {
		t.Errorf("Expected 3 cleared tiles, got %d", len(clearedTiles))
	}

	if !player.IsRackEmpty() {
		t.Errorf("Expected rack to be empty after clearing")
	}

	if player.GetRackSize() != 0 {
		t.Errorf("Expected rack size 0 after clearing, got %d", player.GetRackSize())
	}

	// Verify cleared tiles match original tiles
	for i, tile := range clearedTiles {
		if !tilesEqual(tile, tiles[i]) {
			t.Errorf("Cleared tile %d doesn't match original", i)
		}
	}

	// Test clearing empty rack
	clearedEmpty := player.ClearRack()
	if len(clearedEmpty) != 0 {
		t.Errorf("Expected empty slice when clearing empty rack, got %d tiles", len(clearedEmpty))
	}
}

// TestPlayerStateValidation tests comprehensive player validation
func TestPlayerStateValidation(t *testing.T) {
	// Test valid player
	validPlayer := NewPlayer("valid", "ValidPlayer")
	if err := validPlayer.ValidatePlayer(); err != nil {
		t.Errorf("Valid player should pass validation: %v", err)
	}

	// Test invalid ID
	invalidIDPlayer := &Player{
		ID:       "", // Invalid: empty ID
		Name:     "TestPlayer",
		Rack:     []Tile{},
		Score:    0,
		IsActive: true,
	}
	if err := invalidIDPlayer.ValidatePlayer(); err == nil {
		t.Errorf("Player with empty ID should fail validation")
	}

	// Test invalid name
	invalidNamePlayer := &Player{
		ID:       "test",
		Name:     "", // Invalid: empty name
		Rack:     []Tile{},
		Score:    0,
		IsActive: true,
	}
	if err := invalidNamePlayer.ValidatePlayer(); err == nil {
		t.Errorf("Player with empty name should fail validation")
	}

	// Test rack size exceeding maximum
	oversizedRack := make([]Tile, MaxRackSize+1)
	for i := 0; i <= MaxRackSize; i++ {
		oversizedRack[i] = Tile{Letter: 'A', Points: 1, IsBlank: false}
	}
	oversizedPlayer := &Player{
		ID:       "test",
		Name:     "TestPlayer",
		Rack:     oversizedRack,
		Score:    0,
		IsActive: true,
	}
	if err := oversizedPlayer.ValidatePlayer(); err == nil {
		t.Errorf("Player with oversized rack should fail validation")
	}

	// Test negative score
	negativeScorePlayer := &Player{
		ID:       "test",
		Name:     "TestPlayer",
		Rack:     []Tile{},
		Score:    -10, // Invalid: negative score
		IsActive: true,
	}
	if err := negativeScorePlayer.ValidatePlayer(); err == nil {
		t.Errorf("Player with negative score should fail validation")
	}

	// Test invalid tile with negative points
	invalidTilePlayer := &Player{
		ID:   "test",
		Name: "TestPlayer",
		Rack: []Tile{
			{Letter: 'A', Points: -1, IsBlank: false}, // Invalid: negative points
		},
		Score:    0,
		IsActive: true,
	}
	if err := invalidTilePlayer.ValidatePlayer(); err == nil {
		t.Errorf("Player with invalid tile should fail validation")
	}

	// Test invalid blank tile
	invalidBlankPlayer := &Player{
		ID:   "test",
		Name: "TestPlayer",
		Rack: []Tile{
			{Letter: 'A', Points: 1, IsBlank: true}, // Invalid: blank tile with letter and points
		},
		Score:    0,
		IsActive: true,
	}
	if err := invalidBlankPlayer.ValidatePlayer(); err == nil {
		t.Errorf("Player with invalid blank tile should fail validation")
	}

	// Test non-blank tile without letter
	noLetterPlayer := &Player{
		ID:   "test",
		Name: "TestPlayer",
		Rack: []Tile{
			{Letter: 0, Points: 1, IsBlank: false}, // Invalid: non-blank tile without letter
		},
		Score:    0,
		IsActive: true,
	}
	if err := noLetterPlayer.ValidatePlayer(); err == nil {
		t.Errorf("Player with non-blank tile without letter should fail validation")
	}
}

// TestPlayerStringRepresentation tests the String() method
func TestPlayerStringRepresentation(t *testing.T) {
	player := NewPlayer("test123", "Alice")
	player.SetScore(42)

	// Add some tiles
	tiles := []Tile{
		{Letter: 'A', Points: 1, IsBlank: false},
		{Letter: 'B', Points: 3, IsBlank: false},
		{Letter: 0, Points: 0, IsBlank: true}, // Blank tile
	}
	player.AddTilesToRack(tiles)

	str := player.String()

	// Check that string contains key information
	if str == "" {
		t.Errorf("String representation should not be empty")
	}

	// Basic checks for presence of key data (exact format may vary)
	expectedSubstrings := []string{"test123", "Alice", "42", "3/7"}
	for _, substr := range expectedSubstrings {
		if !contains(str, substr) {
			t.Errorf("String representation should contain '%s', got: %s", substr, str)
		}
	}

	// Test inactive player
	player.SetActive(false)
	str = player.String()
	if !contains(str, "inactive") {
		t.Errorf("String representation should show inactive status")
	}
}

// contains checks if a string contains a substring (helper function)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsAtPosition(s, substr))))
}

func containsAtPosition(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestConcurrentPlayerOperations tests thread safety
func TestConcurrentPlayerOperations(t *testing.T) {
	player := NewPlayer("concurrent", "ConcurrentPlayer")

	// Add initial tiles
	initialTiles := []Tile{
		{Letter: 'A', Points: 1, IsBlank: false},
		{Letter: 'B', Points: 3, IsBlank: false},
		{Letter: 'C', Points: 3, IsBlank: false},
	}
	player.AddTilesToRack(initialTiles)

	const numGoroutines = 10
	const operationsPerGoroutine = 50

	var wg sync.WaitGroup

	// Concurrent score updates
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				player.AddScore(1)
				player.GetScore()
			}
		}()
	}

	// Concurrent rack operations
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				player.GetRackSize()
				player.GetRack()
				player.IsRackFull()
				player.IsRackEmpty()
				player.GetTileCount('A')
				player.GetBlankCount()
				player.GetRackValue()
			}
		}()
	}

	// Concurrent status operations
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				player.IsPlayerActive()
				if id%2 == 0 {
					player.SetActive(true)
				} else {
					player.SetActive(false)
				}
				player.GetName()
				player.GetID()
			}
		}(i)
	}

	wg.Wait()

	// Verify final state is consistent
	expectedScore := numGoroutines * operationsPerGoroutine
	if player.GetScore() != expectedScore {
		t.Errorf("Expected final score %d, got %d", expectedScore, player.GetScore())
	}

	// Verify rack is still valid
	if err := player.ValidatePlayer(); err != nil {
		t.Errorf("Player should be valid after concurrent operations: %v", err)
	}
}

// TestConcurrentRackModification tests thread safety of rack modifications
func TestConcurrentRackModification(t *testing.T) {
	player := NewPlayer("rackmod", "RackModPlayer")

	const numGoroutines = 5
	var wg sync.WaitGroup

	// Concurrent tile additions and removals
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Add tiles
			tiles := []Tile{
				{Letter: rune('A' + id), Points: 1, IsBlank: false},
			}
			player.AddTilesToRack(tiles)

			// Check state
			player.GetRackSize()
			player.HasTile(tiles[0])

			// Remove tiles if possible
			if player.GetRackSize() > 0 {
				player.RemoveTilesFromRack([]int{0})
			}
		}(i)
	}

	wg.Wait()

	// Verify final state is valid
	if err := player.ValidatePlayer(); err != nil {
		t.Errorf("Player should be valid after concurrent rack modifications: %v", err)
	}
}
