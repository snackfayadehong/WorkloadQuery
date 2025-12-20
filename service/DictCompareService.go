package service

import (
	"WorkloadQuery/Interface"
	"WorkloadQuery/controller"
	http2 "WorkloadQuery/http"
	"WorkloadQuery/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DictCompareServiceInstance 导出单例实例
var DictCompareServiceInstance = &DictCompareService{
	ctrl: controller.DictCompareController{},
}

type DictCompareService struct {
	ctrl controller.DictCompareController
}

// HandleCompareRequest 处理前端比对请求
func (s *DictCompareService) HandleCompareRequest(c *gin.Context) {
	res := http2.NewBaseResponse()
	var req struct {
		Keyword string `json:"keyword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		res.Code = 1
		res.Message = "参数解析失败"
		c.JSON(http.StatusOK, res)
	}

	// 1. 获取本地记录列表 (对应 Controller 返回的 *[]model.LocalDictRow)
	locals, err := s.ctrl.GetLocalDictInfo(req.Keyword)

	// 判断是否找到数据
	if err != nil || locals == nil || len(locals) == 0 {
		res.Code = 1
		res.Message = "怡道系统未找到相关材料"
		c.JSON(http.StatusOK, res)
		return
	}

	// 2. 如果存在多个本地 ID，返回列表供用户选择 (解引用指针判断长度)
	if len(locals) > 1 {
		res.Code = 201 // 约定状态码：需要二次选择
		res.Message = "对应怡道系统多个字典信息，请选择具体项进行比对"
		res.Data = locals // 返回整个数组给前端
		c.JSON(http.StatusOK, res)
		return
	}

	// 3. 只有唯一结果时，直接传入该对象指针进行 HIS 比对，避免二次查库
	// 传递具体某一项的地址：&(*locals)[0]
	results, msg, err := s.CompareDictData(&(locals)[0])
	if err != nil {
		res.Code = 1
		res.Message = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	res.Data = results
	res.Message = msg
	c.JSON(http.StatusOK, res)
}

// CompareDictData 核心业务对比逻辑
// 直接接收查好的模型对象，不再接收字符串关键字
func (s *DictCompareService) CompareDictData(local *model.LocalDictRow) ([]model.DictCompareResult, string, error) {

	// 1. 调用 Interface 层的 HIS ICU 通用请求
	// 使用已经查出的 local.Ypdm 构造请求
	icuReq := &Interface.HisIcuRequest{
		Url:     "http://172.21.1.140:8083/integration_inter_icu/wlxt_mis_proc_cx_ypdm",
		ReqData: map[string]string{"ypdm": local.Ypdm},
	}

	respBytes, err := icuReq.CallHisIcuApi()
	if err != nil {
		return nil, "", fmt.Errorf("HIS接口调用失败: %v", err)
	}

	var hisRes Interface.HisDictResponse
	if err := json.Unmarshal(respBytes, &hisRes); err != nil {
		return nil, "", fmt.Errorf("HIS解析失败: %v", err)
	}

	// 2. 停用判断：如果返回长度为 0 代表 HIS 端已停用
	if len(hisRes.Users) == 0 {
		return nil, "HIS系统未找到该材料或已停用", nil
	}
	his := hisRes.Users[0]

	// 3. 定义对比规则 (显示标签 vs 字段Key)
	checkFields := []struct {
		Label string
		Key   string
	}{
		{"通用名", "ypmc"},
		{"规格型号", "ypgg"},
		{"库房单位", "kfdw"},
		{"采购价", "kfcgj"},
		{"零售价", "kflsj"},
		{"库房代码", "kfdm"},
		{"供货单位", "ghdw"},
		{"公司代码", "gsdm"},
	}

	// 4. 执行对比并组装结果
	var results []model.DictCompareResult
	for _, f := range checkFields {
		localVal := s.getReflectVal(local, f.Key)
		hisVal := his[f.Key]

		results = append(results, model.DictCompareResult{
			Label:      f.Label,
			Field:      f.Key,
			LocalValue: localVal,
			HisValue:   hisVal,
			// 统一转为字符串比对，消除浮点数精度干扰
			IsMatch: fmt.Sprintf("%v", localVal) == fmt.Sprintf("%v", hisVal),
		})
	}

	return results, "对比完成", nil
}

// getReflectVal 辅助方法
func (s *DictCompareService) getReflectVal(data *model.LocalDictRow, key string) interface{} {
	switch key {
	case "tymc":
		return data.Ypmc
	case "ypgg":
		return data.Ypgg
	case "kfdw":
		return data.Kfdw
	case "kfcgj":
		return data.Kfcgj
	case "kflsj":
		return data.Kflsj
	case "kfdm":
		return data.Kfdm
	case "ghdw":
		return data.Ghdw
	case "gsdm":
		return data.Gsdm
	default:
		return ""
	}
}
