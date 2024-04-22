package main

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

type scoreSummaryInput struct {
	quiz    PsychometryQuiz
	answers PsychometryAnswers
}

func makeAnswerSection(size int) []int {
	return make([]int, size)
}

func makeAnswerSectionArray(size int) [2][]int {
	return [2][]int{makeAnswerSection(size), makeAnswerSection(size)}
}

func (scoreSummaryInput) Generate(rand *rand.Rand, size int) reflect.Value {
	quiz := PsychometryQuiz{
		EssaySection: "",
		VSections:    makeSectionArray(size),
		QSections:    makeSectionArray(size),
		ESections:    makeSectionArray(size),
	}
	answers := PsychometryAnswers{
		EssaySection: "",
		VSections:    makeAnswerSectionArray(size),
		QSections:    makeAnswerSectionArray(size),
		ESections:    makeAnswerSectionArray(size),
	}

	return reflect.ValueOf(scoreSummaryInput{quiz, answers})
}

func rawOutOfBounds(rawScore int, sections [2]Section) bool {
	return rawScore < 0 || rawScore > len(sections[0].Questions)+len(sections[1].Questions)
}

func uniformOutOfBounds(uniformScore int) bool {
	return uniformScore < 50 || uniformScore > 150
}

func generalInvalid(generalScore [2]int) bool {
	outOfBoundsMin := generalScore[0] < 200 || generalScore[0] > 800
	outOfBoundsMax := generalScore[1] < 200 || generalScore[1] > 800
	outOfBounds := outOfBoundsMin || outOfBoundsMax
	invalidRange := generalScore[0] > generalScore[1]
	return outOfBounds || invalidRange
}

// Test: calculating a score summary always returns a valid score summary within the proper ranges
func TestCalculateScoreSummary_valid(t *testing.T) {
	valid := func(input scoreSummaryInput) bool {
		scoreSummary := CalculateScoreSummary(input.quiz, input.answers)

		// Ensure the raw scores are never greater than their section sizes or less than 0
		vRawOutOfBounds := rawOutOfBounds(scoreSummary.VRaw, input.quiz.VSections)
		qRawOutOfBounds := rawOutOfBounds(scoreSummary.QRaw, input.quiz.QSections)
		eRawOutOfBounds := rawOutOfBounds(scoreSummary.ERaw, input.quiz.ESections)
		if vRawOutOfBounds || qRawOutOfBounds || eRawOutOfBounds {
			return false
		}

		// Ensure the uniform scores are always between 50 and 150
		vUniformOutOfBounds := uniformOutOfBounds(scoreSummary.VUniform)
		qUniformOutOfBounds := uniformOutOfBounds(scoreSummary.QUniform)
		eUniformOutOfBounds := uniformOutOfBounds(scoreSummary.EUniform)
		if vUniformOutOfBounds || qUniformOutOfBounds || eUniformOutOfBounds {
			return false
		}

		// Ensure the general score range is between 200 and 800,
		// and that the minimum score is less than or equal to the maximum score
		multiCategoryInvalid := generalInvalid(scoreSummary.MultiCategoryGeneral)
		mathFocusInvalid := generalInvalid(scoreSummary.MathFocusGeneral)
		languageFocusInvalid := generalInvalid(scoreSummary.LanguageFocusGeneral)
		if multiCategoryInvalid || mathFocusInvalid || languageFocusInvalid {
			return false
		}

		return true
	}

	if err := quick.Check(valid, nil); err != nil {
		t.Error(err)
	}
}
