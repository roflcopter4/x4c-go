package untranslate

/*
#cgo pkg-config: libxml-2.0
#include <string.h>
#include <stdbool.h>
#include <stdio.h>
#include <libxml/HTMLparser.h>
#include <libxml/HTMLtree.h>
#include <libxml/globals.h>
#include <libxml/parser.h>
#include <libxml/parserInternals.h>
#include <libxml/tree.h>
#include <libxml/xmlerror.h>
#include <libxml/xpath.h>
#include <libxml/xpathInternals.h>
#include <libxml/c14n.h>
#include <libxml/xmlschemas.h>

static inline void MY_xmlFree(void *p) {
	xmlFree(p);
}
*/

// func xmlCharToString(s *C.xmlChar) string {
// return C.GoString((*C.char)(unsafe.Pointer(s)))
// }

//========================================================================================

//func init() {
//	spew.Config.Indent = "    "
//}
//
//func Translate(fname string) {
//	b, err := ioutil.ReadFile(fname)
//	if err != nil {
//		panic(err)
//	}
//
//	buf := strings.ReplaceAll(string(b), "\r\n", "\n")
//	tree := parse_buffer(buf)
//
//	// spew.Dump(tree)
//
//	doc := create_xml(tree)
//	out := doc.Dump(true)
//	// out = strings.ReplaceAll(out, "&#10;", "\n")
//	// out := go_to_hell(doc)
//	fmt.Print(out)
//	doc.Free()
//}

// func go_to_hell(doc *XMLdom.Document) string {
//       var (
//             ptr = (*C.xmlDoc)(unsafe.Pointer(doc.Pointer()))
//             i   C.int
//             xc  *C.xmlChar
//       )
//       C.xmlDocDumpMemory(ptr, &xc, &i)
//       defer C.MY_xmlFree(unsafe.Pointer(xc))
//       return xmlCharToString(xc)
// }

//========================================================================================

// func parse_buffer(buf string) ast.AST {
//       p := new(MyParser)
//       p.Buffer = buf
//       p.Pretty = true
//       p.a = ast.NewAst()
//       p.block = p.a.Root()
//       p.cur = p.a.Root()
//
//       lines := strings.Split(p.Buffer, "\n")
//
//       if err := p.Init(); err != nil {
//             die_error(err, lines)
//       }
//
//       err := p.Parse()
//       if err != nil {
//             die_error(err, lines)
//       }
//
//       // p.PrintSyntaxTree()
//       p.Execute()
//       // fmt.Println(p.Tokens())
//
//       return p.a
// }
//
// func getindent(i int) string {
//       return strings.Repeat(" ", i*4)
// }
//
// func die_error(e error, lines []string) {
//       err := e.(*parseError)
//       util.Die(1, "%s", my_parse_error(err, lines))
// }
//
// func my_parse_error(e *parseError, lines []string) string {
//       var (
//             tokens    = []token32{e.max}
//             positions = make([]int, 2*len(tokens))
//             err       = "\n"
//             p         = 0
//       )
//
//       for _, token := range tokens {
//             positions[p] = int(token.begin)
//             p++
//             positions[p] = int(token.end)
//             p++
//       }
//
//       translations := translatePositions(e.p.buffer, positions)
//
//       format := "parse error near \033[34m%v\033[m (line %v symbol %v - line %v symbol %v):\n%v\n"
//
//       for _, token := range tokens {
//             begin, end := int(token.begin), int(token.end)
//             err += fmt.Sprintf(format,
//                   rul3s[token.pegRule],
//                   translations[begin].line, translations[begin].symbol,
//                   translations[end].line, translations[end].symbol,
//                   strconv.Quote(string(e.p.buffer[begin:end])))
//
//             err += fmt.Sprintf("\n"+
//                   "\033[1;31m"+"Parsing Error:"+
//                   "\033[37m"+" near '"+
//                   "\033[36m"+"<HERE>"+
//                   "\033[37m"+"' on line "+
//                   "\033[32m"+"%d"+
//                   "\033[37m"+" col "+
//                   "\033[32m"+"%d"+
//                   "\033[m"+"\n",
//                   translations[begin].line, translations[begin].symbol)
//
//             err += whatever(lines, translations[begin].line-1, translations[begin].symbol+2)
//
//             // err += "\n`" + lines[translations[begin].line-2] + "\n" +
//             //       lines[translations[begin].line-1] + "\n" +
//             //       lines[translations[begin].line] + "\n" +
//             //       lines[translations[begin].line+1] + "\n" +
//             //       lines[translations[begin].line+2] + "`"
//
//       }
//
//       return err
// }
//
// func max(a, b int) int {
//       if a > b {
//             return a
//       } else {
//             return b
//       }
// }
//
// func min(a, b int) int {
//       if a < b {
//             return a
//       } else {
//             return b
//       }
// }
//
// const NLINES = 5
//
// func whatever(lines []string, lineno, column int) string {
//       nlines := min(len(lines), NLINES)
//       ret := ""
//
//       if column > len(lines[lineno])+1 {
//             column = len(lines[lineno]) - 1
//       }
//
//       for i := nlines; i > 0; i-- {
//             if lineno-i < 0 {
//                   continue
//             }
//             ret += "\033[0;33m" + lines[lineno-i] + "\033[0m\n"
//       }
//
//       ret += lines[lineno][:column-1]
//       ret += "\033[1;36m" + "<HERE>" +
//             "\033[0;33m" + lines[lineno][column-1:] +
//             "\033[0m" + "\n"
//
//       if lineno+NLINES+1 > len(lines) {
//             nlines = len(lines) - lineno
//       } else {
//             nlines = NLINES
//       }
//
//       for i := 0; i < nlines; i++ {
//             ret += "\033[0;33m" + lines[lineno+i+1] + "\033[0m\n"
//       }
//
//       return ret
// }
