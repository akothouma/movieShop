package fetch

import (
	"fmt"
	"io"
	"net/http"
)

// Fetch retrieves the response body from the specified URL.
func Fetch(url string) ([]byte, error) {
	body := []byte{}
	var body_err error

	response, fetchErr:= http.Get(url)
	if fetchErr!= nil {
		return nil, fmt.Errorf("error making a get request to the artists api endpoint: %s", artists_err)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, body_err = io.ReadAll(response.Body)
		if body_err != nil {
			return nil, fmt.Errorf("error reading response body: %s", body_err)
		}
	}
	return body, nil
	
}
