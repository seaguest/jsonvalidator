package jsonvalidator_test

import (
	"common/logger"
)

const tmpljson = `{
   "id": "^123456$",
   "data": {
      "token": "Bearer\\s{1}",
      "password": "^123456$"
   }
}`

const srcjson = `{
   "id": "123456",
   "ts": "19226715",
   "data": {
      "token": "Bearer 1|ssds",
      "password": "123456",
      "ok": "yes"
   }
}`

func main() {
	dst, err := ValidateJson(tmpljson, srcjson)

	logger.Error(dst, err)
}
