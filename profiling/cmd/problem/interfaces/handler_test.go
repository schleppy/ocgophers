package interfaces

import (
	"net/http"
	"testing"

	"github.com/emicklei/go-restful"
)

var requestMetricsTable = []struct {
	Endpoint  string
	UserAgent string
	Prefix    string
	Expected  string
}{
	{
		"/testEndpoint",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.81 Safari/537.36",
		"my-prefix",
		"my-prefix.ubuntu.testEndpoint.Mac-OS.Chrome",
	},
	{
		"/basic",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.81 Safari/537.36",
		"my-other-prefix",
		"my-other-prefix.ubuntu.basic.Mac-OS.Chrome",
	},
	{
		"/basic/2",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.81 Safari/537.36",
		"my-other-prefix",
		"my-other-prefix.ubuntu.basic-2.Mac-OS.Chrome",
	},
}

func TestMetricTags(t *testing.T) {
	for _, requestData := range requestMetricsTable {
		_req, _ := http.NewRequest("GET", requestData.Endpoint, nil)
		_req.Header.Set("User-Agent", requestData.UserAgent)
		req := &restful.Request{
			Request: _req,
		}
		out := tagsToMetricName(requestData.Prefix, metricTags(req))
		if out != requestData.Expected {
			t.Errorf("Out: %v does not match Expected: %v", out, requestData.Expected)
		}
	}
}

func BenchmarkTagsToMetricsName(b *testing.B) {
	requestData := requestMetricsTable[0]
	_req, _ := http.NewRequest("GET", requestData.Endpoint, nil)
	_req.Header.Set("User-Agent", requestData.UserAgent)
	req := &restful.Request{
		Request: _req,
	}
	for i := 0; i < b.N; i++ {
		tagsToMetricName(requestData.Prefix, metricTags(req))
	}
}
