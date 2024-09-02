package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

/**
 * @Author: PFinal南丞
 * @Author: lampxiezi@163.com
 * @Date: 2024/9/2
 * @Desc:
 * @Project: pf_util
 */

type DesensitizedType int

const (
	UserId DesensitizedType = iota
	ChineseName
	IdCard
	FixedPhone
	MobilePhone
	ADDRESS
	EMAIL
	PASSWORD
	CarLicense
	BankCard
	IPV4
	IPV6
	FirstMask
	ClearToNull
	ClearToEmpty
)

// 方法注册表
var methodRegistry = map[string]reflect.Value{}

func Desensitized(str string, desensitizedType DesensitizedType) string {
	if str == "" {
		return ""
	}
	newStr := str
	switch desensitizedType {
	case UserId:
		fmt.Println(UserId)
		newStr = strconv.FormatInt(userId(), 10)
	case ChineseName:
		newStr = chineseName(str)
	case IdCard:
		newStr = idCardNum(str, 1, 2)
	case FixedPhone:
		newStr = fixedPhone(str)
	case MobilePhone:
		newStr = mobilePhone(str)
	case ADDRESS:
		newStr = address(str, 8)
	case EMAIL:
		newStr = email(str)
	case PASSWORD:
		newStr = password(str)
	case CarLicense:
		newStr = carLicense(str)
	case BankCard:
		newStr = bankCard(str)
	case IPV4:
		newStr = ipv4(str)
	case IPV6:
		newStr = ipv6(str)
	case FirstMask:
		newStr = firstMask(str)
	case ClearToEmpty:
		newStr = c()
	case ClearToNull:
		newStr = clearToNull()
	default:
	}
	return newStr
}

// init 函数用于初始化注册表
func init() {
	methodRegistry["c"] = reflect.ValueOf(c)
	methodRegistry["userId"] = reflect.ValueOf(userId)
	methodRegistry["chineseName"] = reflect.ValueOf(chineseName)
	methodRegistry["firstMask"] = reflect.ValueOf(firstMask)
	methodRegistry["idCardNum"] = reflect.ValueOf(idCardNum)
	methodRegistry["address"] = reflect.ValueOf(address)
}

func InvokeMethod(methodName string, args ...interface{}) []interface{} {
	if method, found := methodRegistry[methodName]; found {
		// 获取方法的参数类型
		methodType := method.Type()
		numIn := methodType.NumIn()

		// 检查参数数量
		if len(args) != numIn {
			fmt.Println("Argument count mismatch")
			return nil
		}
		// 将 args 转换为 reflect.Value
		in := make([]reflect.Value, len(args))
		for i, arg := range args {
			argValue := reflect.ValueOf(arg)
			expectedType := methodType.In(i)
			// 确保参数类型匹配
			if !argValue.Type().AssignableTo(expectedType) {
				fmt.Printf("Argument type mismatch: expected %v, got %v\n", expectedType, argValue.Type())
				return nil
			}
			in[i] = argValue
		}
		// 调用方法并获取返回值
		results := method.Call(in)
		// 将结果转换为 []interface{}
		var resultInterfaces []interface{}
		for _, result := range results {
			resultInterfaces = append(resultInterfaces, result.Interface())
		}
		return resultInterfaces
	} else {
		return nil
	}
}

func c() string {
	return ""
}

func clearToNull() string {
	return ""
}

/**
 * 【用户id】不对外提供userId
 *
 * @return 脱敏后的主键
 */
func userId() int64 {
	return 0
}

/**
 * 定义了一个first_mask的规则，只显示第一个字符。<br>
 * 脱敏前：123456789；脱敏后：1********。
 *
 * @param str 字符串
 * @return 脱敏后的字符串
 */
func firstMask(str string) string {
	if str == "" {
		return ""
	}
	// 将第一个字符保留，其余字符替换为 '*'
	firstChar := str[0]
	maskedStr := string(firstChar) + strings.Repeat("*", len(str)-1)
	return maskedStr
}

/**
 * 【中文姓名】只显示第一个汉字，其他隐藏为2个星号，比如：李**
 *
 * @param fullName 姓名
 * @return 脱敏后的姓名
 */
