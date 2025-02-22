package types

import (
	"strings"
	"text/template"
)

var templateFuncs = template.FuncMap{
	"split": strings.Split,
}
