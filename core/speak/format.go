package speak

import (
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/apxxxxxxe/Bouyomi/data"
)

type Dialog struct {
	Scope int
	Text  string
}

var tagRep = regexp.MustCompile(`\\_{0,2}[a-zA-Z0-9*!&\-+](\d|\[("([^"]|\\")+?"|([^\]]|\\\])+?)+?\])?`)
var delimRep = regexp.MustCompile(`[！!?？。]`)
var chScopeRep = regexp.MustCompile(`(\\([0h1u])|\\p\[([0-9]+)\])`)

// 元の文字列を処理してからスコープごとに分割する
func SplitDialog(src string, config *data.Config) []Dialog {
	// クイックセクションを削除
	src = deleteQuickSection(src)

	// スコープ切り替え部で分割
	res := splitDialog(src)

	for i := range res {
		// 各テキストを処理
		res[i].Text = clearTags(res[i].Text)
	}

	return res
}

// スコープ切り替えを区切りとしてテキストを分割する
func splitDialog(text string) []Dialog {
	result := []Dialog{}

	if text == "" {
		return result
	} else {
		// 文頭のセリフは\0,\hによる指定がなくても\0
		// 後の処理を共通化するためつけておく
		text = "\\0" + text
	}

	submch := chScopeRep.FindAllStringSubmatch(text, -1)
	sep := chScopeRep.Split(text, -1)
	for i, sub := range submch {
		if len(sub) != 4 {
			log.Println("abnormal submatch:", sub)
			continue
		}

		scope, err := strconv.Atoi(sub[2] + sub[3])
		if err != nil {
			switch sub[2] {
			case "h":
				scope = 0
			case "u":
				scope = 1
			default:
				log.Println("abnormal index:", sub)
				continue
			}
		}

		if d := sep[i+1]; d != "" {
			result = append(result, Dialog{scope, d})
		}
	}

	return result
}

// さくらスクリプトを削除する
func clearTags(src string) string {
	return tagRep.ReplaceAllString(src, "")
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
