// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-server/v6/app"
	"github.com/mattermost/mattermost-server/v6/model"
)

func (api *API) InitInsights() {
	// Reactions
	api.BaseRoutes.InsightsForTeam.Handle("/reactions", api.APISessionRequired(minimumProfessionalLicense(rejectGuests(getTopReactionsForTeamSince)))).Methods("GET")
	api.BaseRoutes.InsightsForUser.Handle("/reactions", api.APISessionRequired(minimumProfessionalLicense(rejectGuests(getTopReactionsForUserSince)))).Methods("GET")

	// Channels
	api.BaseRoutes.InsightsForTeam.Handle("/channels", api.APISessionRequired(minimumProfessionalLicense(rejectGuests(getTopChannelsForTeamSince)))).Methods("GET")
	api.BaseRoutes.InsightsForUser.Handle("/channels", api.APISessionRequired(minimumProfessionalLicense(rejectGuests(getTopChannelsForUserSince)))).Methods("GET")
}

// Top Reactions

func getTopReactionsForTeamSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, err := c.App.GetTeam(c.Params.TeamId)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topReactionList, err := c.App.GetTopReactionsForTeamSince(c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topReactionList)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopReactionsForTeamSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func getTopReactionsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.Params.TeamId = r.URL.Query().Get("team_id")

	// TeamId is an optional parameter
	if c.Params.TeamId != "" {
		if !model.IsValidId(c.Params.TeamId) {
			c.SetInvalidURLParam("team_id")
			return
		}

		team, teamErr := c.App.GetTeam(c.Params.TeamId)
		if teamErr != nil {
			c.Err = teamErr
			return
		}

		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
			c.SetPermissionError(model.PermissionViewTeam)
			return
		}
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topReactionList, err := c.App.GetTopReactionsForUserSince(c.AppContext.Session().UserId, c.Params.TeamId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topReactionList)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopReactionsForUserSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

// Top Channels

func getTopChannelsForTeamSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireTeamId()
	if c.Err != nil {
		return
	}

	team, err := c.App.GetTeam(c.Params.TeamId)
	if err != nil {
		c.Err = err
		return
	}

	if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
		c.SetPermissionError(model.PermissionViewTeam)
		return
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topChannels, err := c.App.GetTopChannelsForTeamSince(c.Params.TeamId, c.AppContext.Session().UserId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})
	if err != nil {
		c.Err = err
		return
	}

	topChannels.PostCountByDuration, err = postCountByDayViewModel(c.App, topChannels, startTime, c.Params.TimeRange, nil)
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topChannels)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopChannelsForTeamSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func getTopChannelsForUserSince(c *Context, w http.ResponseWriter, r *http.Request) {
	c.Params.TeamId = r.URL.Query().Get("team_id")

	// TeamId is an optional parameter
	if c.Params.TeamId != "" {
		if !model.IsValidId(c.Params.TeamId) {
			c.SetInvalidURLParam("team_id")
			return
		}

		team, teamErr := c.App.GetTeam(c.Params.TeamId)
		if teamErr != nil {
			c.Err = teamErr
			return
		}

		if !c.App.SessionHasPermissionToTeam(*c.AppContext.Session(), team.Id, model.PermissionViewTeam) {
			c.SetPermissionError(model.PermissionViewTeam)
			return
		}
	}

	startTime, err := model.GetStartUnixMilliForTimeRange(c.Params.TimeRange)
	if err != nil {
		c.Err = err
		return
	}

	topChannels, err := c.App.GetTopChannelsForUserSince(c.AppContext.Session().UserId, c.Params.TeamId, &model.InsightsOpts{
		StartUnixMilli: startTime,
		Page:           c.Params.Page,
		PerPage:        c.Params.PerPage,
	})

	if err != nil {
		c.Err = err
		return
	}

	topChannels.PostCountByDuration, err = postCountByDayViewModel(c.App, topChannels, startTime, c.Params.TimeRange, &c.AppContext.Session().UserId)
	if err != nil {
		c.Err = err
		return
	}

	js, jsonErr := json.Marshal(topChannels)
	if jsonErr != nil {
		c.Err = model.NewAppError("getTopChannelsForUserSince", "api.marshal_error", nil, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

// postCountByDayViewModel expects a list of channels that are pre-authorized for the given user to view.
func postCountByDayViewModel(app app.AppIface, topChannelList *model.TopChannelList, startTime int64, timeRange string, userID *string) (model.ChannelPostCountByDuration, *model.AppError) {
	if len(topChannelList.Items) == 0 {
		return nil, nil
	}
	var postCountsByDay []*model.DurationPostCount
	channelIDs := topChannelList.ChannelIDs()
	var grouping model.PostCountGrouping
	if timeRange == model.TimeRangeToday {
		grouping = model.PostsByHour
	} else {
		grouping = model.PostsByDay
	}
	postCountsByDay, err := app.PostCountsByDuration(channelIDs, startTime, userID, grouping)
	if err != nil {
		return nil, err
	}
	return model.ToDailyPostCountViewModel(postCountsByDay, startTime, model.TimeRangeToNumberDays(timeRange), channelIDs), nil
}