func chineseName(fullName string) string {
	if fullName == "" {
		return ""
	}

	// Convert the string to a slice of runes to handle Unicode characters properly
	runes := []rune(fullName)

	// Check if there are any characters
	if len(runes) == 0 {
		return ""
	}

	// The first rune
	firstRune := runes[0]

	// Create a builder to construct the masked name
	var maskedName strings.Builder
	maskedName.WriteRune(firstRune) // Write the first character

	// Append the appropriate number of asterisks
	for i := 1; i < len(runes); i++ {
		maskedName.WriteRune('*')
	}

	return maskedName.String()
}

/**
 * 【身份证号】前1位 和后2位
 *
 * @param idCardNum 身份证
 * @param front     保留：前面的front位数；从1开始
 * @param end       保留：后面的end位数；从1开始
 * @return 脱敏后的身份证
 */
func idCardNum(idCardNum string, front, end int) string {
	if idCardNum == "" {
		return ""
	}
	if (front + end) > len(idCardNum) {
		return ""
	}
	if front < 0 || end < 0 {
		return ""
	}
	// 提取前面的字符
	frontPart := idCardNum[:front]
	// 提取后面的字符
	endPart := idCardNum[len(idCardNum)-end:]
	// 生成中间的 * 字符串
	maskedPart := strings.Repeat("*", len(idCardNum)-front-end)

	// 拼接前、后和中间的部分
	return frontPart + maskedPart + endPart
}

/**
 * 【固定电话 前四位，后两位
 *
 * @param num 固定电话
 * @return 脱敏后的固定电话；
 */
func fixedPhone(num string) string {
	if num == "" {
		return ""
	}
	if len(num) < 6 {
		// Handle case where phone number is too short to mask properly
		return num
	}

	// 提取前面的字符
	frontPart := num[:4]
	// 提取后面的字符
	endPart := num[len(num)-2:]
	// 生成中间的 * 字符串
	maskedPart := strings.Repeat("*", len(num)-6) // Adjust length of masked part

	// 拼接前、后和中间的部分
	return frontPart + maskedPart + endPart
}

/**
 * 【手机号码】前三位，后4位，其他隐藏，比如135****2210
 *
 * @param num 移动电话；
 * @return 脱敏后的移动电话；
 */
func mobilePhone(num string) string {
	if num == "" {
		return ""
	}
	// 确保手机号码长度至少为7位（3位前缀 + 4位后缀）
	if len(num) < 7 {
		return num // 返回原始号码（或根据需要处理）
	}

	// 提取前面的字符
	frontPart := num[:3]
	// 提取后面的字符
	endPart := num[len(num)-4:]
	// 生成中间的 * 字符串
	maskedPart := strings.Repeat("*", len(num)-7)

	// 拼接前、后和中间的部分
	return frontPart + maskedPart + endPart
}

/**
 * 【地址】只显示到地区，不显示详细地址，比如：北京市海淀区****
 *
 * @param address       家庭住址
 * @param sensitiveSize 敏感信息长度
 * @return 脱敏后的家庭地址
 */
func address(address string, sensitiveSize int) string {
	if address == "" {
		return ""
	}

	length := utf8.RuneCountInString(address)
	if sensitiveSize < 0 {
		// 处理敏感长度为负数的情况
		return address
	}

	if sensitiveSize >= length {
		// Mask the entire address if sensitiveSize is greater than or equal to the address length
		return strings.Repeat("*", length)
	}

	// Mask the part of the address
	// 获取每个字符的切割位置
	runes := []rune(address)
	maskedPart := strings.Repeat("*", sensitiveSize)
	return string(runes[:length-sensitiveSize]) + maskedPart
}

/**
 * 【电子邮箱】邮箱前缀仅显示第一个字母，前缀其他隐藏，用星号代替，@及后面的地址显示，比如：d**@126.com
 *
 * @param email 邮箱
 * @return 脱敏后的邮箱
 */
func email(email string) string {
	if email == "" {
		return ""
	}

	// Find the position of "@"
	index := strings.Index(email, "@")
	if index <= 1 {
		// If there's no "@" or the "@" is at the beginning, return the email as is
		return email
	}

	// Mask the part before "@"
	localPart := email[:index]
	domainPart := email[index:]

	// Mask the local part except for the first character
	maskedLocalPart := string(localPart[0]) + strings.Repeat("*", len(localPart)-1)

	return maskedLocalPart + domainPart
}

/**
 * 【密码】密码的全部字符都用*代替，比如：******
 *
 * @param password 密码
 * @return 脱敏后的密码
 */
