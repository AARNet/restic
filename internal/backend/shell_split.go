package backend

import (
	"unicode"

	"github.com/restic/restic/internal/errors"
)

// shellSplitter splits a command string into separater arguments. It supports
// single and double quoted strings.
type shellSplitter struct {
	quote    rune
	lastChar rune
}

func (s *shellSplitter) isSplitChar(c rune) bool {
	// only test for quotes if the last char was not a backslash
	if s.lastChar != '\\' {

		// quote ended
		if s.quote != 0 && c == s.quote {
			s.quote = 0
			return true
		}

		// quote starts
		if s.quote == 0 && (c == '"' || c == '\'') {
			s.quote = c
			return true
		}
	}

	s.lastChar = c

	// within quote
	if s.quote != 0 {
		return false
	}

	// outside quote
	return c == '\\' || unicode.IsSpace(c)
}

// SplitShellArgs returns the list of arguments from a shell command string.
func SplitShellArgs(data string) (cmd string, args []string, err error) {
	s := &shellSplitter{}

	// derived from strings.SplitFunc
	fieldStart := -1 // Set to -1 when looking for start of field.
	for i, rune := range data {
		if s.isSplitChar(rune) {
			if fieldStart >= 0 {
				args = append(args, data[fieldStart:i])
				fieldStart = -1
			}
		} else if fieldStart == -1 {
			fieldStart = i
		}
	}
	if fieldStart >= 0 { // Last field might end at EOF.
		args = append(args, data[fieldStart:])
	}

	switch s.quote {
	case '\'':
		return "", nil, errors.New("single-quoted string not terminated")
	case '"':
		return "", nil, errors.New("double-quoted string not terminated")
	}

	if len(args) == 0 {
		return "", nil, errors.New("command string is empty")
	}

	cmd, args = args[0], args[1:]

	return cmd, args, nil
}
