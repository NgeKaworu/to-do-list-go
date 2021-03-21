package resultor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// RetOk 成功处理器
func RetOk(w http.ResponseWriter, data interface{}) {
	res := map[string]interface{}{
		"ok":   true,
		"data": data,
	}
	b, err := json.Marshal(res)
	if err != nil {
		RetFail(w, err.Error())
		return
	}

	fmt.Fprint(w, string(b))
}

// RetFail 失败处理器
func RetFail(w http.ResponseWriter, errMsg string) {
	res := map[string]interface{}{
		"ok":     false,
		"errMsg": errMsg,
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Println(errMsg)
		log.Println(err.Error())
		log.Printf("%+v", res)
		return
	}

	fmt.Fprint(w, string(b))
}
