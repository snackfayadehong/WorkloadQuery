package Interface

import (
	"WorkloadQuery/logger"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const BaseUrl = `http://172.21.1.248:17980/api/bsp-api-engine-others/`

// KLBRReqHeaders 柯林布瑞业务中台请求头
type KLBRReqHeaders struct {
	appId       string
	timestamp   string
	messageId   string
	signature   string
	contentType string
}

// KLBRBaseResponse 柯林布瑞业务中台接口返回
type KLBRBaseResponse struct {
	AckCode      string `json:"ackCode"`
	AckMessage   string `json:"ackMessage"`
	AckMessageID string `json:"ackMessageId"`
}

// ProductChangeData 材料信息变更返回
type ProductChangeData struct {
	Fhxx []struct {
		Sczt   string `json:"sczt"`
		Scsm   string `json:"scsm"`
		Ypspdm string `json:"ypspdm"`
		Ypdm   string `json:"ypdm"`
	} `json:"fhxx"`
}

// DeliveryData 出库信息返回接口
type DeliveryData struct {
	Fhxx []struct {
		Ckdh string `json:"ckdh"`
		Sczt string `json:"sczt"`
		Scsm string `json:"scsm"`
	} `json:"fhxx"`
}

// RefundData 科室退库/入库信息返回 rkfs 02
type RefundData struct {
	Fhxx []struct {
		Rkdh string `json:"rkdh"`
		Sczt string `json:"sczt"`
		Scsm string `json:"scsm"`
	}
}

//type KLBRResPonse[T any] struct {
//	KLBRBaseResponse
//	Data struct {
//		Fhxx []json.RawMessage `json:"fhxx"`
//	} `json:"data"`
//}

// KLBRRequest 柯林布瑞接口请求参数
type KLBRRequest struct {
	Headers *KLBRReqHeaders
	Url     string
	ReqData []byte
}

// NewReqHeaders  KLBRReqHeaders请求头构造函数，根据入参信息生成请求Headers
func NewReqHeaders(serviceCode string) *KLBRReqHeaders {
	reqHeaders := new(KLBRReqHeaders)
	reqHeaders.appId = "HERP"
	reqHeaders.timestamp = strconv.FormatInt(time.Now().UnixMilli(), 10)
	uuid, _ := uuid.NewUUID()
	reqHeaders.messageId = uuid.String()
	reqHeaders.contentType = "json"
	var signStr = fmt.Sprintf("appId=%s&serviceCode=%s&version=%s&timestamp=%v",
		reqHeaders.appId, serviceCode, "1.0", reqHeaders.timestamp)
	reqHeaders.signature = HMACSHA1(signStr)
	return reqHeaders
}

// HMACSHA1 加密转base64
func HMACSHA1(str string) string {
	keyStr := "123456"
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(str))
	//进行base64编码
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return res
}

// KLBRHttpPost 柯林布瑞业务中台接口统一封装方法
func (k *KLBRRequest) KLBRHttpPost() (*[]byte, error) {
	reqData := bytes.NewBuffer(k.ReqData)
	reqBody, err := http.NewRequest("POST", k.Url, reqData)
	if err != nil {
		return nil, err
	}
	logMsg := fmt.Sprintf("\r\n事件:接口请求跟踪\r\n接口地址:%s\r\n入参:%s\r\n%s\r\n", k.Url, string(k.ReqData), logger.LoggerEndStr)
	logger.AsyncLog(logMsg)
	defer reqBody.Body.Close()
	reqBody.Header.Set("Content-Type", "application/json")
	reqBody.Header.Set("appId", k.Headers.appId)
	reqBody.Header.Set("timestamp", k.Headers.timestamp)
	reqBody.Header.Set("messageId", k.Headers.messageId)
	reqBody.Header.Set("signature", k.Headers.signature)
	reqBody.Header.Set("contentType", k.Headers.contentType)
	rep, err := http.DefaultClient.Do(reqBody)
	if err != nil {
		return nil, err
	}
	repBytes, _ := io.ReadAll(rep.Body)
	return &repBytes, nil
}

// ParseResPonse 处理柯林布瑞接口Fhxx字段不固定
//func ParseResPonse[T any](jsonData []byte) (KLBRBaseResponse, []T, error) {
//	var resP KLBRResPonse[T]
//	err := json.Unmarshal(jsonData, &resP)
//	if err != nil {
//		return KLBRBaseResponse{}, nil, err
//	}
//	var fhxxList []T
//	for _, raw := range resP.Data.Fhxx {
//		var fhxx T
//		err = json.Unmarshal(raw, &fhxx)
//		if err != nil {
//			return resP.KLBRBaseResponse, nil, err
//		}
//		fhxxList = append(fhxxList, fhxx)
//	}
//	return resP.KLBRBaseResponse, fhxxList, nil
//}
