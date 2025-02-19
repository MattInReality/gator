package main

import (
	"context"
	"fmt"
	"github.com/MattInReality/gator/internal/database"
	"strconv"
	"time"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.args) > 0 {
		l, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			fmt.Printf("limit value could not be parsed as a number. Defaulting to 2")
		}
		limit = int32(l)
	}
	pa := database.GetPostsForUserParams{
		Limit:  limit,
		UserID: user.ID,
	}
	posts, err := s.db.GetPostsForUser(context.Background(), pa)
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("%s -- Published On: %s\n", post.Title, post.PublishedAt.Time.Format(time.RFC1123Z))
		fmt.Printf("%s\n", post.Description.String)
		fmt.Printf("%s\n", post.Url)
		fmt.Println("-------------------------------------------------------")
	}
	return nil
}
