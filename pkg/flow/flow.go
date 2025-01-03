package flow

import (
	"encoding/json"
	"fmt"
	"iht/utils"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/xshrim/gol"
	"gopkg.in/yaml.v2"
)

type Flow struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

type Step struct {
	Name    string           `json:"name"`
	Type    string           `json:"type"`
	Actions []map[string]any `json:"actions"` // rename, filter, etc.
}

type Filter struct {
	Name string `json:"name"` // 过滤操作名称
	Kind string `json:"kind"` // 过滤操作类型，include, exclude
	Mode string `json:"mode"` // 过滤操作模式，equal, contain, prefix, suffix, regexp
	Expr string `json:"expr"` // 过滤操作表达式, 文本, 正则表达式
	Num  int    `json:"num"`  // 过滤操作次数, 最多匹配次数
}

type Rename struct {
	Name  string `json:"name"`  // 重命名操作名称
	Kind  string `json:"kind"`  // 重命名操作类型，add, delete, replace, transfer, seq, shift
	Mode  string `json:"mode"`  // 重命名操作模式，plain, prefix, suffix, index, regexp
	Expr  string `json:"expr"`  // 重命名操作表达式, 文本, 索引, 正则表达式
	Value string `json:"value"` // 重命名操作值, 追加值, 替换值, 序列值, 偏移值
	Num   int    `json:"num"`   // 重命名操作次数, 追加次数, 替换次数, 序列次数, 偏移第几个
}

func (r Rename) Add(strs []string) []string {
	var output []string
	for _, str := range strs {
		switch r.Mode {
		case "plain":
			num := r.Num
			if r.Num <= 0 {
				num = -1
			}
			idxs := utils.Index(str, r.Expr)
			for i, index := range idxs {
				if i < num || num == -1 {
					idx := len(string([]rune(str)[:index]))
					str = str[:idx+i*len(r.Value)] + r.Value + str[idx+i*len(r.Value):]
				}
			}
			output = append(output, str)
		case "prefix":
			output = append(output, r.Value+str)
		case "suffix":
			output = append(output, str+r.Value)
		case "index":
			var err error
			start := 1
			istrs := []rune(str)
			ranger := strings.Split(r.Expr, ":")
			if ranger[0] != "" {
				start, err = strconv.Atoi(ranger[0])
				if err != nil {
					output = append(output, str)
					break
				}
			}

			if start == 0 {
				start = 1
			}
			if start < 0 {
				start = len(istrs) + 1 + start
			}

			if start > len(istrs) {
				output = append(output, str+r.Value)
				break
			}

			if len(istrs) > 0 {
				start -= 1
			}

			output = append(output, string(istrs[:start])+r.Value+string(istrs[start:]))
		case "regexp":
			num := r.Num
			if r.Num <= 0 {
				num = -1
			}
			if exp, err := regexp.Compile(r.Expr); err == nil {
				idxs := exp.FindAllStringIndex(str, num)
				for i, idx := range idxs {
					str = str[:idx[0]+i*len(r.Value)] + r.Value + str[idx[0]+i*len(r.Value):]
				}
			}
			output = append(output, str)
		default:
			output = append(output, str)
		}
	}

	return output
}

func (r Rename) Delete(strs []string) []string {
	var output []string
	for _, str := range strs {
		switch r.Mode {
		case "plain":
			num := r.Num
			if r.Num <= 0 {
				num = -1
			}
			output = append(output, strings.Replace(str, r.Expr, "", num))
		case "prefix":
			output = append(output, strings.TrimPrefix(str, r.Expr))
		case "suffix":
			output = append(output, strings.TrimSuffix(str, r.Expr))
		case "index":
			var err error
			start := 1
			end := 0
			istrs := []rune(str)
			ranger := strings.Split(r.Expr, ":")
			if ranger[0] != "" {
				start, err = strconv.Atoi(ranger[0])
				if err != nil {
					output = append(output, str)
					break
				}
			}

			if len(ranger) > 1 {
				if ranger[1] == "" {
					end = len(istrs)
				} else {
					end, err = strconv.Atoi(ranger[1])
					if err != nil {
						output = append(output, str)
						break
					}
				}
			} else {
				end = start
			}

			if start == 0 {
				start = 1
			}
			if end == 0 {
				end = 1
			}
			if start < 0 {
				start = len(istrs) + 1 + start
			}
			if end < 0 {
				end = len(istrs) + 1 + end
			}

			if start > len(istrs) {
				output = append(output, str)
				break
			}
			if end > len(istrs) {
				end = len(istrs)
			}

			if start > end {
				start, end = end, start
			}

			start -= 1

			output = append(output, string(istrs[:start])+string(istrs[end:]))
		case "regexp":
			if exp, err := regexp.Compile(r.Expr); err == nil {
				if r.Num <= 0 {
					output = append(output, exp.ReplaceAllString(str, ""))
					break
				}

				offset := 0
				idxs := exp.FindAllStringIndex(str, -1)
				for i, idx := range idxs {
					if i >= r.Num {
						break
					}

					str = str[:idx[0]-offset] + str[idx[1]-offset:]
					offset += idx[1] - idx[0]
				}
			}
			output = append(output, str)
		default:
			output = append(output, str)
		}
	}

	return output
}

