package xml

////////////////////////////////////////////////////////////////////////////
// Porgram: xmlfmt.go
// Purpose: Go XML Beautify from XML string using pure string manipulation
// Reference: https://github.com/go-xmlfmt/xmlfmt
////////////////////////////////////////////////////////////////////////////

import (
	"regexp"
	"strings"
)

var (
	reg = regexp.MustCompile(`<([/!]?)([^>]+?)(/?)>`)
	// NL is the newline string used in XML output, define for DOS-convenient.
	NL = "\r\n"
)

// FormatXML will (purly) reformat the XML string in a readable way, without any rewriting/altering the structure
func FormatXML(xmls, prefix, indent string) string {
	src := regexp.MustCompile(`(?s)>\s+<`).ReplaceAllString(xmls, "><")

	rf := replaceTag(prefix, indent)
	return (prefix + reg.ReplaceAllStringFunc(src, rf))
}

// replaceTag returns a closure function to do 's/(?<=>)\s+(?=<)//g; s(<(/?)([^>]+?)(/?)>)($indent+=$3?0:$1?-1:1;"<$1$2$3>"."\n".("  "x$indent))ge' as in Perl
// and deal with comments as well
func replaceTag(prefix, indent string) func(string) string {
	indentLevel := 0
	frontIsEnd := false
	return func(m string) string {
		// head elem
		if strings.HasPrefix(m, "<?xml") {
			return NL + prefix + strings.Repeat(indent, indentLevel) + m
		}
		// empty elem
		if strings.HasSuffix(m, "/>") {
			defer func() {
				frontIsEnd = true
			}()

			return NL + prefix + strings.Repeat(indent, indentLevel) + m
		}
		// comment elem
		if strings.HasPrefix(m, "<!") {
			return NL + prefix + strings.Repeat(indent, indentLevel) + m
		}

		// end elem
		if strings.HasPrefix(m, "</") {
			indentLevel--
			defer func() {
				frontIsEnd = true
			}()

			if !frontIsEnd {
				return m
			} else {
				return NL + prefix + strings.Repeat(indent, indentLevel) + m
			}
		}

		defer func() {
			indentLevel++
		}()

		frontIsEnd = false
		return NL + prefix + strings.Repeat(indent, indentLevel) + m
	}
}
