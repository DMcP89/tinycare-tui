package utils

// This script should contain a function call GetLatestTweet() that returns the latest tweet from the tinycarebot twitter account using https://github.com/n0madic/twitter-scraper.

import (
	"context"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

func GetLatestTweet() (string, error) {
	scraper := twitterscraper.New()
	tweet := <-scraper.GetTweets(context.Background(), "tinycarebot", 1)
	if tweet.Error != nil {
		return "", tweet.Error
	}
	return tweet.Text, nil
}
