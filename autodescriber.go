package meep

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type meepAutodescriber interface {
	isMeepAutodescriber() *AutodescribingError
}

// note that this function applies if you have a type which embeds this one sans-*, but you have a ref to that type.
func (m *AutodescribingError) isMeepAutodescriber() *AutodescribingError { return m }

var customDescribe map[reflect.Type]func(reflect.Value, io.Writer) = map[reflect.Type]func(reflect.Value, io.Writer){
	reflect.TypeOf(Meep{}):                func(reflect.Value, io.Writer) {},
	reflect.TypeOf(TraceableError{}):      func(reflect.Value, io.Writer) {},
	reflect.TypeOf(CauseableError{}):      func(reflect.Value, io.Writer) {},
	reflect.TypeOf(GroupingError{}):       func(reflect.Value, io.Writer) {},
	reflect.TypeOf(AutodescribingError{}): func(reflect.Value, io.Writer) {},
}

func (m *AutodescribingError) ErrorMessage() string {
	// Check for initialization.
	// Can't do much of use if we didn't get initialized with a selfie reference.
	if m.self == nil {
		panic("meep:uninitialized")
	}
	// Unwind any pointer indirections.
	rv_self := reflect.ValueOf(m.self)
	for rv_self.Kind() == reflect.Ptr {
		rv_self = rv_self.Elem()
	}
	// Start buffering.
	buf := &bytes.Buffer{}
	// Annouce the type info.
	buf.WriteString("Error[")
	buf.WriteString(rv_self.Type().String())
	buf.WriteString("]: ")
	// Iterate over fields.
	nField := rv_self.NumField()
	for i := 0; i < nField; i++ {
		f := rv_self.Field(i)
		if custom, ok := customDescribe[f.Type()]; ok {
			custom(f, buf)
			continue
		}
		buf.WriteString(rv_self.Type().Field(i).Name)
		buf.WriteByte('=')
		buf.WriteString(fmt.Sprintf("%#v", f.Interface()))
		buf.WriteByte(';')
	}
	// That's it.  Return the buffer results.
	return buf.String()
}

/*
	Implements `error`.

	If you're using other mixins, you may want to override this again;
	if you're just using `Autodescriber`, it'll do fine.
*/
func (m *AutodescribingError) Error() string {
	return m.ErrorMessage()
}
