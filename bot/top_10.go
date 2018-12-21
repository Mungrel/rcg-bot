package bot

// Top10 determines the top 10 posts of all time by number of reactions for the year
// and posts it as an album.
func (bot *Bot) Top10() error {
	return nil
}

// Post represents a post from the FB API.
type Post struct{}

// getAllPosts gets all posts from the FB API.
func (bot *Bot) getAllPosts() ([]Post, error) {
	return nil, nil
}

// getTop10Posts sorts the posts by reaction count and returns the top 10.
func getTop10Posts(posts []Post) []Post {
	return nil
}

// postAsAlbum posts the posts to the FB API as an album.
func (bot *Bot) postAsAlbum(posts []Post) error {
	return nil
}
