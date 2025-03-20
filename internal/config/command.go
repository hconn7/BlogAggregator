package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com.hconn7/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

type Command struct {
	Name string
	Args []string
}
type Commands struct {
	CommandMap map[string]func(*State, Command) error
}
type State struct {
	Db         *database.Queries
	CfgPointer *Config
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("No login name entered")

	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("User %s is not registered, please register before you login!", cmd.Args[0])
		os.Exit(1)
	}
	userName := cmd.Args[0]
	s.CfgPointer.SetUser(string(userName))
	fmt.Println("User set!")
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandMap[name] = f
}
func (c *Commands) Run(s *State, cmd Command) error {
	command, ok := c.CommandMap[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	err := command(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func HandlerRegister(state *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		fmt.Println("What is the users name?")
	}
	_, err := state.Db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		fmt.Print("user exists!")
		os.Exit(1)
	}
	_, err = state.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	})
	if err != nil {
		return err
	}
	fmt.Println("User has been registered")
	fmt.Printf("Logging user: %s in", cmd.Args[0])
	HandlerLogin(state, cmd)
	fmt.Printf("User logged in as %s", cmd.Args[0])
	return nil
}
