package server

import (
	"encoding/json"
	"fmt"
	"iht/pkg/cfg"
	"iht/pkg/pan"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/xshrim/gol"
	"github.com/xshrim/gol/tk"
)

func ClientP123() *pan.P123 {
	return &pan.P123{
		Id:     cfg.Conf.P123.Id,
		Secret: cfg.Conf.P123.Secret,
		Token:  cfg.Conf.P123.Token,
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	content := "Hello, World!"
	code := http.StatusOK
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	w.WriteHeader(code)

	fmt.Fprintf(w, "%s", content)
}

func list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reg := regexp.MustCompile(`/list/(.*)/?`)
	matches := reg.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var cid, key string
	if tmp := tk.Jsquery(string(body), ".cid"); tmp != nil {
		cid, _ = tmp.(string)
	}

	if tmp := tk.Jsquery(string(body), ".key"); tmp != nil {
		key, _ = tmp.(string)
	}

	var content string
	p := matches[1]
	switch p {
	case "p123":
		client := ClientP123()

		list, err := client.FetchList(cid, key)
		if err != nil {
			http.Error(w, "Error fetching list", http.StatusInternalServerError)
			return
		}

		if data, err := json.Marshal(list); err != nil {
			http.Error(w, "Error marshaling list", http.StatusInternalServerError)
			return
		} else {
			content = string(data)
		}
		break
	default:
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", content)
}

func rename(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reg := regexp.MustCompile(`/rename/(.*)/?`)
	matches := reg.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	type item struct {
		Fid    string `json:"fid"`
		Name   string `json:"name"`
		Path   string `json:"path"`
		Type   string `json:"type"`
		Status string `json:"status"`
	}

	var items []*item
	if err := json.Unmarshal(body, &items); err != nil {
		http.Error(w, "Error unmarshaling request body", http.StatusBadRequest)
		return
	}

	var rnlist []string
	for _, item := range items {
		rnlist = append(rnlist, fmt.Sprintf("%s|%s", item.Fid, item.Name))
	}

	var content string
	p := matches[1]
	switch p {
	case "p123":
		client := ClientP123()

		result := client.RenameList(rnlist)

		for _, succeed := range result["succeed"] {
			s := strings.Split(succeed, "|")
			for _, item := range items {
				if item.Fid == s[0] {
					item.Status = "1"
					break
				}
			}
		}

		for _, failed := range result["failed"] {
			f := strings.Split(failed, "|")
			for _, item := range items {
				if item.Fid == f[0] {
					item.Status = "0"
					break
				}
			}
		}

		if data, err := json.Marshal(items); err != nil {
			http.Error(w, "Error marshaling list", http.StatusInternalServerError)
			return
		} else {
			content = string(data)
		}

		break
	default:
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", content)
}

func Serve(port int, dir string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))

	http.HandleFunc("/echo", echo)
	http.HandleFunc("/echo/", echo)

	http.HandleFunc("/list/p123", list)
	http.HandleFunc("/list/p123/", list)

	http.HandleFunc("/rename/p123", rename)
	http.HandleFunc("/rename/p123/", rename)

	gol.Infof("Starting server on 0.0.0.0:%d\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		gol.Fatal(err)
	}
}
