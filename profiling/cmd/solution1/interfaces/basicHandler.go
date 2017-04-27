package interfaces

import (
	"github.com/emicklei/go-restful"
	"ocgophers/domain"
	"strconv"
	"fmt"
)

const (
	PathParamID = "id"
)

func NewBasicHandler(responses [][]byte) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/basic").Produces(restful.MIME_JSON, "html/text").Consumes(restful.MIME_JSON)
	ws.Route(ws.GET("/{id}").
		Param(ws.PathParameter(PathParamID, "ID of user to retrieve").DataType("integer")).
		To(basic(responses)).
		Operation("Basic Get").
		Doc("Retrieve user information by _id").
		Writes(domain.UserRecord{}),
	)
	return ws
}

func basic(responses [][]byte) restful.RouteFunction {
	numUsers := len(responses)
	return func(request *restful.Request, response *restful.Response) {
		idStr := request.PathParameter(PathParamID)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.WriteAsJson(map[string]interface{}{"Error": err})
			return
		}
		if id < 0 || id > numUsers - 1 {
			response.WriteAsJson(map[string]interface{}{"Error": fmt.Sprintf("User id %d does not exist", id)})
			return
		}
		response.Write(responses[id])
	}

}