func (r Rename) Replace(strs []string) []string {
	var output []string
	for _, str := range strs {
		switch r.Mode {
		case "plain":
			num := r.Num
			if r.Num <= 0 {
				num = -1
			}
			output = append(output, strings.Replace(str, r.Expr, r.Value, num))
		case "prefix":
			output = append(output, r.Value+strings.TrimPrefix(str, r.Expr))
		case "suffix":
			output = append(output, strings.TrimSuffix(str, r.Expr)+r.Value)
		case "index":
			var err error
			start := 1
			end := 0
			istrs := []rune(str)
			ranger := strings.Split(r.Expr, ":")
			if ranger[0] != "" {
				start, err = strconv.Atoi(ranger[0])
				if err != nil {
					output = append(output, str)
					break
				}
			}

			if len(ranger) > 1 {
				if ranger[1] == "" {
					end = len(istrs)
				} else {
					end, err = strconv.Atoi(ranger[1])
					if err != nil {
						output = append(output, str)
						break
					}
				}
			} else {
				end = start
			}

			if start == 0 {
				start = 1
			}
			if end == 0 {
				end = 1
			}
			if start < 0 {
				start = len(istrs) + 1 + start
			}
			if end < 0 {
				end = len(istrs) + 1 + end
			}

			if start > len(istrs) {
				output = append(output, str)
				break
			}
			if end > len(istrs) {
				end = len(istrs)
			}

			if start > end {
				start, end = end, start
			}

			start -= 1

			output = append(output, string(istrs[:start])+r.Value+string(istrs[end:]))
		case "regexp":
			if exp, err := regexp.Compile(r.Expr); err == nil {
				if r.Num <= 0 {
					output = append(output, exp.ReplaceAllString(str, r.Value))
					break
				}

				offset := 0
				idxs := exp.FindAllStringIndex(str, -1)
				for i, idx := range idxs {
					if i >= r.Num {
						break
					}

					str = str[:idx[0]-offset] + r.Value + str[idx[1]-offset:]
					offset += idx[1] - idx[0] - len(r.Value)
				}
			}
			output = append(output, str)
		default:
			output = append(output, str)
		}
	}
	return output
}

