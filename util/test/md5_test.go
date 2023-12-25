package test

import (
	"fmt"
	"testing"

	"github.com/obse4/goCommon/util"
)

func TestMD5(t *testing.T) {
	res := util.MD5("^-^")
	fmt.Printf("md5 %s\n", res)
}
