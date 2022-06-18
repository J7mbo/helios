package helios

import (
	"github.com/amenzhinsky/go-memexec"
	"reflect"
	"strconv"
)

type HighlightRequest struct {
	ScreenWidth  int     `arg:"--screen-width"`
	ScreenHeight int     `arg:"--screen-height"`
	X            float64 `arg:"--x"`
	Y            float64 `arg:"--y"`
	Width        float64 `arg:"--w"`
	Height       float64 `arg:"--h"`
	Duration     float64 `arg:"--d"`
}

func (r HighlightRequest) AsArgsArray() []string {
	var a []string

	requestReflect := reflect.TypeOf(r)
	reflectVal := reflect.ValueOf(r)
	numFields := requestReflect.NumField()
	i := 0
	for i < numFields {
		fieldName := requestReflect.Field(i).Tag.Get("arg")
		fieldValue := reflectVal.Field(i).Interface()

		a = append(a, fieldName)

		if val, isInt := fieldValue.(int); isInt {
			a = append(a, strconv.Itoa(val))
		}

		if val, isFloat := fieldValue.(float64); isFloat {
			a = append(a, strconv.FormatFloat(val, 'f', -1, 64))
		}

		i++
	}

	return a
}

type Highlighter struct {
	highlighterBinary []byte
}

func NewHighlighter(highlighterBinary []byte) *Highlighter {
	return &Highlighter{
		highlighterBinary: highlighterBinary,
	}
}

func (h *Highlighter) Highlight(r *HighlightRequest) {
	args := r.AsArgsArray()

	exe, err := memexec.New(h.highlighterBinary)
	if err != nil {
		return
	}
	defer func() { _ = exe.Close() }()

	_, _ = exe.Command(args...).Output()
}
