package types

import (
	"errors"
	"fmt"
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
		NewAPI() string
	}
	DefinedTypes []DefinedType
)

func GetDefinedTypes() (DefinedTypes, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cfgPath := filepath.Join(dirname, ".config", "go-mdapi")
	types := []DefinedType{
		internalHTTPTemplate,
	}
	err = filepath.Walk(cfgPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		folderName := filepath.Base(path)
		runTemplatePath := filepath.Join(cfgPath, folderName, "run.tmpl")
		newAPITemplatePath := filepath.Join(cfgPath, folderName, "new_api.md")
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
		runTemplate, err := os.ReadFile(runTemplatePath)
		if err != nil {
			return err
		}

		newAPITemplate, err := os.ReadFile(newAPITemplatePath)
		if err != nil {
			return err
		}
		types = append(types, externalType{
			Name:           folderName,
			RunTemplate:    string(runTemplate),
			NewAPITemplate: string(newAPITemplate),
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
