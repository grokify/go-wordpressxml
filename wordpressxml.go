// wordpressxml provides WordPress XML parser with metadata
package wordpressxml

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/grokify/gotilla/time/timeutil"
)

type WpXml struct {
	Channel        Channel `xml:"channel"`
	CreatorCounts  map[string]int
	CreatorToIndex map[string]int
}

func NewWordpressXml() WpXml {
	theXml := WpXml{}
	return theXml
}

// ReadXml reads a WordPress XML file from the provided path.
func (wpxml *WpXml) ReadXml(filepath string) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(bytes, &wpxml)
	if err != nil {
		return err
	}
	wpxml.inflate()
	return nil
}

func (wpxml *WpXml) inflate() error {
	creatorMap := map[string]int{}
	for i, item := range wpxml.Channel.Items {
		if len(item.Creator) > 0 {
			if _, ok := creatorMap[item.Creator]; ok {
				creatorMap[item.Creator]++
			} else {
				creatorMap[item.Creator] = 1
			}
		}

		item = wpxml.inflateItem(item)
		wpxml.Channel.Items[i] = item
	}
	wpxml.CreatorCounts = creatorMap
	wpxml.inflateAuthors()
	return nil
}

func (wpxml *WpXml) inflateItem(item Item) Item {
	if len(item.Encoded) > 0 && len(item.Encoded[0]) > 0 {
		item.Content = item.Encoded[0]
		item.Encoded[0] = ""
	}
	if len(item.PubDate) > 0 {
		dtRfc3339, err := timeutil.FromTo(item.PubDate, time.RFC1123Z, time.RFC3339)
		if err == nil {
			item.PubDateRfc3339 = dtRfc3339
		}
	}
	return item
}

func (wpxml *WpXml) inflateAuthors() error {
	a2i := wpxml.AuthorsToIndex()
	for i, item := range wpxml.Channel.Items {
		if len(item.Creator) > 0 {
			authorLogin := item.Creator
			if _, ok := a2i[authorLogin]; ok {
				authorIndex := a2i[authorLogin]
				itemThin := ItemThin{Title: item.Title, Index: i}
				if wpxml.Channel.Authors[authorIndex].AuthorArticles == nil {
					wpxml.Channel.Authors[authorIndex].AuthorArticles = []ItemThin{}
				}
				wpxml.Channel.Authors[authorIndex].AuthorArticles = append(wpxml.Channel.Authors[authorIndex].AuthorArticles, itemThin)
			}
		}
	}
	wpxml.CreatorToIndex = a2i
	return nil
}

func (wpxml *WpXml) AuthorsToIndex() map[string]int {
	a2i := map[string]int{}
	for i, author := range wpxml.Channel.Authors {
		a2i[author.AuthorLogin] = i
	}
	return a2i
}

// AuthorForLogin returns the Author object for a given AuthorLogin
// or username.
func (wpxml *WpXml) AuthorForLogin(authorLogin string) (Author, error) {
	a2i := wpxml.CreatorToIndex
	if index, ok := a2i[authorLogin]; ok {
		author := wpxml.Channel.Authors[index]
		return author, nil
	}
	return Author{}, errors.New("Author Not Found")
}

// ArticlesMetaTable generates the data to be written out
// as a CSV.
func (wpxml *WpXml) ArticlesMetaTable() [][]string {
	articles := [][]string{}
	header := []string{"Index", "Date", "Login", "Author", "Title", "Link"}
	articles = append(articles, header)
	a2i := wpxml.AuthorsToIndex()
	for i, item := range wpxml.Channel.Items {
		author, err := wpxml.AuthorForLogin(item.Creator)
		if err != nil {
			panic(err)
		}
		article := []string{strconv.Itoa(i), item.PubDate, item.Creator, author.AuthorDisplayName, item.Title, item.Link}
		articles = append(articles, article)
	}
	wpxml.CreatorToIndex = a2i
	return articles
}

// Writes articles metadata as a CSV file.
func (wpxml *WpXml) WriteMetaCsv(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.WriteAll(wpxml.ArticlesMetaTable())
	return w.Error()
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Authors []Author `xml:"author"`
	Items   []Item   `xml:"item"`
}

// Author is the WordPress XML author object.
type Author struct {
	AuthorId          int        `xml:"author_id"`
	AuthorLogin       string     `xml:"author_login"`
	AuthorEmail       string     `xml:"author_email"`
	AuthorDisplayName string     `xml:"author_display_name"`
	AuthorFirstName   string     `xml:"author_first_name"`
	AuthorLastName    string     `xml:"author_last_name"`
	AuthorArticles    []ItemThin `xml:"-"`
}

// Item is a WordPress XML item which can be a post, page or other object.
type Item struct {
	Id             int        `xml:"post_id"`
	Title          string     `xml:"title"`
	Creator        string     `xml:"creator"`
	Encoded        []string   `xml:"encoded"`
	IsSticky       int        `xml:"is_sticky"`
	Link           string     `xml:"link"`
	PubDate        string     `xml:"pubDate"`
	Description    string     `xml:"description"`
	PostDate       string     `xml:"post_date"`
	PostDateGmt    string     `xml:"post_date_gmt"`
	PostName       string     `xml:"post_name"`
	PostType       string     `xml:"post_type"`
	Status         string     `xml:"status"`
	Categories     []Category `xml:"category"`
	Content        string
	PubDateRfc3339 string
	Comments       []Comment  `xml:"comment"`
}

// ItemThin is a WordPress XML item that is used as additional
// metadata in the Author object.
type ItemThin struct {
	Title string
	Index int
}

type Category struct {
	Domain      string `xml:"domain,attr"`
	DisplayName string `xml:",chardata"`
	UrlSlug     string `xml:"nicename,attr"`
}

type Comment struct {
	Id          int `xml:"comment_id"`
	Parent      int `xml:"comment_parent"`
	Author      string `xml:"comment_author"`
	DateGmt     string `xml:"comment_date_gmt"`
	Content     string `xml:"comment_content"`
}
