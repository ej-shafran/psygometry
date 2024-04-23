package main

type WritingScore struct {
	Linguistic  int
	Content     int
	Explanation string
}

type Scores struct {
	VRaw int
	QRaw int
	ERaw int

	VUniform int
	QUniform int
	EUniform int

	MultiCategoryUniform     int
	VerbalFocusUniform       int
	QuantitativeFocusUniform int

	MultiCategoryGeneral     [2]int
	VerbalFocusGeneral       [2]int
	QuantitativeFocusGeneral [2]int
}

type ScoreSummary struct {
	StaticScores  Scores
	WritingScore  WritingScore
	DynamicScores Scores
}

func multiCategoryUniform(vUniform int, qUniform int, eUniform int) int {
	return (2*vUniform + 2*qUniform + eUniform) / 5
}

func verbalFocusUniform(vUniform int, qUniform int, eUniform int) int {
	return (3*vUniform + qUniform + eUniform) / 5
}

func quantitativeFocusUniform(vUniform int, qUniform int, eUniform int) int {
	return (3*qUniform + vUniform + eUniform) / 5
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

// Converts a score on the uniform scale to the generic PET score scale.
//
// "General PET scores: Each one of the general PET scores is based on a particular weight given to the raw score in each test domain. In the multi-domain score, the verbal reasoning and quantitative reasoning scores are each assigned double the weight of the score in English. In the quantitative-orientated score, the quantitative reasoning score is given triple the weight of each of the other two scores. In the verbal-orientated score the verbal reasoning score is given triple the weight of each of the other two scores." - [nite.org.il]
//
// [nite.org.il]: https://www.nite.org.il/psychometric-entrance-test/scores/calculation/?lang=en
func generalMeasurementRange(score int) [2]int {
	if score <= 50 {
		return [2]int{200, 200}
	}

	if score >= 150 {
		return [2]int{800, 800}
	}

	return measurementRanges[((score-1)/5)*5]
}

func CalculateScoreSummary(quiz PsychometryQuiz, answers PsychometryAnswers) (*ScoreSummary, error) {
	static := calculateStaticScores(quiz, answers)

	writing, err := calculateWritingScore(quiz.WritingSection, answers.WritingSection)
	if err != nil {
		return nil, err
	}

	dynamic := Scores{}

	dynamic.VRaw = writing.Content + writing.Linguistic + static.VRaw
	dynamic.QRaw = static.QRaw
	dynamic.ERaw = static.ERaw

	writingPercent := (writing.Content + writing.Linguistic) * 100 / 12
	staticVPercent := static.VUniform - 50

	dynamic.VUniform = ((writingPercent + (staticVPercent * 3)) / 4) + 50
	dynamic.QUniform = static.QUniform
	dynamic.EUniform = static.EUniform

	dynamic.MultiCategoryUniform = multiCategoryUniform(dynamic.VUniform, dynamic.QUniform, dynamic.EUniform)
	dynamic.VerbalFocusUniform = verbalFocusUniform(dynamic.VUniform, dynamic.QUniform, dynamic.EUniform)
	dynamic.QuantitativeFocusUniform = quantitativeFocusUniform(dynamic.VUniform, dynamic.QUniform, dynamic.EUniform)

	dynamic.MultiCategoryGeneral = generalMeasurementRange(dynamic.MultiCategoryUniform)
	dynamic.VerbalFocusGeneral = generalMeasurementRange(dynamic.VerbalFocusUniform)
	dynamic.QuantitativeFocusGeneral = generalMeasurementRange(dynamic.QuantitativeFocusUniform)

	summary := &ScoreSummary{
		WritingScore:  *writing,
		StaticScores:  static,
		DynamicScores: dynamic,
	}

	return summary, nil
}
