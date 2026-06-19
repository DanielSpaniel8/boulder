package token

import (
	"slices"
)

// lookForToken searches forward from start for a token matching target,
// skipping comments. Returns (found, value, distance).
func LookForToken(tokens []Token, target any, start int) (bool, Token, int) {
	if len(tokens) < start+1 {
		return false, Token{}, 0
	}

	var targets []TokenType
	switch v := target.(type) {
	case TokenType:
		targets = []TokenType{v}
	case []TokenType:
		targets = v
	default:
		return false, Token{}, 0
	}

	skippable := []TokenType{TokenComment, TokenEscapedComment}

	for i := start + 1; i < len(tokens); i++ {
		tok := tokens[i]
		if slices.Contains(targets, tok.Type) {
			return true, tok, i - start
		}
		if !slices.Contains(skippable, tok.Type) {
			break
		}
	}
	return false, Token{}, 0
}
