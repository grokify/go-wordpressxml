package wordpressxml

import (
	"fmt"
	"testing"
	"os"
)

// print a single post
func PrintPost(item Item, author Author) {
	fmt.Println(item.Title)
	fmt.Println(item.Link)
	fmt.Println("by", author.AuthorDisplayName)
	fmt.Println("--------------")
	fmt.Println(item.Content)
	if len(item.Comments) > 0 {
		fmt.Println()
		fmt.Print("COMMENTS:\n\n")
		for _, c := range item.Comments {
			fmt.Println("Author:", c.Author)
            fmt.Println("--------------")
			fmt.Print(c.Content, "\n\n")
		}
	}
	fmt.Print("\n\n")
}


// one global parsed blog export for all tests to operate on
var wp = NewWordpressXml()

// set up the global parsed export and run the tests
func TestMain(m *testing.M) {
    err := wp.ReadXml("./in/wp_export.xml")
    if err != nil {
        panic(err)
    }
    os.Exit(m.Run())
}

// print the first n posts (real posts, not pages)
func TestPrintFirstPosts(*testing.T) {
	n := 2
	wpxml := &wp

	for _, item := range wpxml.Channel.Items {
		author, err := wpxml.AuthorForLogin(item.Creator)
		if err != nil {
			panic(err)
		}
		if item.PostType == "post" {
			PrintPost(item, author)
			n--
		}
		if n == 0 {
			break
		}
	}
}

// find the first post with at least n comments and print it
func TestPrintFirstPostWithMultiComments(*testing.T) {
	n := 7
	wpxml := &wp

	for _, item := range wpxml.Channel.Items {
		author, err := wpxml.AuthorForLogin(item.Creator)
		if err != nil {
			panic(err)
		}
		if item.PostType == "post" && len(item.Comments) >= n {
			PrintPost(item, author)
			break
		}
	}
}
