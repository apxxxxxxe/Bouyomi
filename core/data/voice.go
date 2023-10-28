package data

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type Voice struct {
	BouyomiChanNumber int16
	Name              string
	Language          string
}

type VoiceMap = map[string]int16

const varPath = "yaya_variable.cfg"

var rep = regexp.MustCompile(`^[a-z0-9]+\.voice`)

func ListVoices(noVoiceByDefault bool, japaneseOnly bool) ([]Voice, error) {
	const tokenPath = `SOFTWARE\WOW6432Node\Microsoft\SPEECH\Voices\Tokens`

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, tokenPath, registry.READ)
	if err != nil {
		return nil, err
	}
	defer k.Close()

	sub, err := k.ReadSubKeyNames(0)
	if err != nil {
		return nil, err
	}

	defaultVoice := Voice{0, "棒読みちゃん上で指定(デフォルト)", "0"}
	if noVoiceByDefault {
		defaultVoice = Voice{0, "読み上げなし", "0"}
	}

	voices := []Voice{
		defaultVoice,
		{1, "女性1", "411"},
		{2, "女性2", "411"},
		{3, "男性1", "411"},
		{4, "男性2", "411"},
		{5, "中性", "411"},
		{6, "ロボット", "411"},
		{7, "機械1", "411"},
		{8, "機械2", "411"},
	}
	idx := int16(10001)
	for _, s := range sub {
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, tokenPath+"\\"+s+"\\Attributes", registry.READ)
		if err != nil {
			return nil, err
		}
		defer k.Close()

		name, _, err := k.GetStringValue("Name")
		if err != nil {
			return nil, err
		}

		lang, _, err := k.GetStringValue("Language")
		if err != nil {
			return nil, err
		}

		if !japaneseOnly || lang == "411" || lang == "0" {
			voices = append(voices, Voice{idx, name, lang})
			idx++
		}
	}

	return voices, nil
}

func FindVoice(v VoiceMap, ghost string, scope int) int16 {
	key := fmt.Sprintf("%x", md5.Sum([]byte(ghost+strconv.Itoa(scope)))) + ".voice"
	if voice, ok := v[key]; ok {
		return voice
	} else {
		return 0
	}
}

func LoadVoiceMap() (VoiceMap, error) {
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
