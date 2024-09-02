## 项目介绍

`pf_util` 扩展包 是一个  信息脱敏工具 扩展包

## 效果如下:

```
"黄宗泽"       ->  "黄**"
"18049531999" -> "180****1999"
"123456"      ->  "1*****"
"51343620000320711X" -> "5***************1X"
"51343620000320711X" -> "5134************1X"

```

## 代码

```
package tests

import (
	"fmt"
	"pf_util"
	"testing"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2024/9/2
 * @Desc:
 * @Project: pf_util
 */

func TestDesensitizedUtilUserId(T *testing.T) {
	d := pf_util.DesensitizedUtil{}
	res := d.SetType(0).Desensitized("10000") // 0 不暴露
	fmt.Println(res)
}

func TestDesensitizedUtilName(T *testing.T) {
	d := pf_util.DesensitizedUtil{}
	res := d.SetType(1).Desensitized("黄宗泽") // 0 不暴露
	fmt.Println(res)
}

func TestDesensitizedUtilIDcard(T *testing.T) {
	d := pf_util.DesensitizedUtil{}
	res := d.SetType(2).Desensitized("51343620000320711X") // 5***************1X
	fmt.Println(res)
}

func TestCustomerIDcard(T *testing.T) {
	d := pf_util.DesensitizedUtil{}
	res := d.Method("idCardNum", "51343620000320711X", 4, 2) // 5134************1X
	fmt.Println(res)
}

func TestCustomerChineseName(t *testing.T) {
	d := pf_util.DesensitizedUtil{}
	result := d.Method("chineseName", "黄老板")
	fmt.Println(result)
}

```
