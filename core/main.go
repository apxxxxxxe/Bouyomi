package main

import (
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/apxxxxxxe/Bouyomi/data"
	"github.com/apxxxxxxe/Bouyomi/speak"
)

const defaultBouyomiPath = "null"

func main() {
	var (
		execEngines   bool
		showList      bool
		voice64       int
		ghostName     string
		getHash       string
		getCharaCount string
	)
	flag.BoolVar(&execEngines, "e", false, "execute tools")
	flag.BoolVar(&showList, "l", false, "show available voices")
	flag.IntVar(&voice64, "v", 0, "voice number")
	flag.StringVar(&ghostName, "g", "", "ghost name")
	flag.StringVar(&getHash, "hash", "", "get hash")
	flag.StringVar(&getCharaCount, "count", "", "get character names")
	flag.Parse()

	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile(filepath.Join(filepath.Dir(exePath), "bouyomi.log"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(logFile)
	defer logFile.Close()

	config, err := data.LoadConfig()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	var wg sync.WaitGroup
	if execEngines && config.EnginesPath != nil {
		for _, path := range config.EnginesPath {
			if path == "" {
				continue
			}
			if _, info := os.Stat(path); info != os.ErrNotExist && !data.IsProcExist(filepath.Base(path)) {
				wg.Add(1)
				go func(p string) {
					defer wg.Done()
					if err := data.ExecCommand(p); err != nil {
						log.Fatalf("error: %v\n", err)
					}
				}(path)
			}
		}
		wg.Wait()
		os.Exit(0)
	}

	voices, err := data.ListVoices(*config.NoVoiceByDefault, *config.JapaneseOnly)
	if err != nil {
		log.Fatalf("error: %v\n", err)
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
			log.Fatalf("error: %v\n", err)
		}

		delim := "\u0001"
		res := ""
		for _, line := range f {
			res += strings.Split(line, ",")[1] + delim
		}
		fmt.Print(strings.TrimSuffix(res, delim))
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		err := errors.New("invalid arguments")
		log.Fatalf("error: %v\n", err)
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
		log.Fatalf("error: %v\n", err)
	}

	var rawMsg []byte
	rawMsg = []byte(flag.Arg(0))
	rawMsg = []byte(flag.Arg(0))
	voiceMap, err := data.LoadVoiceMap()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	// 各セリフを読み上げさせる
	for _, dialog := range speak.SplitDialog(string(rawMsg), config) {
		if dialog.Scope == 0 && *config.NoVoiceByDefault {
			continue
		}
		if err := speak.Speak(dialog.Text, data.FindVoice(voiceMap, ghostName, dialog.Scope)); err != nil {
			log.Fatalf("error: %v\n", err)
		}
	}

	os.Exit(0)
}
