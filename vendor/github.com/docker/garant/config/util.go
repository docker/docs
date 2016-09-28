package config

import (
	"strconv"
	"strings"
)

// GetParamString returns the parameter string for the given key.
func (params Parameters) GetParamString(key string) (param string, ok bool) {
	val, ok := params[key]
	if !ok {
		return "", false
	}

	param, ok = val.(string)
	if !ok {
		return "", false
	}

	return param, true
}

// GetParamInt returns the parameter integer for the given key.
func (params Parameters) GetParamInt(key string) (param int, ok bool) {
	if val, ok := params[key]; ok {
		switch param := val.(type) {
		case int:
			return param, true
		case string:
			// May have been from an environment variable.
			if i, err := strconv.Atoi(param); err == nil {
				return i, true
			}
		}
	}

	return 0, false
}

// GetParamBool returns the parameter boolean for the given key.
func (params Parameters) GetParamBool(key string) (param bool, ok bool) {
	if val, ok := params[key]; ok {
		switch param := val.(type) {
		case bool:
			return param, true
		case string:
			return !strings.EqualFold("false", param), true
		case int:
			return param != 0, true
		}
	}

	return false, false
}

// GetParamMap returns the parameter map for the given key.
func (params Parameters) GetParamMap(key string) (param map[string]interface{}, ok bool) {
	val, ok := params[key]
	if !ok {
		return nil, false
	}

	param, ok = val.(map[string]interface{})
	if !ok {
		return nil, false
	}

	return param, true
}

// StringOpt returns the keyed string parameter or uses the default value.
func (params Parameters) StringOpt(key, def string) string {
	if val, ok := params.GetParamString(key); ok {
		return val
	}

	return def
}

// IntOpt returns the keyed integer parameter or uses the default value.
func (params Parameters) IntOpt(key string, def int) int {
	if val, ok := params.GetParamInt(key); ok {
		return val
	}

	return def
}

// BoolOpt returns the keyed boolean parameter or uses the default value.
func (params Parameters) BoolOpt(key string, def bool) bool {
	if val, ok := params.GetParamBool(key); ok {
		return val
	}

	return def
}

// MapOpt returns the keyed map parameter or uses the default value.
func (params Parameters) MapOpt(key string, def map[string]interface{}) map[string]interface{} {
	if val, ok := params.GetParamMap(key); ok {
		return val
	}

	return def
}
