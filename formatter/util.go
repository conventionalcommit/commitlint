package formatter

func truncate(maxSize int, input string) string {
	if len(input) < maxSize {
		return input
	}
	return input[:maxSize-3] + "..."
}
