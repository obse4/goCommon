package test

import (
	"fmt"
	"testing"

	"github.com/obse4/goCommon/util"
)

func TestAESEncode(t *testing.T) {
	aes, err := util.NewAes("HELLOWORLDGOPHER")
	if err != nil {
		fmt.Printf("create aes struct err %s", err.Error())
		t.Fail()
	}

	res, err := aes.EnCode("^-^")

	if err != nil {
		fmt.Printf("encode aes err %s", err.Error())
		t.Fail()
	}

	fmt.Printf("aes encode %s\n", res)
}

func TestAESDecode(t *testing.T) {
	aes, err := util.NewAes("HELLOWORLDGOPHER")
	if err != nil {
		fmt.Printf("create aes struct err %s", err.Error())
		t.Fail()
	}
	res, err := aes.Decode("eLVsUNgHirfIVv+qU5OmKQ==")

	if err != nil {
		fmt.Printf("decode aes err %s", err.Error())
		t.Fail()
	}

	fmt.Printf("aes decode %s\n", res)
}
