package spotting_test

/* Disabled as it requires spots to be available
func TestClient(t *testing.T) {
	cfg := pota.POTAConfig{}
	client := pota.NewPOTAClient(cfg)
	client.Run()
	time.Sleep(2 * time.Second)
	client.Close()
	for s := range client.Spots {
		fmt.Printf("%v\n", s)
	}
}
*/
