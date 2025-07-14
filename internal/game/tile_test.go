package game

import (
	"sync"
	"testing"
)

// TestTileStruct tests the basic Tile struct validation
func TestTileStruct(t *testing.T) {
	tests := []struct {
		name     string
		tile     Tile
		expected string
	}{
		{
			name:     "Regular letter tile",
			tile:     Tile{Letter: 'A', Points: 1, IsBlank: false},
			expected: "A",
		},
		{
			name:     "High value tile",
			tile:     Tile{Letter: 'Q', Points: 10, IsBlank: false},
			expected: "Q",
		},
		{
			name:     "Blank tile",
			tile:     Tile{Letter: 0, Points: 0, IsBlank: true},
			expected: "BLANK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tile.String() != tt.expected {
				t.Errorf("Tile.String() = %v, want %v", tt.tile.String(), tt.expected)
			}
		})
	}
}

// TestTilePointValues tests that tile point values are correct according to Scrabble rules
func TestTilePointValues(t *testing.T) {
	tests := []struct {
		letter rune
		points int
	}{
		// 1 point tiles
		{'A', 1}, {'E', 1}, {'I', 1}, {'O', 1}, {'U', 1},
		{'L', 1}, {'N', 1}, {'S', 1}, {'T', 1}, {'R', 1},
		// 2 point tiles
		{'D', 2}, {'G', 2},
		// 3 point tiles
		{'B', 3}, {'C', 3}, {'M', 3}, {'P', 3},
		// 4 point tiles
		{'F', 4}, {'H', 4}, {'V', 4}, {'W', 4}, {'Y', 4},
		// 5 point tiles
		{'K', 5},
		// 8 point tiles
		{'J', 8}, {'X', 8},
		// 10 point tiles
		{'Q', 10}, {'Z', 10},
		// Blank tiles
		{0, 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.letter), func(t *testing.T) {
			value := GetTileValue(tt.letter)
			if value != tt.points {
				t.Errorf("GetTileValue(%c) = %d, want %d", tt.letter, value, tt.points)
			}
		})
	}

	// Test invalid letter
	if GetTileValue('1') != 0 {
		t.Errorf("GetTileValue('1') should return 0 for invalid letter")
	}
}

// TestTileDistribution tests that the tile bag has the correct distribution
func TestTileDistribution(t *testing.T) {
	// Test the distribution validation function
	if err := ValidateTileDistribution(); err != nil {
		t.Fatalf("ValidateTileDistribution() failed: %v", err)
	}

	// Create a new tile bag and verify distribution
	bag := NewTileBag()

	if bag.RemainingCount() != 100 {
		t.Errorf("NewTileBag() should contain 100 tiles, got %d", bag.RemainingCount())
	}

	// Count tiles by type
	letterCounts := make(map[rune]int)
	blankCount := 0

	// Draw all tiles to count them
	allTiles := bag.DrawTiles(100)

	for _, tile := range allTiles {
		if tile.IsBlank {
			blankCount++
		} else {
			letterCounts[tile.Letter]++
		}
	}

	// Verify letter distribution
	expectedDistribution := map[rune]int{
		'A': 9, 'B': 2, 'C': 2, 'D': 4, 'E': 12, 'F': 2, 'G': 3, 'H': 2,
		'I': 9, 'J': 1, 'K': 1, 'L': 4, 'M': 2, 'N': 6, 'O': 8, 'P': 2,
		'Q': 1, 'R': 6, 'S': 4, 'T': 6, 'U': 4, 'V': 2, 'W': 2, 'X': 1,
		'Y': 2, 'Z': 1,
	}

	for letter, expectedCount := range expectedDistribution {
		if actualCount := letterCounts[letter]; actualCount != expectedCount {
			t.Errorf("Letter %c: expected %d tiles, got %d", letter, expectedCount, actualCount)
		}
	}

	// Verify blank tiles
	if blankCount != 2 {
		t.Errorf("Expected 2 blank tiles, got %d", blankCount)
	}

	// Verify total
	totalCounted := blankCount
	for _, count := range letterCounts {
		totalCounted += count
	}
	if totalCounted != 100 {
		t.Errorf("Total tiles counted: %d, expected 100", totalCounted)
	}
}

