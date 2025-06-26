package errors

var _ error = Const("")

type Const string

func (c Const) Error() string {
	return string(c)
}
