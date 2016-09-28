package types

// Stored in the KV store to represent the controllers in the cluster
type Controller struct {
	Label               string `json:"label"`
	Controller          string `json:"controller"`
	SwarmClassicManager string `json:"manager"`
	EngineProxy         string `json:"proxy"`
	ControllerCA        string `json:"controller_ca"` // They should all be the same, but record each just in case
}

// Helper routine to help cope with changes in the above struct and the Label field
func GetIPsFromControllers(rawControllers []Controller) []string {
	controllers := []string{}
	for _, controller := range rawControllers {
		controllers = append(controllers, controller.Label)
	}
	return controllers
}
