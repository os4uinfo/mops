package msg

import (
    "net"
    "errors"
    "fmt"
    "time"
)

const (
	headerByteCount = 4
)

func read(c net.Conn, bc int) ([]byte, error) {
	bs := make([]byte, bc)
	ix := 0
	for ix < bc {
		n, err := c.Read(bs[ix:])
		if err != nil {
			return nil, err
		}
		ix += n
	}
	return bs, nil
}

func ReadInt(c net.Conn, bc int) (int, error) {
	if bc > 4 {
		return -1, errors.New("ReadInt byteCount too large!")
	}

	bs, err := read(c, bc)
	if err != nil {
		return -1, err
	}

	r := 0
	for i, j := bc-1, 0; i >= 0; i-- {
		r |= int(bs[j]) << uint(i*8)
		j += 1
	}

	return r, nil
}


func ReadBytes(c net.Conn) ([]byte, error) {
	bc, err := ReadInt(c, headerByteCount)
	if err != nil {
		return nil, err
	}

	bs, err := read(c, bc)
	if err != nil {
		return nil, err
	}

	return bs, nil
}


func ReadString(c net.Conn) (string, error) {
	data, err := ReadBytes(c)
	if err != nil {
		return "", err
	}
	return string(data), err
}



func WriteBool(c net.Conn, b bool) error {
	if b {
		return write(c, []byte{1})
	} else {
		return write(c, []byte{0})
	}
}

func write(c net.Conn, bs []byte) error {
	l := len(bs)
	i := 0
	for {
		n, err := c.Write(bs[i:])
		if err != nil {
			return err
		}
		i += n
		if i >= l {
			break
		} else {
			fmt.Println("net.Conn.Write not write all data, tcp sendBuf overflow?", i, l)
			time.Sleep(time.Second)
		}
	}
	return nil
}

func WriteString(c net.Conn, s string) error {
	return WriteBytes(c, []byte(s))
}

func WriteBytes(c net.Conn, data []byte) error {
	// err := WriteInt(c, headerByteCount, len(data)) // this is bug implement
	// if err != nil {
	// 	return err
	// }

	// return write(c, data)

	bc, v := headerByteCount, len(data)

	bs := make([]byte, bc)
	for i, j := bc-1, 0; i >= 0; i-- {
		bs[j] = byte(v >> uint(i*8))
		j += 1
	}

	return write(c, append(bs, data...))
}

func ReadBool(c net.Conn) (bool, error) {
	r, err := read(c, 1)
	if err != nil {
		return false, err
	}
	return r[0] == 1, nil
}
