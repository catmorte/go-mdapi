package types

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/catmorte/go-mdapi/internal/command"
	"github.com/catmorte/go-mdapi/internal/vars"
)

type externalType struct {
	Name           string
	RunTemplate    string
	NewAPITemplate string
	Vars           string
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
	t, err := template.New(d.Name).Funcs(templateFuncs).Parse(d.RunTemplate)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, vrs)
	if err != nil {
		return err
	}
	_, err = command.RunCommand(tpl.String())
	if err != nil {
		return err
	}
	return nil
}

func (d externalType) Compile(vrs vars.Vars) error {
	t, err := template.New(d.Name).Funcs(templateFuncs).Parse(d.RunTemplate)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, vrs)
	if err != nil {
		return err
	}
	fmt.Println(tpl.String())
	return nil
}

func (d externalType) GetVars() []string {
	return strings.Split(d.Vars, "\n")
}
