package config

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com.hconn7/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.Db.MarkFeedFetched(context.Background(), feed.FeedID)
	if err != nil {
		return err
	}

	displayFeeds, _ := FetchFeed(context.Background(), feed.FeedUrl)

	for _, item := range displayFeeds.Channel.Item {

		parsedTime, _ := time.Parse(time.RFC1123, item.PubDate)
		s.Db.CreatePost(context.Background(), database.CreatePostParams{
			ID:    uuid.New(),
			Title: item.Title,
			Url:   item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "", // Only valid if non-empty
			},
			PublishedAt: sql.NullTime{
				Time:  parsedTime,         // Convert string to time
				Valid: item.PubDate == "", // Only valid if non-empty
			},
			FeedID: feed.FeedID,
		})

		fmt.Printf("- %s\n", item.Title)
	}

	return nil
}
