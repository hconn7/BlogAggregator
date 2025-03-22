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

	// DB init
	db, err := sql.Open("postgres", queryString)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	cfg, err := config.ReadDB()
	if err != nil {
		panic(err)
	}
	// Meta struct init
	state := &config.State{CfgPointer: &cfg, Db: dbQueries}
	commands := config.Commands{CommandMap: make(map[string]func(*config.State, config.Command) error)}

	// CMD Registry
	commands.Register("browse", config.MiddlewareLoggedIn(config.HandlerBrowseFeeds))
	commands.Register("unfollow", config.MiddlewareLoggedIn(config.HandlerUnfollowFeed))
	commands.Register("following", config.MiddlewareLoggedIn(config.HandlerGetFeedFollowsForUser))
	commands.Register("follow", config.MiddlewareLoggedIn(config.HandlerFollowFeed))
	commands.Register("feeds", config.HandlerGetFeeds)
	commands.Register("addfeed", config.MiddlewareLoggedIn(config.HandlerAddFeed))
	commands.Register("agg", config.HandlerAgg)
	commands.Register("users", config.HandlerGetUserS)
	commands.Register("reset", config.HandlerReset)
	commands.Register("login", config.HandlerLogin)
	commands.Register("register", config.HandlerRegister)

	// Args Init
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

	// Exist?
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
