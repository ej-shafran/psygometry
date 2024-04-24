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

type Psychometry struct {
	WritingSection string
	VSections      [2]Section
	QSections      [2]Section
	ESections      [2]Section
}

type PsychometryAnswers struct {
	WritingSection string
	VSections      [2][]int
	QSections      [2][]int
	ESections      [2][]int
}

func newPsychometryAnswers(psychometry Psychometry) PsychometryAnswers {
	answers := PsychometryAnswers{
		WritingSection: "",
		VSections: [2][]int{
			make([]int, len(psychometry.VSections[0].Questions)),
			make([]int, len(psychometry.VSections[1].Questions)),
		},
		QSections: [2][]int{
			make([]int, len(psychometry.QSections[0].Questions)),
			make([]int, len(psychometry.QSections[1].Questions)),
		},
		ESections: [2][]int{
			make([]int, len(psychometry.ESections[0].Questions)),
			make([]int, len(psychometry.ESections[1].Questions)),
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

var (
	InvalidKey    = errors.New("invalid key")
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

		if path[0] != "VSections" && path[0] != "QSections" && path[0] != "ESections" {
			return nil, InvalidKey
		}

		if len(path) < 2 {
			return nil, MissingIndex
		}
		sIndex, err := strconv.Atoi(path[1])
		if err != nil {
			return nil, DeformedIndex
		}
		if sIndex != 0 && sIndex != 1 {
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
			return nil, InvalidIndex
		}

		arr[sIndex][qIndex] = value
	}

	return &answers, nil
}

func generateFakeData() Psychometry {
	psychometry := Psychometry{
		WritingSection: "נא לכתוב חיבור על החשיבות של סיפור סיפורים בקולנוע המודרני.",
		VSections: [2]Section{
			{
				Kind: "V",
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
				Kind: "V",
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
		},
		QSections: [2]Section{
			{
				Kind: "Q",
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
				Kind: "Q",
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
		},
		ESections: [2]Section{
			{
				Kind: "E",
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
				Kind: "E",
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
