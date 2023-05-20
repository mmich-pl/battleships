package board_validation

import "testing"

func TestValidateShipPlacement(t *testing.T) {
	testCases := []struct {
		name     string
		ships    []string
		expected bool
	}{
		{
			name: "CommonEdge",
			ships: []string{
				"G8", "H8", "I8", "J8", // Battleship
				"A1", "B1", "C1", // Cruiser
				"A2", "B2", "C2", // Cruiser, common edge
				"D4", "C4", // Destroyer
				"G2", "H2", // Destroyer
				"A8", "B8", // Destroyer
				"C6",  // Submarine
				"E10", // Submarine
				"B10", // Submarine
				"G5",  // Submarine
			},
			expected: false,
		},
		{
			name: "MissingShips",
			ships: []string{
				"A1", "A2", "A3",
				"B1", "B2", "B3", // Cruiser
				"D1", "D2", // Destroyer
				"F1", "F2", // Destroyer
				"H1", // Submarine
				"H3", // Submarine
				"H5", // Submarine
				"H7", // Submarine
			},
			expected: false,
		},
		{
			name: "InvalidColumn",
			ships: []string{
				"I3", "J3", "K3", "L3", // Battleship, invalid column
				"H6", "H7", "H8", // Cruiser
				"F1", "G1", "H1", // Cruiser
				"E5", "F5", // Destroyer
				"A10", "B10", // Destroyer
				"D7", "E7", // Destroyer
				"I4",  // Submarine
				"G10", // Submarine
				"F10", // Submarine
				"C8",  // Submarine
			},
			expected: false,
		},
		{
			name: "InvalidRow",
			ships: []string{
				"A3", "B3", "C3", "D3", // Battleship
				"H6", "H7", "H8", // Cruiser
				"F1", "G1", "H1", // Cruiser
				"E5", "F5", // Destroyer
				"A10", "B10", // Destroyer
				"D7", "E7", // Destroyer
				"I4",  // Submarine
				"G10", // Submarine
				"F11", // Submarine, invalid row
				"C8",  // Submarine
			},
			expected: false,
		},
		{
			name: "DuplicateShipPlacement",
			ships: []string{
				"G8", "H8", "I8", "J8", // Battleship
				"A1", "B1", "C1", // Cruiser
				"J3", "J4", "J5", // Cruiser
				"D4", "C4", // Destroyer
				"C8", "D8", // Destroyer
				"I8", "J8", // Destroyer, duplicate ship placement
				"C6",  // Submarine
				"E10", // Submarine
				"B10", // Submarine
				"G5",  // Submarine
			},
			expected: false,
		},
		{
			name: "ValidPlacement",
			ships: []string{
				"G8", "H8", "I8", "J8", // Battleship
				"A1", "B1", "C1", // Cruiser
				"J3", "J4", "J5", // Cruiser
				"D4", "C4", // Destroyer
				"G2", "H2", // Destroyer
				"A8", "B8", // Destroyer
				"C6",  // Submarine
				"E10", // Submarine
				"B10", // Submarine
				"G5",  // Submarine
			},
			expected: true,
		},
		{
			name: "ValidPlacement",
			ships: []string{
				"A3", "B3", "C3", "D3", // Battleship
				"H6", "H7", "H8", // Cruiser
				"F1", "G1", "H1", // Cruiser
				"E5", "F5", // Destroyer
				"A10", "B10", // Destroyer
				"D7", "D6", // Destroyer
				"I4",  // Submarine
				"G10", // Submarine
				"I10", // Submarine
				"C8",  // Submarine
			},
			expected: true,
		},
		{
			name: "ValidPlacement",
			ships: []string{
				"G2", "H2", "I2", "J2", // Battleship
				"G6", "G7", "G8", // Cruiser
				"D1", "E1", "F1", // Cruiser
				"F3", "F4", // Destroyer
				"H5", "I5", // Destroyer
				"D10", "E10", // Destroyer
				"B2", // Submarine
				"C5", // Submarine
				"D7", // Submarine
				"J9", // Submarine
			},
			expected: true,
		},
		{
			name: "L-shaped battleship",
			ships: []string{
				"B7", "C5", "C6", "C7", // Battleship
				"J2", "J3", "J4", // Cruiser
				"E1", "F1", "G1", // Cruiser
				"A3", "A4", // Destroyer
				"E6", "E7", // Destroyer
				"H9", "I9", // Destroyer
				"H4",  // Submarine
				"I6",  // Submarine
				"F8",  // Submarine
				"B10", // Submarine
			},
			expected: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, message := validateShipPlacement(testCase.ships)
			if result != testCase.expected {
				t.Errorf("Ship placement validation failed for %v. Expected: %v, Got: %v", testCase.ships, testCase.expected, result)
			}
			println(message)
		})
	}
}
