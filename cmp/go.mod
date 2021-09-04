module github.com/SWAN-community/swan-demo-go/cmp

go 1.17

require (
	github.com/SWAN-community/swan-demo-go/common v0.1.1
	github.com/satori/go.uuid v1.2.0
)

require (
	github.com/google/uuid v1.3.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/SWAN-community/swan-demo-go/common => ../common
replace github.com/SWAN-community/swan-demo-go/fod => ../fod
