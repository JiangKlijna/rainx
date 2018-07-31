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
	IsValid() error
	Servers() []ServerSetting
}

// Server Setting Interface
type ServerSetting interface {
	Name() string
	Listen() string
	Locations() []LocationSetting
}

// Location Setting Interface
// Get detailed configuration
type LocationSetting interface {
	Pattern() string
	IsRoot() bool
	Root() string
	IsProxy() bool
	Proxy() string
	IsProxies() bool
	Proxies() []string
	Mode() string
}

// Rainx Setting Implement
type rainxSetting struct {
	data map[string]map[string]interface{}
}

// Rainx Server Setting Implement
type rainxServerSetting struct {
	key  string
	data map[string]interface{}
}

// Rainx Location Setting Implement
type rainxLocatioSetting struct {
	key  string
	data map[string]interface{}
}

// Determine if setting.json is valid
func (s *rainxSetting) IsValid() error {
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
func (s *rainxSetting) Servers() []ServerSetting {
	i := 0
	arr := make([]ServerSetting, len(s.data))
	for k, v := range s.data {
		arr[i] = &rainxServerSetting{k, v}
		i++
	}
	return arr
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

// Get Server Name
func (s rainxServerSetting) Name() string {
	return s.key
}

// Get Addr
func (s rainxServerSetting) Listen() string {
	return s.data["listen"].(string)
}

// Get All of Location
func (s rainxServerSetting) Locations() []LocationSetting {
	return nil
}

// Get Pattern of Location
func (s *rainxLocatioSetting) Pattern() string {
	return s.key
}

// Determine if the [root] is included
func (s *rainxLocatioSetting) IsRoot() bool {
	_, is := s.data["root"].(string)
	return is
}

// Get [root] value
func (s *rainxLocatioSetting) Root() string {
	root, _ := s.data["root"].(string)
	return root
}

// Determine if the [proxy] is included
// Determine if the [proxy] is string
func (s *rainxLocatioSetting) IsProxy() bool {
	_, is := s.data["proxy"].(string)
	return is
}

// Get [proxy] value
func (s *rainxLocatioSetting) Proxy() string {
	proxy, _ := s.data["proxy"].(string)
	return proxy
}

// Determine if the [proxy] is included
// Determine if the [proxy] is map<string, ?>
func (s *rainxLocatioSetting) IsProxies() bool {
	_, is := s.data["proxy"].(map[string]interface{})
	return is
}

// Get [proxy.path] value
func (s *rainxLocatioSetting) Proxies() []string {
	proxies, _ := s.data["proxy"].(map[string]interface{})
	return proxies["path"].([]string)
}

// Get [proxy.mode] value
func (s *rainxLocatioSetting) Mode() string {
	proxies, _ := s.data["proxy"].(map[string]interface{})
	return proxies["mode"].(string)
}
