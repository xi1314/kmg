package kson

import (
	"encoding/binary"
	"encoding/json"
	"testing"
	"time"
)

type TcpPacket struct {
	Time    time.Time
	SrcIp   string
	SrcPort uint16
	DstIp   string
	DstPort uint16
	Seq     uint32
	Ack     uint32
	//Len         int
	//ACK         bool
	//SYN         bool
	//FIN         bool
	//RST         bool
	//Window      uint16 // 2-65535 必须乘 SYN里面的WindowScale值.
	//WindowScale uint8  // 0 表示没有这一项 SYN里面的WindowScale值.
	//TimeStamp   uint32 // tcp的timestamp扩展里面的对方发过来的时间.
	//
	//				   // 其他参数似乎没有什么作用,直接忽略
	//connId string
	//id     string
	//ackId  string
}

func MarshalTcpPacket(o *TcpPacket) (b []byte, err error) {
	timeb, err := o.Time.MarshalBinary() //13.73%
	if err != nil {
		return nil, err
	}
	timeBSize := len(timeb)
	srcIpSize := len(o.SrcIp)
	dstIpSize := len(o.DstIp)
	size := 1 + timeBSize + 1 + srcIpSize + 2 + 1 + dstIpSize + 2 + 4 + 4
	b = make([]byte, size) // 24.44%
	pos := 0
	b[pos] = byte(timeBSize)
	pos++
	copy(b[pos:pos+timeBSize], timeb)
	pos += timeBSize

	b[pos] = byte(srcIpSize)
	pos++
	copy(b[pos:pos+srcIpSize], o.SrcIp)
	pos += srcIpSize

	binary.LittleEndian.PutUint16(b[pos:pos+2], o.SrcPort)
	pos += 2

	b[pos] = byte(dstIpSize)
	pos++
	copy(b[pos:pos+srcIpSize], o.SrcIp)
	pos += dstIpSize

	binary.LittleEndian.PutUint16(b[pos:pos+2], o.DstPort)
	pos += 2

	binary.LittleEndian.PutUint32(b[pos:pos+4], o.Seq)
	pos += 4

	binary.LittleEndian.PutUint32(b[pos:pos+4], o.Ack)
	pos += 4

	return b, nil
}

func UnmarshalTcpPacket(b []byte) (packet *TcpPacket, err error) {
	packet = &TcpPacket{} // 13.73%
	pos := 0
	timeSize := b[pos]
	pos++
	err = packet.Time.UnmarshalBinary(b[pos : pos+int(timeSize)]) //4.39%
	if err != nil {
		return nil, err
	}
	pos += int(timeSize)

	sSize := b[pos]
	packet.SrcIp = string(b[pos : pos+int(sSize)]) //12.48%
	pos += int(sSize)

	packet.SrcPort = binary.LittleEndian.Uint16(b[pos : pos+2])
	pos += 2

	sSize = b[pos]
	packet.DstIp = string(b[pos : pos+int(sSize)])
	pos += int(sSize)

	packet.DstPort = binary.LittleEndian.Uint16(b[pos : pos+2])
	pos += 2

	packet.Seq = binary.LittleEndian.Uint32(b[pos : pos+4])
	pos += 4

	packet.Ack = binary.LittleEndian.Uint32(b[pos : pos+4])
	pos += 4
	return packet, nil
}

/*
BenchmarkJsonMarshal-4            300000              4235 ns/op          28.80 MB/s         600 B/op          7 allocs/op
BenchmarkJsonUnmarshal-4          200000              6197 ns/op          19.68 MB/s         400 B/op          9 allocs/op
BenchmarkKsonMarshal-4          10000000               213 ns/op         205.85 MB/s          64 B/op          2 allocs/op
BenchmarkKsonUnmarshal-4        10000000               223 ns/op         197.04 MB/s          96 B/op          2 allocs/op

*/
func BenchmarkJsonMarshal(b *testing.B) {
	data := TcpPacket{
		Time:  time.Now().In(time.UTC),
		SrcIp: "1.2.3.4",
		DstIp: "1.2.3.4",
	}
	total := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		total += len(buf)
	}
	b.SetBytes(int64(total / b.N))
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	data := TcpPacket{
		Time:  time.Now().In(time.UTC),
		SrcIp: "1.2.3.4",
		DstIp: "1.2.3.4",
	}
	buf, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := json.Unmarshal(buf, &data)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkKsonMarshal(b *testing.B) {
	data := &TcpPacket{
		Time:  time.Now().In(time.UTC),
		SrcIp: "1.2.3.4",
		DstIp: "1.2.3.4",
	}
	total := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf, err := MarshalTcpPacket(data)
		if err != nil {
			panic(err)
		}
		total += len(buf)
	}
	b.SetBytes(int64(total / b.N))
}

func BenchmarkKsonUnmarshal(b *testing.B) {
	data := &TcpPacket{
		Time:  time.Now().In(time.UTC),
		SrcIp: "1.2.3.4",
		DstIp: "1.2.3.4",
	}
	buf, err := MarshalTcpPacket(data)
	if err != nil {
		panic(err)
	}
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := UnmarshalTcpPacket(buf)
		if err != nil {
			panic(err)
		}
	}
}
