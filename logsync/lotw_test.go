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
