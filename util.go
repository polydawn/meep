package meep

import (
	"bufio"
	"io"
)

var (
	tab = []byte{'\t'}
	br  = []byte{'\n'}
)

func indenter(w io.Writer) io.Writer {
	return &rediscipliner{
		edgeDfn: bufio.ScanLines,
		wr:      w,
		prefix:  tab,
		suffix:  br,
	}
}

var _ io.Writer = &rediscipliner{}

type rediscipliner struct {
	edgeDfn bufio.SplitFunc
	wr      io.Writer
	prefix  []byte
	suffix  []byte

	rem []byte
}

func (red *rediscipliner) Write(b []byte) (int, error) {
	red.rem = append(red.rem, b...)
	n := 0
	//fmt.Printf("\n>>>>>\n")
	for len(red.rem) > 0 { // if loop until the buffer is exhausted, or another cond breaks out
		adv, tok, err := red.edgeDfn(red.rem, false)
		//fmt.Printf(">>>\n\tbuf: %q\n\ttok: %q\n\tadv: %d\n", red.rem, tok, adv)
		if err != nil {
			return n, err
		}
		if adv == 0 { // when we no longer have a full chunk, return
			return n, nil
		}
		red.wr.Write(red.prefix)
		red.wr.Write(tok)
		red.wr.Write(red.suffix)
		n += adv
		red.rem = red.rem[adv:]
	}
	return -1, nil
}
