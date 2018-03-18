package stemmer

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestConsonant(t *testing.T) {
	tests := []struct {
		word      string
		consonant []bool // value for each index of word
	}{
		{
			"toy",
			[]bool{true, false, true},
		},
		{
			"syzygy",
			[]bool{true, false, true, false, true, false},
		},
		{
			"yoke",
			[]bool{true, false, true, false},
		},
	}

	for _, test := range tests {
		t.Run(test.word, func(t *testing.T) {
			if len(test.word) != len(test.consonant) {
				t.Errorf("len(word) != len(consonant)")
			}
			for i, expected := range test.consonant {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					if got := Consonant([]byte(test.word), i); got != expected {
						t.Errorf(" = %v, expected %v", got, expected)
					}
				})
			}
		})
	}
}

func TestMeasure(t *testing.T) {
	tests := []struct {
		in  string
		out int
	}{
		{"tr", 0},
		{"ee", 0},
		{"tree", 0},
		{"y", 0},
		{"by", 0},
		{"trouble", 1},
		{"oats", 1},
		{"trees", 1},
		{"ivy", 1},
		{"troubles", 2},
		{"private", 2},
		{"oaten", 2},
		{"orrery", 2},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			if got := Measure([]byte(test.in)); got != test.out {
				t.Errorf(" = %d, expected %d", got, test.out)
			}
		})
	}
}

func TestAll(t *testing.T) {
	type pair struct{ in, out string }
	tests := []struct {
		name   string
		fn     func([]byte) []byte
		params []pair
	}{
		{
			name: "one_a",
			fn:   one_a,
			params: []pair{
				{"caresses", "caress"},
				{"ponies", "poni"},
				{"ties", "ti"},
				{"caress", "caress"},
				{"cats", "cat"},
			},
		},
		{
			name: "one_b",
			fn:   one_b,
			params: []pair{
				{"feed", "feed"},
				{"agreed", "agree"},
				{"plastered", "plaster"},
				{"bled", "bled"},
				{"motoring", "motor"},
				{"sing", "sing"},
				{"motoring", "motor"},
				{"conflated", "conflate"},
				{"troubled", "trouble"},
				{"sized", "size"},
				{"hopping", "hop"},
				{"tanned", "tan"},
				{"failing", "fail"},
				{"filing", "file"},
			},
		},
		{
			name: "one_c",
			fn:   one_c,
			params: []pair{
				{"sky", "sky"},
				{"happy", "happi"},
			},
		},
		{
			name: "two",
			fn:   two,
			params: []pair{
				{"relational", "relate"},
				{"conditional", "condition"},
				{"rational", "rational"},
				{"valenci", "valence"},
				{"hesitanci", "hesitance"},
				{"digitizer", "digitize"},
				{"conformabli", "conformable"},
				{"radicalli", "radical"},
				{"differentli", "different"},
				{"vileli", "vile"},
				{"analogousli", "analogous"},
				{"vietnamization", "vietnamize"},
				{"predication", "predicate"},
				{"operator", "operate"},
				{"feudalism", "feudal"},
				{"decisiveness", "decisive"},
				{"hopefulness", "hopeful"},
				{"callousness", "callous"},
				{"formaliti", "formal"},
				{"sensitiviti", "sensitive"},
				{"sensibiliti", "sensible"},
			},
		},
		{
			name: "three",
			fn:   three,
			params: []pair{
				{"triplicate", "triplic"},
				{"formative", "form"},
				{"formalize", "formal"},
				{"electriciti", "electric"},
				{"electrical", "electric"},
				{"hopeful", "hope"},
				{"goodness", "good"},
			},
		},
		{
			name: "four",
			fn:   four,
			params: []pair{
				{"revival", "reviv"},
				{"allowance", "allow"},
				{"inference", "infer"},
				{"airliner", "airlin"},
				{"gyroscopic", "gyroscop"},
				{"adjustable", "adjust"},
				{"defensible", "defens"},
				{"irritant", "irrit"},
				{"replacement", "replac"},
				{"adjustment", "adjust"},
				{"dependent", "depend"},
				{"adoption", "adopt"},
				{"homologou", "homolog"},
				{"communism", "commun"},
				{"activate", "activ"},
				{"angulariti", "angular"},
				{"homologous", "homolog"},
				{"effective", "effect"},
				{"bowdlerize", "bowdler"},
			},
		},
		{
			name: "five_a",
			fn:   five_a,
			params: []pair{
				{"probate", "probat"},
				{"rate", "rate"},
				{"cease", "ceas"},
			},
		},
		{
			name: "five_b",
			fn:   five_b,
			params: []pair{
				{"controll", "control"},
				{"roll", "roll"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, param := range test.params {
				t.Run(param.in, func(t *testing.T) {
					if got := string(test.fn([]byte(param.in))); got != param.out {
						t.Errorf(" = %q, expected %q", param.in, param.out)
					}
				})
			}
		})
	}
}

func TestVocal(t *testing.T) {
	in := readLines(t, "in.txt")
	out := readLines(t, "out.txt")

	for i, x := range in {
		t.Run(string(x), func(t *testing.T) {
			if got := Stem(x); !bytes.Equal(got, out[i]) {
				t.Errorf(" = %q, expected %q", got, out[i])
			}
		})
	}
}

type fataler interface {
	Fatalf(string, ...interface{})
}

func readLines(f fataler, path string) [][]byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		f.Fatalf("unexpected error opening file: %v", err)
	}
	return bytes.Split(b, []byte("\n"))
}

func BenchmarkVocal(b *testing.B) {
	in := readLines(b, "in.txt")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < len(in); j++ {
			_ = Stem(in[j])
		}
	}
}
