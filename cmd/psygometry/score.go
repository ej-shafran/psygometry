package main

func rawCategoryScore(quizSections [2]Section, answerSections [2][]int) int {
	score := 0

	for i, section := range quizSections {
		for j, question := range section.Questions {
			option := answerSections[i][j]
			if option == question.CorrectOption {
				score += 1
			}
		}
	}

	return score
}

func uniformCategoryScore(quizSections [2]Section, rawScore int) int {
	totalQuestions := len(quizSections[0].Questions) + len(quizSections[1].Questions)
	percent := rawScore * 100 / totalQuestions
	return percent + 50
}

var measurementRanges = map[int][2]int{
	51:  {221, 248},
	56:  {249, 276},
	61:  {277, 304},
	66:  {305, 333},
	71:  {334, 361},
	76:  {362, 389},
	81:  {390, 418},
	86:  {419, 446},
	91:  {447, 474},
	96:  {475, 503},
	101: {504, 531},
	106: {532, 559},
	111: {560, 587},
	116: {588, 616},
	121: {617, 644},
	126: {645, 672},
	131: {673, 701},
	136: {702, 729},
	141: {730, 761},
	145: {762, 795},
}

func generalMeasurementRange(score int) [2]int {
	if score <= 50 {
		return [2]int{200, 200}
	}

	if score >= 150 {
		return [2]int{800, 800}
	}

	return measurementRanges[((score-1)/5)*5]
}

type ScoreSummary struct {
	VRaw int
	QRaw int
	ERaw int

	VUniform int
	QUniform int
	EUniform int

	MultiCategoryUniform int
	LanguageFocusUniform int
	MathFocusUniform     int

	MultiCategoryGeneral [2]int
	LanguageFocusGeneral [2]int
	MathFocusGeneral     [2]int
}

func CalculateScoreSummary(quiz PsychometryQuiz, answers PsychometryAnswers) ScoreSummary {
	summary := ScoreSummary{}

	summary.VRaw = rawCategoryScore(quiz.VSections, answers.VSections)
	summary.QRaw = rawCategoryScore(quiz.QSections, answers.QSections)
	summary.ERaw = rawCategoryScore(quiz.ESections, answers.ESections)

	summary.VUniform = uniformCategoryScore(quiz.VSections, summary.VRaw)
	summary.QUniform = uniformCategoryScore(quiz.QSections, summary.QRaw)
	summary.EUniform = uniformCategoryScore(quiz.ESections, summary.ERaw)

	summary.MultiCategoryUniform = (2*summary.VUniform + 2*summary.QUniform + summary.EUniform) / 5
	summary.LanguageFocusUniform = (3*summary.VUniform + summary.QUniform + summary.EUniform) / 5
	summary.MathFocusUniform = (3*summary.QUniform + summary.VUniform + summary.EUniform) / 5

	summary.MultiCategoryGeneral = generalMeasurementRange(summary.MultiCategoryUniform)
	summary.LanguageFocusGeneral = generalMeasurementRange(summary.LanguageFocusUniform)
	summary.MathFocusGeneral = generalMeasurementRange(summary.MathFocusUniform)

	return summary
}
