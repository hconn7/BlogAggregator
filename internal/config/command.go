package config

import (
	"errors"
	"fmt"

	"github.com.hconn7/BlogAggregator/internal/database"
	"github.com/hconn7/internal/database"
)

type Command struct {
	Name string
	Args []string
}
type Commands struct {
	CommandMap map[string]func(*State, Command) error
}
type State struct {
	db         *database.Queries
	CfgPointer *Config
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("No login name entered")

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
