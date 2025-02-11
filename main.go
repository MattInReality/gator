package main

import (
	"errors"
	"fmt"
	"github.com/MattInReality/gator/internal/config"
	"log"
	"os"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commandHandler func(*state, command) error

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username is required")
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("user has been set")
	return nil
}

type commands struct {
	handlers map[string]commandHandler
}

func (c *commands) register(name string, f commandHandler) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	cmdHndl, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("no handler found for command %s", cmd.name)
	}
	err := cmdHndl(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func NewCommand(name string, args []string) (command, error) {
	return command{name: name, args: args}, nil
}

func main() {
	configFile, err := config.Read()
	if err != nil {
		log.Fatal("config could not be read - check file")
		panic(1)
	}
	appState := state{}
	appState.config = &configFile

	cmds := commands{}
	cmds.handlers = make(map[string]commandHandler)
	cmds.handlers["login"] = handlerLogin

	cmdArgs := os.Args
	if len(cmdArgs) < 2 {
		fmt.Println("not enough arguments were provided")
		os.Exit(1)
	}

	commandName := cmdArgs[1]
	args := cmdArgs[2:]
	cmd, _ := NewCommand(commandName, args)

	err = cmds.run(&appState, cmd)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}
