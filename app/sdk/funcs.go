package sdk

import "strings"

func UnwrapCallback(callback string) (string, string) {
	out := strings.Split(callback, "_")

	return out[0], out[1]
}

func WrapCallback(callbackType, callbackID string) string {
	return callbackType + "_" + callbackID
}