// TestDrawTiles tests the DrawTiles functionality with various scenarios
func TestDrawTiles(t *testing.T) {
	t.Run("Draw normal amount", func(t *testing.T) {
		bag := NewTileBag()
		initialCount := bag.RemainingCount()

		drawn := bag.DrawTiles(7)

		if len(drawn) != 7 {
			t.Errorf("DrawTiles(7) returned %d tiles, expected 7", len(drawn))
		}

		if bag.RemainingCount() != initialCount-7 {
			t.Errorf("Remaining count should be %d, got %d", initialCount-7, bag.RemainingCount())
		}
	})

	t.Run("Draw zero tiles", func(t *testing.T) {
		bag := NewTileBag()
		initialCount := bag.RemainingCount()

		drawn := bag.DrawTiles(0)

		if len(drawn) != 0 {
			t.Errorf("DrawTiles(0) should return empty slice, got %d tiles", len(drawn))
		}

		if bag.RemainingCount() != initialCount {
			t.Errorf("Remaining count should not change when drawing 0 tiles")
		}
	})

	t.Run("Draw negative amount", func(t *testing.T) {
		bag := NewTileBag()
		initialCount := bag.RemainingCount()

		drawn := bag.DrawTiles(-5)

		if len(drawn) != 0 {
			t.Errorf("DrawTiles(-5) should return empty slice, got %d tiles", len(drawn))
		}

		if bag.RemainingCount() != initialCount {
			t.Errorf("Remaining count should not change when drawing negative tiles")
		}
	})

	t.Run("Draw more than available", func(t *testing.T) {
		bag := NewTileBag()

		// Draw most tiles first
		bag.DrawTiles(95)
		remaining := bag.RemainingCount()

		// Try to draw more than remaining
		drawn := bag.DrawTiles(10)

		if len(drawn) != remaining {
			t.Errorf("DrawTiles(10) with only %d tiles should return %d tiles, got %d",
				remaining, remaining, len(drawn))
		}

		if bag.RemainingCount() != 0 {
			t.Errorf("Bag should be empty after drawing all remaining tiles")
		}
	})

	t.Run("Draw from empty bag", func(t *testing.T) {
		bag := NewTileBag()

		// Empty the bag
		bag.DrawTiles(100)

		// Try to draw from empty bag
		drawn := bag.DrawTiles(5)

		if len(drawn) != 0 {
			t.Errorf("DrawTiles from empty bag should return empty slice, got %d tiles", len(drawn))
		}

		if !bag.IsEmpty() {
			t.Errorf("Bag should report as empty")
		}
	})
}

// TestReturnTiles tests the ReturnTiles functionality
func TestReturnTiles(t *testing.T) {
	t.Run("Return drawn tiles", func(t *testing.T) {
		bag := NewTileBag()
		initialCount := bag.RemainingCount()

		// Draw some tiles
		drawn := bag.DrawTiles(7)

		// Return them
		bag.ReturnTiles(drawn)

		if bag.RemainingCount() != initialCount {
			t.Errorf("After returning tiles, count should be %d, got %d",
				initialCount, bag.RemainingCount())
		}
	})

	t.Run("Return empty slice", func(t *testing.T) {
		bag := NewTileBag()
		initialCount := bag.RemainingCount()

		// Return empty slice
		bag.ReturnTiles([]Tile{})

		if bag.RemainingCount() != initialCount {
			t.Errorf("Returning empty slice should not change count")
		}
	})

	t.Run("Return creates new tiles", func(t *testing.T) {
		bag := NewTileBag()

		// Create custom tiles to return
		customTiles := []Tile{
			{Letter: 'X', Points: 8, IsBlank: false},
			{Letter: 'Z', Points: 10, IsBlank: false},
		}

		initialCount := bag.RemainingCount()
		bag.ReturnTiles(customTiles)

		if bag.RemainingCount() != initialCount+len(customTiles) {
			t.Errorf("Returning custom tiles should increase count by %d", len(customTiles))
		}
	})
}

