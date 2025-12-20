package Interface

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// HisIcuRequest HIS ICU直连请求结构
type HisIcuRequest struct {
	Url     string
	ReqData interface{}
}

// CallHisIcuApi 统一调用方法，未来其它ICU接口只需定义不同的ReqData即可
func (k *HisIcuRequest) CallHisIcuApi() ([]byte, error) {
	jsonData, err := json.Marshal(k.ReqData)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(k.Url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// HisDictResponse 字典对比专用的返回结构
type HisDictResponse struct {
	Users []map[string]interface{} `json:"users"`
}
