package errhelp

import "log"

func ThrowPortListenErr(err error) {
	log.Fatalf("Unable to listen at port 5000: %v", err)
}

func ThrowServeErr(err error) {
	log.Fatalf("Unable to listen serve: %v", err)
}
