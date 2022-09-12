package alphabet

import (
	"unicode/utf8"

	"github.com/rs/zerolog/log"
)

// Rack is a machine-friendly representation of a user's rack.
type Rack struct {
	// letArr is an array of letter codes from 0 to MaxAlphabetSize.
	// The blank can go at the MaxAlphabetSize place.
	LetArr     []int
	numLetters uint8
	alphabet   *Alphabet
	repr       string
	// letterIdxs []uint8
}

// NewRack creates a brand new rack structure with an alphabet.
func NewRack(alph *Alphabet) *Rack {
	return &Rack{alphabet: alph, LetArr: make([]int, MaxAlphabetSize+1)}
}

// Hashable returns a hashable representation of this rack, that is
// not necessarily user-friendly.
func (r *Rack) Hashable() string {
	return r.TilesOn().String()
}

// String returns a user-visible version of this rack.
func (r *Rack) String() string {
	return r.TilesOn().UserVisible(r.alphabet)
}

// Copy returns a deep copy of this rack
func (r *Rack) Copy() *Rack {
	n := &Rack{
		numLetters: r.numLetters,
		alphabet:   r.alphabet,
		repr:       r.repr,
	}
	n.LetArr = make([]int, len(r.LetArr))
	copy(n.LetArr, r.LetArr)
	return n
}

func (r *Rack) CopyFrom(other *Rack) {
	r.numLetters = other.numLetters
	r.alphabet = other.alphabet
	r.repr = other.repr
	// These will always be the same size: MaxAlphabetSize + 1
	if r.LetArr == nil {
		r.LetArr = make([]int, MaxAlphabetSize+1)
	}
	copy(r.LetArr, other.LetArr)
}

// RackFromString creates a Rack from a string and an alphabet
func RackFromString(rack string, a *Alphabet) *Rack {
	r := &Rack{}
	r.alphabet = a
	r.setFromStr(rack)
	return r
}

func (r *Rack) setFromStr(rack string) {
	if r.LetArr == nil {
		r.LetArr = make([]int, MaxAlphabetSize+1)
	} else {
		r.Clear()
	}
	for _, c := range rack {
		ml, err := r.alphabet.Val(c)
		if err == nil {
			r.LetArr[ml]++
		} else {
			log.Error().Msgf("Rack has an illegal character: %v", string(c))
		}
	}
	r.numLetters = uint8(utf8.RuneCountInString(rack))
}

// Set sets the rack from a list of machine letters
func (r *Rack) Set(mls []MachineLetter) {
	r.Clear()
	for _, ml := range mls {
		r.LetArr[ml]++
	}
	r.numLetters = uint8(len(mls))
}

func (r *Rack) Clear() {
	// Clear the rack
	for i := 0; i < MaxAlphabetSize+1; i++ {
		r.LetArr[i] = 0
	}
	r.numLetters = 0
}

func (r *Rack) Take(letter MachineLetter) {
	// this function should only be called if there is a letter on the rack
	// it doesn't check if it's there!
	r.LetArr[letter]--
	r.numLetters--
}

func (r *Rack) Has(letter MachineLetter) bool {
	return r.LetArr[letter] > 0
}

func (r *Rack) CountOf(letter MachineLetter) int {
	return r.LetArr[letter]
}

func (r *Rack) Add(letter MachineLetter) {
	r.LetArr[letter]++
	r.numLetters++
}

// TilesOn returns the MachineLetters of the rack's current tiles. It is alphabetized.
func (r *Rack) TilesOn() MachineWord {
	if r.numLetters == 0 {
		return MachineWord([]MachineLetter{})
	}
	letters := make([]MachineLetter, r.numLetters)
	r.NoAllocTilesOn(letters)

	return MachineWord(letters)
}

// NoAllocTilesOn places the tiles in the passed-in slice, and returns the number
// of letters
func (r *Rack) NoAllocTilesOn(letters []MachineLetter) int {
	ct := 0
	numPossibleLetters := r.alphabet.NumLetters()
	var i MachineLetter
	for i = 0; i < MachineLetter(numPossibleLetters); i++ {
		if r.LetArr[i] > 0 {
			for j := 0; j < r.LetArr[i]; j++ {
				letters[ct] = i
				ct++
			}
		}
	}
	if r.LetArr[BlankMachineLetter] > 0 {
		for j := 0; j < r.LetArr[BlankMachineLetter]; j++ {
			letters[ct] = BlankMachineLetter
			ct++
		}
	}
	return ct
}

// ScoreOn returns the total score of the tiles on this rack.
func (r *Rack) ScoreOn(ld *LetterDistribution) int {
	score := 0
	var i MachineLetter
	numPossibleLetters := r.alphabet.NumLetters()
	for i = 0; i < MachineLetter(numPossibleLetters); i++ {
		if r.LetArr[i] > 0 {
			score += ld.Score(i) * r.LetArr[i]
		}
	}
	return score
}

// NumTiles returns the current number of tiles on this rack.
func (r *Rack) NumTiles() uint8 {
	return r.numLetters
}

func (r *Rack) Empty() bool {
	return r.numLetters == 0
}

func (r *Rack) Alphabet() *Alphabet {
	return r.alphabet
}
