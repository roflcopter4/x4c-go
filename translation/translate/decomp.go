package translate

import (
	"github.com/roflcopter4/x4c-go/util/colour"
)

const _WANT_TEXT_NODES = false

var xs_eids = []string{
	// "actions",
	// "append_to_list",
	// "assert",
	// "break",
	// "check_any",
	// "check_value",
	// "create_list",
	// "debug_text",
	// "do_all",
	// "do_any",
	"do_else",
	"do_elseif",
	"do_if",
	// "do_while",
	// "handler",
	// "input_param",
	// "interrupt",
	// "interrupts",
	// "label",
	// "order",
	// "param",
	// "remove_value",
	// "resume",
	// "return",
	// "run_script",
	// "set_value",
	// "add_faction_relation",
	// "add_relation_boost",
	// "attention",
	// "conditions",
	// "create_order",
	// "init",
	// "interrupt_after_time",
	// "match_content",
	// "match_distance",
	// "move_strafe",
	// "on_abort",
	// "position",
	// "replace",
	// "rotation",
	// "save_retval",
	// "set_command",
	// "set_command_action",
	// "set_order_syncpoint_reached",
	// "show_notification",
	// "shuffle_list",
	// "signal_objects",
	// "skill",
	// "substitute_text",
	// "unknown",
	// "wait",
	// "write_to_logbook",
}

var fmtstr = "" +
	colour.BGreen("node:") + " '%v`,  " +
	colour.BGreen("type:") + " `%v',  " +
	colour.BGreen("value:") + " %v"

// func Dumb(doc myxml.DocWrapper, outfp *os.File) error {
//       sort.Strings(xs_eids)
//
//       strlst := make([]string, 0, 1024)
//
//       analyze := func(node XMLtypes.Node) error {
//             var (
//                   str string
//                   err error
//             )
//             switch n := node.(type) {
//             case XMLtypes.Element:
//                   str, err = handle_element(node, n)
//             default:
//                   str, err = handle_default(node)
//             }
//
//             if err != nil {
//                   return err
//             }
//             if str != "" {
//                   strlst = append(strlst, str)
//             }
//
//             return nil
//       }
//
//       if err := doc.Doc().Walk(analyze); err != nil {
//             errors.Wrap(err, "item `XMLtypes.Document' `doc' returned error:")
//             return err
//       }
//
//       for i, s := range strlst {
//             fmt.Fprintf(outfp, "% 4d: %s\n", i, s)
//       }
//
//       return nil
// }

//========================================================================================
// UNUSED !!

// func handle_element(node XMLtypes.Node, elem XMLtypes.Element) (string, error) {
//       var (
//             attributes []XMLtypes.Attribute
//             nn         = elem.NodeName()
//             i          = sort.SearchStrings(xs_eids, elem.NodeName())
//       )
//
//       if attributes, err := elem.Attributes(); err != nil {
//             return "", err
//       }
//
//       if i < len(xs_eids) && xs_eids[i] == nn {
//             // Handle special recognized keywords
//             switch nn {
//             case "do_if":
//                   return fmt.Sprintf("if (%v) {", attributes), nil
//             case "do_elseif":
//                   return fmt.Sprintf("elseif (%v) {", attributes), nil
//             case "do_else":
//                   return fmt.Sprintf("else {"), nil
//             default:
//                   panic("Impossible!")
//             }
//       } else {
//             // handle_generic()
//             // if len(attributes) == 0 {
//             //       // TODO Handle element with no attributes
//             // } else {
//             //       // TODO Handle element with attributes
//             // }
//             //
//             // return fmt.Sprintf(fmtstr, nn, elem.NodeType(), attributes), nil
//       }
// }
//
// func handle_default(node XMLtypes.Node) (string, error) {
//       var str string
//
//       switch node.NodeType() {
//       case XMLclib.TextNode:
//             // Handled differently. This is typically just whitespace.
//             if _WANT_TEXT_NODES {
//                   str = get_node_text(node)
//                   str = strings.ReplaceAll(str, "\r\n", "\\n")
//                   str = strings.ReplaceAll(str, "\n", "\\n")
//                   return fmt.Sprintf("`%v'", str), nil
//             } else {
//                   return "", nil
//             }
//
//       case XMLclib.CommentNode:
//             str = "`" + node.NodeValue() + "'"
//
//       default:
//             panic("unhandled")
//       }
//
//       return fmt.Sprintf(fmtstr, node.NodeName(), node.NodeType(), str), nil
// }
//
// func get_node_text(node XMLtypes.Node) string {
//       str := strings.SplitAfter(node.String(), ">\n")[0]
//       if len(str) > 1 && str[len(str)-1] == '\n' {
//             str = str[:len(str)-1]
//       }
//       return str
// }
