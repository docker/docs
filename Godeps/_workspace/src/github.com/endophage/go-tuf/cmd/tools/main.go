package main

import (
	"fmt"
	"log"

	"github.com/endophage/go-tuf"
	"github.com/endophage/go-tuf/signed"
	"github.com/endophage/go-tuf/store"
	"github.com/endophage/go-tuf/util"
	"github.com/flynn/go-docopt"
)

func main() {
	log.SetFlags(0)

	usage := `usage: tuftools [-h|--help] <command> [<args>...]

Options:
  -h, --help

Commands:
  help         Show usage for a specific command
  meta      Generate metadata from the given file path

See "tuf help <command>" for more information on a specific command
`

	args, _ := docopt.Parse(usage, nil, true, "", true)
	cmd := args.String["<command>"]
	cmdArgs := args.All["<args>"].([]string)

	if cmd == "help" {
		if len(cmdArgs) == 0 { // `tuf help`
			fmt.Println(usage)
			return
		} else { // `tuf help <command>`
			cmd = cmdArgs[0]
			cmdArgs = []string{"--help"}
		}
	}

	if err := runCommand(cmd, cmdArgs); err != nil {
		log.Fatalln("ERROR:", err)
	}
}

type cmdFunc func(*docopt.Args, *tuf.Repo) error

type command struct {
	usage string
	f     cmdFunc
}

var commands = make(map[string]*command)

func register(name string, f cmdFunc, usage string) {
	commands[name] = &command{usage: usage, f: f}
}

func runCommand(name string, args []string) error {
	argv := make([]string, 1, 1+len(args))
	argv[0] = name
	argv = append(argv, args...)

	cmd, ok := commands[name]
	if !ok {
		return fmt.Errorf("%s is not a tuf command. See 'tuf help'", name)
	}

	parsedArgs, err := docopt.Parse(cmd.usage, argv, true, "", true)
	if err != nil {
		return err
	}

	db := util.GetSqliteDB()
	local := store.DBStore(db, "")
	signer := signed.Ed25519{}
	repo, err := tuf.NewRepo(&signer, local, "sha256")
	if err != nil {
		return err
	}
	return cmd.f(parsedArgs, repo)
}
