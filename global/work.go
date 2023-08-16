package global

import (
	"fmt"
	"log"
)

const (
	C_PageID      = "pageID"
	C_ApiURL      = "ApiURL"
	C_VerifyToken = "VerifyToken"
	C_AccessToken = "AccessToken"
)

type IConfig interface {
	Get(key string) string
	Set(key string, value any) error
}

type config struct {
}

var vconfig *config

func InitConfig(filePath string) (IConfig, error) {
	// TODO: init config with filePath， toml, ini, yaml
	vconfig = &config{}
	return vconfig, nil
}

func (p *config) Get(key string) string {
	switch key {
	case C_PageID:
		return "122099998538016269"
	case C_ApiURL:
		return "https://graph.facebook.com/v17.0"
	case C_VerifyToken:
		return "facebook"
	case C_AccessToken:
		return "EAAJu4pzdpeABO1GdHcXFDEudMznO84HL0C6Lk4FOGoisL4JGZAa7DZB8wgR4RtZCkMpQe8KrrP95iZACOocn6bksixNfdSdnVHDl14D4PBPxl0HDARtvGfFUsAy4kjcWjznhwMhvxl9sUJe90xWPUZBJJJMDT2Y1vCfBGjfWAXYrhn5vtwmHOYBcyK00NyX0F8ZA8ZD"
	}
	return ""
}

func (p *config) Set(key string, value any) error {
	return nil
}

func LogPrintln(datas ...interface{}) {
	fmt.Printf("\n")
	log.Println(datas...)
}

func LogPrintf(format string, datas ...interface{}) {
	fmt.Printf("\n")
	log.Printf(format+"\n", datas...)
}
