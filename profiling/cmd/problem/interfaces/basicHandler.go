package interfaces

import (
	"cs/stats"
	"fmt"
	"ocgophers/domain"
	"os"
	"strconv"
	"strings"

	"regexp"

	"github.com/emicklei/go-restful"
	"github.com/varstr/uaparser"
	"github.com/rcrowley/go-metrics"
)

const (
	PathParamID = "id"
)

var MetricsRegistry = metrics.NewRegistry()

func NewBasicHandler(responses domain.Responses) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/basic").Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	ws.Route(ws.GET("/{id}").
		Param(ws.PathParameter(PathParamID, "ID of user to retrieve").DataType("integer")).
		To(basic(responses)).
		Operation("Basic Get User").
		Doc("Retrieve user information by _id").
		Writes(domain.UserRecord{}),
	)
	ws.Route(ws.GET("/").
		To(basicHello).
		Operation("Basic Get Hello World").
		Doc("Print out Hello World"))
	return ws
}

func basicHello(request *restful.Request, response *restful.Response) {
	response.ResponseWriter.Write([]byte("Hello World"))
}

func basic(responses domain.Responses) restful.RouteFunction {
	numUsers := len(responses)
	return func(request *restful.Request, response *restful.Response) {
		idStr := request.PathParameter(PathParamID)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.WriteAsJson(map[string]interface{}{"Error": err})
			return
		}
		if id < 0 || id > numUsers-1 {
			response.WriteAsJson(map[string]interface{}{"Error": fmt.Sprintf("User id %d does not exist", id)})
			return
		}
		response.WriteAsJson(responses[id])
	}
}

func NewCounterHandler(responses domain.Responses) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/counter").Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	ws.Route(ws.GET("/{id}").
		Param(ws.PathParameter(PathParamID, "ID of user to retrieve").DataType("integer")).
		To(basic(responses)).
		Operation("Counter Get User").
		Doc("Retrieve user information by _id").
		Writes(domain.UserRecord{}),
	).Filter(TrackUserMetrics)
	ws.Route(ws.GET("/").
		Param(ws.PathParameter(PathParamID, "ID of user to retrieve").DataType("integer")).
		To(basicHello).
		Operation("Counter Get Hello World").
		Doc("Retrieve user information by _id").
		Writes(domain.UserRecord{}),
	).Filter(TrackUserMetrics)
	return ws
}

func TrackUserMetrics(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	tags := metricTags(request)
	t := stats.NewTimerWithRegistry(MetricsRegistry, tagsToMetricName("requests.response-time", tags))
	timer := t.Start()
	chain.ProcessFilter(request, response)
	timer.Stop()
	stats.GetOrRegisterMeterWithRegistry(MetricsRegistry, tagsToMetricName("requests.handled", tags)).Mark(1)
}

func metricTags(request *restful.Request) map[string]string {
	userBrowser, userOS := userAgentClean(request.Request.UserAgent())

	stats := map[string]string{
		"browser":  userBrowser,
		"os":       userOS,
		"endpoint": request.Request.URL.Path,
	}
	host, err := os.Hostname()
	if err == nil {
		if idx := strings.IndexByte(host, '.'); idx > 0 {
			host = host[:idx]
		}
	}
	stats["host"] = host
	return stats
}

func userAgentClean(userAgent string) (string, string) {

	browser := "Unknown-Browser"
	os := "Unkonwn-OS"
	ua := uaparser.Parse(userAgent)

	if ua.Browser != nil {
		browser = ua.Browser.Name
	}
	if ua.OS != nil {
		os = ua.OS.Name
	}

	return browser, os
}

func tagsToMetricName(prefix string, tags map[string]string) string {
	var keyOrder []string
	if _, ok := tags["host"]; ok {
		keyOrder = append(keyOrder, "host")
	}
	keyOrder = append(keyOrder, "endpoint", "os", "browser")

	parts := []string{prefix}
	for _, k := range keyOrder {
		v, ok := tags[k]
		if !ok || v == "" {
			parts = append(parts, "no-"+k)
			continue
		}
		parts = append(parts, cleanKeyPart(v))
	}

	return strings.Join(parts, ".")
}

var specialChars = regexp.MustCompile(`[{}/\\:\s.]`)


func cleanKeyPart(v string) string {
	return specialChars.ReplaceAllString(strings.TrimPrefix(v, "/"), "-")
}
