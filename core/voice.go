package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const varPath = "yaya_variable.cfg"

type VoiceMap = map[string]int16

func findVoice(v VoiceMap, ghost string, scope int) int16 {
	key := fmt.Sprintf("%x", md5.Sum([]byte(ghost+strconv.Itoa(scope)))) + ".voice"
	if voice, ok := v[key]; ok {
		return voice
	} else {
		return 0
	}
}

var rep = regexp.MustCompile(`^[a-z0-9]+\.voice`)

func loadVoiceMap() (VoiceMap, error) {
	result := VoiceMap{}

	exec, err := os.Executable()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(filepath.Join(filepath.Dir(exec), varPath))
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(b), "\n") {
		if !rep.Match([]byte(line)) {
			continue
		}

		c := strings.Split(line, ",")
		if len(c) < 2 {
			continue
		}

		VoiceNum, err := strconv.Atoi(c[1])
		if err != nil {
			continue
		}

		result[c[0]] = int16(VoiceNum)
	}

	return result, nil
}
