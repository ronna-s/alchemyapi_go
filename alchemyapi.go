package alchemyapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"runtime"
)

type (
	AlchemyAPI struct {
		Endpoints map[string]map[string]string
	}
	alchemy struct {
		api        *AlchemyAPI
		key        string
		baseUrl    string
		httpClient *http.Client
	}
	result map[string]interface{}
)

var api AlchemyAPI

func init() {
	_, currentFilename, _, _ := runtime.Caller(0)
	file, err := ioutil.ReadFile(path.Join(path.Dir(currentFilename), "endpoints.json"))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &api.Endpoints)
	if err != nil {
		panic(err)
	}
}

func New(key string, baseUrl string, httpClient *http.Client) *alchemy {
	return &alchemy{api: &api, key: key, baseUrl: baseUrl, httpClient: httpClient}
}
func (a *alchemy) analyze(action string, flavor string, data string, options ...url.Values) (result, error) {
	var opts url.Values
	if len(options) != 0 {
		opts = options[0]
	} else {
		opts = url.Values{}
	}
	if _, ok := api.Endpoints[action][flavor]; !ok {
		return nil, fmt.Errorf("%s analysis for %s not available", action, flavor)
	}
	opts[flavor] = []string{data}
	return a.Analyze(a.api.Endpoints[action][flavor], opts)
}

func (a *alchemy) Analyze(ep string, options url.Values) (result, error) {
	var (
		v result
	)
	targetUrl := a.baseUrl + ep
	options["apikey"] = []string{a.key}
	options["outputMode"] = []string{"json"}
	response, err := a.httpClient.PostForm(targetUrl, options)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)
	if response != nil {
		response.Body.Close()
	}
	err = json.Unmarshal(body, &v)
	if v["status"] == "ERROR" {
		err = fmt.Errorf("%v", v["statusInfo"])
	}
	return v, err
}

// Calculates the sentiment for text, a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/sentiment-analysis/
// For the docs, please refer to: http://www.alchemyapi.com/api/sentiment-analysis/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// showSourceText -> 0: disabled (default), 1: enabled
// It returns the response as an interface
func (a *alchemy) Sentiment(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("sentiment", flavor, data, options...)
}

// Calculates the targeted sentiment for text, a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/sentiment-analysis/
// For the docs, please refer to: http://www.alchemyapi.com/api/sentiment-analysis/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// target -> the word or phrase to run sentiment analysis on.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// showSourceText	-> 0: disabled, 1: enabled
// It returns the response as an interface
func (a *alchemy) SentimentTargeted(flavor string, data string, target string, options ...url.Values) (result, error) {
	var opts url.Values
	if len(options) != 0 {
		opts = options[0]
	} else {
		opts = url.Values{}
	}
	if target == "" {
		return nil, errors.New("targeted sentiment requires a non-null target")
	}
	opts["target"] = []string{target}
	return a.analyze("sentiment_targeted", flavor, data, opts)
}

// Extracts the entities for text, a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/entity-extraction/
// For the docs, please refer to: http://www.alchemyapi.com/api/entity-extraction/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// disambiguate -> disambiguate entities (i.e. Apple the company vs. apple the fruit). 0: disabled, 1: enabled (default)
// linkedData -> include linked data on disambiguated entities. 0: disabled, 1: enabled (default)
// coreference -> resolve coreferences (i.e. the pronouns that correspond to named entities). 0: disabled, 1: enabled (default)
// quotations -> extract quotations by entities. 0: disabled (default), 1: enabled.
// sentiment -> analyze sentiment for each entity. 0: disabled (default), 1: enabled. Requires 1 additional API transction if enabled.
// showSourceText -> 0: disabled (default), 1: enabled
// maxRetrieve -> the maximum number of entities to retrieve (default: 50)
// It returns the response as an interface
func (a *alchemy) Entities(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("entities", flavor, data, options...)
}

// Extracts the author from a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/author-extraction/
// For the docs, please refer to: http://www.alchemyapi.com/api/author-extraction/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// none
// It returns the response as an interface
func (a *alchemy) Author(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("author", flavor, data, options...)
}

// Extracts the keywords from text, a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/keyword-extraction/
// For the docs, please refer to: http://www.alchemyapi.com/api/keyword-extraction/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// keywordExtractMode -> normal (default), strict
// sentiment -> analyze sentiment for each keyword. 0: disabled (default), 1: enabled. Requires 1 additional API transaction if enabled.
// showSourceText -> 0: disabled (default), 1: enabled.
// maxRetrieve -> the max number of keywords returned (default: 50)
// It returns the response as an interface
func (a *alchemy) Keywords(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("keywords", flavor, data, options...)
}

// Tags the concepts for text, a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/concept-tagging/
// For the docs, please refer to: http://www.alchemyapi.com/api/concept-tagging/
// Available Options:
// maxRetrieve -> the maximum number of concepts to retrieve (default: 8)
// linkedData -> include linked data, 0: disabled, 1: enabled (default)
// showSourceText -> 0:disabled (default), 1: enabled
// It returns the response as an interface

func (a *alchemy) Concepts(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("concepts", flavor, data, options...)
}

// Categorizes the text for text, a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/text-categorization/
// For the docs, please refer to: http://www.alchemyapi.com/api/text-categorization/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// showSourceText -> 0: disabled (default), 1: enabled
// It returns the response as an interface

func (a *alchemy) Category(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("category", flavor, data, options...)
}

