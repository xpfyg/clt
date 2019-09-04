package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

type Commands []*Command

type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool
}

func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (me *Commands) Run() {
	flag.Usage = me.usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		me.usage()
	}

	if args[0] == "help" {
		me.help(args[1:])
		return
	}

	for _, cmd := range []*Command(*me) {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			cmd.Run(cmd, args)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "autocash: unknown subcommand %q\nRun 'autocash help' for usage.\n", args[0])
	os.Exit(2)
}

var usageTemplate = `
commands server - connect everysites!

Usage:

    autocash command [arguments]

The commands are:
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "autocash help [command]" for more information about a command.

Additional help topics:
{{range .}}{{if not .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "autocash help [topic]" for more information about that topic.

`

var helpTemplate = `{{if .Runnable}}usage: autocash {{.UsageLine}}
{{end}}{{.Long}}
`

func (me *Commands) usage() {
	me.tmpl(os.Stdout, usageTemplate, []*Command(*me))
	os.Exit(2)
}

func (me *Commands) tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func (me *Commands) help(args []string) {
	if len(args) == 0 {
		me.usage()
		// not exit 2: succeeded at 'go help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stdout, "usage: autocash help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'autocash help'
	}

	arg := args[0]

	for _, cmd := range []*Command(*me) {
		if cmd.Name() == arg {
			me.tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'go help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stdout, "Unknown help topic %#q.  Run 'autocash help'.\n", arg)
	os.Exit(2) // failed at 'autocash help cmd'
}
