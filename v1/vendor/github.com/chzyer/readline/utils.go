package readline

import (
	"bufio"
	"bytes"
	"strconv"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	isWindows = false
)

// IsTerminal returns true if the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	return terminal.IsTerminal(fd)
}

func MakeRaw(fd int) (*terminal.State, error) {
	return terminal.MakeRaw(fd)
}

func Restore(fd int, state *terminal.State) error {
	err := terminal.Restore(fd, state)
	if err != nil {
		// errno 0 means everything is ok :)
		if err.Error() == "errno 0" {
			err = nil
		}
	}
	return nil
}

func IsPrintable(key rune) bool {
	isInSurrogateArea := key >= 0xd800 && key <= 0xdbff
	return key >= 32 && !isInSurrogateArea
}

// translate Esc[X
func escapeExKey(r rune, reader *bufio.Reader) rune {
	switch r {
	case 'D':
		r = CharBackward
	case 'C':
		r = CharForward
	case 'A':
		r = CharPrev
	case 'B':
		r = CharNext
	case 'H':
		r = CharLineStart
	case 'F':
		r = CharLineEnd
	default:
		if r == '3' && reader != nil {
			d, _, _ := reader.ReadRune()
			if d == '~' {
				r = CharDelete
			} else {
				reader.UnreadRune()
			}
		}
	}
	return r
}

// translate EscX to Meta+X
func escapeKey(r rune) rune {
	switch r {
	case 'b':
		r = MetaBackward
	case 'f':
		r = MetaForward
	case 'd':
		r = MetaDelete
	case CharTranspose:
		r = MetaTranspose
	case CharBackspace:
		r = MetaBackspace
	case CharEsc:

	}
	return r
}

func SplitByLine(start, screenWidth int, rs []rune) []string {
	var ret []string
	buf := bytes.NewBuffer(nil)
	currentWidth := start
	for _, r := range rs {
		w := runes.Width(r)
		currentWidth += w
		buf.WriteRune(r)
		if currentWidth >= screenWidth {
			ret = append(ret, buf.String())
			buf.Reset()
			currentWidth = 0
		}
	}
	ret = append(ret, buf.String())
	return ret
}

// calculate how many lines for N character
func LineCount(screenWidth, w int) int {
	r := w / screenWidth
	if w%screenWidth != 0 {
		r++
	}
	return r
}

func IsWordBreak(i rune) bool {
	switch {
	case i >= 'a' && i <= 'z':
	case i >= 'A' && i <= 'Z':
	case i >= '0' && i <= '9':
	default:
		return true
	}
	return false
}

func GetInt(s []string, def int) int {
	if len(s) == 0 {
		return def
	}
	c, err := strconv.Atoi(s[0])
	if err != nil {
		return def
	}
	return c
}

type RawMode struct {
	state *terminal.State
}

func (r *RawMode) Enter() (err error) {
	r.state, err = MakeRaw(GetStdin())
	return err
}

func (r *RawMode) Exit() error {
	if r.state == nil {
		return nil
	}
	return Restore(GetStdin(), r.state)
}
