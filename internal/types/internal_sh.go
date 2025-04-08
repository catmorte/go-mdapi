package types

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/catmorte/go-mdapi/internal/command"
	"github.com/catmorte/go-mdapi/internal/vars"
)

type internalSh string

const (
	InternalSHScriptField FieldVar = "script"
)

//go:embed internal_sh_new_api.md
var internalShTemplate internalSh

func (d internalSh) GetName() string {
	return "sh"
}

func (d internalSh) NewAPI() string {
	return string(internalShTemplate)
}

func (d internalSh) Run(vrs vars.Vars) error {
	resultDir := vrs.GetResultDir()
	bodyFile := filepath.Join(resultDir, "body")

	script, ok := InternalSHScriptField.Get(vrs)
	if !ok {
		return fmt.Errorf("missing script field")
	}

	body, err := command.RunCommand(script)
	if err != nil {
		return fmt.Errorf("error reading body: %s", err)
	}
	err = os.WriteFile(bodyFile, []byte(body), 0x775)
	if err != nil {
		return fmt.Errorf("error writing body: %w", err)
	}

	return nil
}

func (d internalSh) Compile(vrs vars.Vars) error {
	fmt.Println("not supported for internal commands")
	return nil
}

func (d internalSh) GetVars() []string {
	return []string{
		string(InternalSHScriptField),
	}
}
