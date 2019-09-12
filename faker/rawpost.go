package faker

import "encoding/json"

type rawpost struct {
	Id   int
	Data string
	Raw  string
}

func (rawpost) TableName() string {
	return "rawpost"
}

func newRawpost(f fake) (r rawpost, err error) {
	var data []byte
	data, err = json.Marshal(f)
	if err != nil {
		return
	}
	r = rawpost{Data: string(data), Raw: string(data)}
	return
}
