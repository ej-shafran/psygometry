package main

type Question struct {
	Options       [4]string
	CorrectOption int
}

type Section struct {
	Questions []Question
}

type PsychometryTest struct {
	EssaySection string
	VSections    [2]Section
	QSections    [2]Section
	ESections    [2]Section
}
