package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/MattInReality/gator/internal/config"
	"github.com/MattInReality/gator/internal/database"
	"github.com/google/uuid"
	"log"
	"os"
	"time"
)

import _ "github.com/lib/pq"

type state struct {
	db     *database.Queries
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commandHandler func(*state, command) error

type commandHandlerLoggedIn func(*state, command, database.User) error

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username is required")
	}
	ctx := context.TODO()
	username := cmd.args[0]
	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		os.Exit(1)
	}
	err = s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Println("user has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("username is required to register")
	}
	username := cmd.args[0]
	ctx := context.TODO()
	uuid := uuid.New()
	newUser := database.CreateUserParams{
		ID:        uuid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	user, err := s.db.CreateUser(ctx, newUser)
	if err != nil {
		os.Exit(1)
	}
	err = s.config.SetUser(username)
	if err != nil {
		return err
	}
	fmt.Printf("user was created %v", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	_ = cmd
	ctx := context.TODO()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("error reseting user table: %v", err)
	}
	err = s.config.SetUser("")
	if err != nil {
		return err
	}
	fmt.Printf("reset user table successful")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.TODO()
	currentUser := s.config.CurrentUserName
	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("command agg requires a single argument for the time between requests")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Printf("%v", err)
		}
	}
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

func NewCommands() commands {
	return commands{
		handlers: make(map[string]commandHandler),
	}
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
	dbURL := configFile.DbUrl
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("issue with database connection")
		os.Exit(1)
	}
	dbQueries := database.New(db)
	appState := state{}
	appState.config = &configFile
	appState.db = dbQueries

	cmds := NewCommands()
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", handlerFollowing)
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

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

func middlewareLoggedIn(handler commandHandlerLoggedIn) commandHandler {
	return func(s *state, cmd command) error {
		ctx := context.TODO()
		user, err := s.db.GetUser(ctx, s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
