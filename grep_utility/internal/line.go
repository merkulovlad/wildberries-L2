package internal

type Line struct {
	Number int
	Text   string
	Match  bool
	Printed bool
}

func NewLine(number int, text string, isMatch bool) *Line {
	return &Line{
		Number: number,
		Text:   text,
		Match:  isMatch,
		Printed: false,
	}
}
