package main

import (
	"os"
	"errors"
	"io/ioutil"
	"encoding/json"
)

const defaultJson = `{
	"server0": {
		"listen": "127.0.0.1:80",
		"location": {
            "/": {"proxy": {
                    "path": ["http://127.0.0.1:9090", "http://127.0.0.1:9091"],
                    "mode": "random",
                    "mode": "iphash",
                    "mode": "round"
                }
            },
            "/": {"proxy": "http://127.0.0.1:9090"}
		}
	},
    "server1": {
        "listen": "9090",
        "location": {
            "/": {"root": "html"}
        }
    }
}
`
const indexHtml = `<!DOCTYPE html>
<html>
<head><title>Welcome to rainx!</title></head>
<body>
<h1>Welcome to rainx!</h1>
<p>If you see this page, the rainx web server is successfully installed and working.</p>
<em>Thank you for using rainx.</em>
</body>
</html>
`

// Setting Interface
type Setting interface {
	isValid() error
	servers() []ServerSetting
}

// Server Setting Interface
type ServerSetting interface {
	listen() string
	locations() []LocationSetting
}

// Location Setting Interface
// Get detailed configuration
type LocationSetting interface {
	pattern() string
	isRoot() bool
	root() string
	isProxy() bool
	proxy() string
	isProxies() bool
	proxies() []string
	mode() string
}

// Rainx Setting Implement
type rainxSetting struct {
	data map[string]map[string]interface{}
}

// Determine if setting.json is valid
func (s *rainxSetting) isValid() error {
	for _, v := range s.data {
		if _, is := v["listen"].(string); !is {
			return errors.New("listen must be a string")
		}
		if location, is := v["location"].(map[string]map[string]interface{}); !is {
			return errors.New("location must be a map<string, map<string, ?>>")
		} else {
			for _, loc := range location {
				root, isRoot := loc["root"]
				proxy, isProxy := loc["proxy"]
				if !isRoot && !isProxy {
					return errors.New("location is at least one root and proxy")
				} else if isRoot {
					if _, is := root.(string); !is {
						return errors.New("root must be a string")
					}
				} else if isProxy {
					switch proxy.(type) {
					case string:
						break
					case map[string]interface{}:
						proxies := proxy.(map[string]interface{})
						if _, is := proxies["path"].([]string); !is {
							return errors.New("path must be a string array")
						}
						if mode, is := proxies["mode"].(string); !is {
							return errors.New("mode must be a string")
						} else {
							switch mode {
							case "random":
							case "iphash":
							case "round":
								break
							default:
								return errors.New("mode must be in [random, iphash, round]")
							}
						}
						break
					default:
						return errors.New("proxy must be a string or map<string, ?>")
					}
				}
			}
		}
	}
	return nil
}

// Get all of Server
func (s *rainxSetting) servers() []ServerSetting {
	return nil
}

// New creates a new Setting
func NewSetting(filename string) (Setting, error) {
	var bytes []byte
	var err error
	if !FileExists(filename) {
		bytes = []byte(defaultJson)
		ioutil.WriteFile(filename, bytes, os.ModePerm)
		initHtml()
	} else {
		bytes, err = ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
	}
	data := make(map[string]map[string]interface{})
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &rainxSetting{data}, nil
}

// If the setting.json does not exist
// Ignore errors
func initHtml() {
	if !DirExists("html") {
		os.Mkdir("html", os.ModePerm)
	}
	if !FileExists("html/index.html") {
		ioutil.WriteFile("html/index.html", []byte(indexHtml), os.ModePerm)
	}
}
