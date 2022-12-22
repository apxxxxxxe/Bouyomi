package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	msg := `29d071fc35de88aaee91cb2ac25c18b5.voice,1,",",
b97f0f7b47aa2a3914a647a8a33925a4.voice,10020,",",
1f8faa026b96710b1d059d6345135f21.voice,3,",",`

	for _, line := range strings.Split(msg, "\n") {
		fmt.Println("line:", line)

		if !rep.MatchString(line) {
			fmt.Println("not match:", line)
			continue
		}

		c := strings.Split(line, ",")
		if len(c) < 2 {
			fmt.Println("count:", len(c))
			continue
		}

		VoiceNum, err := strconv.Atoi(c[1])
		if err != nil {
			continue
		}

		fmt.Println(c[0], VoiceNum)
	}
}

func TestSplit(t *testing.T) {
	msg := "\\0あいあいお\\1ほげほげふ\\c[12]"
	dialog := splitDialog(msg)
	for _, d := range dialog {
		fmt.Println(d)
	}
}

func TestList(t *testing.T) {
	voices, err := listVoices(true)
	if err != nil {
		t.Error(err)
	}
	for _, v := range voices {
		fmt.Printf("%v,%v\n", v.BouyomiChanNumber, v.Name)
	}
}

func TestSpeak(t *testing.T) {
	voices, err := listVoices(true)
	if err != nil {
		t.Error(err)
	}
	for _, v := range voices {
		b, err := dec2hex(v.BouyomiChanNumber)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(v, b)
	}

	speak("ほげほげふがふが", 10002)
	speak("ほげほげふがふが", 10003)
	speak("ほげほげふがふが", 1)
	speak("ほげほげふがふが", 0)
}
