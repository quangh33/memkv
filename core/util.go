package core

import (
	"fmt"
	"strconv"
)

func PrintBin(x uint64) {
	fmt.Println(strconv.FormatUint(x, 2))
}
