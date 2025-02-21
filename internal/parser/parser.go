package parser

import (
	"os"
	"regexp"
	"strings"

	"github.com/catmorte/go-mdapi/internal/file"
)

func ParseMarkdownFile(mdPath string) (*file.File, error) {
	bytes, err := os.ReadFile(mdPath)
	if err != nil {
		return nil, err
	}
	s := string(bytes)
	lines := strings.Split(s, "\n")
	f := file.File{}
	for i := 0; i < len(lines); i++ {
		if lines[i] == "## vars" {
			skip, vars := parseVars(lines[i:])
			i += skip
			f.Vars = vars
		} else if lines[i] == "## after" {
			skip, after := parseVars(lines[i+1:])
			i += skip
			f.After = after
		} else if strings.HasPrefix(lines[i], "## type") {
			skip, typ := parseType(lines[i:])
			i += skip
			f.Typ = typ
		}
	}

	return &f, nil
}

func isSection(s string) bool {
	return strings.HasPrefix(s, "## ")
}

func isSubSection(s string) bool {
	return strings.HasPrefix(s, "### ")
}

func parseType(s []string) (int, file.APIType) {
	r := regexp.MustCompile(`## (?P<name>[a-zA-Z0-9_]+)\[(?P<type>[a-zA-Z0-9_]+)\]`)
	matches := r.FindStringSubmatch(s[0])
	typeIndex := r.SubexpIndex("type")
	typ := matches[typeIndex]
	s = s[1:]
	fields := []file.TypedComponent{}
	i := 0
	for ; i < len(s); i++ {
		if isSection(s[i]) {
			break
		}
		if isSubSection(s[i]) {
			skip, v := parseTypedComponent(s[i:])
			i += skip

			fields = append(fields, v)
		}
	}

	return i, file.APIType{
		Typ:    typ,
		Fields: fields,
	}
}

func parseVars(s []string) (int, []file.TypedComponent) {
	s = s[1:]
	vars := []file.TypedComponent{}
	i := 0
	for ; i < len(s); i++ {
		if isSection(s[i]) {
			break
		}
		if isSubSection(s[i]) {
			skip, v := parseTypedComponent(s[i:])
			i += skip
			vars = append(vars, v)
		}
	}

	return i, vars
}

func parseTypedComponent(s []string) (int, file.TypedComponent) {
	r := regexp.MustCompile(`### (?P<name>[a-zA-Z0-9_]+)(?:\[(?P<type>[a-zA-Z0-9_]+)\])?(?::(?P<converter>.+$))?`)
	matches := r.FindStringSubmatch(s[0])
	nameIndex := r.SubexpIndex("name")
	typeIndex := r.SubexpIndex("type")
	converterIndex := r.SubexpIndex("converter")
	varName := matches[nameIndex]
	varType := matches[typeIndex]
	varConvs := matches[converterIndex]

	vals := []file.Value{}
	i := 0
	switch varType {
	case file.ListType:
		skip, multipleVals := parseList(s[1:])
		vals = multipleVals
		i = skip
	case file.TextType:
		skip, singleVal := parseText(s[1:])
		vals = append(vals, singleVal)
		i = skip
	default:
		skip, singleVal := parseText(s[1:])
		vals = append(vals, singleVal)
		i = skip
	}

	convs := strings.Split(varConvs, ":")
	if len(convs) == 1 && convs[0] == "" {
		convs = []string{}
	}

	if varType == "" {
		varType = "text"
	}

	return i, file.TypedComponent{Nam: varName, Typ: varType, Vals: vals, Convs: convs}
}

func parseList(s []string) (int, []file.Value) {
	vals := []file.Value{}
	i := 0
	for ; i < len(s); i++ {
		if isSection(s[i]) || isSubSection(s[i]) {
			break
		}
		if strings.HasPrefix(s[i], "- ") {
			vals = append(vals, file.Value{
				Val: strings.TrimPrefix(s[i], "- "),
				Typ: file.TextType,
			})
		}
	}
	return i, vals
}

func parseText(s []string) (int, file.Value) {
	textType := file.TextType
	found := false
	stringBuilder := strings.Builder{}
	i := 0
	for ; i < len(s); i++ {
		if isSection(s[i]) || isSubSection(s[i]) {
			break
		}
		if strings.HasPrefix(s[i], "```") {
			if !found {
				trimed := strings.TrimSpace(strings.TrimPrefix(s[i], "```"))
				if trimed != "" {
					textType = trimed
				}
				found = true
				continue
			}
			break
		}
		if found {
			stringBuilder.WriteString(s[i])
			stringBuilder.WriteString("\n")
		}
	}
	return i, file.Value{Val: strings.TrimSpace(stringBuilder.String()), Typ: textType}
}
