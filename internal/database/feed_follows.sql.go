// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :many

WITH inserted AS (
    INSERT INTO feed_follows (id, user_id, feed_id)
    VALUES (gen_random_uuid(), $1, $2)
    RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT 
    inserted.id,
    inserted.created_at,
    inserted.updated_at,
    inserted.user_id,
    inserted.feed_id,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted
JOIN users ON users.id = inserted.user_id
JOIN feeds ON feeds.id = inserted.feed_id
`

type CreateFeedFollowParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	UserName  string
	FeedName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) ([]CreateFeedFollowRow, error) {
	rows, err := q.db.QueryContext(ctx, createFeedFollow, arg.UserID, arg.FeedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CreateFeedFollowRow
	for rows.Next() {
		var i CreateFeedFollowRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.UserName,
			&i.FeedName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteFeedFollowByUserAndFeedURL = `-- name: DeleteFeedFollowByUserAndFeedURL :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
AND feed_follows.feed_id = (SELECT id FROM feeds WHERE url = $2)
`

type DeleteFeedFollowByUserAndFeedURLParams struct {
	UserID uuid.UUID
	Url    string
}

func (q *Queries) DeleteFeedFollowByUserAndFeedURL(ctx context.Context, arg DeleteFeedFollowByUserAndFeedURLParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollowByUserAndFeedURL, arg.UserID, arg.Url)
	return err
}

const getFeedByURL = `-- name: GetFeedByURL :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds WHERE url = $1
`

func (q *Queries) GetFeedByURL(ctx context.Context, url string) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByURL, url)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeedFollowsForUser = `-- name: GetFeedFollowsForUser :many
SELECT 
    feed_follows.id,
    feed_follows.created_at,
    feed_follows.updated_at,
    feed_follows.user_id,
    feed_follows.feed_id,
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follows
JOIN users ON users.id = feed_follows.user_id
JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
`

type GetFeedFollowsForUserRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	UserName  string
	FeedName  string
}

func (q *Queries) GetFeedFollowsForUser(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsForUserRow
	for rows.Next() {
		var i GetFeedFollowsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
			&i.UserName,
			&i.FeedName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
