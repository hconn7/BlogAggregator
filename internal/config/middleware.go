package config

import (
	"context"
	"fmt"

	"github.com.hconn7/BlogAggregator/internal/database"
)

// middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.CfgPointer.User)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		return handler(s, cmd, user)
	}
}
