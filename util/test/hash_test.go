package test

import (
	"fmt"
	"testing"

	"github.com/obse4/goCommon/util"
)

func TestHash(t *testing.T) {
	res := util.String2Hash("^-^")

	fmt.Printf("hash old %s\n", res)

	same, err := util.CompareHash(res, "^-^")

	if err != nil {
		fmt.Printf("hash compare err %s\n", err.Error())
		t.Fail()
	}

	fmt.Printf("the hash compare is %v\n", same)
}
