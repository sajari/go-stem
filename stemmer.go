package stemmer

import "bytes"

func Vowel(body []byte, offset int) bool {
	switch body[offset] {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	case 'y':
		if offset == 0 {
			return false
		}
		return offset > 0 && !Vowel(body, offset-1)
	}
	return false
}

func Consonant(body []byte, offset int) bool {
	return !Vowel(body, offset)
}

const (
	vowel_state = iota
	consonant_state
)

func Measure(body []byte) int {
	if len(body) == 0 {
		return 0
	}

	n := 0
	v := Vowel(body, 0)
	for i := 0; i < len(body); i++ {
		c := !Vowel(body, i)
		if !c && !v {
			v = true
		} else if c && v {
			v = false
			n++
		}
	}
	return n
}

func MeasureGT(body []byte, min int) bool {
	if len(body) == 0 {
		return 0 > min
	}

	n := 0
	v := Vowel(body, 0)
	for i := 0; i < len(body); i++ {
		c := !Vowel(body, i)
		if !c && !v {
			v = true
		} else if c && v {
			v = false
			n++
			if n > min {
				return true
			}
		}
	}
	return n > min
}

func hasVowel(body []byte) bool {
	for i := 0; i < len(body); i++ {
		if Vowel(body, i) {
			return true
		}
	}
	return false
}

var (
	suffixSSES = []byte("sses")
	suffixIES  = []byte("ies")
	suffixSS   = []byte("ss")
	suffixS    = []byte("s")
)

func one_a(body []byte) []byte {
	if bytes.HasSuffix(body, suffixSSES) || bytes.HasSuffix(body, suffixIES) {
		return body[:len(body)-2]
	} else if bytes.HasSuffix(body, suffixSS) {
		return body
	} else if bytes.HasSuffix(body, suffixS) {
		return body[:len(body)-1]
	}
	return body
}

func star_o(body []byte) bool {
	size := len(body) - 1
	if size >= 2 && Consonant(body, size-2) && Vowel(body, size-1) && Consonant(body, size) {
		return body[size] != 'w' && body[size] != 'x' && body[size] != 'y'
	}
	return false
}

var (
	suffixAT = []byte("at")
	suffixBL = []byte("bl")
	suffixIZ = []byte("iz")
)

func one_b_a(body []byte) []byte {
	size := len(body)
	if bytes.HasSuffix(body, suffixAT) {
		return append(body, 'e')
	} else if bytes.HasSuffix(body, suffixBL) {
		return append(body, 'e')
	} else if bytes.HasSuffix(body, suffixIZ) {
		return append(body, 'e')
	} else if Consonant(body, size-1) && Consonant(body, size-2) && body[size-1] == body[size-2] {
		if body[size-1] != 'l' && body[size-1] != 's' && body[size-1] != 'z' {
			return body[:size-1]
		}
	} else if star_o(body) && Measure(body) == 1 {
		return append(body, 'e')
	}
	return body
}

var (
	suffixEED = []byte("eed")
	suffixED  = []byte("ed")
	suffixING = []byte("ing")
)

func one_b(body []byte) []byte {
	if bytes.HasSuffix(body, suffixEED) {
		if MeasureGT(body[:len(body)-3], 0) {
			return body[:len(body)-1]
		}
	} else if bytes.HasSuffix(body, suffixED) {
		if hasVowel(body[:len(body)-2]) {
			return one_b_a(body[:len(body)-2])
		}
	} else if bytes.HasSuffix(body, suffixING) {
		if hasVowel(body[:len(body)-3]) {
			return one_b_a(body[:len(body)-3])
		}
	}
	return body
}

var (
	suffixY = []byte("y")
)

func one_c(body []byte) []byte {
	if bytes.HasSuffix(body, suffixY) && hasVowel(body[:len(body)-1]) {
		body[len(body)-1] = 'i'
		return body
	}
	return body
}

var (
	suffixATIONAL = []byte("ational")
	suffixTIONAL  = []byte("tional")
	suffixENCI    = []byte("enci")
	suffixANCI    = []byte("anci")
	suffixIZER    = []byte("izer")
	suffixABLI    = []byte("abli")
	suffixBLI     = []byte("bli")
	suffixALLI    = []byte("alli")
	suffixENTLI   = []byte("entli")
	suffixELI     = []byte("eli")
	suffixOUSLI   = []byte("ousli")
	suffixIZATION = []byte("ization")
	suffixATION   = []byte("ation")
	suffixATOR    = []byte("ator")
	suffixALISM   = []byte("alism")
	suffixIVENESS = []byte("iveness")
	suffixFULNESS = []byte("fulness")
	suffixOUSNESS = []byte("ousness")
	suffixALITI   = []byte("aliti")
	suffixIVITI   = []byte("iviti")
	suffixBILITI  = []byte("biliti")
	suffixLOGI    = []byte("logi")
)

