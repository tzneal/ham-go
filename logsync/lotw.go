package logsync

import (
	"net/http"
	"net/url"
	"time"

	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/util"
)

type LOTWClient struct {
	username string
	password string
}

func NewLOTWClient(username, password string) *LOTWClient {
	return &LOTWClient{
		username: username,
		password: password,
	}
}

func (c LOTWClient) QSLReport(since time.Time) (*adif.Log, error) {
	u, err := url.Parse("https://lotw.arrl.org/lotwuser/lotwreport.adi")
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("login", c.username)
	q.Set("password", c.password)
	q.Set("qso_query", "1")
	q.Set("qso_qsl", "yes") // QSL only
	if !since.IsZero() {
		q.Set("qso_qslsince", since.Format("2006-01-02"))
	}
	u.RawQuery = q.Encode()

	rsp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	r := util.NewSkipReader(rsp.Body, []byte("<PROGRAMID"))

	return adif.Parse(r)
}
