package manifest

import "github.com/gogf/gf/v2/frame/g"

var (
	JWTKey = []byte("hamster-home")
)

// home_network.go配置文件
const (
	HomeNetworkRouterIP       = ""
	HomeNetworkRouterAddress  = ""
	HomeNetworkRouterPassword = ""
)

// node_utils.go配置文件
const (
	XrayBaseUrl      = "ray.sample.top:580" // PS.不加http://
	XrayLoginDataMap = `{"username":"","password":""}`
)

// proxy_network.go配置文件
var (
	XuiBaseUrl      = "http://"
	XuiLoginDataMap = g.Map{
		"username": "",
		"password": "",
	}
)

// nezha 配置文件
const (
	NezhaApiUrl = "http://ip:port/api/v1/server/details?id="
	NezhaApiKey = "token"
)

// uptime_kuma.配置文件

const (
	UptimeKumaApiUrl = "http://ip:port/api/status-page/heartbeat/hamster"
)

// coffee.go配置文件
var (
	ProxyProviderBaseUrl  = "https://*****/api/v1/user/getSubscribe"
	ProxyProviderLoginUrl = "https://*****/api/v1/passport/auth/login"
	ProxyProviderAuthData = g.Map{
		"email":    "",
		"password": "",
	}
)

// docker部分配置文件
var (
	DockerApiUrl  = "http://ip:port/api/endpoints"
	DockerAuthMap = map[string]string{"x-api-key": ""}
)

// home_assistant配置文件
var (
	HomeAssistantBaseUrl = "http://ip:port/api/states/"
	HomeAssistantAuthMap = map[string]string{
		"Authorization": "Bearer ",
	}
)

// cache key
const (
	HomeNetworkCacheKey    = "homeNetwork"
	ProxyNetworkCacheKey   = "proxyNetwork"
	ProxySessionCacheKey   = "proxySession"
	ProxyNodeCacheKey      = "proxyNode"
	ProxySubscribeCacheKey = "proxySubscribe"
	ServerDataCacheKey     = "serverDataList"
	DockerMonitorCacheKey  = "dockerMonitor"
	UptimeCacheKey         = "uptime"
	HaEntitiesCacheKey     = "haEntities"
	XUIUserListCacheKey    = "xuiUserList"
)
