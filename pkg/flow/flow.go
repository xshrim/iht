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
	Index int    `json:"index"`
	Value string `json:"value"`
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
			istrs := []rune(str)
			index := r.Index
			if index < 0 {
				index = len(istrs) + 1 + index
			}
			if index <= 0 {
				index = 1
			} else if index > len(istrs) {
				index = len(istrs) + 1
			}
			index -= 1

			output = append(output, string(istrs[:index])+r.Value+string(istrs[index:]))
		case "regexp":
			index := r.Index
			if r.Index <= 0 {
				index = -1
			}
			if exp, err := regexp.Compile(r.Expr); err == nil {
				idxs := exp.FindAllStringIndex(str, index)
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
			index := r.Index
			if r.Index <= 0 {
				index = -1
			}
			output = append(output, strings.Replace(str, r.Expr, "", index))
		case "prefix":
			output = append(output, strings.TrimPrefix(str, r.Expr))
		case "suffix":
			output = append(output, strings.TrimSuffix(str, r.Expr))
		case "index":
			var err error
			start := 1
			end := 0
			istrs := []rune(str)
			ranger := strings.Split(r.Expr, "-")
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
			}

			if start <= 0 {
				start = 1
			}
			if end == 0 {
				end = start
			}

			if end > len(istrs) {
				end = len(istrs)
			}

			if start > len(istrs) || start > end {
				output = append(output, str)
				break
			}

			start -= 1

			output = append(output, string(istrs[:start])+string(istrs[end:]))
		case "regexp":
			if exp, err := regexp.Compile(r.Expr); err == nil {
				if r.Index <= 0 {
					output = append(output, exp.ReplaceAllString(str, ""))
					break
				}

				offset := 0
				idxs := exp.FindAllStringIndex(str, -1)
				for i, idx := range idxs {
					if i >= r.Index {
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
			index := r.Index
			if r.Index <= 0 {
				index = -1
			}
			output = append(output, strings.Replace(str, r.Expr, r.Value, index))
		case "prefix":
			output = append(output, r.Value+strings.TrimPrefix(str, r.Expr))
		case "suffix":
			output = append(output, strings.TrimSuffix(str, r.Expr)+r.Value)
		case "index":
			var err error
			start := 1
			end := 0
			istrs := []rune(str)
			ranger := strings.Split(r.Expr, "-")
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
			}

			if start <= 0 {
				start = 1
			}
			if end == 0 {
				end = start
			}

			if end > len(istrs) {
				end = len(istrs)
			}

			if start > len(istrs) || start > end {
				output = append(output, str)
				break
			}

			start -= 1

			output = append(output, string(istrs[:start])+r.Value+string(istrs[end:]))
		case "regexp":
			if exp, err := regexp.Compile(r.Expr); err == nil {
				if r.Index <= 0 {
					output = append(output, exp.ReplaceAllString(str, r.Value))
					break
				}

				offset := 0
				idxs := exp.FindAllStringIndex(str, -1)
				for i, idx := range idxs {
					if i >= r.Index {
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
		switch r.Mode {
		case "prefix":
			output = append(output, utils.Seq(r.Expr, i+1)+strings.TrimPrefix(str, r.Expr))
		case "suffix":
			output = append(output, strings.TrimSuffix(str, r.Expr)+r.Value)
		default:
			output = append(output, str)
		}
	}

	return output
}

func (r Rename) Shift(strs []string) []string {
	return strs
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
