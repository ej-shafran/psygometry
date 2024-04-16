package main

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
