// wordpressxml provides WordPress XML parser with metadata
package wordpressxml

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/time/timeutil"
)

type WpXML struct {
	Channel        Channel `xml:"channel"`
	CreatorCounts  map[string]int
	CreatorToIndex map[string]int
}

func NewWpXML() WpXML {
	return WpXML{
		CreatorCounts:  map[string]int{},
		CreatorToIndex: map[string]int{}}
}

// ReadXml reads a WordPress XML file from the provided path.
func (wpxml *WpXML) ReadFile(filepath string) error {
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

func (wpxml *WpXML) inflate() {
	creatorMap := map[string]int{}
	for i, item := range wpxml.Channel.Items {
		if len(item.Creator) > 0 {
			creatorMap[item.Creator]++
		}

		item = wpxml.inflateItem(item)
		wpxml.Channel.Items[i] = item
	}
	wpxml.CreatorCounts = creatorMap
	wpxml.inflateAuthors()
}

func (wpxml *WpXML) inflateItem(item Item) Item {
	if len(item.Encoded) > 0 && len(item.Encoded[0]) > 0 {
		item.Content = item.Encoded[0]
		item.Encoded[0] = ""
	}
	if len(item.PostDate) > 0 {
		dt, err := time.Parse(timeutil.SQLTimestamp, item.PostDate)
		if err == nil {
			item.PostDatetime = dt
		}
	}
	if len(item.PubDate) > 0 {
		dt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			item.PubDatetime = dt
		}
	}
	return item
}

func (wpxml *WpXML) inflateAuthors() {
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
}

func (wpxml *WpXML) AuthorsToIndex() map[string]int {
	a2i := map[string]int{}
	for i, author := range wpxml.Channel.Authors {
		a2i[author.AuthorLogin] = i
	}
	return a2i
}

// AuthorForLogin returns the Author object for a given AuthorLogin
// or username.
func (wpxml *WpXML) AuthorForLogin(authorLogin string) (Author, error) {
	a2i := wpxml.CreatorToIndex
	if index, ok := a2i[authorLogin]; ok {
		author := wpxml.Channel.Authors[index]
		return author, nil
	}
	return Author{}, errors.New("Author Not Found")
}

// ArticlesMetaTable generates the data to be written out
// as a CSV.
func (wpxml *WpXML) ArticlesMetaTable() table.Table {
	tbl := table.NewTable("Articles Metadata")
	tbl.Columns = []string{"Index", "Date", "Login", "Author", "Title", "Link"}
	a2i := wpxml.AuthorsToIndex()
	for i, item := range wpxml.Channel.Items {
		authorDisplayName := ""
		author, err := wpxml.AuthorForLogin(item.Creator)
		if err == nil {
			authorDisplayName = author.AuthorDisplayName
		}
		article := []string{strconv.Itoa(i), item.PubDatetime.Format(time.RFC3339), item.Creator, authorDisplayName, item.Title, item.Link}
		tbl.Rows = append(tbl.Rows, article)
	}
	wpxml.CreatorToIndex = a2i
	return tbl
}

// WriteMetaCsv writes articles metadata as a CSV file.
func (wpxml *WpXML) WriteMetaCsv(filepath string) error {
	tbl := wpxml.ArticlesMetaTable()
	return tbl.WriteCSV(filepath)
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
	AuthorID          int        `xml:"author_id"`
	AuthorLogin       string     `xml:"author_login"`
	AuthorEmail       string     `xml:"author_email"`
	AuthorDisplayName string     `xml:"author_display_name"`
	AuthorFirstName   string     `xml:"author_first_name"`
	AuthorLastName    string     `xml:"author_last_name"`
	AuthorArticles    []ItemThin `xml:"-"`
}

// Item is a WordPress XML item which can be a post, page or other object.
type Item struct {
	ID           int        `xml:"post_id"`
	Title        string     `xml:"title"`
	Creator      string     `xml:"creator"`
	Encoded      []string   `xml:"encoded"`
	IsSticky     int        `xml:"is_sticky"`
	Link         string     `xml:"link"`
	PubDate      string     `xml:"pubDate"`
	Description  string     `xml:"description"`
	PostDate     string     `xml:"post_date"`
	PostDateGmt  string     `xml:"post_date_gmt"`
	PostName     string     `xml:"post_name"`
	PostType     string     `xml:"post_type"`
	Status       string     `xml:"status"`
	Categories   []Category `xml:"category"`
	Comments     []Comment  `xml:"comment"`
	Content      string
	PostDatetime time.Time
	PubDatetime  time.Time
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
	URLSlug     string `xml:"nicename,attr"`
}

type Comment struct {
	ID          int    `xml:"comment_id"`
	Parent      int    `xml:"comment_parent"`
	Author      string `xml:"comment_author"`
	AuthorEmail string `xml:"comment_author_email"`
	AuthorURL   string `xml:"comment_author_url"`
	DateGmt     string `xml:"comment_date_gmt"`
	Content     string `xml:"comment_content"`
	IndentLevel int    `xml:"-"`
}