func (r Rename) Transfer(strs []string) []string {
	var output []string
	for _, str := range strs {
		switch r.Mode {
		case "plain":
			num := r.Num
			if r.Num <= 0 {
				num = -1
			}
			val := r.Expr
			switch r.Value {
			case "lower":
				val = strings.ToLower(val)
			case "upper":
				val = strings.ToUpper(val)
			case "title":
				val = strings.Title(val)
			case "capitalize":
				if len(val) < 2 {
					val = strings.ToUpper(val)
				} else {
					val = strings.ToUpper(val[:1]) + strings.ToLower(val[1:])
				}
			}
			output = append(output, strings.Replace(str, r.Expr, val, num))
		case "prefix":
			val := r.Expr
			switch r.Value {
			case "lower":
				val = strings.ToLower(val)
			case "upper":
				val = strings.ToUpper(val)
			case "title":
				val = strings.Title(val)
			case "capitalize":
				if len(val) < 2 {
					val = strings.ToUpper(val)
				} else {
					val = strings.ToUpper(val[:1]) + strings.ToLower(val[1:])
				}
			}
			output = append(output, val+strings.TrimPrefix(str, r.Expr))
		case "suffix":
			val := r.Expr
			switch r.Value {
			case "lower":
				val = strings.ToLower(val)
			case "upper":
				val = strings.ToUpper(val)
			case "title":
				val = strings.Title(val)
			case "capitalize":
				if len(val) < 2 {
					val = strings.ToUpper(val)
				} else {
					val = strings.ToUpper(val[:1]) + strings.ToLower(val[1:])
				}
			}
			output = append(output, strings.TrimSuffix(str, r.Expr)+val)
		case "index":
			var err error
			start := 1
			end := 0
			istrs := []rune(str)
			ranger := strings.Split(r.Expr, ":")
			if ranger[0] != "" {
				start, err = strconv.Atoi(ranger[0])
				if err != nil {
					output = append(output, str)
					break
				}
			}

			if len(ranger) > 1 {
				if ranger[1] == "" {
					end = len(istrs)
				} else {
					end, err = strconv.Atoi(ranger[1])
					if err != nil {
						output = append(output, str)
						break
					}
				}
			} else {
				end = start
			}

			if start == 0 {
				start = 1
			}
			if end == 0 {
				end = 1
			}
			if start < 0 {
				start = len(istrs) + 1 + start
			}
			if end < 0 {
				end = len(istrs) + 1 + end
			}

			if start > len(istrs) {
				output = append(output, str)
				break
			}
			if end > len(istrs) {
				end = len(istrs)
			}

			if start > end {
				start, end = end, start
			}

			start -= 1

			val := string(istrs[start:end])
			switch r.Value {
			case "lower":
				val = strings.ToLower(val)
			case "upper":
				val = strings.ToUpper(val)
			case "title":
				val = strings.Title(val)
			case "capitalize":
				if len(val) < 2 {
					val = strings.ToUpper(val)
				} else {
					val = strings.ToUpper(val[:1]) + strings.ToLower(val[1:])
				}
			}

			output = append(output, string(istrs[:start])+val+string(istrs[end:]))
		case "regexp":
			if exp, err := regexp.Compile(r.Expr); err == nil {
				num := r.Num
				if r.Num <= 0 {
					num = -1
				}

				idxs := exp.FindAllStringIndex(str, -1)
				for i, idx := range idxs {
					istrs := []rune(str)
					start, end := 0, 0
					if i < num || num == -1 {
						start = idx[0]
						end = idx[1]

						start = utf8.RuneCountInString(str[:start])
						end = utf8.RuneCountInString(str[:end])

						if start == 0 && end == 0 {
							break
						}

						val := string(istrs[start:end])
						switch r.Value {
						case "lower":
							val = strings.ToLower(val)
						case "upper":
							val = strings.ToUpper(val)
						case "title":
							val = strings.Title(val)
						case "capitalize":
							if len(val) < 2 {
								val = strings.ToUpper(val)
							} else {
								val = strings.ToUpper(val[:1]) + strings.ToLower(val[1:])
							}
						}

						str = string(istrs[:start]) + val + string(istrs[end:])
					}
				}
			}
			output = append(output, str)
		default:
			output = append(output, str)
		}
	}
	return output
}

func (r Rename) Seq(strs []string) []string {
	var output []string
	for i, str := range strs {
		seqstr := utils.Seq(r.Value, i+1)
		switch r.Mode {
		case "plain":
			num := r.Num
			if r.Num <= 0 {
				num = -1
			}
			idxs := utils.Index(str, r.Expr)
			for i, index := range idxs {
				if i < num || num == -1 {
					idx := len(string([]rune(str)[:index]))
					str = str[:idx+i*len(seqstr)] + seqstr + str[idx+i*len(seqstr):]
				}
			}
			output = append(output, str)
		case "prefix":
			output = append(output, seqstr+str)
		case "suffix":
			output = append(output, str+seqstr)
		case "index":
			var err error
			start := 1
			istrs := []rune(str)
			ranger := strings.Split(r.Expr, ":")
			if ranger[0] != "" {
				start, err = strconv.Atoi(ranger[0])
				if err != nil {
					output = append(output, str)
					break
				}
			}

			if start == 0 {
				start = 1
			}
			if start < 0 {
				start = len(istrs) + 1 + start
			}

			if start > len(istrs) {
				output = append(output, str+seqstr)
				break
			}

			if len(istrs) > 0 {
				start -= 1
			}

			output = append(output, string(istrs[:start])+seqstr+string(istrs[start:]))
		case "regexp":
			num := r.Num
			if r.Num <= 0 {
				num = -1
			}
			if exp, err := regexp.Compile(r.Expr); err == nil {
				idxs := exp.FindAllStringIndex(str, num)
				for i, idx := range idxs {
					str = str[:idx[0]+i*len(seqstr)] + seqstr + str[idx[0]+i*len(seqstr):]
				}
			}
			output = append(output, str)
		default:
			output = append(output, str)
		}
	}

	return output
}

