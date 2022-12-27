package data

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"sort"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var repCharacter = regexp.MustCompile(`(?m)^(sakura|kero|char[0-9]+)\.name,([^\r\n]*)`)
var repNum = regexp.MustCompile(`\d+`)
var repSjis = regexp.MustCompile(`(?i)shift_jis`)

// descript.txtを読み込んでキャラクター名(\0,\1,...)を返す
func LoadDescript(ghostpath string) ([]string, error) {
	path := filepath.Join(ghostpath, "ghost", "master", "descript.txt")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if repSjis.Match(b) {
		// sjis指定があればデコードして読み直す
		b, err = ioutil.ReadAll(transform.NewReader(bytes.NewBuffer(b), japanese.ShiftJIS.NewDecoder()))
		if err != nil {
			return nil, err
		}
	}

	f := repCharacter.FindAllString(string(b), -1)
	sort.Slice(f, func(i, j int) bool {
		if strings.HasPrefix(string(f[i]), "sakura") {
			// sakura kero : true
			// sakura char : true
			return true
		} else if strings.HasPrefix(string(f[i]), "kero") {
			// kero sakura : false
			// kero char   : true
			return !strings.HasPrefix(string(f[j]), "sakura")
		} else { // strings.HasPrefix(string(f[i]), "char")
			if !strings.HasPrefix(string(f[j]), "char") {
				// char sakura : false
				// char kero   : false
				return false
			} else {
				// char char   : ?
				numI, err := strconv.Atoi(repNum.FindString(f[i]))
				if err != nil {
					panic(err)
				}
				numJ, err := strconv.Atoi(repNum.FindString(f[j]))
				if err != nil {
					panic(err)
				}
				return numI < numJ
			}
		}
	})
	return f, nil
}
