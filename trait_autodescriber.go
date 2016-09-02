package meep

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

// Errors that generate their messages automatically from their fields!
type TraitAutodescribing struct {
	self interface{}
}

type meepAutodescriber interface {
	isMeepAutodescriber() *TraitAutodescribing
}

// note that this function applies if you have a type which embeds this one sans-*, but you have a ref to that type.
func (m *TraitAutodescribing) isMeepAutodescriber() *TraitAutodescribing { return m }

func (m *TraitAutodescribing) ErrorMessage() string {
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
		if consumed, fn := customDescribe(f.Type()); consumed {
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
		inspection := buf.Bytes() // fortunately this is copy free with go slices
		hasTrailingBreak := inspection[len(inspection)-1] == '\n'
		if !hasTrailingBreak {
			buf.WriteByte('\n')
		}
	}
	for _, fn := range custom {
		fn()
		// customs are expected to finish with one trailing \n apiece
	}
}

func customDescribe(typ reflect.Type) (consumed bool, desc func(reflect.Value, io.Writer)) {
	switch typ {
	case reflect.TypeOf(Meep{}):
		return true, func(f reflect.Value, buf io.Writer) {
			describeFields(f, buf.(*bytes.Buffer))
		}
	case reflect.TypeOf(TraitTraceable{}):
		return true, func(f reflect.Value, buf io.Writer) {
			buf = indenter(buf)
			buf.Write([]byte("Stack trace:\n"))
			buf = indenter(buf)
			//buf.(*rediscipliner).prefix = []byte{} // stacks already tab themselves in once
			m := reflect.Indirect(f).Interface().(TraitTraceable)
			m.WriteStack(buf)
		}
	case reflect.TypeOf(TraitCausable{}):
		return true, func(f reflect.Value, buf io.Writer) {
			m := reflect.Indirect(f).Interface().(TraitCausable)
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
		}
	case reflect.TypeOf(TraitAutodescribing{}):
		return true, nil
	default:
		return false, nil
	}
}

/*
	Implements `error`.

	If you're using other mixins, you may want to override this again;
	if you're just using `Autodescriber`, it'll do fine.
*/
func (m *TraitAutodescribing) Error() string {
	return m.ErrorMessage()
}
