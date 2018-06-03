package cli

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"
)

type TableFormatter struct {
	//format string
	w *tabwriter.Writer
}

func newTableFormatter(column ...interface{}) *TableFormatter {
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 1, ' ', 0)
	f := &TableFormatter{w}
	f.Row(column...)
	return f
}

func (f *TableFormatter) Row(v ...interface{}) {
	buf := bytes.NewBufferString("")
	for i, s := range v {
		if i > 0 {
			buf.WriteString("\t")
		}
		buf.WriteString(fmt.Sprint(s))
	}
	buf.WriteString("\n")
	f.w.Write(buf.Bytes())
}

func (f *TableFormatter) Flush() {
	f.w.Flush()
}
