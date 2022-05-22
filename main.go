package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"net"
	"os"
	"regexp"
	"strings"
)

func clearTags(src string) string {
	return regexp.MustCompile(`\\_{0,2}[a-zA-Z0-9*!&](\d|\[("([^"]|\\")+?"|([^\]]|\\\])+?)+?\])?`).ReplaceAllString(src, "")
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

func main() {
	if len(os.Args) != 2 {
		log.Printf("error: %v\n", errors.New("invalid arguments"))
	}

	rawMsg, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	msg := []byte(clearTags(deleteQuickSection(string(rawMsg))))

	if string(msg) == "" {
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

	conn, err := net.Dial("tcp", "localhost:50001")
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	_, err = conn.Write(sData)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	_ = conn.Close()
}
