package clt

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

type ABC struct {
}
type Commands struct {
	CommandList []*Command
	ApiName     string
}
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

	for _, cmd := range []*Command(me.CommandList) {
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

	fmt.Fprintf(os.Stderr, "%s: unknown subcommand %q \n Run '%s help' for usage.\n", me.ApiName, args[0], me.ApiName)
	os.Exit(2)
}

var usageTemplate = `
commands server - connect everysites!

Usage:

    {{.ApiName}} command [arguments]

The commands are:
{{range .CommandList}}{{if .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "{{.ApiName}} help [command]" for more information about a command.

`

var helpTemplate = `{{if .Command.Runnable}}usage: {{.ApiName}} {{.Command.UsageLine}}
{{end}}{{.Command.Long}}
`

func (me *Commands) usage() {
	me.tmpl(os.Stdout, usageTemplate, map[string]interface{}{
		"CommandList": []*Command(me.CommandList),
		"ApiName":     me.ApiName,
	})
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
		fmt.Fprintf(os.Stdout, "usage: %s help command\n\nToo many arguments given.\n", me.ApiName)
		os.Exit(2) // failed at 'me.ApiName help'
	}

	arg := args[0]

	for _, cmd := range []*Command(me.CommandList) {
		if cmd.Name() == arg {

			me.tmpl(os.Stdout, helpTemplate, map[string]interface{}{
				"Command": cmd,
				"ApiName": me.ApiName,
			})
			return
		}
	}

	fmt.Fprintf(os.Stdout, "Unknown help topic %#q.  Run '%s help'.\n", arg, me.ApiName)
	os.Exit(2) // failed at 'me.ApiName help cmd'
}