func two(body []byte) []byte {
	if bytes.HasSuffix(body, suffixATIONAL) {
		if MeasureGT(body[:len(body)-7], 0) {
			return append(body[:len(body)-7], "ate"...)
		}
	} else if bytes.HasSuffix(body, suffixTIONAL) {
		if MeasureGT(body[:len(body)-6], 0) {
			return body[:len(body)-2]
		}
	} else if bytes.HasSuffix(body, suffixENCI) || bytes.HasSuffix(body, suffixANCI) {
		if MeasureGT(body[:len(body)-4], 0) {
			return append(body[:len(body)-1], 'e')
		}
	} else if bytes.HasSuffix(body, suffixIZER) {
		if MeasureGT(body[:len(body)-4], 0) {
			return append(body[:len(body)-4], "ize"...)
		}
	} else if bytes.HasSuffix(body, suffixABLI) {
		if MeasureGT(body[:len(body)-4], 0) {
			return append(body[:len(body)-4], "able"...)
		}
		// To match the published algorithm, delete the following phrase
	} else if bytes.HasSuffix(body, suffixBLI) {
		if MeasureGT(body[:len(body)-3], 0) {
			return append(body[:len(body)-1], 'e')
		}
	} else if bytes.HasSuffix(body, suffixALLI) {
		if MeasureGT(body[:len(body)-4], 0) {
			return append(body[:len(body)-4], "al"...)
		}
	} else if bytes.HasSuffix(body, suffixENTLI) {
		if MeasureGT(body[:len(body)-5], 0) {
			return append(body[:len(body)-5], "ent"...)
		}
	} else if bytes.HasSuffix(body, suffixELI) {
		if MeasureGT(body[:len(body)-3], 0) {
			return append(body[:len(body)-3], "e"...)
		}
	} else if bytes.HasSuffix(body, suffixOUSLI) {
		if MeasureGT(body[:len(body)-5], 0) {
			return append(body[:len(body)-5], "ous"...)
		}
	} else if bytes.HasSuffix(body, suffixIZATION) {
		if MeasureGT(body[:len(body)-7], 0) {
			return append(body[:len(body)-7], "ize"...)
		}
	} else if bytes.HasSuffix(body, suffixATION) {
		if MeasureGT(body[:len(body)-5], 0) {
			return append(body[:len(body)-5], "ate"...)
		}
	} else if bytes.HasSuffix(body, suffixATOR) {
		if MeasureGT(body[:len(body)-4], 0) {
			return append(body[:len(body)-4], "ate"...)
		}
	} else if bytes.HasSuffix(body, suffixALISM) {
		if MeasureGT(body[:len(body)-5], 0) {
			return append(body[:len(body)-5], "al"...)
		}
	} else if bytes.HasSuffix(body, suffixIVENESS) {
		if MeasureGT(body[:len(body)-7], 0) {
			return append(body[:len(body)-7], "ive"...)
		}
	} else if bytes.HasSuffix(body, suffixFULNESS) {
		if MeasureGT(body[:len(body)-7], 0) {
			return append(body[:len(body)-7], "ful"...)
		}
	} else if bytes.HasSuffix(body, suffixOUSNESS) {
		if MeasureGT(body[:len(body)-7], 0) {
			return append(body[:len(body)-7], "ous"...)
		}
	} else if bytes.HasSuffix(body, suffixALITI) {
		if MeasureGT(body[:len(body)-5], 0) {
			return append(body[:len(body)-5], "al"...)
		}
	} else if bytes.HasSuffix(body, suffixIVITI) {
		if MeasureGT(body[:len(body)-5], 0) {
			return append(body[:len(body)-5], "ive"...)
		}
	} else if bytes.HasSuffix(body, suffixBILITI) {
		if MeasureGT(body[:len(body)-6], 0) {
			return append(body[:len(body)-6], "ble"...)
		}
		// To match the published algorithm, delete the following phrase
	} else if bytes.HasSuffix(body, suffixLOGI) {
		if MeasureGT(body[:len(body)-4], 0) {
			return body[:len(body)-1]
		}
	}
	return body
}

var (
	suffixICATE = []byte("icate")
	suffixATIVE = []byte("ative")
	suffixALIZE = []byte("alize")
	suffixICITI = []byte("iciti")
	suffixICAL  = []byte("ical")
	suffixFUL   = []byte("ful")
	suffixNESS  = []byte("ness")
)

func three(body []byte) []byte {
	if bytes.HasSuffix(body, suffixICATE) {
		if MeasureGT(body[:len(body)-5], 0) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixATIVE) {
		if MeasureGT(body[:len(body)-5], 0) {
			return body[:len(body)-5]
		}
	} else if bytes.HasSuffix(body, suffixALIZE) {
		if MeasureGT(body[:len(body)-5], 0) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixICITI) {
		if MeasureGT(body[:len(body)-5], 0) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixICAL) {
		if MeasureGT(body[:len(body)-4], 0) {
			return body[:len(body)-2]
		}
	} else if bytes.HasSuffix(body, suffixFUL) {
		if MeasureGT(body[:len(body)-3], 0) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixNESS) {
		if MeasureGT(body[:len(body)-4], 0) {
			return body[:len(body)-4]
		}
	}
	return body
}

