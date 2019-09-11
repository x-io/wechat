package wapi

import (
	"fmt"
)

// Ticket 请求jsapi_tikcet返回结果
type Ticket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

//GetTicket 从服务器中获取ticket
func GetTicket(accessToken string) (*Ticket, error) {

	url := fmt.Sprintf(TicketURL, accessToken)
	var data Ticket

	if err := getBind(url, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
