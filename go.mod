module github.com/crossworth/daikin

go 1.23.4

require (
	github.com/crossworth/daikin/types v0.0.0-00010101000000-000000000000
	github.com/gosuri/uilive v0.0.4
	github.com/stretchr/testify v1.8.4
	github.com/valyala/fasthttp v1.51.0
)

replace github.com/crossworth/daikin/types => ./types

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
