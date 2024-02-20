package urlbuilder

/*

	Basically just a reimplementation of
	https://github.com/Tyrrrz/DiscordChatExporter/blob/master/DiscordChatExporter.Core/Utils/UrlBuilder.cs



*/

import (
	"net/url"
)

type URLBuilder struct {
	root_uri  string
	arguments map[string]string
}

/*
u.SetRoot("https://google.com/search")
*/
func (u *URLBuilder) SetRoot(root_uri string) *URLBuilder {
	u.root_uri = root_uri
	return u
}

func (u *URLBuilder) AddArgument(key string, value string) *URLBuilder {
	if len(u.arguments) == 0 {
		u.arguments = make(map[string]string)
	}
	u.arguments[key] = value
	return u
}

/*
Return the url rendered as a string

# Code

urlbuilder := urlbuilder.NewURLBuilder("https://google.com/search")

fmt.Println(urlbuilder.BuildString())

urlbuilder.AddArgument("q", "This is a test query with google search!!!")

fmt.Println(urlbuilder.BuildString())

# Returns

https://google.com/search?q=This%20is%20a%20test%20query%20with%20google%20search%21%21%21
*/
func (u URLBuilder) BuildString() string {
	if len(u.arguments) == 0 { // don't even bother
		return u.root_uri
	}

	s := u.root_uri
	s += "?"

	argument_length := len(u.arguments)
	argument_counter := 0
	for key, value := range u.arguments {
		argument_counter++
		s += url.PathEscape(key)
		s += "="
		s += url.PathEscape(value)
		if argument_counter != argument_length {
			s += "&"
		}
	}

	return s
}

/*
Returns a new URLBuilder object given the parameter for the root_uri

# Usage:

u := urlbuilder.NewURLBuilder("https://google.com")

u.AddArgument("q", "Example google search")

fmt.Println(u.BuildString)

# Or Alternatively:

fmt.Println(urlbuilder.NewURLBuilder("https://google.com").AddArgument("q", "Example google search").BuildString())
*/
func NewURLBuilder(root_uri string) *URLBuilder {
	return new(URLBuilder).SetRoot(root_uri)
}
