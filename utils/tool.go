package utils

import (
	"fmt"
	"strings"
)

func int2chinese(num int) string {
	//1、数字为0
	if num == 0 {
		return "零"
	}
	var ans string
	//数字
	szdw := []string{"十", "百", "千", "万", "十万", "百万", "千万", "亿"}
	//数字单位
	sz := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	res := make([]string, 0)

	//数字单位角标
	idx := -1
	for num > 0 {
		//当前位数的值
		x := num % 10
		//2、数字大于等于10
		// 插入数字单位，只有当数字单位角标在范围内，且当前数字不为0 时才有效
		if idx >= 0 && idx < len(szdw) && x != 0 {
			res = append([]string{szdw[idx]}, res...)
		}
		//3、数字中间有多个0
		// 当前数字为0，且后一位也为0 时，为避免重复删除一个零文字
		if x == 0 && len(res) != 0 && res[0] == "零" {
			res = res[1:]
		}
		// 插入数字文字
		res = append([]string{sz[x]}, res...)
		num /= 10
		idx++
	}
	//4、个位数为0
	if len(res) > 1 && res[len(res)-1] == "零" {
		res = res[:len(res)-1]
	}
	//合并字符串
	for i := 0; i < len(res); i++ {
		ans = ans + res[i]
	}
	return ans
}

func int2roman(num int) string {
	//创建映射列表
	numsmap := map[int]string{
		1:    "I",
		4:    "IV",
		5:    "V",
		9:    "IX",
		10:   "X",
		40:   "XL",
		50:   "L",
		90:   "XC",
		100:  "C",
		400:  "CD",
		500:  "D",
		900:  "CM",
		1000: "M",
	}
	//创建整数数组
	numsint := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	results := []string{}
	count := 0
	//进入循环
	for i := 0; i < len(numsint) && num != 0; i++ {
		//判断当前数字是否比map中数值大
		count = num / numsint[i]
		//如果大，则减去当前值
		num = num - count*numsint[i]
		//并记录字符，注意这里用的是for循环
		for count != 0 {
			results = append(results, numsmap[numsint[i]])
			//更新count值
			count--
		}
	}
	return strings.Join(results, "")
}

// intToCircledNumber 将整数转换为带圈的 Unicode 字符。
// 支持 1 到 20 的整数，超出范围返回错误。
func int2cnum(num int) string {
	if num < 1 || num > 20 {
		return ""
	}

	// Unicode 编码从 ① (U+2460) 到 ⑳ (U+2473) 是连续的
	circledNum := rune(0x245F + num)
	return string(circledNum)
}

// intToParenthesizedNumber 将整数转换为带括号的 Unicode 字符
// 支持 1 到 20 的整数，超出范围返回错误。
func int2pnum(num int) string {
	if num < 1 || num > 20 {
		return ""
	}

	// Unicode 编码从 ⑴ (U+2474) 到 ⑳ (U+2487) 是连续的
	parenthesizedNum := rune(0x2473 + num)
	return string(parenthesizedNum)
}

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func Seq(sample string, i int) string {
	if i < 1 {
		return ""
	}

	suffix := ""
	if strings.HasSuffix(sample, ".") {
		suffix = "."
		sample = strings.TrimSuffix(sample, ".")
	}

	switch sample {
	case "1":
		return fmt.Sprintf("%d", i) + suffix
	case "01":
		return fmt.Sprintf("%02d", i) + suffix
	case "001":
		return fmt.Sprintf("%03d", i) + suffix
	case "0001":
		return fmt.Sprintf("%04d", i) + suffix
	case "00001":
		return fmt.Sprintf("%05d", i) + suffix
	case "000001":
		return fmt.Sprintf("%06d", i) + suffix
	case "a":
		return fmt.Sprintf("%c", i+96) + suffix
	case "A":
		return fmt.Sprintf("%c", i+64) + suffix
	case "I":
		return int2roman(i) + suffix
	case "i":
		return strings.ToLower(int2roman(i)) + suffix
	case "一":
		return int2chinese(i) + suffix
	case "①":
		return int2cnum(i) + suffix
	case "⑴":
		return int2pnum(i) + suffix
	case "(1)":
		return fmt.Sprintf("(%d)", i) + suffix
	case "(01)":
		return fmt.Sprintf("(%02d)", i) + suffix
	case "(001)":
		return fmt.Sprintf("(%03d)", i) + suffix
	case "(0001)":
		return fmt.Sprintf("(%04d)", i) + suffix
	case "(00001)":
		return fmt.Sprintf("(%05d)", i) + suffix
	case "(000001)":
		return fmt.Sprintf("(%06d)", i) + suffix
	case "(a)":
		return fmt.Sprintf("(%c)", i+96) + suffix
	case "(A)":
		return fmt.Sprintf("(%c)", i+64) + suffix
	case "(I)":
		return fmt.Sprintf("(%s)", int2roman(i)) + suffix
	case "(i)":
		return fmt.Sprintf("(%s)", strings.ToLower(int2roman(i))) + suffix
	case "(一)":
		return fmt.Sprintf("(%s)", int2chinese(i)) + suffix
	default:
		return fmt.Sprintf("%d", i)
	}
}

func Index(str, sub string) []int {
	var positions []int
	if sub == "" {
		return positions // 如果子串为空，返回空数组
	}

	strs := []rune(str)
	subs := []rune(sub)

	offset := 0
	for i := range strs {
		if i < offset {
			continue
		}

		found := true
		for j := range subs {
			if strs[i+j] != subs[j] {
				found = false
				break
			}
		}
		if found {
			positions = append(positions, i)
			offset = i + 1
		}
	}

	return positions
}

func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}
