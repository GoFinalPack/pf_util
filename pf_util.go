package pf_util

import (
	"fmt"
	"pf_util/utils"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2024/9/2
 * @Desc:
 * @Project: pf_util
 */

type DesensitizedUtil struct {
	utils.DesensitizedType `json:"utils_._desensitized_type,omitempty"`
}

func (d *DesensitizedUtil) SetType(dtype int) *DesensitizedUtil {
	d.DesensitizedType = utils.DesensitizedType(dtype)
	return d
}

func (d *DesensitizedUtil) Desensitized(str string) string {
	return utils.Desensitized(str, d.DesensitizedType)
}

func (d *DesensitizedUtil) Method(methodName string, args ...interface{}) string {
	// 调用 utils.Method 并获取结果
	result := utils.InvokeMethod(methodName, args...)
	if result != nil && len(result) > 0 {
		// 尝试将第一个结果转换为字符串
		if str, ok := result[0].(string); ok {
			return str
		}
		// 使用 fmt.Sprintf 将其他类型转换为字符串
		return fmt.Sprintf("%v", result[0])
	}
	return ""
}
