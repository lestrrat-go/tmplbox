package tmplbox_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/lestrrat/tmplbox"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	f := func(s string) ([]byte, error) {
		return []byte("Hello, World!"), nil
	}

	box := tmplbox.New(tmplbox.AssetSourceFunc(f))
	tmpl, err := box.GetOrCompose("hello.html")
	if !assert.NoError(t, err, "GetOrCompose should succeed") {
		return
	}

	var buf bytes.Buffer
	if !assert.NoError(t, tmpl.Execute(&buf, nil), "Execute() should work") {
		return
	}

	if !assert.Equal(t, "Hello, World!", buf.String()) {
		return
	}
}

func TestComposedTemplates(t *testing.T) {
	f := func(s string) ([]byte, error) {
		t.Logf("Loading %s", s)
		switch s {
		case "hello.html":
			return []byte(`{{ define "content" }}Hello, World!{{ end }}`), nil
		case "base.html":
			return []byte(`{{ define "root" }}content: {{ block "content" . }}{{ end }}{{ end }}`), nil
		}
		return nil, errors.New("not found")
	}

	box := tmplbox.New(tmplbox.AssetSourceFunc(f))
	tmpl, err := box.GetOrCompose("hello.html", "hello.html", "base.html")
	if !assert.NoError(t, err, "GetOrCompose should succeed") {
		return
	}

	var buf bytes.Buffer
	if !assert.NoError(t, tmpl.ExecuteTemplate(&buf, "root", nil), "ExecuteTemplate() should work") {
		return
	}

	if !assert.Equal(t, "content: Hello, World!", buf.String()) {
		return
	}
}
