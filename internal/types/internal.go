package types

import (
	"github.com/catmorte/go-mdapi/internal/vars"
)

type FieldVar string

func (f FieldVar) Get(vrs vars.Vars) (string, bool) {
	v, ok := vrs[string(f)]
	return v, ok
}
