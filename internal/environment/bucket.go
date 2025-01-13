package environment

import (
	"fmt"
)

func GetGamesS3Bucket() string {
	return fmt.Sprintf("%s-games", GetNamespace())
}
