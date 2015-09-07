package kson

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

func MarshalTcpPacketV2(o *TcpPacket) (b []byte, err error) {
	//timeb,err:=o.Time.MarshalBinary() //13.73%
	//if err!=nil{
	//	return nil,err
	//}
	//timeBSize:=len(timeb)
	srcIpSize := len(o.SrcIp)
	dstIpSize := len(o.DstIp)
	//size:=1+timeBSize+1+srcIpSize+2+1+ dstIpSize +2+4+4
	size := 1 + srcIpSize + 2 + 1 + dstIpSize + 2 + 4 + 4
	b = make([]byte, size) // 24.44%
	pos := 0
	//b[pos] = byte(timeBSize)
	//pos++
	//copy(b[pos:pos+timeBSize],timeb)
	//pos+=timeBSize

	b[pos] = byte(srcIpSize)
	pos++
	copy(b[pos:pos+srcIpSize], o.SrcIp)
	pos += srcIpSize

	b[pos] = byte(o.SrcPort)
	b[pos+1] = byte(o.SrcPort >> 8)
	//binary.LittleEndian.PutUint16(b[pos:pos+2],o.SrcPort)
	pos += 2

	b[pos] = byte(dstIpSize)
	pos++
	copy(b[pos:pos+srcIpSize], o.SrcIp)
	pos += dstIpSize

	b[pos] = byte(o.DstPort)
	b[pos+1] = byte(o.DstPort >> 8)
	//binary.LittleEndian.PutUint16(b[pos:pos+2],o.DstPort)
	pos += 2

	b[pos] = byte(o.Seq)
	b[pos+1] = byte(o.Seq >> 8)
	b[pos+2] = byte(o.Seq >> 16)
	b[pos+3] = byte(o.Seq >> 24)
	//binary.LittleEndian.PutUint32(b[pos:pos+4],o.Seq)
	pos += 4

	b[pos] = byte(o.Ack)
	b[pos+1] = byte(o.Ack >> 8)
	b[pos+2] = byte(o.Ack >> 16)
	b[pos+3] = byte(o.Ack >> 24)
	//binary.LittleEndian.PutUint32(b[pos:pos+4],o.Ack)
	pos += 4

	return b, nil
}

func MarshalTcpPacketV2ToBuffer(o *TcpPacket, b []byte) (err error) {
	//timeb,err:=o.Time.MarshalBinary() //13.73%
	//if err!=nil{
	//	return err
	//}
	//timeBSize:=len(timeb)
	srcIpSize := len(o.SrcIp)
	dstIpSize := len(o.DstIp)
	//size:=1+timeBSize+1+srcIpSize+2+1+ dstIpSize +2+4+4
	size := 1 + srcIpSize + 2 + 1 + dstIpSize + 2 + 4 + 4
	if len(b) < size {
		return fmt.Errorf("[MarshalTcpPacketV2ToBuffer] too small buffer")
	}
	pos := 0
	//b[pos] = byte(timeBSize)
	//pos++
	//copy(b[pos:pos+timeBSize],timeb)
	//pos+=timeBSize

	b[pos] = byte(srcIpSize)
	pos++
	copy(b[pos:pos+srcIpSize], o.SrcIp)
	pos += srcIpSize

	b[pos] = byte(o.SrcPort)
	b[pos+1] = byte(o.SrcPort >> 8)
	//binary.LittleEndian.PutUint16(b[pos:pos+2],o.SrcPort)
	pos += 2

	b[pos] = byte(dstIpSize)
	pos++
	copy(b[pos:pos+srcIpSize], o.SrcIp)
	pos += dstIpSize

	b[pos] = byte(o.DstPort)
	b[pos+1] = byte(o.DstPort >> 8)
	//binary.LittleEndian.PutUint16(b[pos:pos+2],o.DstPort)
	pos += 2

	b[pos] = byte(o.Seq)
	b[pos+1] = byte(o.Seq >> 8)
	b[pos+2] = byte(o.Seq >> 16)
	b[pos+3] = byte(o.Seq >> 24)
	//binary.LittleEndian.PutUint32(b[pos:pos+4],o.Seq)
	pos += 4

	b[pos] = byte(o.Ack)
	b[pos+1] = byte(o.Ack >> 8)
	b[pos+2] = byte(o.Ack >> 16)
	b[pos+3] = byte(o.Ack >> 24)
	//binary.LittleEndian.PutUint32(b[pos:pos+4],o.Ack)
	pos += 4

	return nil
}

func UnmarshalTcpPacketV2(b []byte) (packet TcpPacket, err error) {
	packet = TcpPacket{} // 13.73%
	pos := 0
	//timeSize:=b[pos]
	//pos++
	//err = packet.Time.UnmarshalBinary(b[pos:pos+int(timeSize)]) //4.39%
	//if err!=nil{
	//	return packet,err
	//}
	//pos+=int(timeSize)

	sSize := b[pos]
	packet.SrcIp = BytesToString(b[pos : pos+int(sSize)]) //12.48%
	pos += int(sSize)

	packet.SrcPort = binary.LittleEndian.Uint16(b[pos : pos+2])
	pos += 2

	sSize = b[pos]
	packet.DstIp = BytesToString(b[pos : pos+int(sSize)])
	pos += int(sSize)

	packet.DstPort = binary.LittleEndian.Uint16(b[pos : pos+2])
	pos += 2

	packet.Seq = binary.LittleEndian.Uint32(b[pos : pos+4])
	pos += 4

	packet.Ack = binary.LittleEndian.Uint32(b[pos : pos+4])
	pos += 4
	return packet, nil
}

func BytesToString(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{sh.Data, sh.Len, 0}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func BenchmarkKsonV2Marshal(b *testing.B) {
	data := &TcpPacket{
		Time:  time.Now().In(time.UTC),
		SrcIp: "1.2.3.4",
		DstIp: "1.2.3.4",
	}
	total := 0
	buf := make([]byte, 4096)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := MarshalTcpPacketV2ToBuffer(data, buf)
		if err != nil {
			panic(err)
		}
		//total += len(buf)
	}
	b.SetBytes(int64(total / b.N))
}

func BenchmarkKsonV2Unmarshal(b *testing.B) {
	data := &TcpPacket{
		Time:  time.Now().In(time.UTC),
		SrcIp: "1.2.3.4",
		DstIp: "1.2.3.4",
	}
	buf, err := MarshalTcpPacketV2(data)
	if err != nil {
		panic(err)
	}
	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := UnmarshalTcpPacketV2(buf)
		if err != nil {
			panic(err)
		}
	}
}
