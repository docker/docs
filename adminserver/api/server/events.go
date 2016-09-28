package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/manager/schema"
	"github.com/docker/distribution/context"
	"github.com/emicklei/go-restful"
)

func (a *APIServer) getEvents(ctx context.Context, r *restful.Request) responses.APIResponse {
	// Gather request data.
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(false),
		getPagerParams,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	var (
		publishedBefore, publishedAfter *time.Time
		isAdmin                         = false
		queryingUserId                  = ""
	)

	publishedBeforeString := rd.r.QueryParameter("publishedBefore")
	if publishedBeforeString != "" {
		layout := "2006-01-02T15:04:05.000Z"
		parsed, err := time.Parse(layout, publishedBeforeString)
		if err != nil {
			retErr := errors.MakeError(errors.ErrorCodeInvalidParameters, fmt.Sprintf("Error parsing publishedBefore: '%s', publishedBefore should be formatted as: %s", err.Error(), layout))
			return responses.APIError(retErr)
		}
		publishedBefore = &parsed
	}
	publishedAfterString := rd.r.QueryParameter("publishedAfter")
	if publishedAfterString != "" {
		layout := "2006-01-02T15:04:05.000Z"
		parsed, err := time.Parse(layout, publishedAfterString)
		if err != nil {
			retErr := errors.MakeError(errors.ErrorCodeInvalidParameters, fmt.Sprintf("Error parsing publishedAfter: '%s', publishedAfter should be formatted as: %s", err.Error(), layout))
			return responses.APIError(retErr)
		}
		publishedAfter = &parsed
	}

	if !rd.user.IsAnonymous {
		queryingUserId = rd.user.Account.ID

		if rd.user.Account != nil && rd.user.Account.IsAdmin != nil && *rd.user.Account.IsAdmin {
			isAdmin = true
		}
	}

	events, next, err := a.eventMgr.GetEvents(rd.start, rd.limit, publishedBefore, publishedAfter, queryingUserId, rd.r.QueryParameter("actorId"), rd.r.QueryParameter("eventType"), isAdmin)
	if err != nil {
		if err == schema.ErrCannotQueryForOtherUser {
			return responses.APIError(errors.NotAuthorizedError(err.Error()))
		}
		return responses.APIError(errors.InternalError(ctx, err))
	}

	return responses.JSONResponsePage(http.StatusOK, nil, nil, responses.Events{events}, r, next, 0)
}
