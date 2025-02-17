package main

import (
	"context"
	"fmt"
	"github.com/MattInReality/gator/internal/database"
	"github.com/google/uuid"
	"time"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	ctx := context.TODO()
	if len(cmd.args) < 2 {
		return fmt.Errorf("add feed requires 2 arguments: name and url")
	}
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
	}
	feed, err := s.db.CreateFeed(ctx, newFeed)
	if err != nil {
		return err
	}
	ffp := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	_, err = s.db.CreateFeedFollow(ctx, ffp)
	if err != nil {
		return err
	}
	fmt.Printf("%v", feed)
	return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	_, _ = s, cmd
	ctx := context.TODO()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%v\n", feed)
	}
	return nil
}
