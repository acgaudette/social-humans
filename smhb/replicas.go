package smhb

import (
	"log"
	"strconv"
	"strings"
)

var replicas = []string{"0.0.0.0:19138", "0.0.0.0:19139"}

func NextServerIdx(idx int) int {
	return (idx + 1) % len(replicas)
}

func GetAddressAndPort(idx int) (string, int) {
	return ParseAddressAndPort(replicas[idx])
}

func ParseAddressAndPort(entry string) (string, int) {
	pair := strings.Split(entry, ":")
	port, err := strconv.Atoi(pair[1])
	if err != nil {
		log.Printf("%s", err)
	}
	return pair[0], port
}
