package timex

// isBetween returns true if the given value falls between the given
// ceiling and floor
func isBetween(value, floor, ceil int, inclusive bool) bool {
	if inclusive {
		if value >= floor && value <= ceil {
			return true
		}
	}

	if value > floor && value < ceil {
		return true
	}

	return false
}
