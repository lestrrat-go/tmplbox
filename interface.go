package tmplbox

import (
	"io"
	"sync"
)

type Box struct {
	assets       AssetSource
	funcs        FuncMap
	storageMutex sync.RWMutex
	templates    map[string]Template
}

type AssetSource interface {
	Get(string) ([]byte, error)
}

type Template interface {
	DefinedTemplates() string
	Execute(io.Writer, interface{}) error
	ExecuteTemplate(io.Writer, string, interface{}) error
	Name() string
}
