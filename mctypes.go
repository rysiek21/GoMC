package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

func ReadVarInt(reader io.ByteReader) (int, error) {
	var value int
	var position int

	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		value |= int(b&0x7F) << position
		if (b & 0x80) == 0 {
			break
		}
		position += 7
		if position >= 32 {
			return 0, errors.New("Wrong VarInt size")
		}
	}
	return value, nil
}

func MakeVarInt(value int) []byte {
	var final []byte
	for {
		b := value & 0x7F
		value >>= 7
		if value != 0 {
			final = append(final, byte(b|0x80))
		} else {
			final = append(final, byte(b))
			break
		}
	}
	return final
}

func ReadString(reader *bytes.Reader) (string, error) {
	length, err := ReadVarInt(reader)
	if err != nil {
		return "", err
	}

	stringBytes := make([]byte, length)
	_, err = io.ReadFull(reader, stringBytes)
	if err != nil {
		return "", err
	}

	return string(stringBytes), nil
}

func MakeString(value string) []byte {
	var final []byte
	length := len(value)
	final = append(final, MakeVarInt(length)...)
	final = append(final, []byte(value)...)
	return final
}

func ReadUnsignedShort(reader *bytes.Reader) (uint16, error) {
	var ushort uint16
	err := binary.Read(reader, binary.BigEndian, &ushort)
	return ushort, err
}
