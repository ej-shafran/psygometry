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

func makeSection(rand *rand.Rand, size int) Section {
	// At minimum, each section must have one answer
	questions := make([]Question, size+1)
	for i := range questions {
		questions[i] = Question{CorrectOption: rand.Intn(3)}
	}

	kinds := []SectionKind{V, Q, E}
	kind := kinds[rand.Intn(len(kinds))]

	return Section{Kind: kind, Questions: questions, IsCounted: true}
}

func makeSectionArray(rand *rand.Rand, size int) []Section {
	sections := make([]Section, size)

	for i := range sections {
		sections[i] = makeSection(rand, size)
	}

	return sections
}

func (Psychometry) Generate(rand *rand.Rand, size int) reflect.Value {
	psychometry := Psychometry{
		WritingSection: "",
		Sections:       makeSectionArray(rand, size),
	}
	return reflect.ValueOf(psychometry)
}

// Test: parsing an answer form with only a writing section and other irrelevant keys is done successfully
type successValues url.Values

func (successValues) Generate(rand *rand.Rand, size int) reflect.Value {
	form := url.Values{}

	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	for range size {
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
	}

	var writingBuffer bytes.Buffer
	for i := 0; i < size; i++ {
		index := rand.Intn(len(alphabet))
		writingBuffer.WriteString(string(alphabet[index]))
	}

	form.Add("WritingSection", writingBuffer.String())

	return reflect.ValueOf(successValues(form))
}

func TestParsePsychometryAnswers_success(t *testing.T) {
	success := func(psychometry Psychometry, form successValues) bool {
		a, err := ParsePsychometryAnswers(url.Values(form), psychometry)
		return err == nil && a.WritingSection == url.Values(form).Get("WritingSection")
	}

	if err := quick.Check(success, nil); err != nil {
		t.Error(err)
	}
}

// Test: parsing an answer form with keys that have missing indexes errors with `MissingIndex`
type missingIndexValues url.Values

func (missingIndexValues) Generate(rand *rand.Rand, size int) reflect.Value {
	form := url.Values{}

	entries := [][2]string{{"Sections.0", "0"}, {"Sections", "0"}}

	index := rand.Intn(len(entries))
	form.Add(entries[index][0], entries[index][1])

	return reflect.ValueOf(missingIndexValues(form))
}

func TestParsePsychometryAnswers_missingIndex(t *testing.T) {
	missingIndex := func(psychometry Psychometry, form missingIndexValues) bool {
		_, err := ParsePsychometryAnswers(url.Values(form), psychometry)
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

	entries := [][2]string{{"Sections[0][Other]", "0"}, {"Sections[Other][0]", "0"}, {"Sections[0][0]", "Other"}}

	index := rand.Intn(len(entries))
	form.Add(entries[index][0], entries[index][1])

	return reflect.ValueOf(deformedIndexValues(form))
}

func TestParsePsychometryAnswers_deformedIndex(t *testing.T) {
	deformedIndex := func(psychometry Psychometry, form deformedIndexValues) bool {
		_, err := ParsePsychometryAnswers(url.Values(form), psychometry)
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
		{"Sections[0][0]", fmt.Sprint(rand.Int() + 4)},
		{"Sections[0][0]", fmt.Sprint(-1 * rand.Int())},
		{fmt.Sprintf("Sections[%d][0]", rand.Int()+2), "0"},
		{fmt.Sprintf("Sections[%d][0]", -1*rand.Int()), "0"},
	}

	index := rand.Intn(len(entries))
	form.Add(entries[index][0], entries[index][1])

	return reflect.ValueOf(invalidIndexValues(form))
}

func TestParsePsychometryAnswers_invalidIndex(t *testing.T) {
	invalidIndex := func(psychometry Psychometry, form invalidIndexValues) bool {
		_, err := ParsePsychometryAnswers(url.Values(form), psychometry)
		return err == InvalidIndex
	}

	if err := quick.Check(invalidIndex, nil); err != nil {
		t.Error(err)
	}
}
