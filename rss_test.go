package main

import (
	"context"
	"fmt"
	"testing"
)

const xmllink = "https://www.wagslane.dev/index.xml"

func TestRssFeed(t *testing.T) {
	ctx := context.TODO()
	feed, err := fetchFeed(ctx, xmllink)
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Printf("%v", feed)

}
