package interfaces

import (
	"cs/stats"
	"fmt"
	"ocgophers/domain"
	"os"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/varstr/uaparser"
	"bytes"
)

const (
	PathParamID = "id"
)

var _host = getHost()

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
	t := stats.NewTimer(tagsToMetricName("requests.response-time", tags))
	timer := t.Start()
	chain.ProcessFilter(request, response)
	timer.Stop()
	stats.GetOrRegisterMeter(tagsToMetricName("requests.handled", tags)).Mark(1)
}

func metricTags(request *restful.Request) map[string]string {
	userBrowser, userOS := userAgentClean(request.Request.UserAgent())

	stats := map[string]string{
		"browser":  userBrowser,
		"os":       userOS,
		"endpoint": request.Request.URL.Path,
		"host": _host,
	}
	return stats
}

func getHost() string {
	host, err := os.Hostname()
	if err == nil {
		if idx := strings.IndexByte(host, '.'); idx > 0 {
			host = host[:idx]
		}
	}
	return host
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
	var keyOrder = make([]string, 0, 4)
	if _, ok := tags["host"]; ok {
		keyOrder = append(keyOrder, "host")
	}
	keyOrder = append(keyOrder, "endpoint", "os", "browser")

	buf := &bytes.Buffer{}
	buf.WriteString(prefix)
	for _, k := range keyOrder {
		buf.WriteByte('.')
		v, ok := tags[k]
		if !ok || v == "" {
			buf.WriteString("no-")
			buf.WriteString(k)
			continue
		}
		cleanKeyPart(buf, v)
	}

	return buf.String()
}

func cleanKeyPart(buf *bytes.Buffer, v string) {
	for i := 0; i < len(v); i++ {
		switch c := v[i]; c {
		case '[', '{', '}', '/', '\\', ':', ' ', '\t', '.':
			buf.WriteByte('-')
		default:
			buf.WriteByte(c)
		}
	}
}
