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
	reflect.TypeOf(Meep{}): func(f reflect.Value, buf io.Writer) {
		//describeFields(f, buf.(*bytes.Buffer))
	},
	reflect.TypeOf(TraceableError{}): func(f reflect.Value, buf io.Writer) {
		buf = indenter(buf)
		buf.Write([]byte("Stack trace:\n"))
		buf = indenter(buf)
		//buf.(*rediscipliner).prefix = []byte{} // stacks already tab themselves in once
		m := reflect.Indirect(f).Interface().(TraceableError)
		m.WriteStack(buf)
	},
	reflect.TypeOf(CausableError{}): func(f reflect.Value, buf io.Writer) {
		m := reflect.Indirect(f).Interface().(CausableError)
		if m.Cause == nil {
			return
		}
		buf = indenter(buf)
		buf.Write([]byte("Caused by: "))
		// since we're now in multiline mode, we want to wrap up with a br.
		msg := []byte(m.Cause.Error())
		buf.Write(msg)
		if len(msg) == 0 || msg[len(msg)-1] != '\n' {
			buf.Write(br)
		}
	},
	reflect.TypeOf(GroupingError{}): func(f reflect.Value, buf io.Writer) {
		m := reflect.Indirect(f).Interface().(GroupingError)
		if m.Specifically == nil {
			return
		}
		buf = indenter(buf)
		buf.Write([]byte("Specifically: "))
		// since we're now in multiline mode, we want to wrap up with a br.
		msg := []byte(m.Specifically.Error())
		buf.Write(msg)
		if len(msg) == 0 || msg[len(msg)-1] != '\n' {
			buf.Write(br)
		}
	},
	reflect.TypeOf(AutodescribingError{}): nil,
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
	buf.WriteString("]:")
	// Iterate over fields.
	describeFields(rv_self, buf)
	// That's it.  Return the buffer results.
	return buf.String()
}

func describeFields(subject reflect.Value, buf *bytes.Buffer) {
	// Iterate over fields.
	// If we hit any customs, save em; they serialize after other fields.
	nField := subject.NumField()
	havePrintedFields := false
	var custom []func()
	for i := 0; i < nField; i++ {
		f := subject.Field(i)
		// if it's one of the special/multiliners, stack it up for later
		if fn, ok := customDescribe[f.Type()]; ok {
			if fn != nil {
				custom = append(custom, func() { fn(f, buf) })
			}
			continue
		}
		// if it's a regular field, print the field=value pair
		if havePrintedFields == false {
			buf.WriteByte(' ')
		}
		havePrintedFields = true
		buf.WriteString(subject.Type().Field(i).Name)
		buf.WriteByte('=')
		buf.WriteString(fmt.Sprintf("%#v", f.Interface()))
		buf.WriteByte(';')
	}
	// Now go back and let the customs have their say.
	// (If there are any: Start with a clean line; we're now an ML result.)
	if len(custom) > 0 {
		buf.WriteByte('\n')
	}
	for _, fn := range custom {
		fn()
		// customs are expected to finish with one trailing \n apiece
	}
}

/*
	Implements `error`.

	If you're using other mixins, you may want to override this again;
	if you're just using `Autodescriber`, it'll do fine.
*/
func (m *AutodescribingError) Error() string {
	return m.ErrorMessage()
}
