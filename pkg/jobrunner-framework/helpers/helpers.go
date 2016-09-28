package helpers

import (
	"math"
	"strconv"

	"github.com/docker/orca/enzi/api"
	"github.com/emicklei/go-restful"
)

// XXX: copied from ucp/enzi/api/server/internal/helpers

// PageParams gets the page start offset and page results limit from the given
// request using query parameters with the given names.
func PageParams(r *restful.Request, startParamName, limitParamName string) (start string, limit uint) {
	start = r.QueryParameter(startParamName)
	limitStr := r.QueryParameter(limitParamName)

	parsedLimit, _ := strconv.ParseUint(limitStr, 10, 32)
	if parsedLimit == 0 {
		parsedLimit = api.DefaultPerPageLimit
	}

	limit = uint(parsedLimit)
	if limit > api.MaxPerPageLimit {
		limit = api.MaxPerPageLimit
	}

	return start, limit
}

// ParseOffsetAsUint is used for paging when the offset should be a numbered
// offset rather than a string identifie
func ParseOffsetAsUint(offsetStr string) uint {
	offset, _ := strconv.ParseUint(offsetStr, 10, 32)
	if offset > math.MaxUint32 {
		offset = math.MaxUint32
	}

	return uint(offset)
}
