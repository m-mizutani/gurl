package gurl

type Method int

const (
	MethodGET Method = iota + 1
	MethodHEAD
	MethodPOST
	MethodPUT
	MethodDELETE
	MethodCONNECT
	MethodOPTIONS
	MethodTRACE
	MethodPATCH
)

var methodMap = map[Method]string{
	MethodGET:     "GET",
	MethodHEAD:    "HEAD",
	MethodPOST:    "POST",
	MethodPUT:     "PUT",
	MethodDELETE:  "DELETE",
	MethodCONNECT: "CONNECT",
	MethodOPTIONS: "OPTIONS",
	MethodTRACE:   "TRACE",
	MethodPATCH:   "PATCH",
}

func Get(uri string, options ...Option) error {
	return New(MethodGET, uri, options...).Emit()
}

func Head(uri string, options ...Option) error {
	return New(MethodHEAD, uri, options...).Emit()
}

func Post(uri string, options ...Option) error {
	return New(MethodPOST, uri, options...).Emit()
}

func Put(uri string, options ...Option) error {
	return New(MethodPUT, uri, options...).Emit()
}

func Delete(uri string, options ...Option) error {
	return New(MethodDELETE, uri, options...).Emit()
}

func Connect(uri string, options ...Option) error {
	return New(MethodCONNECT, uri, options...).Emit()
}

func Options(uri string, options ...Option) error {
	return New(MethodOPTIONS, uri, options...).Emit()
}

func Trace(uri string, options ...Option) error {
	return New(MethodTRACE, uri, options...).Emit()
}

func Patch(uri string, options ...Option) error {
	return New(MethodPATCH, uri, options...).Emit()
}
