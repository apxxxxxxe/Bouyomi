package main

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/apxxxxxxe/Bouyomi/data"
	"github.com/apxxxxxxe/Bouyomi/speak"
)

func main() {
	var (
		showList      bool
		isBase64      bool
		voice64       int
		ghostName     string
		getHash       string
		getCharaCount string
	)
	flag.BoolVar(&showList, "l", false, "show available voices")
	flag.BoolVar(&isBase64, "base64", false, "input base64 message")
	flag.IntVar(&voice64, "v", 0, "voice number")
	flag.StringVar(&ghostName, "g", "", "ghost name")
	flag.StringVar(&getHash, "hash", "", "get hash")
	flag.StringVar(&getCharaCount, "count", "", "get character names")
	flag.Parse()

	config, err := data.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	voices, err := data.ListVoices(config.JapaneseOnly)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if showList {
		// 利用可能な声質のリストを返す
		delim := "\u0001"
		p := ""
		for _, v := range voices {
			p += fmt.Sprintf("%v,%v%s", v.BouyomiChanNumber, v.Name, delim)
		}
		fmt.Println(strings.TrimSuffix(p, delim))
		os.Exit(0)
	} else if getHash != "" {
		// 文字列をmd5エンコードして返す
		fmt.Printf("%x", md5.Sum([]byte(getHash)))
		os.Exit(0)
	} else if getCharaCount != "" {
		// 指定ゴーストのdescript.txtの名前指定からキャラクター数をカウントして返す
		f, err := data.LoadDescript(getCharaCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		delim := "\u0001"
		res := ""
		for _, line := range f {
			res += strings.Split(line, ",")[1] + delim
		}
		fmt.Printf(strings.TrimSuffix(res, delim))
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		err := errors.New("invalid arguments")
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	voice := int16(voice64)
	isValidVoice := false
	for _, v := range voices {
		if v.BouyomiChanNumber == voice {
			isValidVoice = true
		}
	}
	if !isValidVoice {
		err := errors.New("invalid voice number")
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var rawMsg []byte
	if isBase64 {
		// 文字化け防止のためbase64で渡されたセリフをデコードする
		rawMsg, err = base64.StdEncoding.DecodeString(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	} else {
		rawMsg = []byte(flag.Arg(0))
	}

	voiceMap, err := data.LoadVoiceMap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// 各セリフを読み上げさせる
	for _, dialog := range speak.SplitDialog(string(rawMsg), config) {
		if err := speak.Speak(dialog.Text, data.FindVoice(voiceMap, ghostName, dialog.Scope)); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
