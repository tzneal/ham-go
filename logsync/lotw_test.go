package logsync_test

/* Disabled as we need a valid username/password to test
func TestLOTW(t *testing.T) {
	c := logsync.NewLOTWClient("", "")
	var qslSince time.Time
	log, err := c.QSLReport(qslSince)
	if err != nil {
		t.Fatalf("error retrieving log: %s", err)
	}
	for rec := range log.Records {
		fmt.Println(rec)
	}
}
*/

/* Disabled as we don't want to upload QSOs during testing
func TestLOTWSave(t *testing.T) {
	alog, err := adif.ParseFile("/home/todd/termlog/Apr_2020.adif")
	if err != nil {
		t.Fatalf("error parsing log: %s", err)
	}
	rec := alog.Records[0]

	c := logsync.NewLOTWClient("", "", "/usr/local/bin/tqsl")
	if err := c.UploadQSO(rec); err != nil {
		t.Errorf("error uploading records: %s", err)
	}
}

*/
