package oldfart

import "oldfartscounter/internal/environment"

type (
	OldFart struct {
		SteamId SteamId `json:"steamId" dynamodbav:"steamId"`
		Person  Person  `json:"person" dynamodbav:"person"`
		Nick    string  `json:"nick" dynamodbav:"nick"`
	}

	Person struct {
		FirstName  string `json:"firstName" dynamodbav:"firstName"`
		LastName   string `json:"lastName" dynamodbav:"lastName"`
		TelegramId string `json:"telegramId" dynamodbav:"telegramId"`
	}

	SteamId string
)

func (o *OldFart) TableName() string {
	return environment.GetOldFartsPeopleDDBTable()
}
