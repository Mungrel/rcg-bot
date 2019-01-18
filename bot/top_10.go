package bot

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"
)

const top10Year = 2018

// Top10 determines the top 10 posts of all time by number of reactions for the year
// and posts it as an album.
func (bot *Bot) Top10() error {
	allPosts, err := bot.getAllPosts()
	if err != nil {
		return err
	}

	top10Posts := getTop10Posts(allPosts)

	return bot.postTop10Posts(top10Posts)
}

// Post represents a simplified summary of an FB Post.
type Post struct {
	ID             string
	CreatedTime    *time.Time
	Message        string
	Permalink      string
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
	Permalink       string `json:"permalink_url"`
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
		Permalink:      p.Permalink,
		TotalReactions: p.Reactions.Summary.TotalCount,
	}
}

const (
	feedURL       = "/feed"
	postsURL      = "/posts"
	postsPageSize = "100"
)

// getAllPosts gets all posts from the FB API.
func (bot *Bot) getAllPosts() ([]Post, error) {
	fbPosts := []fbPost{}

	since := time.Date(top10Year, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(top10Year, 12, 31, 23, 59, 59, 0, time.UTC)

	params := url.Values{}
	params.Add("date_format", "U")
	params.Add("fields", "id,created_time,message,permalink_url,reactions.limit(0).summary(1)")
	params.Add("limit", postsPageSize)
	params.Add("since", strconv.FormatInt(since.Unix(), 10))
	params.Add("until", strconv.FormatInt(until.Unix(), 10))

	var response postResponse
	err := bot.fbClient.Get(postsURL, params, &response)
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

// getTop10Posts sorts the posts by reaction count in descending order and returns the first 10.
func getTop10Posts(posts []Post) []Post {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].TotalReactions > posts[j].TotalReactions
	})

	return posts[:10]
}

// postAsAlbum posts the posts to the FB API as an album.
func (bot *Bot) postTop10Posts(posts []Post) error {
	message := "Top 10 page posts of " + strconv.Itoa(top10Year) + "\n\n"
	for _, post := range posts {
		line := fmt.Sprintf("%s\n%d reactions\n\n", post.Permalink, post.TotalReactions)
		message += line
	}

	fmt.Print(message)

	params := url.Values{}
	params.Add("message", message)

	return bot.fbClient.Post(feedURL, params)
	return nil
}
