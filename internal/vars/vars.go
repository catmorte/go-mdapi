package vars

import "strings"

const (
	currentDir  = "CURDIR"
	currentFile = "CURFILE"
	resultDir   = "RESULTDIR"
)

type Vars map[string]string

func ReplacePatterns(text string, allFields map[string]string) string {
	for variableName, value := range allFields {
		placeholder := "{{" + variableName + "}}"
		text = strings.ReplaceAll(text, placeholder, value)
	}
	return text
}

func (v Vars) GetCurrentDir() string {
	return v[currentDir]
}

func (v Vars) GetCurrentFile() string {
	return v[currentFile]
}

func (v Vars) GetResultDir() string {
	return v[resultDir]
}

func (v Vars) SetCurrentDir(dir string) {
	v[currentDir] = dir
}

func (v Vars) SetCurrentFile(file string) {
	v[currentFile] = file
}

func (v Vars) SetResultDir(dir string) {
	v[resultDir] = dir
}
