package main

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		showList      bool
		isBase64      bool
		getHash       bool
		getCharaCount bool
		voice64       int
		ghostName     string
	)
	flag.BoolVar(&showList, "l", false, "show available voices")
	flag.BoolVar(&isBase64, "base64", false, "input base64 message")
	flag.BoolVar(&getHash, "hash", false, "show hash")
	flag.BoolVar(&getCharaCount, "count", false, "count characters")
	flag.IntVar(&voice64, "v", 0, "voice number")
	flag.StringVar(&ghostName, "g", "", "ghost name")
	flag.Parse()

	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	voices, err := listVoices(config.JapaneseOnly)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if showList {
		delim := "\u0001"
		p := ""
		for _, v := range voices {
			p += fmt.Sprintf("%v,%v%s", v.BouyomiChanNumber, v.Name, delim)
		}
		fmt.Println(strings.TrimSuffix(p, delim))
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		err := errors.New("invalid arguments")
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if getHash {
		// arg[0]: ハッシュ化したい文字列
		// 文字列をmd5エンコードして返す
		fmt.Printf("%x", md5.Sum([]byte(flag.Arg(0))))
		os.Exit(0)
	}

	if getCharaCount {
		// arg[0]: ゴーストのフルパス
		// 指定ゴーストのdescript.txtの名前指定からキャラクター数をカウントして返す

		f, err := loadDescript(flag.Arg(0))
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

	voiceMap, err := loadVoiceMap()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	baseText := deleteQuickSection(string(rawMsg))

	// 各セリフを読み上げさせる
	for _, dialog := range splitDialog(baseText) {
		msg := processNoWordSentence(clearTags(dialog.Text), config)
		voice := findVoice(voiceMap, ghostName, dialog.Scope)
		if err := speak(msg, voice); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
