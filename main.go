// use for test algorithm
package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/luanngominh/secure-tranfer-file/util"
)

var (
	pub  string
	priv string
)

func init() {
	data, err := base64.StdEncoding.DecodeString(os.Getenv("PRIVATE"))
	if err != nil {
		panic(err)
	}
	priv = string(data)

	data, err = base64.StdEncoding.DecodeString(os.Getenv("PUBLIC"))
	if err != nil {
		panic(err)
	}
	pub = string(data)
}

func main() {
	mess := []byte("meo con xinh xan, say meo meo")
	key := []byte("meoconxinhxinh")

	encryptMsg, _ := util.EncryptWithKey(mess, key)
	msg, _ := util.DecryptWithKey(encryptMsg, key)
	fmt.Println(string(msg))
}
