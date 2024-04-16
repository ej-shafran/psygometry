package main

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
)

type Question struct {
	Content       string
	Options       [4]string
	CorrectOption int
}

type Section struct {
	Kind      string
	Questions []Question
}

type PsychometryTest struct {
	EssaySection string
	VSections    [2]Section
	QSections    [2]Section
	ESections    [2]Section
}

type PsychometryAnswers struct {
	EssaySection string
	VSections    [2][]int
	QSections    [2][]int
	ESections    [2][]int
}

func newPsychometryAnswers(test PsychometryTest) PsychometryAnswers {
	answers := PsychometryAnswers{
		EssaySection: "",
		VSections: [2][]int{
			make([]int, len(test.VSections[0].Questions)),
			make([]int, len(test.VSections[1].Questions)),
		},
		QSections: [2][]int{
			make([]int, len(test.QSections[0].Questions)),
			make([]int, len(test.QSections[1].Questions)),
		},
		ESections: [2][]int{
			make([]int, len(test.ESections[0].Questions)),
			make([]int, len(test.ESections[1].Questions)),
		},
	}

	for _, s := range answers.VSections {
		for j := range s {
			s[j] = -1
		}
	}
	for _, s := range answers.QSections {
		for j := range s {
			s[j] = -1
		}
	}
	for _, s := range answers.ESections {
		for j := range s {
			s[j] = -1
		}
	}

	return answers
}

func ParsePsychometryAnswers(form url.Values, test PsychometryTest) (*PsychometryAnswers, error) {
	answers := newPsychometryAnswers(test)

	r := regexp.MustCompile("[\\][.]+")

	for key := range form {
		path := r.Split(key, -1)

		if len(path) < 1 {
			return nil, errors.New("deformed key")
		}

		if path[0] == "EssaySection" {
			answers.EssaySection = form.Get(key)
			continue
		}

		if path[0] != "VSections" && path[0] != "QSections" && path[0] != "ESections" {
			return nil, errors.New("invalid key")
		}

		if len(path) < 2 {
			return nil, errors.New("missing section index")
		}
		sIndex, err := strconv.Atoi(path[1])
		if err != nil {
			return nil, errors.New("deformed section index")
		}
		if sIndex != 0 && sIndex != 1 {
			return nil, errors.New("invalid section index")
		}

		if len(path) < 3 {
			return nil, errors.New("missing question index")
		}
		qIndex, err := strconv.Atoi(path[2])
		if err != nil {
			return nil, errors.New("deformed question index")
		}

		rawValue := form.Get(key)
		var value int
		if rawValue == "" {
			value = -1
		} else {
			value, err = strconv.Atoi(rawValue)
			if err != nil {
				return nil, errors.New("deformed option index")
			}
		}
		if value < 0 || value > 4 {
			return nil, errors.New("invalid option index")
		}

		var arr [2][]int
		switch path[0] {
		case "VSections":
			arr = answers.VSections
			break
		case "QSections":
			arr = answers.QSections
			break
		case "ESections":
			arr = answers.ESections
			break
		default:
			panic("Invariant")
		}

		if qIndex < 0 || qIndex >= len(arr[sIndex]) {
			return nil, errors.New("invalid question index")
		}

		arr[sIndex][qIndex] = value
	}

	return &answers, nil
}

func generateFakeData() PsychometryTest {
	test := PsychometryTest{
		EssaySection: "Please write an essay on the importance of storytelling in modern cinema.",
		VSections: [2]Section{
			{
				Kind: "V",
				Questions: []Question{
					{
						Content:       "Who played the lead role in the movie 'Inception'?",
						Options:       [4]string{"Leonardo DiCaprio", "Brad Pitt", "Tom Hanks", "Johnny Depp"},
						CorrectOption: 0,
					},
					{
						Content:       "Which movie was not directed by Christopher Nolan?",
						Options:       [4]string{"Inception", "Legally Blonde", "Interstellar", "Shutter Island"},
						CorrectOption: 1,
					},
				},
			},
			{
				Kind: "V",
				Questions: []Question{
					{
						Content:       "Who is the author of the 'Game of Thrones' book series?",
						Options:       [4]string{"J.K. Rowling", "Stephen King", "George R.R. Martin", "J.R.R. Tolkien"},
						CorrectOption: 2,
					},
					{
						Content:       "Which book series features a character named Harry Potter?",
						Options:       [4]string{"Harry Potter", "Lord of the Rings", "Game of Thrones", "The Hunger Games"},
						CorrectOption: 0,
					},
				},
			},
		},
		QSections: [2]Section{
			{
				Kind: "Q",
				Questions: []Question{
					{
						Content:       "Which Avenger is known for his green appearance and incredible strength?",
						Options:       [4]string{"Iron Man", "Captain America", "Thor", "Hulk"},
						CorrectOption: 3,
					},
					{
						Content:       "Who portrayed the character of Black Widow in the Marvel Cinematic Universe?",
						Options:       [4]string{"Scarlett Johansson", "Gal Gadot", "Angelina Jolie", "Jennifer Lawrence"},
						CorrectOption: 0,
					},
				},
			},
			{
				Kind: "Q",
				Questions: []Question{
					{
						Content:       "Which band is known for the song 'Bohemian Rhapsody'?",
						Options:       [4]string{"The Beatles", "Led Zeppelin", "Queen", "Pink Floyd"},
						CorrectOption: 2,
					},
					{
						Content:       "Which movie is often referred to as 'the greatest film ever made'?",
						Options:       [4]string{"The Godfather", "Pulp Fiction", "Goodfellas", "Scarface"},
						CorrectOption: 0,
					},
				},
			},
		},
		ESections: [2]Section{
			{
				Kind: "E",
				Questions: []Question{
					{
						Content:       "Who painted the famous artwork 'Starry Night'?",
						Options:       [4]string{"Monet", "Van Gogh", "Picasso", "Da Vinci"},
						CorrectOption: 1,
					},
					{
						Content:       "Which composer is known as 'The Genius'?",
						Options:       [4]string{"Mozart", "Beethoven", "Bach", "Chopin"},
						CorrectOption: 0,
					},
				},
			},
			{
				Kind: "E",
				Questions: []Question{
					{
						Content:       "Who won the Academy Award for Best Actress for her role in 'Black Swan'?",
						Options:       [4]string{"Meryl Streep", "Cate Blanchett", "Julianne Moore", "Natalie Portman"},
						CorrectOption: 3,
					},
					{
						Content:       "Which director is known for his epic films like 'Schindler's List' and 'Saving Private Ryan'?",
						Options:       [4]string{"Steven Spielberg", "Martin Scorsese", "Quentin Tarantino", "Christopher Nolan"},
						CorrectOption: 0,
					},
				},
			},
		},
	}
	return test
}
