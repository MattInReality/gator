package main

import (
	"context"
	"fmt"
	"github.com/MattInReality/gator/internal/database"
	"github.com/google/uuid"
	"time"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("follow command requires a url argument")
	}
	ctx := context.TODO()
	url := cmd.args[0]
	feed, err := s.db.FindFeedFromURL(ctx, url)
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
	feedFollow, err := s.db.CreateFeedFollow(ctx, ffp)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", feedFollow)
	return nil
}

func handlerFollowing(s *state, cmd command) error {
	ctx := context.TODO()
	// TODO: middleware brief suggests this should be quering the user table.
	// my query must be different to theirs.
	feeds, err := s.db.GetFeedFollowsForUser(ctx, s.config.CurrentUserName)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("unfollow feed requires 1 argument - the url to unfollow")
	}
	ctx := context.TODO()
	ufp := database.UnfollowFeedParams{
		Url:    cmd.args[0],
		UserID: user.ID,
	}
	err := s.db.UnfollowFeed(ctx, ufp)
	if err != nil {
		return err
	}
	return nil
}
