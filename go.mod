module pdflabs

go 1.25.0

// require github.com/dslipak/pdf v0.0.2
// require github.com/signintech/gopdf v0.36.0

// require (
// 	github.com/phpdave11/gofpdi v1.0.14-0.20211212211723-1f10f9844311 // indirect
// 	github.com/pkg/errors v0.8.1 // indirect
// )

require (
	github.com/dslipak/pdf v0.0.2
	github.com/pkg/errors v0.8.1
	github.com/signintech/gopdf v0.36.0
)

replace github.com/signintech/gopdf => ./gopdf

replace github.com/dslipak/pdf => ./pdf

replace github.com/pkg/errors => ./errors
