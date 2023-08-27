package core

import (
	"fmt"
	"strconv"
)

func PrintBin(x uint64) {
	fmt.Println(fmt.Sprintf("0b%s", strconv.FormatUint(x, 2)))
}
