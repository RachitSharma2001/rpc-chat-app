package errhelp

import "log"

func ThrowPortListenErr(err error) {
	log.Fatalf("Unable to listen at port 5000: %v", err)
}

func ThrowServeErr(err error) {
	log.Fatalf("Unable to listen serve: %v", err)
}

func ThrowDialErr(err error) {
	log.Fatalf("Error while dialing port 5000: %v", err)
}

func ThrowCreateStreamErr(err error) {
	log.Fatalf("Error while creating stream: %v", err)
}

func ThrowReceiveMsgErr(err error) {
	log.Fatalf("Error while recieving messages: %v", err)
}
