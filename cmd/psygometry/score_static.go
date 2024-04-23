package main

// Calculates raw scores for a certain domain.
//
// "Raw scores on the multiple-choice sections: Each correct answer is worth one point. The number of correct answers in each domain is equal to the raw score in that domain." - [nite.org.il]
//
// [nite.org.il]: https://www.nite.org.il/psychometric-entrance-test/scores/calculation/?lang=en
func rawCategoryScore(psychometrySections [2]Section, answerSections [2][]int) int {
	score := 0

	for i, section := range psychometrySections {
		for j, question := range section.Questions {
			option := answerSections[i][j]
			if option == question.CorrectOption {
				score += 1
			}
		}
	}

	return score
}

// Converts a raw score to a uniform scale for a certain domain (based on the percentage of correct answers).
//
// Note: this only considers multiple-choice sections, not the writing score.
//
// "Scores in each of the three test domains: In order to compare the scores of examinees who took different versions of the test or who took the test in different languages or on different dates, the raw scores for the writing task and the raw scores for the multiple-choice sections in each of the three test domains are converted to a uniform scale. The verbal reasoning score includes the score on the writing task, which is weighted at 25%. The scale for scores in each of the three domains is from 50 to 150." - [nite.org.il]
//
// [nite.org.il]: https://www.nite.org.il/psychometric-entrance-test/scores/calculation/?lang=en
func uniformCategoryScore(psychometrySections [2]Section, rawScore int) int {
	totalQuestions := len(psychometrySections[0].Questions) + len(psychometrySections[1].Questions)
	percent := rawScore * 100 / totalQuestions
	return percent + 50
}

func calculateStaticScores(psychometry Psychometry, answers PsychometryAnswers) Scores {
	scores := Scores{}

	scores.VRaw = rawCategoryScore(psychometry.VSections, answers.VSections)
	scores.QRaw = rawCategoryScore(psychometry.QSections, answers.QSections)
	scores.ERaw = rawCategoryScore(psychometry.ESections, answers.ESections)

	scores.VUniform = uniformCategoryScore(psychometry.VSections, scores.VRaw)
	scores.QUniform = uniformCategoryScore(psychometry.QSections, scores.QRaw)
	scores.EUniform = uniformCategoryScore(psychometry.ESections, scores.ERaw)

	scores.MultiCategoryUniform = multiCategoryUniform(scores.VUniform, scores.QUniform, scores.EUniform)
	scores.VerbalFocusUniform = verbalFocusUniform(scores.VUniform, scores.QUniform, scores.EUniform)
	scores.QuantitativeFocusUniform = quantitativeFocusUniform(scores.VUniform, scores.QUniform, scores.EUniform)

	scores.MultiCategoryGeneral = generalMeasurementRange(scores.MultiCategoryUniform)
	scores.VerbalFocusGeneral = generalMeasurementRange(scores.VerbalFocusUniform)
	scores.QuantitativeFocusGeneral = generalMeasurementRange(scores.QuantitativeFocusUniform)

	return scores
}
