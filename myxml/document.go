package myxml

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	XML "github.com/lestrrat-go/libxml2"
	XMLtypes "github.com/lestrrat-go/libxml2/types"
	XMLxsd "github.com/lestrrat-go/libxml2/xsd"
	"github.com/pkg/errors"

	"github.com/roflcopter4/x4c/util"
)

type DocWrapper interface {
	Doc() XMLtypes.Document
	Schema() *XMLxsd.Schema
	Name() *fileName
	ScriptName() string
	GetSchema() error
	SetSchema(*XMLxsd.Schema)
	ValidateSchema() error
	Free()
}

type fileName struct {
	base string
	path string
	full string
}

type schema struct {
	xml_schema *XMLxsd.Schema
	url        string
	filename   string
}

// Implements DocWrapper
type document struct {
	doc        XMLtypes.Document
	schema     schema
	filename   fileName
	scriptname string
	lines      []string
}

// Create a new DocWrapper from a filename
func New_Document(fname string) (DocWrapper, error) {
	var (
		err      error
		xml_file *os.File
		doc      XMLtypes.Document
		lines    []string
		fileinfo *fileName
	)

	if fileinfo, err = get_filename_info(fname); err != nil {
		return nil, errors.Wrap(err, "Failed to resolve path")
	}

	docstring := badly_escape_line_breaks(fileinfo.full)
	// fmt.Print(1, "%s", docstring)
	// os.Exit(1)

	// if xml_file, err = os.Open(fileinfo.full); err != nil {
	//       return nil, err
	// }
	// if doc, err = XML.ParseReader(xml_file); err != nil {
	//       return nil, errors.WithMessage(err, "XML parse failed")
	// }
	if doc, err = XML.ParseString(docstring); err != nil {
		return nil, errors.WithMessage(err, "XML parse failed")
	}

	doc.MakeMortal()
	xml_file.Close()

	if lines, err = read_lines(fileinfo.full); err != nil {
		return nil, errors.WithStack(err)
	}

	docwrap := &document{
		schema:     schema{},
		filename:   *fileinfo,
		doc:        doc,
		scriptname: "",
		lines:      lines,
	}

	return docwrap, nil
}

func get_filename_info(fname string) (*fileName, error) {
	var err error
	if fname, err = filepath.Abs(fname); err != nil {
		return nil, errors.Wrap(err, "Failed to resolve path")
	}
	return &fileName{
		base: filepath.Base(fname),
		path: filepath.Dir(fname),
		full: fname,
	}, nil
}

// Free the document and schema if applicable
func (d *document) Free() {
	if d.doc != nil {
		d.doc.Free()
		d.doc = nil
	}
	if d.schema.xml_schema != nil {
		// d.schema.references--
		// if d.schema.references == 0 {
		//       d.schema.xml_schema.Free()
		// }
		d.schema.xml_schema.Free()
		d.schema.xml_schema = nil
	}
}

//========================================================================================
// Boring accessor methods

// -------------
// type document

// Return the XML document
func (d *document) Doc() XMLtypes.Document {
	return d.doc
}

// Return the filename struct
func (d *document) Name() *fileName {
	return &d.filename
}

// Return the md or aiscript name
func (d *document) ScriptName() string {
	return d.scriptname
}

// Return the compiled xsd schema
func (d *document) Schema() *XMLxsd.Schema {
	return d.schema.xml_schema
}

// -------------
// type fileName

// Return the base filename (as with basename(1))
func (fname *fileName) Base() string {
	return fname.base
}

// Return the file path (as with dirname(1))
func (fname *fileName) Path() string {
	return fname.path
}

// Return the full file path
func (fname *fileName) Full() string {
	return fname.full
}

//========================================================================================
// Util

// Reads a whole file into memory and returns a slice of its lines.
func read_lines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Can you tell that this is mostly directly pasted from a C program?
func badly_escape_line_breaks(fname string) string {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		util.DieE(1, err)
	}

	var (
		orig         = string(b)
		repl         = ""
		esc          = false
		is_attribute = false
		is_comment   = false
		is_string    = false
	)

	for i, l := 0, utf8.RuneCountInString(orig); i < l; i++ {
	skip_increment:
		ch := orig[i]

		if is_comment {
			// Check for the end of the comment
			if strings.HasPrefix(orig[i:], "-->") {
				is_comment = false
				repl += "-->"
				// I think this is clearer than `i+=2; continue` and
				// then having the for loop increment again
				i += 3
				goto skip_increment
			} else {
				repl += string(ch)
			}

			continue
		}

		switch ch {
		case '<':
			if strings.HasPrefix(orig[i+1:], "!--") {
				is_comment = true
				repl += "<!--"
				i += 4
				goto skip_increment
			}

		case '\\':
			esc = true
			repl += string(ch)
			continue

		case '\'':
			if is_attribute && !esc {
				is_string = !is_string
			}

		case '"':
			if !is_string {
				is_attribute = !is_attribute
			}

		case '\r':
			i++
			ch = orig[i]
			if ch != '\n' {
				util.Die(1, "Stray \\r in file...")
			}
			fallthrough
		case '\n':
			if is_attribute {
				repl += "&#xA;"
				esc = false
				continue
			}
		}

		repl += string(ch)
		esc = false
	}

	return repl
}
