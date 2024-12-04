
See https://github.com/golang/go/issues/70681 for further information/discussion.

Seems `text/template` behaves different if a map is defined other than `map[string]string`

The expected behavior is that if a map key is missing the template should not print/substitute 
anything if `.Option` is called with `"missingkey=zero"`, not printing `<no value>`

Package `html/template` on the other hand behaves correct.

See docs for both libraries...

https://pkg.go.dev/text/template#Template.Option
https://pkg.go.dev/html/template#Template.Option
