package input

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/hyprstereo/go-dao/utils/template/ft"
)

func Read(format string, in string) (i *Input, score float64, err error) {
	i = &Input{}
	score, err = i.Read(format, strings.NewReader(in))
	return
}

func NewInput() (i *Input) {
	i = &Input{}
	return
}

type Input struct {
	line     string
	fmtValue string
	vars     map[string]*Var
}

func (i *Input) Read(format string, r io.Reader) (score float64, err error) {
	i.vars = make(map[string]*Var)
	words := Split(format)
	matchers := []string{}
	cnt := 0
	//vars := make([]any, 0)
	out, _ := ioutil.ReadAll(r)
	i.line = string(out)
	outs := SplitArgs(string(out))
	var scores = 0
	for x := 0; x < len(words); x++ {
		if isVar(words[x]) {
			v := extractVar(words[x], cnt)
			if v.expectedKind == Array {
				fmt.Println("kindof", KindString(kindOf(outs[x])))
				if kindOf(outs[x]) == String {
					v.Value = outs[x]
				} else {
					v.Value = outs[x].([]any)[0]
				}
				scores++
			} else {
				if v.expectedKind == kindOf(outs[x]) {
					score++
				}
				v.Value = outs[x]
				//v.Value = "expected " + KindString(v.expectedKind) + ", got " + KindString(kindOf(out[x]))
			}
			i.vars[v.Name] = v
			matchers = append(matchers, v.fmtValue)

			cnt++
		} else {
			if words[x] == outs[x] {
				scores++
			}
			matchers = append(matchers, words[x])
		}
	}

	i.fmtValue = strings.Join(matchers, " ")
	score = float64(scores) / float64(len(words))
	return
}

// func (i *Input) ReadFMT(format string, r io.Reader) (err error) {
// 	i.vars = make(map[string]*Var)
// 	words := strings.Split(format, " ")
// 	matchers := []string{}
// 	cnt := 0
// 	vars := make([]any, 0)
// 	out, _ := ioutil.ReadAll(r)
// 	i.line = string(out)
// 	//outs := SplitArgs(string(out))

// 	for x := 0; x < len(words); x++ {
// 		if isVar(words[x]) {
// 			v := extractVar(words[x], cnt)
// 			i.vars[v.Name] = v
// 			vr := KindValue(v.expectedKind)
// 			vars = append(vars, &vr)
// 			matchers = append(matchers, v.fmtValue)

// 			cnt++
// 		} else {
// 			matchers = append(matchers, words[x])
// 		}
// 	}
// 	i.fmtValue = strings.Join(matchers, " ")
// 	fmt.Println(i.fmtValue)
// 	fmt.Fscanf(strings.NewReader(string(out)), i.fmtValue, vars...)
// 	for i, v := range vars {
// 		fmt.Println(i, pointerValue(vars[i], 0), v)
// 	}

// 	return
// }

func (i *Input) Matches() (vars []*Var, score float64) {
	scored := 0
	for _, v := range i.vars {
		if v.Value != nil {
			vars = append(vars, v)
			scored++
		}
	}
	score = float64(scored) / float64(len(i.vars))
	return
}

func (i *Input) Line() (str string) {
	str = i.fmtValue
	return
}

func (i *Input) Get(n string) (v *Var) {
	v = i.vars[n]
	return
}

func (i *Input) Values() (v []any) {
	for _, val := range i.vars {
		v = append(v, val.Value)
	}
	return
}

func (i *Input) Printf(w io.Writer, format string) (int, error) {
	str := ft.ExecuteString(format, "{{", "}}", i.All())
	return w.Write([]byte(str))
}

func (i *Input) Sprintf(format string) (out string) {
	out = ft.ExecuteString(format, "{{", "}}", i.All())
	return
}

func (i *Input) All() (v map[string]any) {
	o := map[string]any{}
	for n, val := range i.vars {
		o[n] = val.Value
	}
	v = o
	return
}

func isVar(v string) (ok bool) {
	if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
		ok = true
	}
	return
}

func extractVar(v string, pos int) (nv *Var) {
	v = strings.TrimPrefix(v, "${")
	v = strings.TrimSuffix(v, "}")
	if strings.Contains(v, ":") {
		toks := strings.Split(v, ":")
		knd := StringToKind(toks[1])
		nv = &Var{
			Name:         toks[0],
			Pos:          pos,
			expectedKind: knd,
			fmtValue:     KindFmtSymbol(knd),
		}
	} else {
		nv = &Var{
			Name:         v,
			Pos:          pos,
			expectedKind: Any,
			fmtValue:     KindFmtSymbol(Any),
		}
	}
	return
}

func Split(value string) (res []string) {
	canSplit := true
	inString := false
	lvl := 0
	tags := strings.FieldsFunc(value, func(r rune) bool {
		if r == '{' || r == '[' {
			lvl++
			canSplit = false
		} else if r == '}' || r == ']' {
			if lvl > 0 {
				lvl--
			}
			canSplit = lvl == 0
		}
		if r == '\'' || r == '"' {
			inString = !inString

		}

		if canSplit {
			return (r == ' ' || r == ',' || r == ':') && !inString
		} else {
			return false
		}
	})
	for _, t := range tags {
		if strings.TrimSpace(t) != "" {
			res = append(res, t)
		}
	}
	return
}

func SplitArgs(value string) (res []any) {
	tags := Split(value)
	res = make([]any, 0)
	for _, t := range tags {
		if out, er := expr.Eval(t, nil); er == nil {
			res = append(res, out)
		} else {
			res = append(res, t)
		}
	}
	return
}

func getLastIndexOf(value string, char rune, ignoreChar rune, start int) (pos int) {
	isIgnored := false
	for x := start; x < len(value); x++ {
		cr := rune(value[x])
		if cr == ignoreChar {
			isIgnored = true
		}
		if cr == char {
			if !isIgnored {
				pos = x
				return
			} else {
				isIgnored = false
			}
		}
	}
	return
}

func getRangeOf(value string, startChar rune, endChar rune, start int) (pos []int) {
	val := value[start:]
	first := strings.IndexRune(val, startChar)
	last := getLastIndexOf(val, endChar, startChar, first) + 1
	if last <= 1 {
		last = strings.LastIndex(val, string(endChar)) + 1
	}
	pos = []int{start + first, start + last}
	return
}
