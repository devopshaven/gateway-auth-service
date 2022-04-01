package pb

import (
	"errors"
	"strings"
)

var ErrNoContentType = errors.New("no content type in request")

func (x *AuthorizeRequest) ContentType() (string, error) {
	for _, h := range x.RequestHeaders {
		if strings.ToLower(h.Name) == "content-type" {
			return h.Values[0], nil
		}
	}

	return "", ErrNoContentType
}
