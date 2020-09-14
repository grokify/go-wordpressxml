// hugo converts Wordpress XML to Hugo posts.
package hugo

import (
	"strings"
	"time"

	"github.com/grokify/go-wordpressxml"
)

type Post struct {
	Author   string
	Title    string
	Date     time.Time
	Tags     []string
	URL      string
	DisqusID string
	Draft    bool
	Body     string
}

func Convert(wpxml *wordpressxml.WpXml) []Post {
	posts := []Post{}
	for _, item := range wpxml.Channel.Items {
		item.PostType = strings.ToLower(strings.TrimSpace(item.PostType))
		item.Status = strings.ToLower(strings.TrimSpace(item.Status))
		post := Post{
			Title: strings.TrimSpace(item.Title),
			Body:  item.Content}
		if item.PostType != "post" {
			continue
		}
		if item.Status == "draft" {
			post.Draft = true
			post.Date = item.PostDatetime
		} else {
			post.Date = item.PubDatetime
		}
		author, err := wpxml.AuthorForLogin(item.Creator)
		if err == nil {
			post.Author = strings.TrimSpace(author.AuthorDisplayName)
		}
		posts = append(posts, post)
	}
	return posts
}
