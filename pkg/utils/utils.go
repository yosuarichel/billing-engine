package utils

import (
	"hash/fnv"
	"math"
	"net"
	"os"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/yosuarichel/billing-engine/pkg/shared/domain"
)

func GetEnvWithKey(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func GetEnv() string {
	return GetEnvWithKey(APP_ENV, "")
}

func GetAppType() string {
	return GetEnvWithKey(APP_TYPE, "")
}

func GetAppName() string {
	return GetEnvWithKey(APP_NAME, "")
}

func GetMachineID() (int, error) {
	hostname, err := os.Hostname()
	klog.Info("HOSTNAME ====> ", hostname)
	if err != nil {
		return 0, err
	}
	h := fnv.New32a()
	h.Write([]byte(hostname))
	return int(h.Sum32() % 65535), nil
}

// Paginate generic function
func Paginate[T any](list []T, pageNum, pageSize int64) ([]T, int64) {
	if len(list) == 0 || list == nil {
		return []T{}, 0
	}
	total := int64(len(list))

	if pageNum <= 0 || pageSize <= 0 {
		return list, total
	}

	startIndex := (pageNum - 1) * pageSize
	if startIndex >= total {
		return []T{}, total
	}

	endIndex := startIndex + pageSize
	if endIndex > total {
		endIndex = total
	}

	return list[startIndex:endIndex], total
}

func CalculateTotalAmount(principal int64, interestRate float64, termWeeks int) int64 {
	years := math.Ceil(float64(termWeeks) / domain.WeeksInYear)
	interestExact := float64(principal) * interestRate * years
	totalExact := float64(principal) + interestExact

	return int64(math.Floor(totalExact))
}

func GetLocalIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "127.0.0.1"
}
