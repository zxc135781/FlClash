package main

import (
	"encoding/json"
	"unsafe"
)

type Action struct {
	Id     string      `json:"id"`
	Method Method      `json:"method"`
	Data   interface{} `json:"data"`
}

type ActionResult struct {
	Id       string      `json:"id"`
	Method   Method      `json:"method"`
	Data     interface{} `json:"data"`
	Code     int         `json:"code"`
	callback unsafe.Pointer
}

func (result ActionResult) Json() ([]byte, error) {
	data, err := json.Marshal(result)
	return data, err
}

func (result ActionResult) success(data interface{}) {
	result.Code = 0
	result.Data = data
	result.send()
}

func (result ActionResult) error(data interface{}) {
	result.Code = -1
	result.Data = data
	result.send()
}

func requireString(data interface{}, result ActionResult) (string, bool) {
	s, ok := data.(string)
	if !ok {
		result.error("invalid data type: expected string")
		return "", false
	}
	return s, true
}

func requireBool(data interface{}, result ActionResult) (bool, bool) {
	b, ok := data.(bool)
	if !ok {
		result.error("invalid data type: expected bool")
		return false, false
	}
	return b, ok
}

func handleAction(action *Action, result ActionResult) {
	switch action.Method {
	case initClashMethod:
		paramsString, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		result.success(handleInitClash(paramsString))
		return
	case getIsInitMethod:
		result.success(handleGetIsInit())
		return
	case forceGcMethod:
		handleForceGC()
		result.success(true)
		return
	case shutdownMethod:
		result.success(handleShutdown())
		return
	case validateConfigMethod:
		path, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		result.success(handleValidateConfig(path))
		return
	case updateConfigMethod:
		dataStr, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		result.success(handleUpdateConfig([]byte(dataStr)))
		return
	case setupConfigMethod:
		dataStr, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		result.success(handleSetupConfig([]byte(dataStr)))
		return
	case getProxiesMethod:
		result.success(handleGetProxies())
		return
	case changeProxyMethod:
		data, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		handleChangeProxy(data, func(value string) {
			result.success(value)
		})
		return
	case getTrafficMethod:
		data, ok := requireBool(action.Data, result)
		if !ok {
			return
		}
		result.success(handleGetTraffic(data))
		return
	case getTotalTrafficMethod:
		data, ok := requireBool(action.Data, result)
		if !ok {
			return
		}
		result.success(handleGetTotalTraffic(data))
		return
	case resetTrafficMethod:
		handleResetTraffic()
		result.success(true)
		return
	case asyncTestDelayMethod:
		data, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		handleAsyncTestDelay(data, func(value string) {
			result.success(value)
		})
		return
	case getConnectionsMethod:
		result.success(handleGetConnections())
		return
	case closeConnectionsMethod:
		result.success(handleCloseConnections())
		return
	case resetConnectionsMethod:
		result.success(handleResetConnections())
		return
	case getConfigMethod:
		path, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		config, err := handleGetConfig(path)
		if err != nil {
			result.error(err)
			return
		}
		result.success(config)
		return
	case closeConnectionMethod:
		id, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		result.success(handleCloseConnection(id))
		return
	case getExternalProvidersMethod:
		result.success(handleGetExternalProviders())
		return
	case getExternalProviderMethod:
		externalProviderName, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		result.success(handleGetExternalProvider(externalProviderName))
		return
	case updateGeoDataMethod:
		paramsString, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		var params = map[string]string{}
		err := json.Unmarshal([]byte(paramsString), &params)
		if err != nil {
			result.success(err.Error())
			return
		}
		geoType := params["geo-type"]
		geoName := params["geo-name"]
		handleUpdateGeoData(geoType, geoName, func(value string) {
			result.success(value)
		})
		return
	case updateExternalProviderMethod:
		providerName, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		handleUpdateExternalProvider(providerName, func(value string) {
			result.success(value)
		})
		return
	case sideLoadExternalProviderMethod:
		paramsString, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		var params = map[string]string{}
		err := json.Unmarshal([]byte(paramsString), &params)
		if err != nil {
			result.success(err.Error())
			return
		}
		providerName := params["providerName"]
		data := params["data"]
		handleSideLoadExternalProvider(providerName, []byte(data), func(value string) {
			result.success(value)
		})
		return
	case startLogMethod:
		handleStartLog()
		result.success(true)
		return
	case stopLogMethod:
		handleStopLog()
		result.success(true)
		return
	case startListenerMethod:
		result.success(handleStartListener())
		return
	case stopListenerMethod:
		result.success(handleStopListener())
		return
	case getCountryCodeMethod:
		ip, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		handleGetCountryCode(ip, func(value string) {
			result.success(value)
		})
		return
	case getMemoryMethod:
		handleGetMemory(func(value string) {
			result.success(value)
		})
		return
	case crashMethod:
		result.success(true)
		handleCrash()
		return
	case deleteFile:
		path, ok := requireString(action.Data, result)
		if !ok {
			return
		}
		handleDelFile(path, result)
		return
	default:
		nextHandle(action, result)
	}
}
