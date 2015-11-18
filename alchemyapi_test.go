package alchemyapi

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"testing"
)

type Assert struct {
	*testing.T
}

func NewAssert(t *testing.T) *Assert {
	return &Assert{t}
}

func (ast *Assert) Equal(expected, actual interface{}, logs ...interface{}) {
	ast.equalAssert(false, true, expected, actual, logs...)
}
func (ast *Assert) equalAssert(fatal bool, isEqual bool, expected, actual interface{}, logs ...interface{}) {
	expected = normalizeValue(expected)
	actual = normalizeValue(actual)
	if isEqual != (reflect.DeepEqual(expected, actual)) {
		_, file, line, _ := runtime.Caller(2)
		ast.Logf("Caller: %v:%d", file, line)
		if len(logs) > 0 {
			ast.Log(logs...)
		} else {
			if isEqual {
				ast.Log("Values not equal")
			} else {
				ast.Log("Values equal")
			}
		}
		ast.Log("Expected: ", expected)
		ast.Log("Actual: ", actual)
		ast.FailNow()
	}
}
func (ast *Assert) NotNil(value interface{}, logs ...interface{}) {
	ast.nilAssert(false, false, value, logs...)
}

func (ast *Assert) nilAssert(fatal bool, isNil bool, value interface{}, logs ...interface{}) {
	if isNil != (value == nil || reflect.ValueOf(value).IsNil()) {
		_, file, line, _ := runtime.Caller(2)
		ast.Logf("Caller: %v:%d", file, line)
		if len(logs) > 0 {
			ast.Log(logs...)
		} else {
			if isNil {
				ast.Log("value is not nil:", value)
			} else {
				ast.Log("value is nil")
			}
		}
		ast.FailNow()
	}
}

func normalizeValue(value interface{}) interface{} {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(val.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int()
	case reflect.Float32, reflect.Float64:
		return val.Float()
	case reflect.Complex64, reflect.Complex128:
		return val.Complex()
	case reflect.String:
		return val.String()
	case reflect.Bool:
		return val.Bool()
	case reflect.Slice:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return val.Bytes()
		}
	}
	return value
}

func TestAlchemy(t *testing.T) {
	b, _ := ioutil.ReadFile("apikey.txt")
	key := string(b)
	a := New(key, "http://access.alchemyapi.com/calls", &http.Client{})
	assert := NewAssert(t)

	testText := "Bob broke my heart, and then made up this silly sentence to test the Ruby SDK"
	testHtml := "<html><head><title>The best SDK Test | AlchemyAPI</title></head><body><h1>Hello World!</h1><p>My favorite language is Ruby</p></body></html>"
	testUrl := "http://www.nytimes.com/2013/07/13/us/politics/a-day-of-friction-notable-even-for-a-fractious-congress.html?_r=0"

	response, err := a.Entities("text", testText)
	assert.Equal(response["status"], "OK")
	response, err = a.Entities("html", testHtml)
	assert.Equal(response["status"], "OK")
	response, err = a.Entities("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Entities("random", testText)
	assert.NotNil(err)
	response, err = a.Keywords("text", testText)
	assert.Equal(response["status"], "OK")
	response, err = a.Keywords("html", testHtml)
	assert.Equal(response["status"], "OK")
	response, err = a.Keywords("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Keywords("random", testText)
	assert.NotNil(err)
	response, err = a.Concepts("text", testText)
	assert.Equal(response["status"], "OK")
	response, err = a.Concepts("html", testHtml)
	assert.Equal(response["status"], "OK")
	response, err = a.Concepts("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Concepts("random", testText)
	assert.NotNil(err)
	response, err = a.Sentiment("text", testText)
	assert.Equal(response["status"], "OK")
	response, err = a.Sentiment("html", testHtml)
	assert.Equal(response["status"], "OK")
	response, err = a.Sentiment("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Sentiment("random", testText)
	assert.NotNil(err)
	response, err = a.SentimentTargeted("text", testText, "heart")
	assert.Equal(response["status"], "OK")
	response, err = a.SentimentTargeted("html", testHtml, "language")
	assert.Equal(response["status"], "OK")
	response, err = a.SentimentTargeted("url", testUrl, "Congress")
	assert.Equal(response["status"], "OK")
	response, err = a.SentimentTargeted("text", testText, "")
	assert.NotNil(err)
	response, err = a.SentimentTargeted("random", testUrl, "Congress")
	assert.NotNil(err)
	response, err = a.Text("text", testText)
	assert.NotNil(err)
	response, err = a.Text("html", testHtml)
	assert.Equal(response["status"], "OK")
	response, err = a.Text("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.TextRaw("text", testText)
	assert.NotNil(err)
	response, err = a.TextRaw("html", testHtml)
	assert.Equal(response["status"], "OK")
	response, err = a.TextRaw("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Author("text", testText)
	assert.NotNil(err)
	response, err = a.Author("html", testHtml)
	assert.NotNil(err)
	response, err = a.Author("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Title("text", testText)
	assert.NotNil(err)
	response, err = a.Title("html", testHtml)
	assert.Equal(response["status"], "OK")
	response, err = a.Title("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Relations("text", testText)
	assert.Equal(response["status"], "OK")
	response, err = a.Relations("html", testHtml)
	assert.Equal(response["status"], "OK")
	response, err = a.Relations("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Relations("random", testText)
	assert.NotNil(err)
	response, err = a.Category("text", testText)
	assert.Equal(response["status"], "OK")
	response, err = a.Category("html", testHtml, map[string][]string{"url": []string{"test"}})
	assert.Equal(response["status"], "OK")
	response, err = a.Category("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Category("random", testText)
	assert.NotNil(err)
	response, err = a.Feeds("text", testText)
	assert.NotNil(err)
	response, err = a.Feeds("html", testHtml, map[string][]string{"url": []string{"test"}})
	assert.Equal(response["status"], "OK")
	response, err = a.Feeds("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Microformats("text", testText)
	assert.NotNil(err)
	response, err = a.Microformats("html", testHtml, map[string][]string{"url": []string{"test"}})
	assert.Equal(response["status"], "OK")
	response, err = a.Microformats("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Taxonomy("text", testText)
	assert.Equal(response["status"], "OK")
	response, err = a.Taxonomy("url", testUrl)
	assert.Equal(response["status"], "OK")
	response, err = a.Taxonomy("html", testHtml, map[string][]string{"url": []string{"test"}})
	assert.Equal(response["status"], "OK")
	response, err = a.Taxonomy("random", testText)
	assert.NotNil(err)
	response, err = a.Combined("html", testHtml, map[string][]string{"url": []string{"test"}})
	assert.NotNil(err)
	response, err = a.Combined("text", testText)
	assert.Equal(response["status"], "OK")
	response, err = a.Combined("url", testUrl)
	assert.Equal(response["status"], "OK")
}
