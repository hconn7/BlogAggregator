package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
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

func HandlerFollowFeed(s *State, cmd Command, user database.User) error {
	URL := cmd.Args[0]
	feed, err := s.Db.GetFeedByURL(context.Background(), URL)
	if err != nil {
		return err
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}
	fmt.Printf("Feed followed! User Name: %s, Feed Name: %s\n", user, feed.Name)
	return nil
}
func HandlerGetFeedFollowsForUser(s *State, cmd Command, user database.User) error {

	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
		fmt.Println(feed.UserName)
	}
	return nil
}
func HandlerGetFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println(feeds)
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("please provide time between requests (e.g. 1s, 1m, 1h)")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(time.Second * 10)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
func HandlerGetUserS(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.CfgPointer.User {
			fmt.Printf("* %s (current)", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
func HandlerBrowseFeeds(s *State, cmd Command, user database.User) error {
	limit := cmd.Args[0]
	limitInt, err := strconv.Atoi(limit)

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limitInt)})

	fmt.Printf("Fetching posts for UserID: %s with limit: %d and:\n", user.ID, limitInt)
	if err != nil {
		return err
	}

	fmt.Println(posts)

	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	name := cmd.Args[0]
	URL := cmd.Args[1]
	s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       URL,
		UserID:    user.ID,
	})
	feed, err := FetchFeed(context.Background(), URL)
	if err != nil {
		return err
	}
	feedParams, err := s.Db.GetFeedByURL(context.Background(), URL)
	if err != nil {
		return err
	}
	_, err = s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{UserID: user.ID, FeedID: feedParams.ID})
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("No login name entered")

	}
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		fmt.Printf("User %s is not registered, please register before you login!\n", cmd.Args[0])
		os.Exit(1)
	}
	userName := cmd.Args[0]
	s.CfgPointer.SetUser(string(userName))
	fmt.Printf("--Current User %s\n--", userName)
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
func HandlerReset(state *State, cmd Command) error {
	err := state.Db.DeleteUsers(context.Background())
	if err != nil {
		println("Can't reset users, internal issue")
		os.Exit(1)
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

func HandlerUnfollowFeed(s *State, cmd Command, user database.User) error {
	URL := cmd.Args[0]
	err := s.Db.DeleteFeedFollowByUserAndFeedURL(context.Background(), database.DeleteFeedFollowByUserAndFeedURLParams{Url: URL, UserID: user.ID})
	if err != nil {
		return err
	}
	return nil
}
