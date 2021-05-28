package tests

import (
	"testing"

	"github.com/yuedun/zhuque/util"
)

func TestJWT(t *testing.T) {
	token, err := util.CreateToken("yuedun", "sec")
	if err != nil {
		t.Log(">>>>>>>>err", err)
	}
	t.Log(token)
}

func TestJWT2string(t *testing.T) {
	tokenstring, err := util.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjI2OTAyNTQsIm9yaWdfaWF0IjoxNjIxODI2MjU0LCJ1c2VyX2lkIjo1LCJ1c2VybmFtZSI6ImhhbGUuaHVvIn0.DWfif4lIYMIErHVd1GtihuIZOgycuR6SU8A_rNxNUp4", "JWTSecret")
	if err != nil {
		t.Log(">>>>>>>>err", err)
	}
	t.Log(tokenstring)
}
