package logsync

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/tzneal/ham-go/adif"
	"github.com/tzneal/ham-go/util"
)

type LOTWClient struct {
	username string
	password string
	tqslPath string
}

var ErrNoRecordsUploaded = errors.New("No records to upload")

// NewLOTWClient constructs a new LOTW client that can fetch QSLs using the REST API and upload
// signed QSLs using the tqsl command
func NewLOTWClient(username, password string, tqslPath string) *LOTWClient {
	return &LOTWClient{
		username: username,
		password: password,
		tqslPath: tqslPath,
	}
}

// QSLReport fetches the LOTW QSLs that have been submitted since the time given
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

	buf, _ := ioutil.ReadAll(rsp.Body)
	ioutil.WriteFile("/tmp/lotw", buf, 0644)
	r := util.NewSkipReader(bytes.NewReader(buf), []byte("<PROGRAMID"))

	return adif.Parse(r)
}

// UploadQSOs uploads a batch of QSOs to LOTW using the tqsl command
func (c LOTWClient) UploadQSOs(records []adif.Record) error {
	alog := adif.NewLog()
	alog.AddRecords(records)
	tf, err := ioutil.TempFile("", "lotwupload")
	if err != nil {
		return err
	}
	alog.Normalize()
	if err := alog.Write(tf); err != nil {
		return err
	}
	tf.Close()
	defer os.Remove(tf.Name())

	cmd := exec.Command(c.tqslPath, "-a", "compliant", // don't send duplicates
		"--nodate", // don't ask about dates
		"--upload", // upload
		"--batch",  // and exit
		tf.Name(),
	)

	op, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running tqsl: %s %s", err, string(op))
	}

	return nil
}

// UploadQSOs uploads a single QSO to LOTW using the tqsl command
func (c LOTWClient) UploadQSO(record adif.Record) error {
	return c.UploadQSOs([]adif.Record{record})
}
