package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

func clearTags(src string) string {
	return regexp.MustCompile(`\\_{0,2}[a-zA-Z0-9*!&](\d|\[("([^"]|\\")+?"|([^\]]|\\\])+?)+?\])?`).ReplaceAllString(src, "")
}

func deleteQuickSection(src string) string {
	var result string

	min := func(x, y int) int {
		if x > y {
			return y
		} else {
			return x
		}
	}

	getStartPoint := func(s string) (int, int) {
		var (
			result int
			length int
		)
		quickStart := "\\![quicksection,true]"
		quick := "\\_q"
		index_true := strings.Index(s, quickStart)
		index := strings.Index(s, quick)

		if index_true != -1 && index != -1 {
			if index < index_true {
				result = index
				length = len(quick)
			} else {
				result = index_true
				length = len(quickStart)
			}
		} else if index_true == -1 && index != -1 {
			result = index
			length = len(quick)
		} else if index_true != -1 && index == -1 {
			result = index_true
			length = len(quickStart)
		} else {
			result = -1
			length = 0
		}

		return result, length
	}

	getEndPoint := func(s string) int {
		var result int
		quickEnd := "\\![quicksection,false]"
		quick := "\\_q"
		index_end := strings.Index(s, quickEnd)
		index := strings.Index(s, quick)

		if index_end != -1 && index != -1 {
			result = min(index_end+len(quickEnd), index+len(quick))
		} else if index_end == -1 && index != -1 {
			result = index + len(quick)
		} else if index_end != -1 && index == -1 {
			result = index_end + len(quickEnd)
		} else {
			result = -1
		}

		return result
	}

	isquick := false
	for {
		if !isquick {
			i, l := getStartPoint(src)
			if i == -1 {
				result += src
				break
			}
			result += src[:i]
			src = src[i+l:]
			isquick = true
		} else {
			i := getEndPoint(src)
			if i == -1 {
				break
			}
			src = src[i:]
			isquick = false
		}
	}

	return result
}

func main() {

	if len(os.Args) != 2 {
		panic(errors.New("invalid arguments"))
	}

	rawMsg, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		panic(err)
	}
	msg := []byte(clearTags(deleteQuickSection(string(rawMsg))))

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
		panic(err)
	}

	_, err = conn.Write(sData)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	_ = conn.Close()
}
