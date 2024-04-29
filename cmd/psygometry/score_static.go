package main

// Calculates raw scores for a certain domain.
//
// "Raw scores on the multiple-choice sections: Each correct answer is worth one point. The number of correct answers in each domain is equal to the raw score in that domain." - [nite.org.il]
//
// [nite.org.il]: https://www.nite.org.il/psychometric-entrance-test/scores/calculation/?lang=en
func rawCategoryScore(psychometrySections []Section, answerSections [][]int) int {
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
func uniformCategoryScore(psychometrySections []Section, rawScore int) int {
	totalQuestions := 0
	for _, section := range psychometrySections {
		totalQuestions += len(section.Questions)
	}

	if totalQuestions == 0 {
		return 50
	}

	percent := rawScore * 100 / totalQuestions
	return percent + 50
}

func calculateStaticScores(psychometry Psychometry, answers PsychometryAnswers) Scores {
	scores := Scores{}

	scores.VRaw = rawCategoryScore(psychometry.GetSections(V), answers.GetSections(psychometry, V))
	scores.QRaw = rawCategoryScore(psychometry.GetSections(Q), answers.GetSections(psychometry, Q))
	scores.ERaw = rawCategoryScore(psychometry.GetSections(E), answers.GetSections(psychometry, E))

	scores.VUniform = uniformCategoryScore(psychometry.GetSections(V), scores.VRaw)
	scores.QUniform = uniformCategoryScore(psychometry.GetSections(Q), scores.QRaw)
	scores.EUniform = uniformCategoryScore(psychometry.GetSections(E), scores.ERaw)

	scores.MultiCategoryUniform = multiCategoryUniform(scores.VUniform, scores.QUniform, scores.EUniform)
	scores.VerbalFocusUniform = verbalFocusUniform(scores.VUniform, scores.QUniform, scores.EUniform)
	scores.QuantitativeFocusUniform = quantitativeFocusUniform(scores.VUniform, scores.QUniform, scores.EUniform)

	scores.MultiCategoryGeneral = generalMeasurementRange(scores.MultiCategoryUniform)
	scores.VerbalFocusGeneral = generalMeasurementRange(scores.VerbalFocusUniform)
	scores.QuantitativeFocusGeneral = generalMeasurementRange(scores.QuantitativeFocusUniform)

	return scores
}
