package main

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
