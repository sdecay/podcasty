// Morphs the goose-generated model to a format we want to return

package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/sdecay/podcasty/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func dbUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func dbFeedtoFeed(feed database.Feed) Feed {
	return Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		URL:       feed.Url,
		UserID:    feed.UserID,
	}
}

func dbFeedstoFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}

	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, dbFeedtoFeed(dbFeed))
	}

	return feeds
}

type Follow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func dbFollowToFollow(dbFollow database.Follow) Follow {
	return Follow{
		ID:        dbFollow.ID,
		CreatedAt: dbFollow.CreatedAt,
		UpdatedAt: dbFollow.UpdatedAt,
		UserID:    dbFollow.UserID,
		FeedID:    dbFollow.FeedID,
	}
}

func dbFollowedToFollowed(dbFollowed []database.Follow) []Follow {
	follows := []Follow{}

	for _, dbFeed := range dbFollowed {
		follows = append(follows, dbFollowToFollow(dbFeed))
	}

	return follows
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func dbPostToPost(dbPost database.Post) Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
	}
}

func dbPostsToPosts(dbPosts []database.Post) []Post {
	posts := []Post{}

	for _, dbPost := range dbPosts {
		posts = append(posts, dbPostToPost(dbPost))
	}

	return posts
}
