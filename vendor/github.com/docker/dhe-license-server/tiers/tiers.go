package tiers

import "errors"

// License types. The type dictates whether the DHE instance will make license check requests
const (
	Online  = "Online"
	Offline = "Offline"
	Hourly  = "Hourly"
)

// License tiers
// Refer to https://docker.atlassian.net/wiki/pages/viewpage.action?pageId=10093025
const (
	Trial       = "Trial"       // Always Online
	Starter     = "Starter"     // Always Online
	Evaluation  = "Evaluation"  // Always Offline
	Team        = "Team"        // Always Offline
	Production  = "Production"  // Always Offline
	HourlyAWS   = "HourlyAWS"   // Always Hourly
	HourlyAzure = "HourlyAzure" // Always Hourly
)

func GetTypeFromTier(tier string) (string, error) {
	switch tier {
	case Trial:
		return Online, nil
	case Starter:
		return Online, nil
	case Evaluation:
		return Offline, nil
	case Team:
		return Offline, nil
	case Production:
		return Offline, nil
	case HourlyAWS:
		return Hourly, nil
	case HourlyAzure:
		return Hourly, nil
	case "":
		return Offline, nil
	default:
		return "", errors.New("Must provide a valid license tier")
	}
}
