package smhb

import (
	"errors"
	"strconv"
	"strings"
)

var replicas = []string{
	"localhost:19138",
	"localhost:19139",
}

func NextReplicaIndex(i int) int {
	return (i + 1) % len(replicas)
}

func GetReplicaAddress(i int) (string, int, error) {
	if i >= len(replicas) {
		return "", -1, errors.New("invalid replica index")
	}

	pair := strings.Split(replicas[i], ":")
	port, err := strconv.Atoi(pair[1])

	if err != nil {
		return "", -1, err
	}

	return pair[0], port, nil
}
