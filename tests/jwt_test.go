package tests

import (
	"fmt"
	"testing"

	"github.com/yuedun/zhuque/util"
)

func TestJWT(t *testing.T) {
	token, err := util.CreateToken("yuedun", "sec")
	if err != nil {
		fmt.Println(">>>>>>>>err", err)
	}
	fmt.Println(token)
}

func TestJWT2string(t *testing.T) {
	tokenstring, err := util.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDY0NDU5NzcsInVpZCI6Inl1ZWR1biJ9.sTFWe9ZLprCLS5luOm43o-Xh65lE-vujTyXXhcSP7ic", "sec")
	if err != nil {
		fmt.Println(">>>>>>>>err", err)
	}
	fmt.Println(tokenstring)
}
