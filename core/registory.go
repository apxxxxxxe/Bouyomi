package main

import (
	"golang.org/x/sys/windows/registry"
)

type Voice struct {
	BouyomiChanNumber int16
	Name              string
	Language          string
}

func listVoices(japaneseOnly bool) ([]Voice, error) {
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

	voices := []Voice{
		{0, "棒読みちゃん上で指定(デフォルト)", "0"},
		{1, "女性1", "411"},
		{2, "女性2", "411"},
		{3, "男性1", "411"},
		{4, "男性2", "411"},
		{5, "中性", "411"},
		{6, "ロボット", "411"},
		{7, "機械1", "411"},
		{8, "機械1", "411"},
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
