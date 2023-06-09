package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/katallaxie/autobot/util"
)

const adviceURL = "https://api.adviceslip.com/advice"

// AdvicesTag is the intent tag for its module
var AdvicesTag = "advices"

// AdvicesReplacer replaces the pattern contained inside the response by a random advice from the api
// specified by the adviceURL.
// See modules/modules.go#Module.Replacer() for more details.
func AdvicesReplacer(locale, entry, response, _ string) (string, string) {
	resp, err := http.Get(adviceURL) //nolint:noctx
	if err != nil {
		responseTag := "no advices"
		return responseTag, util.GetMessage(locale, responseTag)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		responseTag := "no advices"
		return responseTag, util.GetMessage(locale, responseTag)
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	advice := result["slip"].(map[string]interface{})["advice"]

	return AdvicesTag, fmt.Sprintf(response, advice)
}
