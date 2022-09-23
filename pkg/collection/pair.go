package collection

type StringPair struct {
	First  string
	Second string
}

// ZipString match the elements of two arrays
// on the same index into a StringPair.
// Pitfall: DO remember to call index++.
func ZipString(a, b []string) []StringPair {
	var pairs []StringPair

	var aIndex, bIndex int
	var aLen = len(a)
	var bLen = len(b)

	for aIndex < aLen && bIndex < bLen {
		pairs = append(pairs, StringPair{
			First:  a[aIndex],
			Second: b[bIndex],
		})
		aIndex++
		bIndex++
	}

	// Handle dangling elements if a has more elements than b
	for aIndex < aLen {
		pairs = append(pairs, StringPair{
			First:  a[aIndex],
			Second: "",
		})
		// DON'T forget this!
		// It causes infinite loop if this block is run
		// and you forgot this.
		aIndex++
	}

	// Handle dangling elements if b has more elements than a.
	for bIndex < bLen {
		pairs = append(pairs, StringPair{
			First:  "",
			Second: b[bIndex],
		})

		// DON'T forget this!
		// It causes infinite loop if this block is run
		// and you forgot this.
		bIndex++
	}

	return pairs
}
