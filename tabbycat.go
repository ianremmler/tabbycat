// Package tabbycat is a wrapper around text/tabwriter which ignores the width
// of any text matching a given regular expression.  This can be used to
// properly tabbify text containing non-printing ANSI terminal control codes,
// for example.  Note that text/tabwriter's ability to filter HTML tags is
// always enabled since this mechanism is exploited to achieve this package's
// purpose.
package tabbycat

import (
	"bytes"
	"io"
	"regexp"
	"text/tabwriter"
)

// Writer wraps text/tabwriter, ignoring the width of specified text when
// tabbifying.
type Writer struct {
	tabwriter *tabwriter.Writer
	output    io.Writer
	ignore    *regexp.Regexp
	tagged    *regexp.Regexp
	buf       *bytes.Buffer
}

// NewWriter allocates and initializes a new tabbycat.Writer.
// The parameters are the same as for the Init function.
func NewWriter(ignore string, output io.Writer, minWidth, tabWidth, padding int, padChar byte, flags uint) *Writer {
	return new(Writer).Init(ignore, output, minWidth, tabWidth, padding, padChar, flags)
}

// Init initializes a Writer.  The ignore parameters is a regular expression
// that specifies which text to ignore when tabbifying.  The other parameters
// are identical to those of text/tabwriter, with the exception of the
// FilterHTML flag being always enabled.
func (w *Writer) Init(ignore string, output io.Writer, minWidth, tabWidth, padding int, padChar byte, flags uint) *Writer {
	var err error
	if w.ignore, err = regexp.Compile(ignore); err != nil {
		return nil
	}
	if w.tagged, err = regexp.Compile("<!--(" + ignore + ")-->"); err != nil {
		return nil
	}
	w.output = output
	w.buf = &bytes.Buffer{}
	w.tabwriter = tabwriter.NewWriter(w.buf, minWidth, tabWidth, padding, padChar, flags|tabwriter.FilterHTML)

	return w
}

// Write wraps text to be ignored in HTML comment tags to prevent its width
// from being considered, then writes it to the underlying tabwriter.
func (w *Writer) Write(buf []byte) (int, error) {
	return w.tabwriter.Write(w.ignore.ReplaceAll(buf, []byte("<!--$0-->")))
}

// Flush flushes the underlying tabwriter, extracts ignored text from the HTML
// comment tags, and writes the result.
func (w *Writer) Flush() error {
	if err := w.tabwriter.Flush(); err != nil {
		return err
	}
	_, err := w.output.Write(w.tagged.ReplaceAll(w.buf.Bytes(), []byte("$1")))

	return err
}
