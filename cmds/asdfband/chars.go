package main

type (
	// CellPos is just a single, square-ish block rendered on screen.
	CellPos struct {
		X, Y int
	}

	// CellChar is a latin character made out of cells in a 4x6 grid.
	CellChar struct {
		Cells []CellPos
	}
)

// AvailableChars returns a list of all characters with a known render pattern.
func AvailableChars() []byte {
	// note (bs): I suspect this should be implicitly derivable from the set of
	// known chars by having better, more intelligent underlying data structures.
	chars := []byte{}
	for c := byte('A'); c <= 'Z'; c++ {
		chars = append(chars, c)
	}
	for c := byte('0'); c <= '9'; c++ {
		chars = append(chars, c)
	}
	return chars
}

// GetCellChar will return the set of cells needed to render the given
// character, if it's available. Otherwise, returns an empty set of cells.
func GetCellChar(c byte) CellChar {
	switch c {
	case 'A':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{0, 2}, {1, 2}, {2, 2}, {3, 2},
				{0, 3}, {3, 3},
				{0, 4}, {3, 4},
				{0, 5}, {3, 5},
			},
		}

	case 'B':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {2, 0},
				{0, 1}, {1, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {1, 3}, {2, 3},
				{0, 4}, {3, 4},
				{0, 5}, {1, 5}, {2, 5},
			},
		}

	case 'C':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{0, 2},
				{0, 3},
				{0, 4}, {3, 4},
				{1, 5}, {2, 5},
			},
		}

	case 'D':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {2, 0},
				{0, 1}, {1, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {3, 3},
				{0, 4}, {3, 4},
				{0, 5}, {1, 5}, {2, 5},
			},
		}

	case 'E':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {1, 0}, {2, 0}, {3, 0},
				{0, 1},
				{0, 2},
				{0, 3}, {1, 3}, {2, 3},
				{0, 4},
				{0, 5}, {1, 5}, {2, 5}, {3, 5},
			},
		}

	case 'F':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {1, 0}, {2, 0}, {3, 0},
				{0, 1},
				{0, 2}, {1, 2}, {2, 2},
				{0, 3},
				{0, 4},
				{0, 5},
			},
		}

	case 'G':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{0, 2},
				{0, 3}, {2, 3}, {3, 3},
				{0, 4}, {3, 4},
				{1, 5}, {2, 5},
			},
		}

	case 'H':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {3, 0},
				{0, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {1, 3}, {2, 3}, {3, 3},
				{0, 4}, {3, 4},
				{0, 5}, {3, 5},
			},
		}

	case 'I':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {1, 0}, {2, 0},
				{1, 1},
				{1, 2},
				{1, 3},
				{1, 4},
				{0, 5}, {1, 5}, {2, 5},
			},
		}

	case 'J':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {1, 0}, {2, 0}, {3, 0},
				{2, 1},
				{2, 2},
				{2, 3},
				{0, 4}, {2, 4},
				{1, 5},
			},
		}

	case 'K':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {3, 0},
				{0, 1}, {3, 1},
				{0, 2}, {2, 2},
				{0, 3}, {1, 3},
				{0, 4}, {2, 4},
				{0, 5}, {3, 5},
			},
		}

	case 'L':
		return CellChar{
			Cells: []CellPos{
				{0, 0},
				{0, 1},
				{0, 2},
				{0, 3},
				{0, 4},
				{0, 5}, {1, 5}, {2, 5}, {3, 5},
			},
		}

	case 'M':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {3, 0},
				{0, 1}, {2, 1}, {3, 1},
				{0, 2}, {1, 2}, {2, 2}, {3, 2},
				{0, 3}, {3, 3},
				{0, 4}, {3, 4},
				{0, 5}, {3, 5},
			},
		}

	case 'N':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {3, 0},
				{0, 1}, {3, 1},
				{0, 2}, {1, 2}, {3, 2},
				{0, 3}, {2, 3}, {3, 3},
				{0, 4}, {3, 4},
				{0, 5}, {3, 5},
			},
		}

	case 'O':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {3, 3},
				{0, 4}, {3, 4},
				{1, 5}, {2, 5},
			},
		}

	case 'P':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {2, 0},
				{0, 1}, {1, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {1, 3}, {2, 3},
				{0, 4},
				{0, 5},
			},
		}

	case 'Q':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {3, 3},
				{1, 4}, {2, 4},
				{3, 5},
			},
		}

	case 'R':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {2, 0},
				{0, 1}, {1, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {1, 3}, {2, 3},
				{0, 4}, {3, 4},
				{0, 5}, {3, 5},
			},
		}

	case 'S':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{1, 2},
				{2, 3},
				{0, 4}, {3, 4},
				{1, 5}, {2, 5},
			},
		}

	case 'T':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {1, 0}, {2, 0},
				{1, 1},
				{1, 2},
				{1, 3},
				{1, 4},
				{1, 5},
			},
		}

	case 'U':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {3, 0},
				{0, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {3, 3},
				{0, 4}, {3, 4},
				{1, 5}, {2, 5},
			},
		}

	case 'V':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {3, 0},
				{0, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {3, 3},
				{0, 4}, {2, 4},
				{1, 5},
			},
		}

	case 'W':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {3, 0},
				{0, 1}, {3, 1},
				{0, 2}, {3, 2},
				{0, 3}, {1, 3}, {3, 3},
				{0, 4}, {1, 4}, {2, 4}, {3, 4},
				{0, 5}, {3, 5},
			},
		}

	case 'X':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {3, 0},
				{0, 1}, {3, 1},
				{1, 2}, {2, 2},
				{1, 3}, {2, 3},
				{0, 4}, {3, 4},
				{0, 5}, {3, 5},
			},
		}

	case 'Y':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {3, 0},
				{0, 1}, {3, 1},
				{1, 2}, {2, 2},
				{1, 3},
				{1, 4},
				{1, 5},
			},
		}

	case 'Z':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {1, 0}, {2, 0}, {3, 0},
				{3, 1},
				{2, 2},
				{1, 3},
				{0, 4},
				{0, 5}, {1, 5}, {2, 5}, {3, 5},
			},
		}

	case '0':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{0, 2}, {2, 2}, {3, 2},
				{0, 3}, {1, 3}, {3, 3},
				{0, 4}, {3, 4},
				{1, 5}, {2, 5},
			},
		}

	case '1':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{2, 1},
				{2, 2},
				{2, 3},
				{2, 4},
				{0, 5}, {1, 5}, {2, 5}, {3, 5},
			},
		}

	case '2':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{3, 2},
				{2, 3},
				{1, 4},
				{0, 5}, {1, 5}, {2, 5}, {3, 5},
			},
		}

	case '3':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{2, 2},
				{1, 3}, {2, 3},
				{3, 4},
				{1, 5}, {2, 5},
			},
		}

	case '4':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {2, 0},
				{0, 1}, {2, 1},
				{0, 2}, {2, 2},
				{1, 3}, {2, 3}, {3, 3},
				{2, 4},
				{2, 5},
			},
		}

	case '5':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {1, 0}, {2, 0}, {3, 0},
				{0, 1},
				{0, 2}, {1, 2}, {2, 2},
				{3, 3},
				{3, 4},
				{0, 5}, {1, 5}, {2, 5},
			},
		}

	case '6':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{0, 2},
				{0, 3}, {1, 3}, {2, 3},
				{0, 4}, {3, 4},
				{1, 5}, {2, 5},
			},
		}

	case '7':
		return CellChar{
			Cells: []CellPos{
				{0, 0}, {1, 0}, {2, 0}, {3, 0},
				{3, 1},
				{2, 2},
				{1, 3},
				{0, 4},
				{0, 5},
			},
		}

	case '8':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{0, 2}, {3, 2},
				{1, 3}, {2, 3},
				{0, 4}, {3, 4},
				{1, 5}, {2, 5},
			},
		}

	case '9':
		return CellChar{
			Cells: []CellPos{
				{1, 0}, {2, 0},
				{0, 1}, {3, 1},
				{0, 2}, {3, 2},
				{1, 3}, {2, 3}, {3, 3},
				{3, 4},
				{2, 5},
			},
		}

	default:
		return CellChar{
			Cells: []CellPos{},
		}
	}
}
