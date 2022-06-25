package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-ps"
)

type Config struct {
	NoWordPhrase string `json:"NoWordPhrase"`
}

func initConfig(path string) (*Config, error) {
	config := &Config{
		NoWordPhrase: "んん",
	}

	fp, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	err = json.NewEncoder(fp).Encode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func loadConfig() (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(filepath.Dir(exePath), "config.json")

	fp, err := os.Open(path)
	if err != nil {
		return initConfig(path)
	}
	defer fp.Close()

	var config Config
	if err := json.NewDecoder(fp).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func isProcExist(name string) bool {
	var result bool

	processes, err := ps.Processes()

	if err != nil {
		os.Exit(1)
	}

	result = false
	for _, p := range processes {
		if p.Executable() == name {
			result = true
		}
	}

	return result
}

var tagRep = regexp.MustCompile(`\\_{0,2}[a-zA-Z0-9*!&\-+](\d|\[("([^"]|\\")+?"|([^\]]|\\\])+?)+?\])?`)
var noWordRep = regexp.MustCompile(`^[….]+$`)

func clearTags(src string) string {
	return tagRep.ReplaceAllString(src, "")
}

func processNoWordSentence(src string, config *Config) string {
	var result string
	ary := strings.Split(src, "。")
	for i, s := range ary {
		if noWordRep.MatchString(s) {
			s = config.NoWordPhrase
		}
		result += s
		if i != len(ary)-1 && s != "" {
			result += "。"
		}
	}
	return result
}

func deleteQuickSection(src string) string {
	getMinIndex := func(arr []int) int {
		result := -1
		min := math.MaxInt
		for i, n := range arr {
			if n >= 0 && min > n {
				result = i
				min = n
			}
		}
		return result
	}

	getStartPoint := func(s string) (int, int) {
		tags := []string{"\\![quicksection,1]", "\\![quicksection,true]", "\\_q"}

		var indexes []int
		for _, tag := range tags {
			indexes = append(indexes, strings.Index(s, tag))
		}

		minIndex := getMinIndex(indexes)

		var result int
		var length int
		if minIndex != -1 {
			result = indexes[minIndex]
			length = len(tags[minIndex])
		} else {
			result = -1
			length = 0
		}

		return result, length
	}

	getEndPoint := func(s string) int {
		tags := []string{"\\![quicksection,0]", "\\![quicksection,false]", "\\_q"}

		var indexes []int
		for _, tag := range tags {
			indexes = append(indexes, strings.Index(s, tag))
		}

		minIndex := getMinIndex(indexes)

		var result int
		if minIndex != -1 {
			result = indexes[minIndex] + len(tags[minIndex])
		} else {
			result = -1
		}

		return result
	}

	var result string
	isQuick := false
	for {
		if !isQuick {
			i, l := getStartPoint(src)
			if i == -1 {
				result += src
				break
			}
			result += src[:i]
			src = src[i+l:]
			isQuick = true
		} else {
			i := getEndPoint(src)
			if i == -1 {
				break
			}
			src = src[i:]
			isQuick = false
		}
	}

	return result
}

func establishTCP() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:50001")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func writeConn(conn net.Conn, sData []byte) error {
	_, err := conn.Write(sData)
	if err != nil {
		return err
	}
	return nil
}

func closeConn(conn net.Conn) {
	_ = conn.Close()
}

func main() {
	if len(os.Args) != 2 {
		err := errors.New("invalid arguments")
		log.Printf("error: %v\n", err)
	}

	config, err := loadConfig()
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	rawMsg, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	msg := []byte(processNoWordSentence(clearTags(deleteQuickSection(string(rawMsg))), config))

	if string(msg) == "" || !isProcExist("BouyomiChan.exe") {
		return
	}

	msg_length := uint32(len(msg))
	iCommand := []byte{1, 0}
	iSpeed := []byte{255, 255}
	iTone := []byte{255, 255}
	iVolume := []byte{255, 255}
	iVoice := []byte{0, 0}
	bCode := []byte{0}

	h := msg_length
	bMsgLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(bMsgLength, h)

	sData := append(iCommand, iSpeed...)
	sData = append(sData, iTone...)
	sData = append(sData, iVolume...)
	sData = append(sData, iVoice...)
	sData = append(sData, bCode...)
	sData = append(sData, bMsgLength...)
	sData = append(sData, msg...)

	conn, err := establishTCP()
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	err = writeConn(conn, sData)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	closeConn(conn)
}
