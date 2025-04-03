module htmx-templ

go 1.23.2

require (
	github.com/a-h/templ v0.2.793
	github.com/alexedwards/scs/v2 v2.8.0
	github.com/fullpipe/icu-mf v0.99.99
	golang.org/x/text v0.23.0
)

replace github.com/fullpipe/icu-mf => ../..

require (
	github.com/alecthomas/participle/v2 v2.1.4 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
