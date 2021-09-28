package util

type Optional struct {
	val interface{}
}

func NewOptional(val interface{}) Optional {
	return Optional{val}
}

func (op Optional) OrElseGet(fallback func() interface{}) interface{} {
	if op.val != nil {
		return op.val
	}
	return fallback()
}
