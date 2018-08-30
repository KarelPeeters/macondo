// blank_challenges has utilities for generating racks with blanks
// that have 1 or more solutions.
package anagrammer

import (
	"fmt"
	"log"

	"github.com/domino14/macondo/gaddag"
	"github.com/domino14/macondo/lexicon"
)

const BlankCharacter = '?'

// GenerateBlanks - Generate a list of blank word challenges given the
// parameters in args.
func GenerateBlanks(args *BlankChallengeArgs, dawg gaddag.SimpleDawg) (
	[]*Question, int) {
	var dist lexicon.LetterDistribution
	if args.Lexicon == "FISE09" {
		dist = lexicon.SpanishLetterDistribution()
	} else {
		dist = lexicon.EnglishLetterDistribution()
	}
	tries := 0
	// Handle 2-blank challenges at the end.
	// First gen 1-blank challenges.
	answerMap := make(map[string]bool)
	questions := make([]*Question, args.NumQuestions)
	qIndex := 0

	// try tries to generate challenges. It returns an error if it fails
	// to generate a challenge with too many or too few answers, or if
	// an answer has already been generated.
	try := func(nBlanks int) (*Question, error) {
		rack := genRack(dist, args.WordLength, nBlanks)
		tries++
		answers := Anagram(string(rack), dawg, ModeExact)
		if len(answers) == 0 || len(answers) > args.MaxSolutions {
			// Try again!
			return nil, fmt.Errorf("Too many or few answers! %v %v",
				len(answers), string(rack))
		}
		for _, answer := range answers {
			if answerMap[answer] {
				return nil, fmt.Errorf("Duplicate answer %v!", answer)
			}
		}
		for _, answer := range answers {
			answerMap[answer] = true
		}
		w := lexicon.Word{Word: string(rack), Dist: dist}
		return &Question{Q: w.MakeAlphagram(), A: answers}, nil
	}

	for qIndex < args.NumQuestions-args.Num2Blanks {
		q, tryErr := try(1)
		if tryErr != nil {
			log.Println("[DEBUG]", tryErr)
			continue
		}
		questions[qIndex] = q
		qIndex++
	}
	for qIndex < args.NumQuestions {
		q, tryErr := try(2)
		if tryErr != nil {
			continue
		}
		questions[qIndex] = q
		qIndex++
	}
	log.Println(tries, "tries")
	return questions, len(answerMap)
}

// genRack - Generate a random rack using `dist` and with `blanks` blanks.
func genRack(dist lexicon.LetterDistribution, wordLength, blanks int) []rune {
	bag := dist.MakeBag()
	// it's a bag of runes.
	rack := make([]rune, wordLength)
	idx := 0
	draw := func(avoidBlanks bool) rune {
		var tile rune
		if avoidBlanks {
			for tile, _ = bag.Draw(); tile == BlankCharacter; {
				tile, _ = bag.Draw()
			}
		} else {
			tile, _ = bag.Draw()
		}
		return tile
	}
	for idx < wordLength-blanks {
		// Avoid blanks on draw if user specifies a number of blanks.
		rack[idx] = draw(blanks != 0)
		idx++
	}
	for ; idx < wordLength; idx++ {
		rack[idx] = BlankCharacter
	}
	return rack
}
