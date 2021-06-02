package httpoption

type Validator interface {
	Validate(value interface{}) error
}

// FuncValidator is an adapter to allow the use of ordinary functions as
// validators. If f is a function with the appropriate signature,
// Func(f) is a Validator that calls f.
type FuncValidator func(value interface{}) error

// Validate calls fv(value).
func (fv FuncValidator) Validate(value interface{}) error {
	return fv(value)
}
