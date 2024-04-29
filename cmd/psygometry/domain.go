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

type SectionKind string

const (
	V SectionKind = "V"
	Q SectionKind = "Q"
	E SectionKind = "E"
)

type Section struct {
	Kind      SectionKind
	Index     int
	IsCounted bool
	Questions []Question
}

type Psychometry struct {
	WritingSection string
	Sections       []Section
}

func (p *Psychometry) GetSections(kind SectionKind) []Section {
	sections := []Section{}

	for _, section := range p.Sections {
		if section.Kind == kind && section.IsCounted {
			sections = append(sections, section)
		}
	}

	return sections
}

type PsychometryAnswers struct {
	WritingSection string
	Sections       [][]int
}

func (a *PsychometryAnswers) GetSections(psychometry Psychometry, kind SectionKind) [][]int {
	answerSections := [][]int{}

	for i, section := range psychometry.Sections {
		if section.Kind == kind && section.IsCounted {
			answerSections = append(answerSections, a.Sections[i])
		}
	}

	return answerSections
}

func newPsychometryAnswers(psychometry Psychometry) PsychometryAnswers {
	answerSections := make([][]int, len(psychometry.Sections))
	for i, section := range psychometry.Sections {
		answerSections[i] = make([]int, len(section.Questions))
		for j := range answerSections[i] {
			answerSections[i][j] = -1
		}
	}

	answers := PsychometryAnswers{
		WritingSection: "",
		Sections:       answerSections,
	}
	return answers
}

var (
	MissingIndex  = errors.New("missing index")
	DeformedIndex = errors.New("deformed index")
	InvalidIndex  = errors.New("invalid index")
)

func ParsePsychometryAnswers(form url.Values, psychometry Psychometry) (*PsychometryAnswers, error) {
	answers := newPsychometryAnswers(psychometry)

	r := regexp.MustCompile("[\\][.]+")

	for key := range form {
		path := r.Split(key, -1)

		if path[0] == "WritingSection" {
			answers.WritingSection = form.Get(key)
			continue
		}

		if path[0] != "Sections" {
			continue
		}

		if len(path) < 2 {
			return nil, MissingIndex
		}
		sIndex, err := strconv.Atoi(path[1])
		if err != nil {
			return nil, DeformedIndex
		}
		if sIndex < 0 || sIndex >= len(psychometry.Sections) {
			return nil, InvalidIndex
		}

		if len(path) < 3 {
			return nil, MissingIndex
		}
		qIndex, err := strconv.Atoi(path[2])
		if err != nil {
			return nil, DeformedIndex
		}

		rawValue := form.Get(key)
		var value int
		if rawValue == "" {
			value = -1
		} else {
			value, err = strconv.Atoi(rawValue)
			if err != nil {
				return nil, DeformedIndex
			}
		}
		if value < 0 || value > 4 {
			return nil, InvalidIndex
		}

		if qIndex < 0 || qIndex >= len(answers.Sections[sIndex]) {
			return nil, InvalidIndex
		}

		answers.Sections[sIndex][qIndex] = value
	}

	return &answers, nil
}

func generateFakeData() Psychometry {
	psychometry := Psychometry{
		WritingSection: "נא לכתוב חיבור על החשיבות של סיפור סיפורים בקולנוע המודרני.",
		Sections: []Section{
			{
				Kind:      V,
				Index:     0,
				IsCounted: true,
				Questions: []Question{
					{
						Content:       "מי משחק את הדמות הראשית בסרט 'ההסתערות'?",
						Options:       [4]string{"ליאונרדו דיקפריו", "בראד פיט", "טום הנקס", "ג'וני דפ"},
						CorrectOption: 0,
					},
					{
						Content:       "איזה סרט לא נבחר על ידי כריסטופר נולן?",
						Options:       [4]string{"ההסתערות", "בלונדינית משפטית", "בין הכוכבים", "אי הצנום"},
						CorrectOption: 1,
					},
				},
			},
			{
				Kind:      V,
				Index:     1,
				IsCounted: true,
				Questions: []Question{
					{
						Content:       "מי הוא המחבר של סדרת הספרים 'משחקי הכס'?",
						Options:       [4]string{"ג'יי. קי. רואלינג", "סטיבן קינג", "ג'ורג' אר.אר. מרטין", "ג'יי.אר.אר. טולקין"},
						CorrectOption: 2,
					},
					{
						Content:       "איזו סדרת ספרים כוללת דמות בשם 'הארי פוטר'?",
						Options:       [4]string{"הארי פוטר", "אדון הטבעות", "משחקי הכס", "המשחקים של הרעב"},
						CorrectOption: 0,
					},
				},
			},
			{
				Kind:      Q,
				Index:     2,
				IsCounted: true,
				Questions: []Question{
					{
						Content:       "איזה אבנג'ר מכונה בגלל המראה הירוק שלו והכוח המדהים שלו?",
						Options:       [4]string{"איירון מן", "קפטן אמריקה", "תור", "האלק"},
						CorrectOption: 3,
					},
					{
						Content:       "מי מגלם את הדמות של נרייט שחורה ביקום הסרטים המרובע של מארו?",
						Options:       [4]string{"סקרלט יוהנסון", "גל גדות", "אנג'לינה ג'ולי", "ג'ניפר לורנס"},
						CorrectOption: 0,
					},
				},
			},
			{
				Kind:      Q,
				Index:     3,
				IsCounted: true,
				Questions: []Question{
					{
						Content:       "איזה להקה מפורסמת בשיר 'בוהמיאן ראפסודיה'?",
						Options:       [4]string{"הביטלס", "לד זפלין", "קווין", "פלוויד הוויד"},
						CorrectOption: 2,
					},
					{
						Content:       "איזה סרט לעיתים קרוא 'הסרט הגדול ביותר שנעשה אי פעם'?",
						Options:       [4]string{"הקרוטונאי", "פיקדון דמים", "בראש ובראש", "פנים שטוחות"},
						CorrectOption: 0,
					},
				},
			},
			{
				Kind:      E,
				Index:     4,
				IsCounted: true,
				Questions: []Question{
					{
						Content:       "מי צייר את היצירה המפורסמת 'לילה כוכבי'?",
						Options:       [4]string{"מונה", "ואן גוך", "פיקאסו", "דה וינצ'י"},
						CorrectOption: 1,
					},
					{
						Content:       "איזה מלחין מוכר כ 'הגאון'?",
						Options:       [4]string{"מוצארט", "בטהובן", "באך", "שופין"},
						CorrectOption: 0,
					},
				},
			},
			{
				Kind:      E,
				Index:     5,
				IsCounted: true,
				Questions: []Question{
					{
						Content:       "מי זכתה בפרס אוסקר לשחקנית הטובה ביותר על תפקידה ב'ברבור שחור'?",
						Options:       [4]string{"מריל סטריפ", "קייט בלנשט", "ג'וליאן מור", "נטלי פורטמן"},
						CorrectOption: 3,
					},
					{
						Content:       "איזה במאי ידוע בסרטיו האפיים כמו 'רשימת שינדלר' ו'שמור פרטי'?",
						Options:       [4]string{"סטיבן שפילברג", "מרטין סקורסזה", "קוונטין טרנטינו", "כריסטופר נולן"},
						CorrectOption: 0,
					},
				},
			},
		},
	}
	return psychometry
}
