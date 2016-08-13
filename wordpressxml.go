package wordpressxml

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
)

type WpXml struct {
	Channel        Channel `xml:"channel"`
	CreatorCounts  map[string]int
	CreatorToIndex map[string]int
}

func NewWordpressXml() WpXml {
	xml := WpXml{}
	return xml
}

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
		if len(item.Encoded) > 0 && len(item.Encoded[0]) > 0 {
			content := item.Encoded[0]
			wpxml.Channel.Items[i].Content = content
			item.Encoded[0] = ""
		}
	}
	wpxml.CreatorCounts = creatorMap
	wpxml.inflateAuthors()
	return nil
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

func (wpxml *WpXml) AuthorForLogin(authorLogin string) (Author, error) {
	a2i := wpxml.CreatorToIndex
	if index, ok := a2i[authorLogin]; ok {
		author := wpxml.Channel.Authors[index]
		return author, nil
	}
	return Author{}, errors.New("Author Not Found")
}

func (wpxml *WpXml) ArticlesMetaTable() [][]string {
	articles := [][]string{}
	header := []string{"Index", "Date", "Login", "Author", "Title"}
	articles = append(articles, header)
	a2i := wpxml.AuthorsToIndex()
	for i, item := range wpxml.Channel.Items {
		author, err := wpxml.AuthorForLogin(item.Creator)
		if err != nil {
			panic(err)
		}
		article := []string{strconv.Itoa(i), item.PubDate, item.Creator, author.AuthorDisplayName, item.Title}
		articles = append(articles, article)
	}
	wpxml.CreatorToIndex = a2i
	return articles
}

func (wpxml *WpXml) WriteCsv(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	articles := wpxml.ArticlesMetaTable()

	w := csv.NewWriter(file)
	w.WriteAll(articles)
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

type Author struct {
	AuthorId          int    `xml:"author_id"`
	AuthorLogin       string `xml:"author_login"`
	AuthorEmail       string `xml:"author_email"`
	AuthorDisplayName string `xml:"author_display_name"`
	AuthorFirstName   string `xml:"author_first_name"`
	AuthorLastName    string `xml:"author_last_name"`
	AuthorArticles    []ItemThin
}

type Item struct {
	Title       string   `xml:"title"`
	Creator     string   `xml:"creator"`
	Encoded     []string `xml:"encoded"`
	Content     string
	IsSticky    int    `xml:"is_sticky"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	PostDate    string `xml:"post_date"`
	PostDateGmt string `xml:"post_date_gmt"`
	PostName    string `xml:"post_name"`
	PostType    string `xml:"post_type"`
	Status      string `xml:"status"`
}

type ItemThin struct {
	Title string
	Index int
}
