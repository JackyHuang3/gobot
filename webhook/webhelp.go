package webhook

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const c_params_parsed = "parsed"

func WebParams(r *http.Request) url.Values {
	ret := make(url.Values, 0)
	if r.Header.Get(c_params_parsed) != "true" {
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				LogPrintf("%s %s.ParseMultipartForm: %v", r.Method, r.URL.Path, err.Error())
				return ret
			}
		} else {
			if err := r.ParseForm(); err != nil {
				LogPrintf("%s %s.ParseForm: %v", r.Method, r.URL.Path, err.Error())
				return ret
			}
		}
	}

	for k := range r.URL.Query() {
		ret.Add(k, r.URL.Query().Get(k))
	}
	for k := range r.Header {
		ret.Add(k, r.Header.Get(k))
	}
	for k := range r.Form {
		if ret.Get(k) == "" {
			ret.Add(k, r.Form.Get(k))
		}
	}
	for k := range r.PostForm {
		if ret.Get(k) == "" {
			ret.Add(k, r.PostForm.Get(k))
		}
	}
	r.Header.Add(c_params_parsed, "true")
	return ret
}

func WebResponse(w http.ResponseWriter, r *http.Request, data []byte, status int) {
	if status != http.StatusOK {
		LogPrintf("WebResponse error: %v, status: %d", string(data), status)
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.WriteHeader(status)
	w.Write(data)
}

// http.Request
func WebGetByFunc(url string, fHandle func(buf []byte) error) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if fHandle != nil {
		return fHandle(raw)
	}
	return nil
}

func WebPostFormByFunc(url string, form url.Values, fHandle func(buf []byte) error) error {
	resp, err := http.PostForm(url, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if fHandle != nil {
		return fHandle(raw)
	}
	return nil
}

func WebPostByFunc(url string, response io.Reader, fHandle func(buf []byte) error) error {
	resp, err := http.Post(url, "application/json", response)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if fHandle != nil {
		return fHandle(raw)
	}
	return nil
}
