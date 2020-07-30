package jsonvalidator

import (
	"common/logger"
	"testing"
)

const tmpljson = `
{
  "code": "int|eq=0",
  "id": "re|^123456$",
  "list": [
    "int|eq=1",
    "int|eq=2"
  ],
  "data": {
    "token": "re|Bearer\\s{1}",
    "password": "re|^123456$"
  }
}
`

const srcjson = `
{
  "code": 0,
  "list": [
    1,
    2
  ],
  "id": "123456",
  "ts": "19226715",
  "data": {
    "token": "Bearer 1|ssds",
    "password": "123456",
    "ok": "yes"
  }
}
`

func TestValidate(t *testing.T) {
	err := Validate([]byte(tmpljson), []byte(srcjson))
	logger.Error(err)
}
