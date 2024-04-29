package main

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

type scoreSummaryInput struct {
	psychometry Psychometry
	answers     PsychometryAnswers
}

func makeAnswerSection(rand *rand.Rand, size int) []int {
	answerSection := make([]int, size)

	for i := range answerSection {
		answerSection[i] = rand.Intn(3)
	}

	return answerSection
}

func makeAnswerSectionArray(rand *rand.Rand, size int) [][]int {
	answerSections := make([][]int, size)

	for i := range answerSections {
		answerSections[i] = makeAnswerSection(rand, size)
	}

	return answerSections
}

func (scoreSummaryInput) Generate(rand *rand.Rand, size int) reflect.Value {
	psychometry := Psychometry{
		WritingSection: "",
		Sections:       make([]Section, size),
	}
	answers := PsychometryAnswers{
		WritingSection: "",
		Sections:       makeAnswerSectionArray(rand, size),
	}

	return reflect.ValueOf(scoreSummaryInput{psychometry, answers})
}

func rawOutOfBounds(rawScore int, sections []Section) bool {
	totalQuestions := 0
	for _, section := range sections {
		totalQuestions += len(section.Questions)

	}

	return rawScore < 0 || rawScore > totalQuestions
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
		scoreSummary := calculateStaticScores(input.psychometry, input.answers)

		// Ensure the raw scores are never greater than their section sizes or less than 0
		vRawOutOfBounds := rawOutOfBounds(scoreSummary.VRaw, input.psychometry.GetSections(V))
		qRawOutOfBounds := rawOutOfBounds(scoreSummary.QRaw, input.psychometry.GetSections(Q))
		eRawOutOfBounds := rawOutOfBounds(scoreSummary.ERaw, input.psychometry.GetSections(E))
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
		mathFocusInvalid := generalInvalid(scoreSummary.QuantitativeFocusGeneral)
		languageFocusInvalid := generalInvalid(scoreSummary.VerbalFocusGeneral)
		if multiCategoryInvalid || mathFocusInvalid || languageFocusInvalid {
			return false
		}

		return true
	}

	if err := quick.Check(valid, nil); err != nil {
		t.Error(err)
	}
}
