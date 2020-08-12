package jsonvalidator

import (
	"fmt"
	"testing"

	_ "github.com/go-playground/validator"
)

const tmpljson = `
{
  "int": "int|eq=0",
  "string": "string|eq=0",
  "re": "re|^123456$",
  "list1": "list|eq=2",
  "list2": [
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
  "int": 0,
  "string": "0",
  "re": "123456",
  "list1": [
    1,
    2
  ],
  "list2": [
    1,
    2
  ],
  "ts": "19226715",
  "data": {
    "token": "Bearer 1|ssds",
    "password": "123456",
    "ok": "yes"
  }
}
`

func TestValidate(t *testing.T) {
	err := Validate([]byte(srcjson), []byte(tmpljson))
	fmt.Println(err)

}
