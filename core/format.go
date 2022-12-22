package main

import (
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Dialog struct {
	Scope int
	Text  string
}

var tagRep = regexp.MustCompile(`\\_{0,2}[a-zA-Z0-9*!&\-+](\d|\[("([^"]|\\")+?"|([^\]]|\\\])+?)+?\])?`)
var noWordRep = regexp.MustCompile(`^[…‥.]+$`)
var delimRep = regexp.MustCompile(`[！!?？。]`)
var chScopeRep = regexp.MustCompile(`(\\([01])|\\c\[([0-9]+)\])`)

// スコープ切り替えを区切りとしてテキストを分割する
func splitDialog(text string) []Dialog {
	result := []Dialog{}

	submch := chScopeRep.FindAllStringSubmatch(text, -1)
	sep := chScopeRep.Split(text, -1)
	for i := range submch {
		sub := submch[i]

		if len(sub) != 4 {
			log.Println("abnormal submatch:", sub)
			continue
		}

		scope, err := strconv.Atoi(sub[2] + sub[3])
		if err != nil {
			log.Println("abnormal index:", sub)
			continue
		}

		result = append(result, Dialog{scope, sep[i+1]})
	}

	return result
}

// さくらスクリプトを削除する
func clearTags(src string) string {
	return tagRep.ReplaceAllString(src, "")
}

// ……。など読み上げのない文を代替文に置き換える
func processNoWordSentence(src string, config *Config) string {
	var result string

	punctuations := []string{}
	for _, p := range delimRep.FindAllStringSubmatch(src, -1) {
		punctuations = append(punctuations, p[0])
	}

	ary := delimRep.Split(src, -1)
	for i, s := range ary {
		if s == "" {
			continue
		}

		if noWordRep.MatchString(s) {
			s = config.NoWordPhrase
		}

		result += s
		if i < len(punctuations) {
			result += punctuations[i]
		}
	}
	return result
}

// クイックセクションを削除する
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
