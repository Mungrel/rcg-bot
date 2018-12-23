package bot

import "time"

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
	Data []fbPost `json:"data"`
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
	var response postResponse
	err := bot.fbClient.Get(postsURL, &response)
	if err != nil {
		return nil, err
	}

	posts := make([]Post, 0, len(response.Data))
	for _, fbPost := range response.Data {
		posts = append(posts, fbPost.toPost())
	}

	return posts, nil
}

// getTop10Posts sorts the posts by reaction count and returns the top 10.
func getTop10Posts(posts []Post) []Post {
	return nil
}

// postAsAlbum posts the posts to the FB API as an album.
func (bot *Bot) postAsAlbum(posts []Post) error {
	return nil
}
