package common

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type ConvertResponse struct {
	Code   string `json:"code"`
	Desc   string `json:"desc"`
	FileID string `json:"fileid"`
}
type ConversionResult struct {
	Code     string `json:"code"`
	Desc     string `json:"desc"`
	FileID   string `json:"fileid"`
	Stat     string `json:"stat"`
	OutFile  string `json:"outfile"`
	HtmlFile string `json:"htmlfile,omitempty"` // omitempty 表示如果 htmlfile 字段为空，则在解析 JSON 时忽略该字段
	Title    string `json:"title,omitempty"`    // 同上
	Tag      string `json:"tag,omitempty"`      // 同上
}

func GenerateSign(baseURL, appKey string) string {
	// 在URL最后拼接appkey
	signedURL := fmt.Sprintf("%s&appkey=%s", baseURL, appKey)

	// 将拼接appkey后的url进行md5编码
	hasher := md5.New()
	hasher.Write([]byte(signedURL))
	sign := hex.EncodeToString(hasher.Sum(nil))

	return sign
}

func ConvertModel(infileURL, outType, appID, appKey string) ([]byte, error) {
	// 构造请求URL
	baseURL := "https://open.3dwhere.com/api/add"
	params := url.Values{}
	params.Add("appid", appID)
	params.Add("infile", infileURL)
	params.Add("outtype", outType)
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// 生成签名
	sign := GenerateSign(requestURL, appKey)

	// 在请求URL中加入签名
	requestURLWithSign := fmt.Sprintf("%s&sign=%s", requestURL, sign)
	println("url = ", requestURLWithSign)
	// 发起请求
	resp, err := http.Get(requestURLWithSign)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 这里简化处理，直接返回响应体字符串
	// 实际应用中，你需要解析返回的JSON，提取其中的outfile字段
	return body, nil
}

func QueryConversionResult(fileID, appID, appKey string, timeout int) (string, error) {
	endTime := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		if time.Now().After(endTime) {
			// 超时退出
			return "", fmt.Errorf("request model conversion timeout")
		}

		baseURL := "https://open.3dwhere.com/api/query"
		params := url.Values{}
		params.Add("fileid", fileID)
		params.Add("appid", appID)
		requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		sign := GenerateSign(requestURL, appKey)
		requestURLWithSign := fmt.Sprintf("%s&sign=%s", requestURL, sign)
		// 发起查询请求
		resp, err := http.Get(requestURLWithSign)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// 解析响应体
		var result ConversionResult
		err = json.Unmarshal([]byte(body), &result)
		if err != nil {
			fmt.Println("Error parsing conversion result:", err)
			return "", nil
		}
		println("result.Stat: ", result.Stat)
		// 检查转换状态
		if result.Code == "1" {
			if result.Stat == "10" {
				// 转换成功且存在转换后的文件
				return result.OutFile, nil
			} else if result.Stat == "-1" {
				// 转换失败
				return "", fmt.Errorf("conversion failed, desc: %s", result.Desc)
			}
			// 其他状态，等待转换完成
		} else {
			// 返回错误信息或其他必要的处理
			return "", fmt.Errorf("conversion failed or still processing, desc: %s", result.Desc)
		}

		// 如果没有成功也没有失败，等待5秒后重试
		time.Sleep(5 * time.Second)
	}
}

// ConvertAndQueryModel 使用指定的文件URI和输出类型请求模型转换，并查询转换结果
func ConvertAndQueryModel(fileUri, outType string, timeoutSec int) (string, error) {
	appID := "y8wDUEW6ry467o36"                  // 你的AppID
	appKey := "1e603d09deadf33ea28789db9a9315d2" // 你的AppKey

	// 调用模型转换API
	// response, err := ConvertModel(fileUri, outType, appID, appKey)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to request model conversion: %v", err)
	// }

	// // 解析转换请求响应
	// var resp ConvertResponse
	// err = json.Unmarshal(response, &resp)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to parse conversion response: %v", err)
	// }
	// if resp.Code != "1" {
	// 	return "", fmt.Errorf("conversion request failed: %s", resp.Desc)
	// }

	// 等待转换处理完成
	// time.Sleep(60 * time.Second) // 根据实际情况调整等待时间

	// 查询转换结果
	// outfile, err := QueryConversionResult(resp.FileID, appID, appKey, timeoutSec)
	outfile, err := QueryConversionResult("k0DRt1sqMv2Jo6vk", appID, appKey, timeoutSec)
	if err != nil {
		return "", err
	}

	return outfile, nil
}
