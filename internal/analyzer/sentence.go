package analyzer

import (
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/katallaxie/autobot/internal/network"
	"github.com/katallaxie/autobot/pkg/modules"
	"github.com/katallaxie/autobot/util"
	"go.uber.org/zap"

	gocache "github.com/patrickmn/go-cache"
)

// A Sentence represents simply a sentence with its content as a string
type Sentence struct {
	Locale  string
	Content string
}

// Result contains a predicted value with its tag and its value
type Result struct {
	Tag   string
	Value float64
}

var userCache = gocache.New(5*time.Minute, 5*time.Minute)

// DontUnderstand contains the tag for the don't understand messages
const DontUnderstand = "don't understand"

type analyzer struct{}

// Analyze ...
func (a *analyzer) Analyze(locale, content string) Sentence {
	s := Sentence{
		Locale:  locale,
		Content: content,
	}
	s.arrange()

	return s
}

// Opt ...
type Opt func(*analyzer)

// NewAnalyzer ...
func NewAnalyzer() *analyzer {
	return &analyzer{}
}

// NewSentence returns a Sentence object where the content has been arranged
func NewSentence(locale, content string) (sentence Sentence) {
	sentence = Sentence{
		Locale:  locale,
		Content: content,
	}
	sentence.arrange()

	return
}

// PredictTag classifies the sentence with the model
func (sentence Sentence) PredictTag(neuralNetwork network.Network, confidence float64) string {
	words, classes, _ := Organize(sentence.Locale)

	// Predict with the model
	predict := neuralNetwork.Predict(sentence.WordsBag(words))

	// Enumerate the results with the intent tags
	resultsTag := make([]Result, 0)
	for i, result := range predict {
		if i >= len(classes) {
			continue
		}
		resultsTag = append(resultsTag, Result{classes[i], result})
	}

	// Sort the results in ascending order
	sort.Slice(resultsTag, func(i, j int) bool {
		return resultsTag[i].Value > resultsTag[j].Value
	})

	LogResults(sentence.Locale, sentence.Content, resultsTag)

	if resultsTag[0].Value <= confidence {
		return DontUnderstand
	}

	return resultsTag[0].Tag
}

// RandomizeResponse takes the entry message, the response tag and the token and returns a random
// message from res/datasets/intents.json where the triggers are applied
func RandomizeResponse(locale, entry, tag, token string) (string, string) {
	if tag == DontUnderstand {
		return DontUnderstand, util.GetMessage(locale, tag)
	}

	// Append the modules intents to the intents from res/datasets/intents.json
	intents := append(SerializeIntents(locale), SerializeModulesIntents(locale)...)

	for _, intent := range intents {
		if intent.Tag != tag {
			continue
		}

		// Reply a "don't understand" message if the context isn't correct
		cacheTag, _ := userCache.Get(token)
		if intent.Context != "" && cacheTag != intent.Context {
			return DontUnderstand, util.GetMessage(locale, DontUnderstand)
		}

		// Set the actual context
		userCache.Set(token, tag, gocache.DefaultExpiration)

		// Choose a random response in intents
		response := intent.Responses[0]
		if len(intent.Responses) > 1 {
			rand.Seed(time.Now().UnixNano())
			response = intent.Responses[rand.Intn(len(intent.Responses))]
		}

		// And then apply the triggers on the message
		return modules.ReplaceContent(locale, tag, entry, response, token)
	}

	return DontUnderstand, util.GetMessage(locale, DontUnderstand)
}

// Calculate send the sentence content to the neural network and returns a response with the matching tag
func (sentence Sentence) Calculate(cache gocache.Cache, neuralNetwork network.Network, token string, confidence float64) (string, string) {
	tag, found := cache.Get(sentence.Content)

	// Predict tag with the neural network if the sentence isn't in the cache
	if !found {
		tag = sentence.PredictTag(neuralNetwork, confidence)
		cache.Set(sentence.Content, tag, gocache.DefaultExpiration)
	}

	return RandomizeResponse(sentence.Locale, sentence.Content, tag.(string), token)
}

// LogResults print in the console the sentence and its tags sorted by prediction
func LogResults(locale, entry string, results []Result) {
	// If NO_LOGS is present, then don't print the given messages
	if os.Getenv("NO_LOGS") == "1" {
		return
	}

	matches := make([]interface{}, 0)
	for _, result := range results {
		// Arbitrary choice of 0.004 to have less tags to show
		if result.Value < 0.004 {
			continue
		}

		matches = append(matches, result.Tag, result.Value)
	}

	zap.S().Infow("matches sentences", matches...)
}
