package token

import (
	"fmt"
	"strings"
)

type TokenType int

const (
	TokenNone TokenType = iota
	TokenError
	TokenName
	TokenNumber
	TokenBoolean
	TokenEnumMember
	TokenString
	TokenChunk
	TokenMessageStart
	TokenMessageEnd
	TokenCollectionStart
	TokenCollectionEnd
	TokenComment
	TokenEscapedComment
	TokenTrigger
	TokenLuaBlock
	TokenLuaBlockType
	TokenUnprintable
	TokenRenameIndicator
	TokenEOF
)

var TokenNames = map[TokenType]string{
	TokenNone:            "none",
	TokenError:           "errortoken",
	TokenName:            "name",
	TokenNumber:          "number",
	TokenBoolean:         "bool",
	TokenEnumMember:      "enum_member",
	TokenString:          "string",
	TokenChunk:           "chunk",
	TokenMessageStart:    "message_start",
	TokenMessageEnd:      "message_end",
	TokenCollectionStart: "collection_start",
	TokenCollectionEnd:   "collection_end",
	TokenComment:         "comment",
	TokenEscapedComment:  "escaped_comment",
	TokenTrigger:         "trigger",
	TokenLuaBlock:        "lua_block",
	TokenLuaBlockType:    "lua_block_type",
	TokenUnprintable:     "unprintable",
	TokenRenameIndicator: "rename_indicator",
	TokenEOF:             "end_of_file",
}

type Token struct {
	Type      TokenType
	Value     string
	StartLine int
	StartChar int
	EndLine   int
	EndChar   int
	FromHook  bool
	HookName  string
}

func Tokenize(content []byte) (tokens []Token) {
	l := &lexer{
		data:   content,
		len:    len(content),
		str:    string(content),
		tokens: make([]Token, 0, len(content)/18),
		line:   1,
		char:   1,
	}
	for state := normalState; state != nil; {
		state = state(l)
	}
	return l.tokens
}

type stateFn func(*lexer) stateFn

var (
	isNameStart = [256]bool{}
	isNameChar  = [256]bool{}
	isLower     = [256]bool{}
)

func init() {
	for i := range isNameStart {
		c := byte(i)
		isNameStart[c] = (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_'
		isNameChar[c] = isNameStart[c] || (c >= '0' && c <= '9')
		isLower[c] = c >= 'a' && c <= 'z'
	}
}

func normalState(l *lexer) stateFn {
	for {
		l.start = l.pos
		l.startline = l.line
		l.startchar = l.char
		char := l.peek()
		if char == -1 {
			break
		}
		if isNonPrintable(char) {
			l.move(1)
			l.split(TokenUnprintable)
		}
		if isNameStart[char] {
			return nameState
		}
		if char == '#' {
			l.move(1)
			l.start++
			return commentState
		}
		if (char == '-' || char == '/' || char == ';') &&
			l.pos+1 < l.len && int(l.data[l.pos+1]) == char {
			l.move(2)
			l.start += 2
			return commentState
		}
		if char >= '0' && char <= '9' || char == '+' || char == '-' || char == '.' {
			return numberState
		}
		if char == '\'' {
			l.terminator = '\''
			return stringState
		}
		if char == '"' {
			l.terminator = '"'
			return stringState
		}
		if char == '{' {
			l.move(1)
			l.split(TokenMessageStart)
			l.move(-1)
		}
		if char == '}' {
			l.move(1)
			l.split(TokenMessageEnd)
			l.move(-1)
		}
		if char == '[' {
			l.move(1)
			l.split(TokenCollectionStart)
			l.move(-1)
		}
		if char == ']' {
			l.move(1)
			l.split(TokenCollectionEnd)
			l.move(-1)
		}
		if char == '$' {
			l.move(1)
			l.start++
			l.startchar++
			return multilineStringState
		}
		if char == '@' {
			l.move(1)
			l.start++
			return triggerState
		}
		if char == '>' {
			l.move(1)
			l.split(TokenRenameIndicator)
			l.move(-1)
		}
		if char == '`' {
			return chunkState
		}
		if l.next() == -1 {
			break
		}
	}
	return nil
}

func nameState(l *lexer) stateFn {
	l.acceptRun(nameChars)
	name := l.str[l.start:l.pos]
	tokenType := TokenName
	switch {
	case name == "true" || name == "false":
		tokenType = TokenBoolean
	case name == "NaN":
		tokenType = TokenNumber
	case len(name) > 1 && name != "HEALTH_TYPE" && isAllUpper(name):
		tokenType = TokenEnumMember
	case name == "N", name == "S", name == "E", name == "W":
		tokenType = TokenEnumMember
	}
	l.split(tokenType)
	return normalState
}

func commentState(l *lexer) stateFn {
	idx := strings.IndexByte(l.str[l.pos:], '\n')
	if idx == -1 {
		l.pos = l.len - 1
	} else {
		l.move(idx)
	}
	l.split(TokenComment)
	return normalState
}

func numberState(l *lexer) stateFn {
	l.accept("+-")
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun(digits)
	}
	l.accept("d")
	l.split(TokenNumber)
	return normalState
}

