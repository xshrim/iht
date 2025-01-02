package flow

import (
	"encoding/json"
	"fmt"
	"iht/utils"
	"os"
	"regexp"
	"strconv"
	"strings"

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

type Rename struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Mode  string `json:"mode"`
	Expr  string `json:"expr"`
	Value string `json:"value"`
	Num   int    `json:"num"`
}

func (r Rename) Add(strs []string) []string {
	var output []string
	for _, str := range strs {
		switch r.Mode {
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
			// istrs := []rune(str)
			// index := r.Index
			// if index < 0 {
			// 	index = len(istrs) + 1 + index
			// }
			// if index <= 0 {
			// 	index = 1
			// } else if index > len(istrs) {
			// 	index = len(istrs) + 1
			// }
			// index -= 1

			// output = append(output, string(istrs[:index])+r.Value+string(istrs[index:]))
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

func (r Rename) Seq(strs []string) []string {
	var output []string
	for i, str := range strs {
		seqstr := utils.Seq(r.Value, i+1)
		switch r.Mode {
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

			strs := []rune(str)

			start, end := 0, 0
			for idx, index := range utils.Index(str, r.Expr) {
				if idx+1 == num {
					start = index
					end = index + len([]rune(r.Expr))
					break
				}
			}
			if start == 0 && end == 0 {
				output = append(output, str)
				break
			}

			sub := strs[start:end]
			tmp := append([]rune{}, strs[:start]...)
			tmp = append(tmp, strs[end:]...)

			index := start + offset
			if index < 0 {
				index = 0
			}
			if index > len(tmp) {
				index = len(tmp)
			}

			output = append(output, string(tmp[:index])+string(sub)+string(tmp[index:]))
		}
	}
	return output
}

func Load(fname string) (*Flow, error) {
	var flow *Flow
	if dataBytes, err := os.ReadFile(fname); err == nil {
		_ = yaml.Unmarshal(dataBytes, &flow)
	}

	if flow.Name == "" {
		return nil, fmt.Errorf("flow name is empty")
	}

	gol.Infof("flow %s loaded\n", flow.Name)

	return flow, nil
}

func (f *Flow) Run(obj any) error {
	gol.Infof("run flow %s\n", f.Name)
	for _, step := range f.Steps {
		gol.Infof("start step %s\n", step.Name)
		switch step.Type {
		case "rename":
			strs, _ := obj.([]string)
			for _, action := range step.Actions {
				actionBytes, err := json.Marshal(action)
				if err != nil {
					return err
				}
				var rename Rename
				err = json.Unmarshal(actionBytes, &rename)
				if err != nil {
					return err
				}

				switch rename.Kind {
				case "add":
					strs = rename.Add(strs)
				case "delete":
					strs = rename.Delete(strs)
				case "replace":
					strs = rename.Replace(strs)
				case "seq":
					strs = rename.Seq(strs)
				case "shift":
					strs = rename.Shift(strs)
				default:
					return fmt.Errorf("unsupported rename kind %s\n", rename.Kind)
				}

				fmt.Printf("rename %s -> %s\n", rename.Kind, strs)
			}
		default:
			return fmt.Errorf("unsupported step type %s\n", step.Type)
		}
	}
	return nil
}
