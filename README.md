# jsonvalidator
A simple json validator based on json template which support go-playground/validator and regular expression.

### Installation

`go get github.com/seaguest/jsonvalidator`

### Usage

``` 
package main

import (
	"fmt"
	
	"github.com/seaguest/jsonvalidator"
)

const tmpljson = `{
   "id": "string|eq=123456",
   "data": {
      "token": "re|Bearer\\s{1}",
      "password": "re|^123456$"
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
	dst, err := jsonvalidator.Validate([]byte(tmpljson), []byte(srcjson))

	fmt.Println(dst, err)
}

```
