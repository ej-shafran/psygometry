package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"testing"
	"testing/quick"
)

func makeSection(size int) Section {
	return Section{Kind: "", Questions: make([]Question, size)}
}

func makeSectionArray(size int) [2]Section {
	return [2]Section{makeSection(size), makeSection(size)}
}

func (PsychometryQuiz) Generate(rand *rand.Rand, size int) reflect.Value {
	quiz := PsychometryQuiz{
		EssaySection: "",
		VSections:    makeSectionArray(size),
		QSections:    makeSectionArray(size),
		ESections:    makeSectionArray(size),
	}
	return reflect.ValueOf(quiz)
}

// Test: parsing an answer form with only an essay is done successfully
func TestParsePsychometryAnswers_success(t *testing.T) {
	success := func(quiz PsychometryQuiz, essay string) bool {
		form := url.Values{}
		form.Add("EssaySection", essay)

		a, err := ParsePsychometryAnswers(form, quiz)
		return err == nil && a.EssaySection == essay
	}

	if err := quick.Check(success, nil); err != nil {
		t.Error(err)
	}
}

// Test: parsing an answer form with a random key errors with `InvalidKey`
type invalidKeyValues url.Values

func (invalidKeyValues) Generate(rand *rand.Rand, size int) reflect.Value {
	form := url.Values{}

	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var keyBuffer bytes.Buffer
	for i := 0; i < size; i++ {
		index := rand.Intn(len(alphabet))
		keyBuffer.WriteString(string(alphabet[index]))
	}

	var valueBuffer bytes.Buffer
	for i := 0; i < size; i++ {
		index := rand.Intn(len(alphabet))
		valueBuffer.WriteString(string(alphabet[index]))
	}

	form.Add(keyBuffer.String(), valueBuffer.String())

	return reflect.ValueOf(invalidKeyValues(form))
}

func TestParsePsychometryAnswers_invalidKey(t *testing.T) {
	invalidKey := func(quiz PsychometryQuiz, form invalidKeyValues) bool {
		_, err := ParsePsychometryAnswers(url.Values(form), quiz)
		return err == InvalidKey
	}

	if err := quick.Check(invalidKey, nil); err != nil {
		t.Error(err)
	}
}

// Test: parsing an answer form with keys that have missing indexes errors with `MissingIndex`
type missingIndexValues url.Values

func (missingIndexValues) Generate(rand *rand.Rand, size int) reflect.Value {
	form := url.Values{}

	entries := [][2]string{{"VSections.0", "0"}, {"VSections", "0"}}

	index := rand.Intn(len(entries))
	form.Add(entries[index][0], entries[index][1])

	return reflect.ValueOf(missingIndexValues(form))
}

func TestParsePsychometryAnswers_missingIndex(t *testing.T) {
	missingIndex := func(quiz PsychometryQuiz, form missingIndexValues) bool {
		_, err := ParsePsychometryAnswers(url.Values(form), quiz)
		return err == MissingIndex
	}

	if err := quick.Check(missingIndex, nil); err != nil {
		t.Error(err)
	}
}

// Test: parsing an answer form with keys that have non-numerical values as indexes errors with `DeformedIndex`
type deformedIndexValues url.Values

func (deformedIndexValues) Generate(rand *rand.Rand, size int) reflect.Value {
	form := url.Values{}

	entries := [][2]string{{"VSections[0][Other]", "0"}, {"VSections[Other][0]", "0"}, {"VSections[0][0]", "Other"}}

	index := rand.Intn(len(entries))
	form.Add(entries[index][0], entries[index][1])

	return reflect.ValueOf(deformedIndexValues(form))
}

func TestParsePsychometryAnswers_deformedIndex(t *testing.T) {
	deformedIndex := func(quiz PsychometryQuiz, form deformedIndexValues) bool {
		_, err := ParsePsychometryAnswers(url.Values(form), quiz)
		return err == DeformedIndex
	}

	if err := quick.Check(deformedIndex, nil); err != nil {
		t.Error(err)
	}
}

// Test: parsing an answer form with keys that have out-of-range indexes errors with `InvalidIndex`
type invalidIndexValues url.Values

func (invalidIndexValues) Generate(rand *rand.Rand, size int) reflect.Value {
	form := url.Values{}

	entries := [][2]string{
		{"VSections[0][0]", fmt.Sprint(rand.Int() + 4)},
		{"VSections[0][0]", fmt.Sprint(-1 * rand.Int())},
		{fmt.Sprintf("VSections[%d][0]", rand.Int()+2), "0"},
		{fmt.Sprintf("VSections[%d][0]", -1*rand.Int()), "0"},
	}

	index := rand.Intn(len(entries))
	form.Add(entries[index][0], entries[index][1])

	return reflect.ValueOf(invalidIndexValues(form))
}

func TestParsePsychometryAnswers_invalidIndex(t *testing.T) {
	invalidIndex := func(quiz PsychometryQuiz, form invalidIndexValues) bool {
		_, err := ParsePsychometryAnswers(url.Values(form), quiz)
		return err == InvalidIndex
	}

	if err := quick.Check(invalidIndex, nil); err != nil {
		t.Error(err)
	}
}
