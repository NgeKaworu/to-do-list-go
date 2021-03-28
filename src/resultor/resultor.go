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
		RetFail(w, err)
		return
	}

	fmt.Fprint(w, string(b))
}

// RetFail 失败处理器
func RetFail(w http.ResponseWriter, e error) {
	errMsg := e.Error()

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

// RetOkWithTotal 成功处理器
func RetOkWithTotal(w http.ResponseWriter, data interface{}, total int64) {
	res := map[string]interface{}{
		"ok":    true,
		"data":  data,
		"total": total,
	}
	b, err := json.Marshal(res)
	if err != nil {
		RetFail(w, err)
		return
	}

	fmt.Fprint(w, string(b))
}
