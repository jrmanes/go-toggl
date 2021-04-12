package kit

// split return divided strings
func Split(tosplit string, sep rune) []string {
	var fields []string

	last := 0
	for i, c := range tosplit {
		if c == sep {
			// Found the separator, append a slice
			fields = append(fields, string(tosplit[last:i]))
			last = i + 1
		}
	}

	// last field
	fields = append(fields, string(tosplit[last:]))

	return fields
}
