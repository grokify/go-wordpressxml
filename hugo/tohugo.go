// hugo converts Wordpress XML to Hugo posts.
package hugo

import (
	"fmt"
	"strings"
	"time"

	wxr "github.com/frankbille/go-wxr-import"
	wordpressxml "github.com/grokify/go-wordpressxml"
	"github.com/grokify/mogo/type/stringsutil"
)

type Post struct {
	Title    string
	Author   string
	Date     time.Time
	Tags     []string
	URL      string
	DisqusID string
	Draft    bool
	Body     string
}

type WxrConverter struct {
	URLFunc      func(*wxr.Item) string
	DisqusIDFunc func(*wxr.Item) string
}

func URLFuncDefaultFunc(wxrItem *wxr.Item) string {
	if wxrItem == nil {
		return ""
	}
	return wxrItem.Link
}

func DisqusIDDefaultFunc(wxrItem *wxr.Item) string {
	if wxrItem == nil {
		return ""
	}
	return wxrItem.PostName
}

func (wc *WxrConverter) ConvertPosts(wxrData wxr.Wxr) ([]Post, error) {
	posts := []Post{}

	for _, channel := range wxrData.Channels {
		for _, item := range channel.Items {
			post, err := wc.ConvertPost(item)
			if err != nil {
				return posts, err
			}
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func (wc *WxrConverter) ConvertPost(wxrItem wxr.Item) (Post, error) {
	post := Post{
		Title:  wxrItem.Title,
		Author: wxrItem.Creator,
		Tags:   ConvertWxrItemCategories(wxrItem.Categories),
		Body:   wxrItem.Content}
	postTypes := map[string]int{"page": 1, "post": 1}
	_, ok := postTypes[wxrItem.PostType]
	if !ok {
		return post, fmt.Errorf("E_UNHANDLED_POST_TYPE [%s]", wxrItem.PostType)
	}
	switch wxrItem.Status {
	case "publish":
		post.Date = wxrItem.PubDate
	case "draft":
		post.Draft = true
		post.Date = wxrItem.PostDate
	default:
		return post, fmt.Errorf("E_STATUS_NOT_FOUND [%s]", wxrItem.Status)
	}
	if wc.URLFunc == nil {
		post.URL = URLFuncDefaultFunc(&wxrItem)
	}
	if wc.URLFunc == nil {
		post.DisqusID = DisqusIDDefaultFunc(&wxrItem)
	}

	return post, nil
}

func ConvertWxrItemCategories(cats []wxr.ItemCategory) []string {
	out := []string{}
	for _, cat := range cats {
		out = append(out, cat.DisplayName)
	}
	return stringsutil.SliceCondenseSpace(out, true, false)
}

func Convert(wpxml *wordpressxml.WpXML) []Post {
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