// Extracts the relations for text, a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/relation-extraction/
// For the docs, please refer to: http://www.alchemyapi.com/api/relation-extraction/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// sentiment -> 0: disabled (default), 1: enabled. Requires one additional API transaction if enabled.
// keywords -> extract keywords from the subject and object. 0: disabled (default), 1: enabled. Requires one additional API transaction if enabled.
// entities -> extract entities from the subject and object. 0: disabled (default), 1: enabled. Requires one additional API transaction if enabled.
// requireEntities -> only extract relations that have entities. 0: disabled (default), 1: enabled.
// sentimentExcludeEntities -> exclude full entity name in sentiment analysis. 0: disabled, 1: enabled (default)
// disambiguate -> disambiguate entities (i.e. Apple the company vs. apple the fruit). 0: disabled, 1: enabled (default)
// linkedData -> include linked data with disambiguated entities. 0: disabled, 1: enabled (default).
// coreference -> resolve entity coreferences. 0: disabled, 1: enabled (default)
// showSourceText -> 0: disabled (default), 1: enabled.
// maxRetrieve -> the maximum number of relations to extract (default: 50, max: 100)
// It returns the response as an interface

func (a *alchemy) Relations(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("relations", flavor, data, options...)
}

// Detects the language for text, a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/api/language-detection/
// For the docs, please refer to: http://www.alchemyapi.com/products/features/language-detection/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// none
// It returns the response as an interface

func (a *alchemy) Language(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("language", flavor, data, options...)
}

// Extracts the cleaned text (removes ads, navigation, etc.) for text, a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/text-extraction/
// For the docs, please refer to: http://www.alchemyapi.com/api/text-extraction/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// useMetadata -> utilize meta description data, 0: disabled, 1: enabled (default)
// extractLinks -> include links, 0: disabled (default), 1: enabled.
// It returns the response as an interface
func (a *alchemy) Text(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("text", flavor, data, options...)
}

// Extracts the raw text (includes ads, navigation, etc.) for a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/text-extraction/
// For the docs, please refer to: http://www.alchemyapi.com/api/text-extraction/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// none
// It returns the response as an interface

func (a *alchemy) TextRaw(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("text_raw", flavor, data, options...)
}

// Extracts the title for a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/text-extraction/
// For the docs, please refer to: http://www.alchemyapi.com/api/text-extraction/
// INPUT:
// flavor -> which version of the call, i.e. text, url or html.
// data -> the data to analyze, either the text, the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// useMetadata -> utilize title info embedded in meta data, 0: disabled, 1: enabled (default)
// It returns the response as an interface
func (a *alchemy) Title(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("title", flavor, data, options...)
}

// Parses the microformats for a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/microformats-parsing/
// For the docs, please refer to: http://www.alchemyapi.com/api/microformats-parsing/
// INPUT:
// flavor -> which version of the call, i.e.  url or html.
// data -> the data to analyze, either the the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// none
// It returns the response as an interface

func (a *alchemy) Microformats(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("microformats", flavor, data, options...)
}

// Detects the RSS/ATOM feeds for a URL or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/feed-detection/
// For the docs, please refer to: http://www.alchemyapi.com/api/feed-detection/

// INPUT:
// flavor -> which version of the call, i.e.  url or html.
// data -> the data to analyze, either the the url or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.

// Available Options:
// none
// It returns the response as an interface
func (a *alchemy) Feeds(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("feeds", flavor, data, options...)
}

// Categorizes the text for a URL, text or HTML.
// For an overview, please refer to: http://www.alchemyapi.com/products/features/text-categorization/
// For the docs, please refer to: http://www.alchemyapi.com/api/taxonomy/
// INPUT:
// flavor -> which version of the call, i.e.  url, text or html.
// data -> the data to analyze, either the the url, text or html code.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// showSourceText -> 0: disabled (default), 1: enabled.
// It returns the response as an interface
func (a *alchemy) Taxonomy(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("taxonomy", flavor, data, options...)
}

// Combined call (see options below for available extractions) for a URL or text.
// INPUT:
// flavor -> which version of the call, i.e.  url or text.
// data -> the data to analyze, either the the url or text.
// options -> various parameters that can be used to adjust how the API works, see below for more info on the available options.
// Available Options:
// extract -> VALUE,VALUE,VALUE,... (possible VALUEs: page-image,entity,keyword,title,author,taxonomy,concept,relation,doc-sentiment)
// extractMode -> (only applies when 'page-image' VALUE passed to 'extract' option)
// 		trust-metadata: less CPU-intensive, less accurate
// 		always-infer: more CPU-intensive, more accurate
// disambiguate -> whether to disambiguate detected entities, 0: disabled, 1: enabled (default)
// linkedData -> whether to include Linked Data content links with disambiguated entities, 0: disabled, 1: enabled (default). disambiguate must be enabled to use this.
// coreference -> whether to he/she/etc coreferences into detected entities, 0: disabled, 1: enabled (default)
// quotations -> whether to enable quotations extraction, 0: disabled (default), 1: enabled
// sentiment -> whether to enable entity-level sentiment analysis, 0: disabled (default), 1: enabled. Requires one additional API transaction if enabled.
// showSourceText -> 0: disabled (default), 1: enabled.
// maxRetrieve -> maximum number of named entities to extract (default: 50)
// It returns the response as an interface
func (a *alchemy) Combined(flavor string, data string, options ...url.Values) (result, error) {
	return a.analyze("combined", flavor, data, options...)
}