var (
	suffixAL    = []byte("al")
	suffixANCE  = []byte("ance")
	suffixENCE  = []byte("ence")
	suffixER    = []byte("er")
	suffixIC    = []byte("ic")
	suffixABLE  = []byte("able")
	suffixIBLE  = []byte("ible")
	suffixANT   = []byte("ant")
	suffixEMENT = []byte("ement")
	suffixMENT  = []byte("ment")
	suffixENT   = []byte("ent")
	suffixION   = []byte("ion")
	suffixOU    = []byte("ou")
	suffixISM   = []byte("ism")
	suffixATE   = []byte("ate")
	suffixITI   = []byte("iti")
	suffixOUS   = []byte("ous")
	suffixIVE   = []byte("ive")
	suffixIZE   = []byte("ize")
)

func four(body []byte) []byte {
	if bytes.HasSuffix(body, suffixAL) {
		if MeasureGT(body[:len(body)-2], 1) {
			return body[:len(body)-2]
		}
	} else if bytes.HasSuffix(body, suffixANCE) {
		if MeasureGT(body[:len(body)-4], 1) {
			return body[:len(body)-4]
		}
	} else if bytes.HasSuffix(body, suffixENCE) {
		if MeasureGT(body[:len(body)-4], 1) {
			return body[:len(body)-4]
		}
	} else if bytes.HasSuffix(body, suffixER) {
		if MeasureGT(body[:len(body)-2], 1) {
			return body[:len(body)-2]
		}
	} else if bytes.HasSuffix(body, suffixIC) {
		if MeasureGT(body[:len(body)-2], 1) {
			return body[:len(body)-2]
		}
	} else if bytes.HasSuffix(body, suffixABLE) {
		if MeasureGT(body[:len(body)-4], 1) {
			return body[:len(body)-4]
		}
	} else if bytes.HasSuffix(body, suffixIBLE) {
		if MeasureGT(body[:len(body)-4], 1) {
			return body[:len(body)-4]
		}
	} else if bytes.HasSuffix(body, suffixANT) {
		if MeasureGT(body[:len(body)-3], 1) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixEMENT) {
		if MeasureGT(body[:len(body)-5], 1) {
			return body[:len(body)-5]
		}
	} else if bytes.HasSuffix(body, suffixMENT) {
		if MeasureGT(body[:len(body)-4], 1) {
			return body[:len(body)-4]
		}
	} else if bytes.HasSuffix(body, suffixENT) {
		if MeasureGT(body[:len(body)-3], 1) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixION) {
		if MeasureGT(body[:len(body)-3], 1) {
			if len(body) > 4 && (body[len(body)-4] == 's' || body[len(body)-4] == 't') {
				return body[:len(body)-3]
			}
		}
	} else if bytes.HasSuffix(body, suffixOU) {
		if MeasureGT(body[:len(body)-2], 1) {
			return body[:len(body)-2]
		}
	} else if bytes.HasSuffix(body, suffixISM) {
		if MeasureGT(body[:len(body)-3], 1) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixATE) {
		if MeasureGT(body[:len(body)-3], 1) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixITI) {
		if MeasureGT(body[:len(body)-3], 1) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixOUS) {
		if MeasureGT(body[:len(body)-3], 1) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixIVE) {
		if MeasureGT(body[:len(body)-3], 1) {
			return body[:len(body)-3]
		}
	} else if bytes.HasSuffix(body, suffixIZE) {
		if MeasureGT(body[:len(body)-3], 1) {
			return body[:len(body)-3]
		}
	}
	return body
}

var (
	suffixE = []byte("e")
)

func five_a(body []byte) []byte {
	if bytes.HasSuffix(body, suffixE) && MeasureGT(body[:len(body)-1], 1) {
		return body[:len(body)-1]
	} else if bytes.HasSuffix(body, suffixE) && Measure(body[:len(body)-1]) == 1 && !star_o(body[:len(body)-1]) {
		return body[:len(body)-1]
	}
	return body
}

func five_b(body []byte) []byte {
	size := len(body)
	if MeasureGT(body, 1) && Consonant(body, size-1) && Consonant(body, size-2) && body[size-1] == body[size-2] && body[size-1] == 'l' {
		return body[:len(body)-1]
	}
	return body
}

// Stem computes the stemmed version of a single word, assumed to be lowercase and not
// contain trailing/leading whitespace.
func Stem(word []byte) []byte {
	if len(word) > 2 {
		return five_b(five_a(four(three(two(one_c(one_b(one_a(word))))))))
	}
	return word
}
