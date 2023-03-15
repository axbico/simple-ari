package cli

/// ///

type Cmd struct {
	name        string
	usage       string
	description string
	conf        string
	example     []string
	exec        func(*Cli, ...string) (string, error)
	pcli        *Cli
}

/// ///

func (cmd *Cmd) Name(name string) *Cmd {
	cmd.name = name
	return cmd
}

func (cmd *Cmd) Usage(usage string) *Cmd {
	cmd.usage = usage
	return cmd
}

func (cmd *Cmd) Description(description string) *Cmd {
	cmd.description = description
	return cmd
}

func (cmd *Cmd) AddExample(example string) *Cmd {
	cmd.example = append(cmd.example, example)
	return cmd
}

func (cmd *Cmd) Exec(exec func(*Cli, ...string) (string, error)) *Cmd {
	cmd.exec = exec
	return cmd
}

func (cmd *Cmd) Conf(conf string) *Cmd {
	cmd.conf = conf
	return cmd
}

func (cmd *Cmd) Cli(pcli *Cli) *Cmd {
	cmd.pcli = pcli
	return cmd
}

/// ///

func Help() *Cmd {
	var help *Cmd = new(Cmd)

	help.Name("help")
	help.Usage("help [command]...")
	help.Description("Print commands usage information")

	help.Exec(
		func(cl *Cli, args ...string) (string, error) {

			if len(args) > 1 {
				for _, name := range args[1:] {
					if command, ok := cl.commands[name]; ok {
						cl.Print("[" + command.name + "]")
						if len(command.usage) > 0 {
							cl.Print("Usage: " + command.usage)
						}
						if len(command.description) > 0 {
							cl.Print(command.description)
						}
						if len(command.example) > 0 {
							cl.Print("")
							cl.Print("Examples:")
							for _, example := range command.example {
								cl.Print("  " + example)
							}
						}
						cl.Print("")
					} else {
						cl.Print("Command `" + name + "` is not defined")
					}
				}
			} else {
				cl.Print("Usage: " + cl.commands["help"].usage)
				cl.Print(cl.commands["help"].description)
				cl.Print("")
				cl.Print("List of available commands:")
				for name := range cl.commands {
					cl.Print("  " + name)
				}
				cl.Print("")
			}

			return "", nil
		},
	)

	return help
}

/// ///

func Exit() *Cmd {
	var cmd *Cmd = new(Cmd)

	cmd.Name("exit")

	cmd.Exec(
		func(cl *Cli, args ...string) (string, error) {
			cl.Kill()

			return "", nil
		},
	)

	return cmd
}

/// ///
