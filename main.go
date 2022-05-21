package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		panic(errors.New("invalid arguments"))
	}

	msg, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		panic(err)
	}
	msg_length := uint32(len(msg))
	iCommand := []byte{1, 0}
	iSpeed := []byte{255, 255}
	iTone := []byte{255, 255}
	iVolume := []byte{255, 255}
	iVoice := []byte{0, 0}
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
	sData = append(sData, msg...)

	conn, err := net.Dial("tcp", "localhost:50001")
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(sData)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	_ = conn.Close()
}
