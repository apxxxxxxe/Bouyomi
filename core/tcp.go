package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/mitchellh/go-ps"
)

// 棒読みちゃんに文字列を読み上げてもらう
func speak(msg string, voice int16) error {
	if msg == "" || !isProcExist("BouyomiChan.exe") {
		return nil
	}

	bMsg := []byte(msg)

	msg_length := uint32(len(bMsg))
	iCommand := []byte{1, 0}
	iSpeed := []byte{255, 255}
	iTone := []byte{255, 255}
	iVolume := []byte{255, 255}
	iVoice, err := dec2hex(voice)
	if err != nil {
		return err
	}
	bCode := []byte{0}

	h := msg_length
	bMsgLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(bMsgLength, h)

	sData := append(iCommand, iSpeed...)
	sData = append(sData, iTone...)
	sData = append(sData, iVolume...)
	sData = append(sData, iVoice...)
	sData = append(sData, bCode...)
	sData = append(sData, bMsgLength...)
	sData = append(sData, bMsg...)

	// 棒読みちゃんと通信を確立
	conn, err := net.Dial("tcp", "localhost:50001")
	if err != nil {
		return err
	}

	// 読み上げデータを送信
	_, err = conn.Write(sData)
	if err != nil {
		return err
	}

	_ = conn.Close()

	return nil
}

// 声質番号をバイト列に変換 ex)10001 -> {17, 39}
func dec2hex(d int16) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint16(d))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes(), nil
}

// プロセスの存在を返す
func isProcExist(name string) bool {
	var result bool

	processes, err := ps.Processes()

	if err != nil {
		os.Exit(1)
	}

	result = false
	for _, p := range processes {
		if p.Executable() == name {
			result = true
		}
	}

	return result
}
