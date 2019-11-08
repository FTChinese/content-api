package models

type StringPair struct {
	First  string
	Second string
}

// ZipString match the elements of two arrays
// on the same index into a StringPair
func ZipString(a, b []string) []StringPair {
	var pairs []StringPair

	var aIndex, bIndex int
	var aLen = len(a)
	var bLen = len(b)

	for aIndex < aLen && bIndex < bLen {
		p := StringPair{
			First:  a[aIndex],
			Second: b[bIndex],
		}

		pairs = append(pairs, p)
		aIndex++
		bIndex++
	}

	for aIndex < aLen {
		p := StringPair{
			First:  a[aIndex],
			Second: "",
		}
		pairs = append(pairs, p)
		aIndex++
	}

	for bIndex < bLen {
		p := StringPair{
			First:  "",
			Second: b[bIndex],
		}
		pairs = append(pairs, p)
	}

	return pairs
}
