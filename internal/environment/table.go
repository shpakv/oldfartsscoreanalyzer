package environment

import "fmt"

func GetOldFartsPeopleDDBTable() string {
	return fmt.Sprintf("%s-old-farts-people", GetNamespace())
}
