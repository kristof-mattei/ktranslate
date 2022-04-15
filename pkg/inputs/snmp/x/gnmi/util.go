package gnmi

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	gnmipb "github.com/openconfig/gnmi/proto/gnmi"
)

// Taken from https://github.com/openconfig/ygot to avoid annoying import.

// StringToStructuredPath takes a string representing a path, and converts it to
// a gnmi.Path, using the PathElem element message that is defined in gNMI 0.4.0.
func stringToStructuredPath(path string) (*gnmipb.Path, error) {
	parts := pathStringToElements(path)

	gpath := &gnmipb.Path{}
	for _, p := range parts {
		name, kv, err := extractKV(p)
		if err != nil {
			return nil, fmt.Errorf("error parsing path %s: %v", path, err)
		}
		gpath.Elem = append(gpath.Elem, &gnmipb.PathElem{
			Name: name,
			Key:  kv,
		})
	}
	return gpath, nil
}

// extractKV extracts key value predicates from the input string in. It returns
// the name of the element, a map keyed by key name with values of the predicates
// specified. It removes escape characters from keys and values where they are
// specified.
func extractKV(in string) (string, map[string]string, error) {
	var inEscape, inKey, inValue bool
	var name, currentKey string
	var buf bytes.Buffer
	keys := map[string]string{}

	for _, ch := range in {
		switch {
		case ch == '[' && !inEscape && !inValue && inKey:
			return "", nil, fmt.Errorf("received an unescaped [ in key of element %s", name)
		case ch == '[' && !inEscape && !inKey:
			inKey = true
			if len(keys) == 0 {
				if buf.Len() == 0 {
					return "", nil, errors.New("received a value when the element name was null")
				}
				name = buf.String()
				buf.Reset()
			}
			continue
		case ch == ']' && !inEscape && !inKey:
			return "", nil, fmt.Errorf("received an unescaped ] when not in a key for element %s", buf.String())
		case ch == ']' && !inEscape:
			inKey = false
			inValue = false
			if err := addKey(keys, name, currentKey, buf.String()); err != nil {
				return "", nil, err
			}
			buf.Reset()
			currentKey = ""
			continue
		case ch == '\\' && !inEscape:
			inEscape = true
			continue
		case ch == '=' && inKey && !inEscape && !inValue:
			currentKey = buf.String()
			buf.Reset()
			inValue = true
			continue
		}

		buf.WriteRune(ch)
		inEscape = false
	}

	if len(keys) == 0 {
		name = buf.String()
	}

	if len(keys) != 0 && buf.Len() != 0 {
		// In this case, we have trailing garbage following the key.
		return "", nil, fmt.Errorf("trailing garbage following keys in element %s, got: %v", name, buf.String())
	}

	if strings.Contains(name, " ") {
		return "", nil, fmt.Errorf("invalid space character included in element name '%s'", name)
	}

	return name, keys, nil
}

// PathStringToElements splits the string s, which represents a gNMI string
// path into its constituent elements. It does not parse keys, which are left
// unchanged within the path - but removes escape characters from element
// names. The path returned omits any leading or trailing empty elements when
// splitting on the / character.
func pathStringToElements(path string) []string {
	parts := splitPath(path)
	// Remove leading empty element
	if len(parts) > 0 && parts[0] == "" {
		parts = parts[1:]
	}
	// Remove trailing empty element
	if len(parts) > 0 && path[len(path)-1] == '/' {
		parts = parts[:len(parts)-1]
	}
	return parts
}

// addKey adds key k with value v to the key's map. The key, value pair is specified
// to be for an element named e.
func addKey(keys map[string]string, e, k, v string) error {
	switch {
	case strings.Contains(k, " "):
		return fmt.Errorf("received an invalid space in element %s key name '%s'", e, k)
	case e == "":
		return fmt.Errorf("received null element value with key and value %s=%s", k, v)
	case k == "":
		return fmt.Errorf("received null key name for element %s", e)
	case v == "":
		return fmt.Errorf("received null value for key %s of element %s", k, e)
	}
	keys[k] = v
	return nil
}

// SplitPath splits path across unescaped /.
// Any / inside square brackets are ignored.
func splitPath(path string) []string {
	var parts []string
	var buf bytes.Buffer

	var inKey, inEscape bool

	var ch rune
	for _, ch = range path {
		switch {
		case ch == '[' && !inEscape:
			inKey = true
		case ch == ']' && !inEscape:
			inKey = false
		case ch == '\\' && !inEscape && !inKey:
			inEscape = true
			continue
		case ch == '/' && !inEscape && !inKey:
			parts = append(parts, buf.String())
			buf.Reset()
			continue
		}

		buf.WriteRune(ch)
		inEscape = false
	}

	if buf.Len() != 0 || (len(path) != 1 && ch == '/') {
		parts = append(parts, buf.String())
	}

	return parts
}
