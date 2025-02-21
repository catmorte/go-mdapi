package types

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/catmorte/go-mdapi/internal/command"
	"github.com/catmorte/go-mdapi/internal/vars"
)

type externalType struct {
	Name           string
	RunTemplate    string
	NewAPITemplate string
}

func (d externalType) GetName() string {
	return d.Name
}

func (d externalType) NewAPI() string {
	return d.NewAPITemplate
}

type TemplateData struct {
	Vars map[string]string
}

func (d externalType) Run(vrs vars.Vars) error {
	td := TemplateData{
		Vars: vrs,
	}

	t, err := template.ParseGlob(d.RunTemplate)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, td)
	if err != nil {
		return err
	}
	output, err := command.RunCommand(tpl.String())
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