func stringState(l *lexer) stateFn {
	l.move(1)
	l.start++
	escape := false
	for {
		char := l.data[l.pos]
		if char == byte(l.terminator) && !escape {
			break
		}
		switch char {
		case '\n':
			return l.errorf("string not terminated @%d\n", l.pos)
		case '\\':
			escape = !escape
		default:
			escape = false
		}
		l.move(1)
		if l.pos == l.len {
			return l.errorf("string not terminated @%d\n", l.pos)
		}
	}
	l.char++
	l.split(TokenString)
	l.char--
	l.move(1)
	return normalState
}

func multilineStringState(l *lexer) stateFn {
	// TODO check if it's a template reference
	for {
		lineEnd := strings.IndexByte(l.str[l.pos:], '\n')
		// check if it's a template reference

		chunkEnd := strings.Index(l.str[l.pos:], "$end")
		if chunkEnd == -1 {
			return l.errorf("unterminated string @%s\n", l.line)
		}
		if chunkEnd < lineEnd {
			l.pos += chunkEnd + 1
			l.char += chunkEnd + 1
			break
		}
		l.char = lineEnd + 1
		l.pos += lineEnd + 1
		l.line++
	}
	l.pos--
	l.char--
	l.split(TokenChunk)
	l.pos += 4
	l.char += 4
	return normalState
}

func triggerState(l *lexer) stateFn {
	l.acceptRun(lowercaseChars)
	l.split(TokenTrigger)
	return normalState
}

func chunkState(l *lexer) stateFn {
	if strings.HasPrefix(l.str[l.pos:], "``") {
		l.move(3)
		l.start = l.pos
		l.startchar = l.pos + 1
		idx := strings.IndexByte(l.str[l.pos:], '\n')
		if idx == -1 {
			return l.errorf("chunk ends on header line @%d\n", l.pos)
		}
		l.move(idx)
		l.split(TokenLuaBlockType)
		l.move(1)
		l.start = l.pos
		l.startline++
		l.char = 1
		l.startchar = 1
		for {
			lineEnd := strings.IndexByte(l.str[l.pos:], '\n')
			if lineEnd == -1 {
				return l.errorf("unterminated string @%s\n", l.line)
			}
			line := l.str[l.pos : l.pos+lineEnd]
			if line == "```" {
				break
			}
			l.char = lineEnd + 1
			l.pos += lineEnd + 1
			l.line++
		}
		l.pos--
		l.split(TokenLuaBlock)
		l.move(3)
		l.line++
	} else {
		l.move(1)
	}
	return normalState
}

type lexer struct {
	data       []byte
	len        int
	str        string
	start      int
	startline  int
	startchar  int
	pos        int
	tokens     []Token
	line       int
	char       int
	terminator rune
}

func (l *lexer) split(t TokenType) {
	l.tokens = append(l.tokens, Token{
		Type:      t,
		Value:     l.str[l.start:l.pos],
		StartLine: l.startline,
		StartChar: l.startchar,
		EndLine:   l.line,
		EndChar:   l.char - 1,
	})
	l.start = l.pos
}

func (l *lexer) next() (c int) {
	if l.pos >= l.len {
		return -1
	}
	c = int(l.data[l.pos])
	if c == '\n' {
		l.line++
		l.startline++
		l.char = 1
		l.startchar = 1
	} else {
		l.char++
		l.startchar++
	}
	l.pos++
	return
}

func (l *lexer) peek() (c int) {
	if l.pos >= l.len {
		return -1
	}
	return int(l.data[l.pos])
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) move(amount int) {
	l.pos += amount
	l.char += amount
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, rune(l.peek())) {
		l.move(1)
		return true
	}
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, rune(l.peek())) {
		l.move(1)
	}
}

func (l *lexer) errorf(format string, args ...any) stateFn {
	l.tokens = append(l.tokens, Token{
		Type:      TokenError,
		Value:     fmt.Sprintf(format, args...),
		StartLine: l.line,
		StartChar: l.char,
		EndLine:   l.line,
		EndChar:   l.char + 1,
	})
	return nil
}

const nameStartChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
const nameChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_0123456789"
const lowercaseChars = "abcdefghijklmnopqrstuvwxyz"

func isAllUpper(s string) bool {
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			return false
		}
	}
	return true
}

func isNonPrintable(char int) bool {
	return ((char < 32) && !(char >= 9 && char <= 12)) || char >= 127
}
