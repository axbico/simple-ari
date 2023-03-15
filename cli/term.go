package cli

// #include <sys/ioctl.h>
// #include <termios.h>
// #include <unistd.h>
//
// typedef struct winsize winsize;
// typedef struct termios termios;
//
// void go_ioctl(int i, unsigned long l, winsize * ws) { ioctl(i, l, ws); }
//
import "C"

import (
	"fmt"
	"syscall"
)

const TIOCGWINSZ C.ulong = 0x5413

/// ///

const BUFFER_SIZE = 64

const (
	PROMPT_INDICATOR = "> "
	CSI              = "\033["
)

const (
	ALTSCRBUFF_ENABLE  = CSI + "?1049h"
	ALTSCRBUFF_DISABLE = CSI + "?1049l"
	DECTCEM_ENABLE     = CSI + "?25h"
	DECTCEM_DISABLE    = CSI + "?25l"
	CNL_1              = CSI + "E"
	CPL_1              = CSI + "1F"
	CUP_0_0            = CSI + "H"
	ED_0               = CSI + "0J"
	EL_2               = CSI + "2K"
	CURSOR_ENABLE      = CSI + "7m" + CSI + "1m"
	CURSOR_DISABLE     = CSI + "0m"
)

/// ///

type term struct {
	cursor           int
	deferSubroutines func()
	inputBuffer      string
	windowHeight     uint16
	promptIndicator  string
}

/// ///

func initiateTerminal() *term {

	var t *term = new(term)

	var ws C.winsize
	C.go_ioctl(0, TIOCGWINSZ, &ws)
	t.windowHeight = uint16(ws.ws_row)

	var canonical C.termios
	C.tcgetattr(0, &canonical)

	noncanonical := canonical
	noncanonical.c_lflag &^= syscall.ICANON | syscall.ECHO
	C.tcsetattr(0, C.TCSANOW, &noncanonical)

	fmt.Print(ALTSCRBUFF_ENABLE + CUP_0_0 + DECTCEM_DISABLE)

	t.deferSubroutines = func() {
		fmt.Print(DECTCEM_ENABLE + ALTSCRBUFF_DISABLE + ED_0)

		C.tcsetattr(0, C.TCSANOW, &canonical)
	}

	t.promptIndicator = PROMPT_INDICATOR

	return t
}

/// ///

func (t *term) inputBufferInsert(input string) {
	if len(t.inputBuffer) != BUFFER_SIZE {
		t.inputBuffer = t.inputBuffer[:t.cursor] + input + t.inputBuffer[t.cursor:]
		t.cursor++
	}
}

func (t *term) pushLeftDelete() {
	if t.cursor != 0 {
		t.cursor--
		t.inputBuffer = t.inputBuffer[:t.cursor] + t.inputBuffer[t.cursor+1:]
	}
}

func (t *term) pushRightDelete() {
	if t.cursor < len(t.inputBuffer) {
		t.inputBuffer = t.inputBuffer[:t.cursor] + t.inputBuffer[t.cursor+1:]
	}
}

/// ///

func (t *term) readEscapeSequence() {
	var input []byte = make([]byte, 1)
	if _, err := syscall.Read(0, input); err == nil {
		switch input[0] {
		case 0x5b: // CSI sequence
			if _, err := syscall.Read(0, input); err == nil {
				switch input[0] {
				case 0x44: // left arrow
					if t.cursor > 0 {
						t.cursor--
					}
				case 0x43: // right arrow
					if t.cursor < len(t.inputBuffer) {
						t.cursor++
					}
				case 0x33: // del
					if _, err := syscall.Read(0, input); err == nil && input[0] == 0x7e {
						t.pushRightDelete()
					}
				}
			}
		}
	}
}

/// ///

func (t *term) resetInputBuffer() {
	t.inputBuffer = ""
	t.cursor = 0
}

func (t *term) flushInputBufferPrompt() {

	formatBuffer := t.inputBuffer + " "

	fmt.Print(
		CPL_1 + CNL_1 + EL_2 + t.promptIndicator +
			formatBuffer[:t.cursor] +
			CURSOR_ENABLE + string(formatBuffer[t.cursor]) + CURSOR_DISABLE +
			formatBuffer[t.cursor:][1:],
	)
}

func (t *term) printLine(line string) {

	fmt.Print(CPL_1 + CNL_1 + EL_2 + line + "\n")

	t.flushInputBufferPrompt()
}

/// ///