func (r Rename) Shift(strs []string) []string {
	var output []string
	offset, err := strconv.Atoi(r.Value)
	if err != nil {
		return output
	}
	for _, str := range strs {
		switch r.Mode {
		case "plain":
			num := r.Num
			if r.Num <= 0 {
				num = 1
			}

			istrs := []rune(str)

			start, end := 0, 0
			for i, idx := range utils.Index(str, r.Expr) {
				if i+1 == num {
					start = idx
					end = idx + len([]rune(r.Expr))
					break
				}
			}
			if start == 0 && end == 0 {
				output = append(output, str)
				break
			}

			sub := istrs[start:end]
			tmp := append([]rune{}, istrs[:start]...)
			tmp = append(tmp, istrs[end:]...)

			index := start + offset
			if index < 0 {
				index = 0
			}
			if index > len(tmp) {
				index = len(tmp)
			}

			output = append(output, string(tmp[:index])+string(sub)+string(tmp[index:]))
		case "prefix":
			istrs := []rune(str)
			cnt, err := strconv.Atoi(r.Expr)
			if err != nil && strings.HasPrefix(str, r.Expr) {
				cnt = len([]rune(r.Expr))
			}
			if cnt <= 0 {
				output = append(output, str)
				break
			}

			if cnt > len(istrs) {
				cnt = len(istrs)
			}
			start := 0
			end := start + cnt
			sub := istrs[start:end]
			tmp := append([]rune{}, istrs[:start]...)
			tmp = append(tmp, istrs[end:]...)

			index := start + offset
			if index < 0 {
				index = 0
			}
			if index > len(tmp) {
				index = len(tmp)
			}

			output = append(output, string(tmp[:index])+string(sub)+string(tmp[index:]))
		case "suffix":
			istrs := []rune(str)
			cnt, err := strconv.Atoi(r.Expr)
			if err != nil && strings.HasSuffix(str, r.Expr) {
				cnt = len([]rune(r.Expr))
			}
			if cnt <= 0 {
				output = append(output, str)
				break
			}

			if cnt > len(istrs) {
				cnt = len(istrs)
			}
			end := len(istrs)
			start := end - cnt
			sub := istrs[start:end]
			tmp := append([]rune{}, istrs[:start]...)
			tmp = append(tmp, istrs[end:]...)

			index := start + offset
			if index < 0 {
				index = 0
			}
			if index > len(tmp) {
				index = len(tmp)
			}

			output = append(output, string(tmp[:index])+string(sub)+string(tmp[index:]))
		case "index":
			start := 1
			end := 0

			istrs := []rune(str)
			ranger := strings.Split(r.Expr, ":")
			if ranger[0] != "" {
				start, err = strconv.Atoi(ranger[0])
				if err != nil {
					output = append(output, str)
					break
				}
			}

			if len(ranger) > 1 {
				if ranger[1] == "" {
					end = len(istrs)
				} else {
					end, err = strconv.Atoi(ranger[1])
					if err != nil {
						output = append(output, str)
						break
					}
				}
			} else {
				end = start
			}

			if start == 0 {
				start = 1
			}
			if end == 0 {
				end = 1
			}
			if start < 0 {
				start = len(istrs) + 1 + start
			}
			if end < 0 {
				end = len(istrs) + 1 + end
			}

			if start > len(istrs) {
				output = append(output, str)
				break
			}
			if end > len(istrs) {
				end = len(istrs)
			}

			if start > end {
				start, end = end, start
			}

			start -= 1

			sub := istrs[start:end]
			tmp := append([]rune{}, istrs[:start]...)
			tmp = append(tmp, istrs[end:]...)

			index := start + offset
			if index < 0 {
				index = 0
			}
			if index > len(tmp) {
				index = len(tmp)
			}

			output = append(output, string(tmp[:index])+string(sub)+string(tmp[index:]))
		case "regexp":
			if exp, err := regexp.Compile(r.Expr); err == nil {
				num := r.Num
				if r.Num <= 0 {
					num = 1
				}

				istrs := []rune(str)

				start, end := 0, 0
				idxs := exp.FindAllStringIndex(str, -1)
				for i, idx := range idxs {
					if i+1 == num {
						start = idx[0]
						end = idx[1]
						break
					}
				}
				start = utf8.RuneCountInString(str[:start])
				end = utf8.RuneCountInString(str[:end])

				if start == 0 && end == 0 {
					output = append(output, str)
					break
				}

				sub := istrs[start:end]
				tmp := append([]rune{}, istrs[:start]...)
				tmp = append(tmp, istrs[end:]...)

				index := start + offset
				if index < 0 {
					index = 0
				}
				if index > len(tmp) {
					index = len(tmp)
				}

				output = append(output, string(tmp[:index])+string(sub)+string(tmp[index:]))
			} else {
				output = append(output, str)
			}
		default:
			output = append(output, str)
		}
	}
	return output
}

