package edcoder

import (
	"bufio"
	"errors"
	"io"
	"os"
)

func readFile(fileName string) (string, error) {
	fi, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer fi.Close()

	c, err := readConf(fi)
	if nil != err {
		return "", err
	}
	return c, nil
}

func readConf(fi io.Reader) (string, error) {
	var chunks []byte
	r := bufio.NewReader(fi)
	buf := make([]byte, 4096)
	if nil == buf {
		return "", errors.New("make []byte error")
	}
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if 0 == n {
			break
		}

		chunks = append(chunks, buf[:n]...)
	}

	return string(chunks), nil
}
