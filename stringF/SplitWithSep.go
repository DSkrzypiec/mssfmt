package stringF

import "bytes"

// SplitWithSep function splits the given string by some predefined characters
// but in outcome separators are also included.
// This function is used for split the "important" part of a word but also not
// loosing also important whitespaces and special characters.
//
// Examples:
// SplitWithSep("x,max()") returns ["x", ",", "max", "(", ")"]
// SplitWithSep("min, 	max()") returns ["min", ",", "\t", "max", "(", ")"]
//
// It is especially useful while formatting keywords. In case when keywords are
// connected with some special characters we want to format only the keyword but
// we don't want to loose the prior content.
func SplitWithSep(w string) []string {
	result := make([]string, 0, len(w))
	var c rune
	var buff bytes.Buffer

	for id, char := range w {
		c = rune(char)
		if buff.Len() == 0 && (c == ',' || c == '\t' || c == '\n' ||
			c == ' ' || c == '(' || c == ')' || c == '\r') {
			result = append(result, string(c))
			continue
		}

		if c == ',' || c == '\t' || c == '\n' || c == ' ' ||
			c == '(' || c == ')' || c == '\r' {
			result = append(result, buff.String(), string(c))
			buff.Reset()
			continue
		}

		if id == len(w)-1 { // last character which is not a separator
			buff.WriteRune(rune(char))
			result = append(result, buff.String())
			buff.Reset()
			continue
		}

		buff.WriteRune(rune(char))
	}

	return result
}
