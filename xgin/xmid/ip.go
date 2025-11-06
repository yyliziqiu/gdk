package xmid

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yyliziqiu/gdk/xerr"
	"github.com/yyliziqiu/gdk/xgin"
	"github.com/yyliziqiu/gdk/xgin/xresp"
	"github.com/yyliziqiu/gdk/xlog"
)

var ErrInvalidIpFormat = errors.New("invalid ip format")

var _ipRanges []uint32

func CheckIp(ips []string) gin.HandlerFunc {
	for _, ip := range ips {
		r, err := parseIpRange(ip)
		if err != nil {
			xlog.Errorf("Parse ip failed, ip: %s, error: %v.", ip, err)
			continue
		}
		_ipRanges = append(_ipRanges, r)
	}

	return func(ctx *gin.Context) {
		// ip := ctx.RemoteIP() // 这个 ip 地址是客户端与服务器建立 tcp 连接时的原始 ip
		// ip := ctx.ClientIP() // 先检查 X-Forwarded-For、X-Real-IP 等 header，如果没有代理信息或代理不受信任则退回到 RemoteIP()

		ip := ctx.RemoteIP()

		iv, err := ip2int(ip)
		if err != nil {
			if logger := xgin.GetLogger(); logger != nil {
				logger.Warnf("Parse ip failed, ip: %s, error: %v.", ip, err)
			}
			xresp.AbortError(ctx, xerr.ForbiddenIp)
			return
		}

		for _, r := range _ipRanges {
			if iv&r == r {
				return
			}
		}

		xresp.AbortError(ctx, xerr.ForbiddenIp)
	}
}

func ip2int(ip string) (uint32, error) {
	ipv4 := net.ParseIP(ip).To4()
	if len(ipv4) == 0 {
		return 0, ErrInvalidIpFormat
	}

	i := uint32(0)
	for _, b := range ipv4 {
		i = i << 8
		i |= uint32(b)
	}

	return i, nil
}

func parseIpRange(ip string) (uint32, error) {
	im := strings.Split(ip, "/")
	if len(im) == 0 || len(im) > 2 {
		return 0, ErrInvalidIpFormat
	}

	mk := 32
	if len(im) == 2 {
		i, err := strconv.Atoi(im[1])
		if err != nil {
			return 0, ErrInvalidIpFormat
		}
		ip = im[0]
		mk = i
	}

	iv, err := ip2int(ip)
	if err != nil {
		return 0, err
	}

	return iv & (uint32(0xffffffff) << (32 - mk)), nil
}
