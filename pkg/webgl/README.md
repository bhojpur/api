# Bhojpur API - WebGL Access Library

The [GopherJS](https://github.com/gopherjs/gopherjs) bindings for [WebGL 1.0](https://www.khronos.org/registry/webgl/specs/latest/1.0/) context.

## Example

Look into the source code of `internal/web/main.go`:

```go
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/bhojpur/api/pkg/webgl"
)

func main() {
	document := js.Global.Get("document")
	canvas := document.Call("createElement", "canvas")
	document.Get("body").Call("appendChild", canvas)

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false

	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		js.Global.Call("alert", "Error: "+err.Error())
	}

	gl.ClearColor(0.8, 0.3, 0.01, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
```

And, the landing page (i.e., `index.html` file)

```html
<html><body><script src="main.js"></script></body></html>
```

To produce `main.js` file, run `gopherjs build internal/main.go`.