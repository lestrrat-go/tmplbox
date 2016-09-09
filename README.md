# go-tmplbox

[![Build Status](https://travis-ci.org/lestrrat/go-tmplbox.svg?branch=master)](https://travis-ci.org/lestrrat/go-tmplbox)

[![GoDoc](https://godoc.org/github.com/lestrrat/go-tmplbox?status.svg)](https://godoc.org/github.com/lestrrat/go-tmplbox)

Yet Another (text|html)/template wrapper

# SYNOPSIS

```html
<!-- base.html -->
{{ define "root" }}
<html>
<body>
{{ block "content" }}{{ end }}
</body>
</html>
{{ end }}
```

```html
<!-- index.html -->
{{ define "content" }}Hello, World!{{ end }}
```

```go
// assuming you're using go-bindata's Asset function to
// retrieve template sources here.
var box := tmplbox.New(tmplbox.AssetSourceFunc(Asset))

func indexHandler(w http.ResponseWriter, r *http.Response) {
    // Assuming you have a dependency between index.html and
    // base.html (see above), GetOrCompose will automatically
    // compile and merge all of the specified templates.
    t, _ := box.GetOrCompose("index.html", "base.html")

    // Will output 
    // <!-- base.html -->
    // <html>
    // <body>
    // Hello, World!
    // </body>
    // </html>
    t.ExecuteTemplate(w, "root", nil)
}
```

# FEATURES

* Plays well with go-bindata
* Easy to setup template dependency to emulate "extends", without another big library
