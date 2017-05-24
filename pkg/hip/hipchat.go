package hip

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func checkReponseCode(r *http.Response, err error) error {
	if err != nil {
		return err
	}
	if r.StatusCode/100 != 2 {
		responseBody, _ := ioutil.ReadAll(r.Body)
		_ = r.Body.Close()
		return errors.New(fmt.Sprintf("Unexpected http code %d (%s) %s", r.StatusCode, r.Status, responseBody))
	}
	return nil
}
