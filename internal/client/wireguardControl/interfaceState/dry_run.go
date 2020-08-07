package interfaceState

import "log"

type dry struct{}

func CreateDummy() *dry {
	return &dry{}
}

func (*dry) DeleteInterface() error {
	log.Println("CALL TO DeleteInterface() - dry run")
	return nil
}

func (*dry) CreateOrUpdateInterface(options InterfaceOptions) error {
	log.Println("CALL TO CreateOrUpdateInterface() - dry run")
	return nil
}
