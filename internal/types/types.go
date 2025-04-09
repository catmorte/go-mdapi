package types

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/catmorte/go-mdapi/internal/vars"
)

var ErrNotExist = errors.New("no type defined")

type (
	DefinedType interface {
		GetName() string
		Run(vars.Vars) error
		Compile(vars.Vars) error
		NewAPI() string
		GetVars() []string
	}
	DefinedTypes []DefinedType
)

func GetDefinedTypes(cfgPath string) (DefinedTypes, error) {
	types := InternalTypes()

	info, err := os.Lstat(cfgPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return types, nil
		}
		return nil, err
	}

	if info.Mode()&os.ModeSymlink != 0 {
		cfgPath, err = filepath.EvalSymlinks(cfgPath)
		if err != nil {
			return nil, err
		}
	}

	err = filepath.Walk(cfgPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		folderName := filepath.Base(path)
		runTemplatePath := filepath.Join(cfgPath, folderName, "run.tmpl")
		newAPITemplatePath := filepath.Join(cfgPath, folderName, "new_api.md")
		varsPath := filepath.Join(cfgPath, folderName, "vars")
		_, err = os.Stat(runTemplatePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		_, err = os.Stat(newAPITemplatePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		_, err = os.Stat(varsPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		runTemplate, err := os.ReadFile(runTemplatePath)
		if err != nil {
			return err
		}

		newAPITemplate, err := os.ReadFile(newAPITemplatePath)
		if err != nil {
			return err
		}

		vars, err := os.ReadFile(varsPath)
		if err != nil {
			return err
		}
		types = append(types, externalType{
			Name:           folderName,
			RunTemplate:    string(runTemplate),
			NewAPITemplate: string(newAPITemplate),
			Vars:           string(vars),
		})
		return nil
	})

	return types, err
}

func (dts DefinedTypes) FindByName(name string) (DefinedType, error) {
	for _, v := range dts {
		if v.GetName() == name {
			return v, nil
		}
	}
	return nil, ErrNotExist
}

func InternalTypes() []DefinedType {
	return []DefinedType{
		internalHTTPTemplate,
		internalShTemplate,
	}
}
