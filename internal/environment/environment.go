package environment

import "os"

func LookupVariable(key string) (value string, exists bool) {
	return os.LookupEnv(key)
}

func GetVariable(key string, defaultValue ...string) (value string) {
	value, exists := LookupVariable(key)
	if !exists {
		if len(defaultValue) > 0 {
			value = defaultValue[0]
		}
	}
	return value
}
