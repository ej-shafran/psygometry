package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const charactersPerLine = 40
const minimumLines = 25
const maximumLines = 50

func writingOutOfBounds(writing string) *WritingScore {
	count := strings.Count(writing, "")

	if count < minimumLines*charactersPerLine {
		return &WritingScore{
			Linguistic:  0,
			Content:     0,
			Explanation: fmt.Sprintf("הכתיבה הייתה מתחת למספר השורות המינימלי של %d.", minimumLines),
		}
	}

	if count > maximumLines*charactersPerLine {
		return &WritingScore{
			Linguistic:  0,
			Content:     0,
			Explanation: fmt.Sprintf("הכתיבה הייתה מעל למספר השורות המקסימלי של %d.", maximumLines),
		}
	}

	return nil
}

var scoreSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"linguistic":  {Type: genai.TypeInteger},
		"content":     {Type: genai.TypeInteger},
		"explanation": {Type: genai.TypeString},
	},
	Required: []string{"linguistic", "content", "explanation"},
}

var calculateWritingScoreFunc = genai.FunctionDeclaration{
	Name:       "CalculateWritingScore",
	Parameters: scoreSchema,
}

const writingScorePrompt = `
Please return JSON grading this essay, based on its prompt and the following rules.

The prompt (between the "-----" delimiters):
-----
%s
-----

The rules:

- The "linguistic" field must be a score between 0 and 6 grading the essay's grammar, spelling, and linguistic level. The essay should be in Hebrew. If it is not, this field should be 0.
- The "content" field must be a score between 0 and 6 grading the essay's coherency, structure, and critical thinking as it relates to the prompt. The essay should be in Hebrew. If it is not, this field should be 0.
- The "explanation" field must be a textual explanation of why you have chosen the two grades listed above. It should be in Hebrew.

Here is the essay (between the "-----" delimiters):
-----
%s
-----
`

var malformedRegexp = regexp.MustCompile("-{5}")

func parseGeminiInt(v any) (int, bool) {
	kind := reflect.TypeOf(v).Kind()
	if kind == reflect.String {
		ret, err := strconv.Atoi(v.(string))
		if err != nil {
			return 0, false
		}
		return ret, true
	}

	if kind == reflect.Float32 {
		return int(v.(float32)), true
	}

	if kind == reflect.Float64 {
		return int(v.(float64)), true
	}

	if kind == reflect.Int {
		return v.(int), true
	}

	return 0, false
}

func calculateWritingScore(prompt string, writing string) (*WritingScore, error) {
	outOfBounds := writingOutOfBounds(writing)
	if outOfBounds != nil {
		return outOfBounds, nil
	}

	if malformedRegexp.MatchString(writing) {
		return nil, errors.New("malformed content")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("missing api key")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	model.Tools = []*genai.Tool{
		{FunctionDeclarations: []*genai.FunctionDeclaration{&calculateWritingScoreFunc}},
	}
	model.ToolConfig = &genai.ToolConfig{
		FunctionCallingConfig: &genai.FunctionCallingConfig{Mode: genai.FunctionCallingAny},
	}

	response, err := model.GenerateContent(ctx, genai.Text(fmt.Sprintf(writingScorePrompt, prompt, writing)))
	if err != nil {
		return nil, err
	}

	log.Println("gemini response = ")
	responseJson, err := json.Marshal(response)
	log.Println(string(responseJson))

	data, ok := response.Candidates[0].Content.Parts[0].(genai.FunctionCall)
	if !ok {
		return nil, errors.New("invalid gemini response")
	}
	explanation, ok := data.Args["explanation"].(string)
	if !ok {
		return nil, errors.New("invalid gemini response")
	}
	linguistic, ok := parseGeminiInt(data.Args["linguistic"])
	if !ok {
		return nil, errors.New("invalid gemini response")
	}
	content, ok := parseGeminiInt(data.Args["content"])
	if !ok {
		return nil, errors.New("invalid gemini response")
	}
	writingScore := &WritingScore{
		Linguistic:  linguistic,
		Content:     content,
		Explanation: explanation,
	}

	return writingScore, nil
}
