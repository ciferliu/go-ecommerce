package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// RequestClient represents the request come from which client.
// Every HTTP request must be sent with a 'request_client' header,
// and with the header, it is possible to serve multi clients on a same HTTP endpoint.
//Supported clients are:[admin, merchant_web, merchant_ios, merchant_andriod, consumer_web, consumer_ios, consumer_andriod]
type RequestClient int8

const (
	ADMIN            RequestClient = 11
	MERCHANT_WEB     RequestClient = 21
	MERCHANT_IOS     RequestClient = 22
	MERCHANT_ANDRIOD RequestClient = 23
	CONSUMER_WEB     RequestClient = 31
	CONSUMER_IOS     RequestClient = 32
	CONSUMER_ANDRIOD RequestClient = 33
)

func (c *RequestClient) IsAdmin() bool {
	return *c > 10 && *c < 20
}

func (c *RequestClient) IsMerchant() bool {
	return *c > 20 && *c < 30
}

func (c *RequestClient) IsConsumer() bool {
	return *c > 30 && *c < 40
}

func GetRequestClient(r *http.Request) (rc RequestClient, err error) {
	client := strings.ToLower(r.Header.Get("request_client"))
	if len(client) == 0 {
		err = errors.New("there is no request_client in request header")
		return
	}

	switch client {
	case "admin":
		rc = ADMIN
	case "merchant_web":
		rc = MERCHANT_WEB
	case "merchant_ios":
		rc = MERCHANT_IOS
	case "merchant_andriod":
		rc = MERCHANT_ANDRIOD
	case "consumer_web":
		rc = CONSUMER_WEB
	case "consumer_ios":
		rc = CONSUMER_IOS
	case "consumer_andriod":
		rc = CONSUMER_ANDRIOD
	default:
		err = fmt.Errorf("request client [%s] is not supported", client)
	}
	return
}