func (f *Filter) Include(strs []string) []string {
	var output []string
	num := f.Num
	if f.Num <= 0 {
		num = -1
	}
	for idx, str := range strs {
		if idx >= num && num != -1 {
			output = append(output, str)
			continue
		}
		switch f.Mode {
		case "equal":
			if f.Expr == str {
				output = append(output, str)
			}
		case "contain":
			if strings.Contains(str, f.Expr) {
				output = append(output, str)
			}
		case "prefix":
			if strings.HasPrefix(str, f.Expr) {
				output = append(output, str)
			}
		case "suffix":
			if strings.HasSuffix(str, f.Expr) {
				output = append(output, str)
			}
		case "regexp":
			if exp, err := regexp.Compile(f.Expr); err == nil {
				if exp.MatchString(str) {
					output = append(output, str)
				}
			} else {
				output = append(output, str)
			}
		default:
			output = append(output, str)
		}
	}

	return output
}

func (f *Filter) Exclude(strs []string) []string {
	var output []string
	num := f.Num
	if f.Num <= 0 {
		num = -1
	}
	for idx, str := range strs {
		if idx >= num && num != -1 {
			output = append(output, str)
			continue
		}
		switch f.Mode {
		case "equal":
			if f.Expr != str {
				output = append(output, str)
			}
		case "contain":
			if !strings.Contains(str, f.Expr) {
				output = append(output, str)
			}
		case "prefix":
			if !strings.HasPrefix(str, f.Expr) {
				output = append(output, str)
			}
		case "suffix":
			if !strings.HasSuffix(str, f.Expr) {
				output = append(output, str)
			}
		case "regexp":
			if exp, err := regexp.Compile(f.Expr); err == nil {
				if !exp.MatchString(str) {
					output = append(output, str)
				}
			} else {
				output = append(output, str)
			}
		default:
			output = append(output, str)
		}
	}

	return output
}

func Load(fname string) (*Flow, error) {
	var flow *Flow
	if dataBytes, err := os.ReadFile(fname); err == nil {
		_ = yaml.Unmarshal(dataBytes, &flow)
	} else {
		_ = yaml.Unmarshal([]byte(fname), &flow)
	}

	if flow.Name == "" {
		return nil, fmt.Errorf("flow name is empty")
	}

	gol.Infof("flow %s loaded\n", flow.Name)

	return flow, nil
}

func (f *Flow) Run(strs []string) ([]string, error) {
	gol.Infof("run flow %s\n", f.Name)
	// strs, _ := obj.([]string)
	for _, step := range f.Steps {
		gol.Infof("start step %s\n", step.Name)
		switch step.Type {
		case "filter":
			for _, action := range step.Actions {
				actionBytes, err := json.Marshal(action)
				if err != nil {
					return strs, err
				}
				var filter Filter
				err = json.Unmarshal(actionBytes, &filter)
				if err != nil {
					return strs, err
				}

				switch filter.Kind {
				case "include":
					strs = filter.Include(strs)
				case "exclude":
					strs = filter.Exclude(strs)
				default:
					return strs, fmt.Errorf("unsupported filter kind %s\n", filter.Kind)
				}

				gol.Debugf("filter %s -> %s\n", filter.Kind, strs)
			}
		case "rename":
			for _, action := range step.Actions {
				actionBytes, err := json.Marshal(action)
				if err != nil {
					return strs, err
				}
				var rename Rename
				err = json.Unmarshal(actionBytes, &rename)
				if err != nil {
					return strs, err
				}

				switch rename.Kind {
				case "add":
					strs = rename.Add(strs)
				case "delete":
					strs = rename.Delete(strs)
				case "replace":
					strs = rename.Replace(strs)
				case "transfer":
					strs = rename.Transfer(strs)
				case "seq":
					strs = rename.Seq(strs)
				case "shift":
					strs = rename.Shift(strs)
				default:
					return strs, fmt.Errorf("unsupported rename kind %s\n", rename.Kind)
				}

				gol.Debugf("rename %s -> %s\n", rename.Kind, strs)
			}

		default:
			return strs, fmt.Errorf("unsupported step type %s\n", step.Type)
		}
	}
	return strs, nil
}
