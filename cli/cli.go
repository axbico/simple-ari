package cli

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

/// ///

type Cli struct {
	title              string
	commands           map[string]*Cmd
	backgroundCommands []*Cmd
	instream           chan string
	outstream          chan string
	terminate          chan bool
	verbose            bool
	tty                *term
}

/// ///

func (ocli *Cli) AddCommand(cmd *Cmd) *Cli {
	if ocli.commands == nil {
		ocli.commands = make(map[string]*Cmd)
	}

	if _, ok := ocli.commands[cmd.name]; ok {
		panic("Command already defined")
	}

	ocli.commands[cmd.name] = cmd

	return ocli
}

func (ocli *Cli) AddBackgroundCommand(cmd *Cmd) *Cli {
	if cmd.name != "" {
		panic("`" + cmd.name + "` can't run as background command")
	}

	ocli.backgroundCommands = append(ocli.backgroundCommands, cmd)

	return ocli
}

/// ///

func (ocli *Cli) Title(title string) *Cli {
	ocli.title = title
	return ocli
}

func (ocli *Cli) Verbose(v bool) *Cli {
	ocli.verbose = v
	return ocli
}

func (ocli *Cli) Print(output string) *Cli {
	ocli.outstream <- output
	return ocli
}

/// ///

func (ocli *Cli) waitExecute() {
	defer func() {
		ocli.tty.deferSubroutines()
	}()

	for {
		stdin := <-ocli.instream
		args := strings.Split(stdin, " ")

		if _, ok := ocli.commands[args[0]]; !ok {
			ocli.Print("Command `" + args[0] + "` doesn't exist or wrong use, enter `help` to get list of commands and usage")
			continue
		}

		var execArgs []string
		execArgs = append(execArgs, ocli.commands[args[0]].conf)
		execArgs = append(execArgs, args[1:]...)

		go func(ocli *Cli, execArgs []string) {
			start := time.Now()
			if ret, err := ocli.commands[args[0]].exec(ocli, execArgs...); err != nil {
				ocli.Print("!!ERROR[" + args[0] + "]: " + err.Error())
			} else {
				ocli.Print("@" + args[0] + " =" + fmt.Sprintf("%.4fs", time.Since(start).Seconds()) + "] " + ret)
			}
		}(ocli, execArgs)
	}
}

/// ///

func (ocli *Cli) Kill() {
	ocli.terminate <- true
}

/// ///

func (ocli *Cli) waitStdin() {
	defer func() {
		ocli.tty.deferSubroutines()
	}()

	var input []byte = make([]byte, 1)

	for {
		syscall.Read(0, input)

		switch input[0] {
		case '\r', '\n': // cr lf
			ocli.Print("[" + time.Now().Format("15:04:05.00") + "*CLI] " + ocli.tty.inputBuffer)
			ocli.instream <- ocli.tty.inputBuffer

			ocli.tty.resetInputBuffer()

		case '\b', 0x7f: // backspace
			ocli.tty.pushLeftDelete()
		case '\033', 0x5b: // escape sequence start
			ocli.tty.readEscapeSequence()
		default:
			ocli.tty.inputBufferInsert(string(input))
		}

		ocli.tty.flushInputBufferPrompt()
	}
}

/// ///

func (ocli *Cli) waitStdout() {
	defer func() {
		ocli.tty.deferSubroutines()
	}()

	for {
		ocli.tty.printLine(<-ocli.outstream)
	}
}

/// ///

func Run(ocli *Cli) {
	defer func() {
		ocli.tty.deferSubroutines()
	}()

	ocli.instream = make(chan string)
	ocli.outstream = make(chan string)
	ocli.terminate = make(chan bool)

	sigints := make(chan os.Signal, 1)
	signal.Notify(sigints, syscall.SIGINT, syscall.SIGTERM)

	ocli.tty = initiateTerminal()

	ocli.tty.promptIndicator = ocli.title + PROMPT_INDICATOR

	go ocli.waitStdout()

	go func() {
		<-sigints
		ocli.Kill()
	}()

	ocli.Print("ARI Prompt")
	ocli.Print("^^^^^^^^^^")
	ocli.Print("")

	for _, command := range ocli.backgroundCommands {
		go command.exec(ocli, command.conf)
	}

	go ocli.waitStdin()

	go ocli.waitExecute()

	<-ocli.terminate
}

/// ///
