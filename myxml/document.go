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
		err       error
		xml_file  *os.File
		doc       XMLtypes.Document
		docstring string
		lines     []string
		fileinfo  *fileName
	)

	if fileinfo, err = get_filename_info(fname); err != nil {
		return nil, errors.Wrap(err, "Failed to resolve path")
	}

	// docstring := badly_escape_line_breaks(fileinfo.full)

	if lines, docstring, err = badly_escape_line_breaks(fileinfo.full); err != nil {
		return nil, err
	}

	// util.Eprintf("%s\n\n", docstring)
	// os.Exit(1)

	// if doc, err = XML.ParseReader(xml_file); err != nil {
	//       return nil, errors.WithMessage(err, "XML parse failed")
	// }
	if doc, err = XML.ParseString(docstring); err != nil {
		// docstring += "hi"
		// if doc, err = XML.ParseString(strings.Join(lines, "")); err != nil {
		return nil, errors.WithMessage(err, "XML parse failed")
	}

	doc.MakeMortal()
	xml_file.Close()

	// if lines, err = read_lines(fileinfo.full); err != nil {
	//       return nil, err
	// }

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
	if fname == "-" {
		return &fileName{"<STDIN>", "<STDIN>", "<STDIN>"}, nil
	} else {
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
func (d *document) Doc() XMLtypes.Document { return d.doc }

// Return the filename struct
func (d *document) Name() *fileName { return &d.filename }

// Return the md or aiscript name
func (d *document) ScriptName() string { return d.scriptname }

// Return the compiled xsd schema
func (d *document) Schema() *XMLxsd.Schema { return d.schema.xml_schema }

// -------------
// type fileName

// Return the base filename (as with basename(1))
func (fname *fileName) Base() string { return fname.base }

// Return the file path (as with dirname(1))
func (fname *fileName) Path() string { return fname.path }

// Return the full file path
func (fname *fileName) Full() string { return fname.full }

//========================================================================================
// Util

// Reads a whole file into memory and returns a slice of its lines.
func read_lines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
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
func badly_escape_line_breaks(fname string) ([]string, string, error) {
	var (
		b   []byte
		err error
	)

	if fname == "<STDIN>" {
		if b, err = ioutil.ReadAll(os.Stdin); err != nil {
			return nil, "", errors.WithStack(err)
		}
	} else {
		if b, err = ioutil.ReadFile(fname); err != nil {
			return nil, "", errors.WithStack(err)
		}
	}

	var (
		orig = strings.ReplaceAll(string(b), "\r\n", "\n")
		repl = ""
		// esc          = false
		is_attribute = false
		is_comment   = false
		// is_string    = false
	)

	for i, l := 0, utf8.RuneCountInString(orig); i < l; i++ {
		ch := orig[i]

		if is_comment {
			// Check for the end of the comment
			if strings.HasPrefix(orig[i:], "-->") {
				is_comment = false
				repl += "-->"
				i += 3 - 1
				continue
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
				i += 4 - 1
				continue
			}

		// case '\\':
		//       esc = true
		//       repl += string(ch)
		//       continue

		//case '\'':
		//	if is_attribute && !esc {
		//		is_string = !is_string
		//	}

		case '"':
			// if !is_string {
			is_attribute = !is_attribute
			// }

		case '\n':
			if is_attribute {
				repl += "&#xA;"
				// esc = false
				continue
			}
		}

		repl += string(ch)
		// esc = false
	}

	return strings.Split(orig, "\n"), repl, nil
}
