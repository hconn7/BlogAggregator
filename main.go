package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com.hconn7/BlogAggregator/internal/config"
	_ "github.com/lib/pq"
)

const queryString = "postgres://henry:@localhost:5432/gator"

func main() {
	db, err := sql.Open("postgres", quequeryString)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	cfg, err := config.ReadDB()
	if err != nil {
		panic(err)
	}
	state := &config.State{CfgPointer: &cfg}

	commands := config.Commands{CommandMap: make(map[string]func(*config.State, config.Command) error)}

	commands.Register("login", config.HandlerLogin)

	args := os.Args
	if len(args) < 2 {
		os.Exit(1)
		fmt.Println("Error! need more arguments")
	}
	var command config.Command
	cmdName := args[1]
	cmdUser := args[2:]
	command.Name = cmdName
	command.Args = cmdUser
	_, ok := commands.CommandMap[cmdName]
	if ok {
		err = commands.Run(state, command)
		if err != nil {
			os.Exit(1)
			fmt.Println("Username required")
		}
	} else {
		os.Exit(1)
		fmt.Println("Command doesn't exist!")
	}
}
