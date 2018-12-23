package bot

import (
	"fmt"
	"sort"
	"time"
)

// Top10 determines the top 10 posts of all time by number of reactions for the year
// and posts it as an album.
func (bot *Bot) Top10() error {
	allPosts, err := bot.getAllPosts()
	if err != nil {
		return err
	}

	top10Posts := getTop10Posts(allPosts)

	return bot.postAsAlbum(top10Posts)
}

// Post represents a simplified summary of an FB Post.
type Post struct {
	ID             string
	CreatedTime    *time.Time
	Message        string
	TotalReactions int
}

// postResponse represents a response from the FB API for post reaction summaries
type postResponse struct {
	Data   []fbPost `json:"data"`
	Paging paging   `json:"paging"`
}

type fbPost struct {
	ID              string `json:"id"`
	CreatedTimeUnix int64  `json:"created_time"`
	Message         string `json:"message"`
	Reactions       struct {
		Summary struct {
			TotalCount int `json:"total_count"`
		} `json:"summary"`
	} `json:"reactions"`
}

type paging struct {
	Next *string `json:"next"`
}

func (p *fbPost) toPost() Post {
	t := time.Unix(p.CreatedTimeUnix, 0)
	return Post{
		ID:             p.ID,
		CreatedTime:    &t,
		Message:        p.Message,
		TotalReactions: p.Reactions.Summary.TotalCount,
	}
}

const postsURL = "/posts"

// getAllPosts gets all posts from the FB API.
func (bot *Bot) getAllPosts() ([]Post, error) {
	fbPosts := []fbPost{}

	var response postResponse
	err := bot.fbClient.Get(postsURL, &response)
	if err != nil {
		return nil, err
	}

	fbPosts = append(fbPosts, response.Data...)

	for response.Paging.Next != nil {
		fmt.Println("Fetching next page")
		// Need to declare a fresh struct here so that Next will be nil if not in response
		var tmpResponse postResponse
		err := bot.fbClient.GetAbsoluteURL(*response.Paging.Next, &tmpResponse)
		if err != nil {
			return nil, err
		}

		response = tmpResponse

		fbPosts = append(fbPosts, response.Data...)
	}

	fmt.Printf("\nNumber of posts: %d\n\n", len(fbPosts))
	posts := make([]Post, 0, len(fbPosts))
	for _, fbPost := range fbPosts {
		posts = append(posts, fbPost.toPost())
	}

	return posts, nil
}

// getTop10Posts sorts the posts by reaction count and returns the top 10.
func getTop10Posts(posts []Post) []Post {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].TotalReactions > posts[j].TotalReactions
	})

	return posts[:10]
}

// postAsAlbum posts the posts to the FB API as an album.
func (bot *Bot) postAsAlbum(posts []Post) error {
	return nil
}
