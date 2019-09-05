# jsonvalidator
A simple json validator based on json template which is in regular expression.

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
	dst, err := jsonvalidator.ValidateJson(tmpljson, srcjson)

	fmt.Println(dst, err)
}

```
