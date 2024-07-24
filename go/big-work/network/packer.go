package network

import (
	"encoding/binary"
	"io"
)

type NormalPacker struct {
	ByteOrder binary.ByteOrder
}

func NewNormalPacker(order binary.ByteOrder) *NormalPacker {
	return &NormalPacker{
		ByteOrder: order,
	}
}

// Pack |data 长度|id|data|
func (p *NormalPacker) Pack(message *Message) ([]byte, error) {
	buffer := make([]byte, 8+8+len(message.Data))
	// 大小端 CPU小端 网络大端

	p.ByteOrder.PutUint64(buffer[:8], uint64(len(buffer)))
	p.ByteOrder.PutUint64(buffer[8:16], message.Id)
	copy(buffer[16:], message.Data)

	return buffer, nil
}

func (p *NormalPacker) Unpack(reader io.Reader) (*Message, error) {
	//fmt.Println("UNPPPP")
	/*
		err := reader.(*net.TCPConn).SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			return nil, err
		}
	*/

	buffer := make([]byte, 8+8)
	_, err := io.ReadFull(reader, buffer)
	if err != nil {
		return nil, err
	}
	totalLen := p.ByteOrder.Uint64(buffer[:8])
	id := p.ByteOrder.Uint64(buffer[8:])
	dataLen := totalLen - 8 - 8
	//fmt.Println("totalLen:", totalLen, " id:", id, " dataLen:", dataLen)

	data := make([]byte, dataLen)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Id:   id,
		Data: data,
	}

	return msg, nil
}
