package meep

import (
	"fmt"
	"reflect"
)

type meepAutodescriber interface {
	isMeepAutodescriber() *AutodescribingError
}

// note that this function applies if you have a type which embeds this one sans-*, but you have a ref to that type.
func (m *AutodescribingError) isMeepAutodescriber() *AutodescribingError { return m }

var dontDescribe map[reflect.Type]struct{} = map[reflect.Type]struct{}{
	reflect.TypeOf(Meep{}):                struct{}{},
	reflect.TypeOf(StackableError{}):      struct{}{},
	reflect.TypeOf(CauseableError{}):      struct{}{},
	reflect.TypeOf(GroupingError{}):       struct{}{},
	reflect.TypeOf(AutodescribingError{}): struct{}{},
}

func (m *AutodescribingError) ErrorMessage() string {
	if m.self == nil {
		// can't do much of use if we didn't get initialized with a selfie reference.
		panic("uninitialized")
	}
	rv_self := reflect.ValueOf(m.self)
	for rv_self.Kind() == reflect.Ptr {
		rv_self = rv_self.Elem()
	}
	rv_self.CanInterface()
	nField := rv_self.NumField()
	content := ""
	for i := 0; i < nField; i++ {
		f := rv_self.Field(i)
		if _, ok := dontDescribe[f.Type()]; ok {
			continue
		}
		content += rv_self.Type().Field(i).Name + "=" + fmt.Sprintf("%#v", f.Interface()) + "; "
	}
	return fmt.Sprintf(
		"Error[%s]: %s",
		rv_self.Type(),
		content,
	)
}

/*
	Implements `error`.

	If you're using other mixins, you may want to override this again;
	if you're just using `Autodescriber`, it'll do fine.
*/
func (m *AutodescribingError) Error() string {
	return m.ErrorMessage()
}
