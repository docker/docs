package server

import (
	"fmt"
	"strconv"

	"github.com/emicklei/go-restful"

	"github.com/docker/dhe-deploy/adminserver/api/common"
)

func pagerParams(r *restful.Request) (string, uint) {
	start := r.QueryParameter("start")
	limitStr := r.QueryParameter("limit")
	limit, _ := strconv.ParseUint(limitStr, 10, 32)
	if limit == 0 {
		limit = common.DefaultPerPageLimit
	}
	return start, uint(limit)
}

func pagerNext(r *restful.Request, next uint) string {
	if next == 0 {
		return ""
	}
	url := r.Request.URL
	values := url.Query()
	values.Set("start", fmt.Sprintf("%d", next))
	url.RawQuery = values.Encode()
	return url.String()
}
