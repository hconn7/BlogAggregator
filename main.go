package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com.hconn7/BlogAggregator/internal/config"
	"github.com.hconn7/BlogAggregator/internal/database"
	_ "github.com/lib/pq"
)

const queryString = "postgres://henry:@localhost:5432/gator?sslmode=disable"

func main() {
	db, err := sql.Open("postgres", queryString)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	cfg, err := config.ReadDB()
	if err != nil {
		panic(err)
	}
	state := &config.State{CfgPointer: &cfg, Db: dbQueries}

	commands := config.Commands{CommandMap: make(map[string]func(*config.State, config.Command) error)}

	commands.Register("login", config.HandlerLogin)
	commands.Register("register", config.HandlerRegister)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error! need more arguments")
		os.Exit(1)
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

			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Command doesn't exist!")
		os.Exit(1)
	}
}
