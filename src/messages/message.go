package messages

type Message struct {
	Id 	int 	`json:"id"`
	Data 	string 	`json:"data"`
	Type	string  `json:"type"`
}
