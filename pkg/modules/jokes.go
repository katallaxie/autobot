package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/katallaxie/autobot/util"
)

const jokeURL = "https://official-joke-api.appspot.com/random_joke"

// JokesTag is the intent tag for its module
var JokesTag = "jokes"

// Joke represents the response from the joke api
type Joke struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

// JokesReplacer replaces the pattern contained inside the response by a random joke from the api
// specified in jokeURL.
// See modules/modules.go#Module.Replacer() for more details.
func JokesReplacer(locale, entry, response, _ string) (string, string) {
	resp, err := http.Get(jokeURL) //nolint:noctx
	if err != nil {
		responseTag := "no jokes"
		return responseTag, util.GetMessage(locale, responseTag)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		responseTag := "no jokes"
		return responseTag, util.GetMessage(locale, responseTag)
	}

	joke := &Joke{}

	err = json.Unmarshal(body, joke)
	if err != nil {
		responseTag := "no jokes"
		return responseTag, util.GetMessage(locale, responseTag)
	}

	jokeStr := joke.Setup + " " + joke.Punchline

	return JokesTag, fmt.Sprintf(response, jokeStr)
}
