package core

import (
	"errors"
	"fmt"
)

func getType(te uint8) uint8 {
	return te & 0b11110000
}

func getEncoding(te uint8) uint8 {
	return te & 0b00001111
}

func assertType(te uint8, t uint8) error {
	if getType(te) != t {
		return errors.New(fmt.Sprintf("operation is not permitted on type %d", getType(te)))
	}
	return nil
}

func assertEncoding(te uint8, e uint8) error {
	if getEncoding(te) != e {
		return errors.New(fmt.Sprintf("operation is not permitted on encoding %d", getEncoding(te)))
	}
	return nil
}
