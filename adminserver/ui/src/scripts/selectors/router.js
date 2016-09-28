'use strict';

export const locationSelector = (state, props) => props.location;

export const getTeamNameFromRouter = (state, props) => props.params.team;
export const getOrgNameFromRouter = (state, props) => props.params.org;
