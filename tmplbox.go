package tmplbox

import (
	htemplate "html/template"
	"strings"
	ttemplate "text/template"

	"github.com/pkg/errors"
)

func New(assets AssetSource) *Box {
	return &Box{
		assets:    assets,
		templates: make(map[string]Template),
	}
}

type FuncMap map[string]interface{}

// Funcs specifies the set of template functions that will be
// applied to all templates that are compiled
func (b *Box) Funcs(fm FuncMap) *Box {
	b.funcs = fm
	return b
}

// Compose creates a template instance using the templates
// specified in `name` and `deps`.
func (b *Box) Compose(name string, deps ...string) (Template, error) {
	var t Template
	if strings.HasSuffix(name, ".html") {
		t = newHTML(name, b.funcs)
	} else {
		t = newText(name, b.funcs)
	}

	tmpls := append([]string{name}, deps...)
	for _, tname := range tmpls {
		tmp, err := b.Compile(tname)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to compile template %s", tname)
		}
		t, err = addParseTree(t, tmp)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to add parse tree for template %s", tname)
		}
	}
	return t, nil
}

func definedTemplates(t Template) []string {
	s := t.DefinedTemplates()
	l := strings.Split(s[len(`; defined templates are: `):], ", ")
	for i, n := range l {
		l[i] = n[1 : len(n)-1]
	}

	return l
}

func addParseTree(t1, t2 Template) (Template, error) {
	var err error
	switch t1.(type) {
	case *htemplate.Template:
		t1h := t1.(*htemplate.Template)
		t2h := t2.(*htemplate.Template)
		for _, n := range definedTemplates(t2) {
			if tmp := t2h.Lookup(n); tmp != nil {
				t1h, err = t1h.AddParseTree(n, tmp.Tree)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to add template %s", n)
				}
			}
		}
		return t1h, nil
	case *ttemplate.Template:
		t1h := t1.(*ttemplate.Template)
		t2h := t2.(*ttemplate.Template)
		for _, n := range definedTemplates(t2) {
			if tmp := t2h.Lookup(n); tmp != nil {
				t1h, err = t1h.AddParseTree(n, tmp.Tree)
				if err != nil {
					return nil, errors.Wrapf(err, "failed to add template %s", n)
				}
			}
		}
		return t1h, nil
	default:
		return nil, errors.New("unknown template type")
	}
}

type newFunc func(string) Template

func newHTML(s string, fm FuncMap) *htemplate.Template {
	t := htemplate.New(s)
	if fm != nil {
		t = t.Funcs(htemplate.FuncMap(fm))
	}
	return t
}
func newText(s string, fm FuncMap) *ttemplate.Template {
	t := ttemplate.New(s)
	if fm != nil {
		t = t.Funcs(ttemplate.FuncMap(fm))
	}
	return t
}

type compileFunc func(string, []byte, FuncMap) (Template, error)

func compileHTML(s string, b []byte, fm FuncMap) (Template, error) {
	return newHTML(s, fm).Parse(string(b))
}
func compileText(s string, b []byte, fm FuncMap) (Template, error) {
	return newText(s, fm).Parse(string(b))
}

// Compile compiles the template specified the given name.
// If the template name has a ".html" suffix, html/template is used.
// Otherwise text/template is assumed.
func (b *Box) Compile(name string) (Template, error) {
	var compileFn compileFunc
	if strings.HasSuffix(name, ".html") {
		compileFn = compileHTML
	} else {
		compileFn = compileText
	}

	buf, err := b.assets.Get(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get source for specified template")
	}

	return compileFn(name, buf, b.funcs)
}

// Get returns the template associated with asset `name`.
// An error is returned if the template does not exist.
func (b *Box) Get(name string) (Template, error) {
	b.storageMutex.RLock()
	t, ok := b.templates[name]
	b.storageMutex.RUnlock()

	if !ok {
		return nil, errors.New("template not found")
	}
	return t, nil
}

func (b *Box) Set(name string, t Template) error {
	b.storageMutex.Lock()
	b.templates[name] = t
	b.storageMutex.Unlock()
	return nil
}

// GetOrCompose is like Get, except that if the template
// named by `name` does not exist already, it will call
// Compose to generate it. Otherwise this methods turns a
// previously cached copy.
func (b *Box) GetOrCompose(name string, deps ...string) (Template, error) {
	t, err := b.Get(name)
	if err == nil {
		return t, nil
	}

	t, err = b.Compose(name, deps...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compose new template")
	}

	b.Set(name, t)
	return t, nil
}
