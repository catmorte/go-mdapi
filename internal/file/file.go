package file

import (
	"fmt"

	"github.com/catmorte/go-mdapi/internal/command"
	"github.com/catmorte/go-mdapi/internal/converters"
	varsPkg "github.com/catmorte/go-mdapi/internal/vars"
)

const (
	TextType             = "text"
	ListType             = "list"
	FileListType         = "file_list"
	AbsoluteFileListType = "abs_file_list"
	ScriptType           = "script"
)

var typesDescriptions = map[string]string{
	TextType:             "simple text type withing ``` ```",
	ListType:             "one of the values in md list format (- value)",
	FileListType:         "file path from .md dir from which the value will be selected (the value won't be computed, value is the real path to the file)",
	AbsoluteFileListType: "absolute file path from which the value will be selected (the value won't be computed, value is the real path to the file)",
	ScriptType:           "same as text, but the content will be executed in sh",
}

func GetSupportedTypes() []string {
	return []string{TextType, ListType, ScriptType, FileListType, AbsoluteFileListType}
}

func GetTypeDescription(key string) (string, error) {
	d, ok := typesDescriptions[key]
	if !ok {
		return "", fmt.Errorf("unknown type %s", key)
	}

	return d, nil
}

type (
	File struct {
		Dir   string
		Nam   string
		Vars  TypedComponents
		Typ   APIType
		After TypedComponents
	}
	APIType struct {
		Typ    string
		Fields TypedComponents
	}
	TypedComponent struct {
		Nam   string
		Typ   string
		Convs []string
		Vals  []Value
	}
	Value struct {
		Val string
		Typ string
	}
	TypedComponents []TypedComponent
)

func (f File) Compute(vars varsPkg.Vars) (varsPkg.Vars, error) {
	err := f.Vars.Compute(vars, false)
	if err != nil {
		return nil, err
	}
	return vars, nil
}

func (ts TypedComponents) Compute(vars varsPkg.Vars, forceCompute bool) error {
	var err error
	for _, v := range ts {
		val, ok := vars[v.Nam]
		if !ok || forceCompute {
			val, err = v.Compute(vars)
			if err != nil {
				return err
			}
		} else {
			err = v.Validate(val)
			if err != nil {
				return err
			}
		}

		val, err := converters.Convert(val, v.Convs)
		if err != nil {
			return fmt.Errorf("failed to convert value %s: %w", val, err)
		}
		vars[v.Nam] = val
	}
	return nil
}

func (t TypedComponent) Compute(vars varsPkg.Vars) (string, error) {
	var val string
	var err error
	switch t.Typ {
	case TextType:
		fallthrough
	case ListType:
		val = varsPkg.ReplacePatterns(t.Vals[0].Val, vars)
	case ScriptType:
		val = varsPkg.ReplacePatterns(t.Vals[0].Val, vars)
		val, err = command.RunCommand(val)
		if err != nil {
			return "", fmt.Errorf("failed to run command %s: %w", val, err)
		}
	}
	return converters.Convert(val, t.Convs)
}

func (t TypedComponent) Validate(val string) error {
	switch t.Typ {
	case ListType:
		found := false
		for _, varVal := range t.Vals {
			if varVal.Val == val {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown value %s", val)
		}
	}
	return nil
}

func (f File) GetVarByName(name string) (TypedComponent, bool) {
	for _, v := range f.Vars {
		if v.Nam == name {
			return v, true
		}
	}
	return TypedComponent{}, false
}
