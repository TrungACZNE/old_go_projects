package cliparser

import "fmt"

func putCaret(orgString, errMsg string, pos int) error {
	formatStr := fmt.Sprintf("%%s\n%%s\n%%%vs", pos+1)
	return fmt.Errorf(formatStr, errMsg, orgString, "^")
}

func Parse(params string) (string, []string, error) {
	var result []string

	token := ""
	openParen := false
	openParenPos := -1
	for i := 0; i < len(params); i++ {
		ch := params[i]
		switch ch {
		case '\\':
			i++
			if i >= len(params) {
				return "", nil, putCaret(params, `Unexpected "\" `, i)
			} else {
				ch := params[i]
				switch ch {
				case '"':
					token += "\""
				case '\\':
					token += "\\"
				default:
					return "", nil, putCaret(params, `Unexpected "\" `, i)
				}
			}
		case '"':
			if openParen {
				result = append(result, token)
				token = ""
				openParen = false
			} else {
				if token != "" {
					return "", nil, putCaret(params, `Missing space before open parenthesis?`, i)
				}
				openParen = true
				openParenPos = i
			}
		case ' ':
			if openParen == true {
				token += string(ch)
			} else if token != "" {
				result = append(result, token)
				token = ""
			}
		default:
			token += string(ch)
		}
	}

	if openParen {
		return "", nil, putCaret(params, `Unterminated open parenthesis`, openParenPos)
	}

	if token != "" {
		result = append(result, token)
	}

	switch len(result) {
	case 0:
		return "", nil, fmt.Errorf("Could not parse command name")
	case 1:
		return result[0], []string{}, nil
	default:
		return result[0], result[1:], nil
	}
}

func ParseCommandless(params string) ([]string, error) {
	command, args, err := Parse(params)
	return append([]string{command}, args...), err
}
