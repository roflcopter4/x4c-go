package myxml

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/roflcopter4/x4c/config"

	XMLtypes "github.com/lestrrat-go/libxml2/types"
	XMLxsd "github.com/lestrrat-go/libxml2/xsd"
)

func (d *document) GetSchema() error {
	var (
		sch_fname string
		sch       *XMLxsd.Schema
		err       error
	)

	if sch_fname, err = d.get_schema_file(); err != nil {
		return err
	}
	if sch, err = XMLxsd.ParseFromFile(sch_fname); err != nil {
		return err
	}

	d.schema.filename = sch_fname
	d.schema.xml_schema = sch

	return nil
}

func (d *document) SetSchema(sch *XMLxsd.Schema) {
	d.schema = schema{xml_schema: sch}
}

func (d *document) get_schema_file() (string, error) {
	var (
		fname string = ""
		found bool   = false
	)

	err := d.doc.Walk(
		func(node XMLtypes.Node) error {
			switch n := node.(type) {
			case XMLtypes.Element:
				attributes, err := n.Attributes()
				if err != nil {
					return err
				}
				for _, attr := range attributes {
					switch attr.NodeName() {
					case "xsi:noNamespaceSchemaLocation":
						fname = attr.NodeValue()
						found = true
					case "name":
						switch n.NodeName() {
						case "aiscript", "mdscript":
							d.scriptname = attr.NodeValue()
						}
					}
				}
			}

			if found {
				return errors.New("")
			}
			return nil
		},
	)

	switch {
	case err != nil:
		return "", errors.New(fmt.Sprintf("item `XMLtypes.Document' `docwrap.d' returned error:\n%v", err))
	case !found:
		return "", errors.New("Failed to find XML Schema file name")
	case fname == "":
		return "", errors.New("Filename is empty")
	}

	var ret string

	if _, err = os.Stat(config.Ini_Data["libraries"]); err == nil {
		ret, err = filepath.Abs(filepath.Join(config.Ini_Data["libraries"], fname))
	} else {
		ret, err = filepath.Abs(filepath.Join(d.filename.path, fname))
	}
	if err != nil {
		return "", err
	}
	if _, err = os.Stat(ret); err != nil {
		return "", errors.New(fmt.Sprintf("Schema file does not exist:\n%v", err))
	}

	return ret, nil
}

func (d *document) ValidateSchema() error {
	return d.schema.validate(d.doc)
}

func (s schema) validate(doc XMLtypes.Document) error {
	// if !s.initialized {
	//       return errors.New("Uninitialized schema")
	// }
	if s.xml_schema == nil {
		return errors.New("Invalid schema")
	}

	if err := s.xml_schema.Validate(doc); err != nil {
		switch errtype := err.(type) {
		default:
			panic(err)

		case XMLxsd.SchemaValidationError:
			str := fmt.Sprintf("%v\n", errtype)
			for i, scherr := range errtype.Errors() {
				str += fmt.Sprintf("  %3d:\t%v\n", i+1, scherr)
			}
			return errors.New(str)
		}
	}

	return nil

}