// TestRemainingCount tests the RemainingCount method accuracy
func TestRemainingCount(t *testing.T) {
	bag := NewTileBag()

	// Test initial count
	if count := bag.RemainingCount(); count != 100 {
		t.Errorf("Initial RemainingCount() = %d, want 100", count)
	}

	// Test after drawing tiles
	for i := 0; i < 10; i++ {
		drawn := bag.DrawTiles(i + 1)
		expectedRemaining := 100 - (i+1)*(i+2)/2 // Sum of 1+2+...+(i+1)

		if actualRemaining := bag.RemainingCount(); actualRemaining != expectedRemaining {
			t.Errorf("After drawing %d batches, RemainingCount() = %d, want %d",
				i+1, actualRemaining, expectedRemaining)
		}

		if len(drawn) != min(i+1, expectedRemaining+len(drawn)) {
			t.Errorf("DrawTiles(%d) returned %d tiles", i+1, len(drawn))
		}
	}
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestIsEmpty tests the IsEmpty method
func TestIsEmpty(t *testing.T) {
	bag := NewTileBag()

	// Bag should not be empty initially
	if bag.IsEmpty() {
		t.Errorf("New bag should not be empty")
	}

	// Empty the bag
	bag.DrawTiles(100)

	// Now it should be empty
	if !bag.IsEmpty() {
		t.Errorf("Bag should be empty after drawing all tiles")
	}

	// Add tiles back
	bag.ReturnTiles([]Tile{{Letter: 'A', Points: 1, IsBlank: false}})

	// Should not be empty anymore
	if bag.IsEmpty() {
		t.Errorf("Bag should not be empty after returning tiles")
	}
}

// TestGetTileQuantity tests the GetTileQuantity function
func TestGetTileQuantity(t *testing.T) {
	// Test known quantities
	if GetTileQuantity('A') != 9 {
		t.Errorf("GetTileQuantity('A') = %d, want 9", GetTileQuantity('A'))
	}

	if GetTileQuantity('E') != 12 {
		t.Errorf("GetTileQuantity('E') = %d, want 12", GetTileQuantity('E'))
	}

	if GetTileQuantity('Q') != 1 {
		t.Errorf("GetTileQuantity('Q') = %d, want 1", GetTileQuantity('Q'))
	}

	// Test blank tiles
	if GetTileQuantity(0) != 2 {
		t.Errorf("GetTileQuantity(0) = %d, want 2", GetTileQuantity(0))
	}

	// Test invalid letter
	if GetTileQuantity('1') != 0 {
		t.Errorf("GetTileQuantity('1') = %d, want 0", GetTileQuantity('1'))
	}
}

// TestGetAllTileInfo tests the GetAllTileInfo function
func TestGetAllTileInfo(t *testing.T) {
	info := GetAllTileInfo()

	// Check that we have info for all letters plus blank
	expectedLetters := 26 + 1 // 26 letters + blank (rune 0)
	if len(info) != expectedLetters {
		t.Errorf("GetAllTileInfo() returned info for %d letters, want %d", len(info), expectedLetters)
	}

	// Check specific values
	if info['A'].Quantity != 9 || info['A'].Points != 1 {
		t.Errorf("Info for 'A': quantity=%d points=%d, want quantity=9 points=1",
			info['A'].Quantity, info['A'].Points)
	}

	if info[0].Quantity != 2 || info[0].Points != 0 {
		t.Errorf("Info for blank: quantity=%d points=%d, want quantity=2 points=0",
			info[0].Quantity, info[0].Points)
	}
}

// TestConcurrentTileBag tests thread-safety of the TileBag
func TestConcurrentTileBag(t *testing.T) {
	bag := NewTileBag()

	const numGoroutines = 10
	const tilesPerGoroutine = 5

	var wg sync.WaitGroup
	drawnTiles := make([][]Tile, numGoroutines)

	// Launch multiple goroutines to draw tiles concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			drawnTiles[index] = bag.DrawTiles(tilesPerGoroutine)
		}(i)
	}

	wg.Wait()

	// Verify that tiles were drawn correctly
	totalDrawn := 0
	allDrawnTiles := make([]Tile, 0)

	for _, tiles := range drawnTiles {
		totalDrawn += len(tiles)
		allDrawnTiles = append(allDrawnTiles, tiles...)
	}

	// Should have drawn exactly the requested amount
	expectedTotal := numGoroutines * tilesPerGoroutine
	if totalDrawn != expectedTotal {
		t.Errorf("Concurrent drawing: got %d tiles, want %d", totalDrawn, expectedTotal)
	}

	// Remaining count should be correct
	expectedRemaining := 100 - expectedTotal
	if remaining := bag.RemainingCount(); remaining != expectedRemaining {
		t.Errorf("After concurrent drawing: %d tiles remaining, want %d", remaining, expectedRemaining)
	}

	// Test concurrent return
	wg = sync.WaitGroup{}
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			bag.ReturnTiles(drawnTiles[index])
		}(i)
	}

	wg.Wait()

	// After returning all tiles, count should be back to 100
	if bag.RemainingCount() != 100 {
		t.Errorf("After concurrent returning: %d tiles, want 100", bag.RemainingCount())
	}
}

// TestConcurrentRemainingCount tests that RemainingCount is thread-safe
func TestConcurrentRemainingCount(t *testing.T) {
	bag := NewTileBag()

	const numGoroutines = 20
	counts := make([]int, numGoroutines)

	var wg sync.WaitGroup

	// Multiple goroutines calling RemainingCount concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			counts[index] = bag.RemainingCount()
		}(i)
	}

	wg.Wait()

	// All should return the same count (100) since no tiles are being drawn
	for i, count := range counts {
		if count != 100 {
			t.Errorf("Goroutine %d: RemainingCount() = %d, want 100", i, count)
		}
	}
}

// TestConcurrentMixedOperations tests mixed concurrent operations
func TestConcurrentMixedOperations(t *testing.T) {
	bag := NewTileBag()

	var wg sync.WaitGroup

	// Mix of different operations
	operations := []func(){
		func() { bag.DrawTiles(1) },
		func() { bag.DrawTiles(2) },
		func() { bag.RemainingCount() },
		func() { bag.IsEmpty() },
		func() { bag.ReturnTiles([]Tile{{Letter: 'X', Points: 8}}) },
	}

	// Run operations concurrently
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			op := operations[index%len(operations)]
			op()
		}(i)
	}

	wg.Wait()

	// Bag should still be in a valid state
	remaining := bag.RemainingCount()
	if remaining < 0 {
		t.Errorf("Invalid remaining count after concurrent operations: %d", remaining)
	}
}
