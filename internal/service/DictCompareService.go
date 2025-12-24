package service

import (
	"SupperSystem/internal/controller"
	"SupperSystem/internal/model"
	http2 "SupperSystem/pkg/http"
	"SupperSystem/pkg/integration"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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
	// 分流
	kLen := len(req.Keyword)
	if kLen != 6 && kLen != 14 {
		res.Code = 1
		res.Message = "材料代码或产品ID不正确"
		c.JSON(http.StatusOK, res)
		return
	}
	_, err := strconv.Atoi(req.Keyword)
	if err != nil {
		res.Code = 1
		res.Message = "非法字符串!"
		c.JSON(http.StatusOK, res)
		return
	}
	isIdQuery := kLen == 6
	// 1. 获取本地记录列表 (对应 Controller 返回的 *[]model.LocalDictRow)
	locals, err := s.ctrl.GetLocalDictInfo(req.Keyword, isIdQuery)

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

	res.Data = gin.H{
		"ProductInfoID": locals[0].ProductInfoID,
		"ypdm":          locals[0].Ypdm,
		"results":       results,
	}
	res.Message = msg
	c.JSON(http.StatusOK, res)
}

// CompareDictData 核心业务对比逻辑
// 直接接收查好的模型对象，不再接收字符串关键字
func (s *DictCompareService) CompareDictData(local *model.LocalDictRow) ([]model.DictCompareResult, string, error) {

	// 1. 调用 integration 层的 HIS ICU 通用请求
	// 使用已经查出的 local.Ypdm 构造请求
	icuReq := &integration.HisIcuRequest{
		Url:     "http://172.21.1.140:8083/integration_inter_icu/wlxt_mis_proc_cx_ypdm",
		ReqData: map[string]string{"xmbh": local.Ypdm},
	}
	respBytes, err := icuReq.CallHisIcuApi()
	if err != nil {
		return nil, "", fmt.Errorf("HIS接口调用失败: %v", err)
	}
	// 转utf8
	utf8Bytes, err := GbkToUtf8(respBytes)
	if err != nil {
		utf8Bytes = respBytes
	}
	// 清理 "\"
	cleanStr := string(utf8Bytes)
	cleanStr = strings.ReplaceAll(cleanStr, "\\", "\\\\")
	// 还原被过度转义的引号
	cleanStr = strings.ReplaceAll(cleanStr, "\\\\\"", "\\\"")
	// 处理末尾可能存在的垃圾字符（截取到最后一个 '}'）
	lastBrace := strings.LastIndex(cleanStr, "}")
	if lastBrace == -1 {
		cleanStr = cleanStr[:lastBrace+1]
	}
	var hisRes integration.HisDictResponse
	if err := json.Unmarshal([]byte(cleanStr), &hisRes); err != nil {
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
		{"产品名称", "ypmc"},
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
		rawLocal := s.getReflectVal(local, f.Key)
		RawHis := his[f.Key]
		// 清洗首位空格
		localVal := strings.TrimSpace(fmt.Sprintf("%v", rawLocal))
		hisVal := strings.TrimSpace(fmt.Sprintf("%v", RawHis))
		// null值处理
		if localVal == "<nil>" {
			localVal = ""
		}
		if hisVal == "<nil>" {
			hisVal = ""
		}
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
	case "ypmc":
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

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// CleanHISJson 智能清洗 HIS 接口返回的“脏”数据
func CleanHISJson(raw []byte) []byte {
	// 1. 处理斜杠：只转义后面不是合法转义字符的斜杠
	// 我们先将已有的合法转义占位，防止被误伤
	s := string(raw)
	s = strings.ReplaceAll(s, "\\\"", "###QUOTE###") // 保护 \"
	s = strings.ReplaceAll(s, "\\n", "###N###")      // 保护 \n
	s = strings.ReplaceAll(s, "\\r", "###R###")      // 保护 \r
	s = strings.ReplaceAll(s, "\\t", "###T###")      // 保护 \t

	// 2. 将剩下的单斜杠全部转义为双斜杠
	s = strings.ReplaceAll(s, "\\", "\\\\")

	// 3. 还原保护的字符
	s = strings.ReplaceAll(s, "###QUOTE###", "\\\"")
	s = strings.ReplaceAll(s, "###N###", "\\n")
	s = strings.ReplaceAll(s, "###R###", "\\r")
	s = strings.ReplaceAll(s, "###T###", "\\t")

	// 4. 【关键】处理末尾垃圾字符：截取到最后一个 '}' 或 ']'
	lastBrace := strings.LastIndex(s, "}")
	lastBracket := strings.LastIndex(s, "]")
	endIdx := lastBrace
	if lastBracket > endIdx {
		endIdx = lastBracket
	}
	if endIdx != -1 {
		s = s[:endIdx+1]
	}

	return []byte(s)
}

//// 测试
//func (s *DictCompareService) CompareDictDataCs(local *model.LocalDictRow) ([]model.DictCompareResult, string, error) {
//	// --- 测试模式：模拟 HIS 接口返回数据 ---
//	mockHisJson := `{
//    "users": [
//        {
//            "lbdm": "09",
//            "sfwjkcl": false,
//            "ydcldm": "1770293273",
//            "sccjdm": "1825",
//            "pym2": null,
//            "ypbm2": null,
//            "yppp": "浙江新亚医疗科技股份",
//            "kfpfj": 30.0000,
//            "sfwwhp1": null,
//            "ypzczh_xq": 1826812800000,
//            "kfzhl": 1,
//            "tymc": "口外正畸牵引装置",
//            "kfcgj": 30.0000,
//            "cljflx": "1",
//            "kfdm": "2095",
//            "gnzdl": "",
//            "sccj": "浙江新亚医疗科技股份有限公司",
//            "ypbwm": "C07110914800003041210000007",
//            "lsh": "/",
//            "ypzczh": "浙械注准20172171224",
//            "zxzhl": 1,
//            "kfdw": "套        ",
//            "pym1": null,
//            "lrrq": 1744078711000,
//            "zjm": "kwzjqyzz       ",
//            "jxbm": "17",
//            "cctj": "常温",
//            "yplb": "0305      ",
//            "ypbz": "副",
//            "bz": "名称：caiwu2;IP:172.21.67.5",
//            "ypmc": "口外正畸牵引装置",
//            "ypdm": "03050000003951",
//            "kflsj": 30.0000,
//            "gsdm": "1877",
//            "sybz": "1",
//            "ypgg": "J型钩|轻力型-Ⅱ型  65mm          2*1",
//            "zxlsj": 30.0000,
//            "pym": "kwzjqyzz       ",
//            "ghdw": "成都登思特医疗器械有限公司",
//            "zxcgj": 30.0000,
//            "ypbm": "口外正畸牵引装置              ",
//            "zxdw": "套      ",
//            "ypbm1": null
//        }
//    ]
//}`
//
//	var hisRes integration.HisDictResponse
//	if err := json.Unmarshal([]byte(mockHisJson), &hisRes); err != nil {
//		return nil, "", fmt.Errorf("测试数据解析失败: %v", err)
//	}
//	// 3. 停用判断
//	if len(hisRes.Users) == 0 {
//		return nil, "HIS系统未找到该材料或已停用", nil
//	}
//	his := hisRes.Users[0]
//
//	// 4. 定义对比规则（确保 Key 与 hisres.json 中的字段一致）
//	checkFields := []struct {
//		Label string
//		Key   string
//	}{
//		{"产品名称", "ypmc"}, // 对应 "tymc": "脑压板"
//		{"规格型号", "ypgg"}, // 对应 "ypgg": "脑压板"
//		{"库房单位", "kfdw"}, // 对应 "kfdw": "件        "
//		{"采购价", "kfcgj"}, // 对应 "kfcgj": 101.4000
//		{"库房代码", "kfdm"}, // 对应 "kfdm": "2095"
//		{"公司代码", "gsdm"},
//	}
//
//	// 5. 执行对比
//	var results []model.DictCompareResult
//	for _, f := range checkFields {
//		localVal := s.getReflectVal(local, f.Key)
//		hisVal := his[f.Key]
//
//		results = append(results, model.DictCompareResult{
//			Label:      f.Label,
//			Field:      f.Key,
//			LocalValue: localVal,
//			HisValue:   hisVal,
//			// 重点：使用 fmt.Sprintf 消除 float64 精度差异（如 101.4000 与 101.4）
//			IsMatch: fmt.Sprintf("%v", localVal) == fmt.Sprintf("%v", hisVal),
//		})
//	}
//
//	return results, "测试比对完成 (Mock 数据)", nil
//}