func password(password string) string {
	if password == "" {
		return ""
	}
	return strings.Repeat("*", len(password))
}

/**
 * 【中国车牌】车牌中间用*代替
 * eg1：null       -》 ""
 * eg1：""         -》 ""
 * eg3：苏D40000   -》 苏D4***0
 * eg4：陕A12345D  -》 陕A1****D
 * eg5：京A123     -》 京A123     如果是错误的车牌，不处理
 *
 * @param carLicense 完整的车牌号
 * @return 脱敏后的车牌
 */
func carLicense(carLicense string) string {
	if carLicense == "" {
		return ""
	}

	// Check if the car license plate is valid
	if len(carLicense) != 7 && len(carLicense) != 8 {
		// Invalid license plate
		return carLicense
	}

	// Define how many characters to mask
	var start, end int
	if len(carLicense) == 7 {
		start, end = 3, 6
	} else if len(carLicense) == 8 {
		start, end = 3, 7
	}

	// Extract parts of the license plate
	prefix := carLicense[:start]
	middle := carLicense[start:end]
	suffix := carLicense[end:]

	// Mask the middle part
	maskedMiddle := strings.Repeat("*", len(middle))

	// Concatenate the parts
	return prefix + maskedMiddle + suffix
}

/**
 * 【银行卡号脱敏】由于银行卡号长度不定，所以只展示前4位，后面的位数根据卡号决定展示1-4位
 * 例如：
 * <pre>{@code
 *      1. "1234 2222 3333 4444 6789 9"    ->   "1234 **** **** **** **** 9"
 *      2. "1234 2222 3333 4444 6789 91"   ->   "1234 **** **** **** **** 91"
 *      3. "1234 2222 3333 4444 678"       ->    "1234 **** **** **** 678"
 *      4. "1234 2222 3333 4444 6789"      ->    "1234 **** **** **** 6789"
 *  }</pre>
 *
 * @param bankCardNo 银行卡号
 * @return 脱敏之后的银行卡号
 */
func bankCard(bankCardNo string) string {
	if bankCardNo == "" {
		return bankCardNo
	}

	// Remove all spaces from the bank card number
	bankCardNo = strings.Join(strings.Fields(bankCardNo), "")

	// If the length of the card number is less than 9, return it as is
	if len(bankCardNo) < 9 {
		return bankCardNo
	}

	length := len(bankCardNo)
	// Calculate the number of digits to show at the end
	endLength := length % 4
	if endLength == 0 {
		endLength = 4
	}
	// Calculate the length of the masked part
	midLength := length - 4 - endLength

	var buf strings.Builder
	// Write the first 4 digits
	buf.WriteString(bankCardNo[:4])
	for i := 0; i < midLength; i++ {
		if i%4 == 0 {
			buf.WriteString(" ")
		}
		buf.WriteString("*")
	}
	// Write the last part of the card number
	buf.WriteString(" ")
	buf.WriteString(bankCardNo[length-endLength:])
	return buf.String()
}

/**
 * IPv4脱敏，如：脱敏前：192.0.2.1；脱敏后：192.*.*.*。
 *
 * @param ipv4 IPv4地址
 * @return 脱敏后的地址
 */
func ipv4(ipv4 string) string {
	if ipv4 == "" {
		return ""
	}

	// Split the IPv4 address into octets
	octets := strings.Split(ipv4, ".")

	// Check if the address has exactly 4 octets
	if len(octets) != 4 {
		// If it's not a valid IPv4 address, return as is
		return ipv4
	}

	// Construct the masked address
	return octets[0] + ".*.*.*"
}

/**
 * IPv6脱敏，如：脱敏前：2001:0db8:86a3:08d3:1319:8a2e:0370:7344；脱敏后：2001:*:*:*:*:*:*:*
 *
 * @param ipv6 IPv6地址
 * @return 脱敏后的地址
 */
func ipv6(ipv6 string) string {
	if ipv6 == "" {
		return ""
	}

	// Split the IPv6 address into segments
	segments := strings.Split(ipv6, ":")

	// Check if the address has exactly 8 segments
	if len(segments) != 8 {
		// If it's not a valid IPv6 address, return as is
		return ipv6
	}

	// Construct the masked address
	return segments[0] + ":*:*:*:*:*:*:*"
}
