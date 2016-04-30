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
	reflect.TypeOf(Meep{}):           func(reflect.Value, io.Writer) {},
	reflect.TypeOf(TraceableError{}): func(reflect.Value, io.Writer) {},
	reflect.TypeOf(CauseableError{}): func(f reflect.Value, buf io.Writer) {
		buf.Write([]byte("\n\tCaused by: "))
		buf.Write([]byte(f.FieldByName("Cause").Interface().(error).Error()))
	},
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
	// If we hit any customs, save em; they serialize after other fields.
	nField := rv_self.NumField()
	var custom []func()
	for i := 0; i < nField; i++ {
		f := rv_self.Field(i)
		if fn, ok := customDescribe[f.Type()]; ok {
			custom = append(custom, func() { fn(f, buf) })
			continue
		}
		buf.WriteString(rv_self.Type().Field(i).Name)
		buf.WriteByte('=')
		buf.WriteString(fmt.Sprintf("%#v", f.Interface()))
		buf.WriteByte(';')
	}
	// Now go back and let the customs have their say.
	for _, fn := range custom {
		fn()
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
