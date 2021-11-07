package utils

import (
	"encoding/binary"
	"io"
	"net"
	"sync"

	"github.com/pkg/errors"
)

type TCPDialer struct {
	lock                 sync.RWMutex
	conn                 net.Conn
	plannedDisconnect    func(*TCPDialer)
	unexpectedDisconnect func(*TCPDialer, error)
}

var ErrConnectionClosed = errors.New("connection closed")

// PlannedDisconnect 预料中的断开连接
// 如调用 Close() Connect()
func (t *TCPDialer) PlannedDisconnect(f func(*TCPDialer)) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.plannedDisconnect = f
}

// UnexpectedDisconnect 未预料的断开连接
func (t *TCPDialer) UnexpectedDisconnect(f func(*TCPDialer, error)) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.unexpectedDisconnect = f
}

func (t *TCPDialer) Connect(addr *net.TCPAddr) error {
	t.Close()
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return errors.Wrap(err, "dial tcp error")
	}
	t.lock.Lock()
	defer t.lock.Unlock()
	t.conn = conn
	return nil
}

func (t *TCPDialer) Write(buf []byte) error {
	if conn := t.getConn(); conn != nil {
		_, err := conn.Write(buf)
		if err != nil {
			t.unexpectedClose(err)
			return ErrConnectionClosed
		}
		return nil
	}

	return ErrConnectionClosed
}

func (t *TCPDialer) ReadBytes(len int) ([]byte, error) {
	buf := make([]byte, len)
	if conn := t.getConn(); conn != nil {
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			// time.Sleep(time.Millisecond * 100) // 服务器会发送offline包后立即断开连接, 此时还没解析, 可能还是得加锁
			t.unexpectedClose(err)
			return nil, ErrConnectionClosed
		}
		return buf, nil
	}

	return nil, ErrConnectionClosed
}

func (t *TCPDialer) ReadInt32() (int32, error) {
	b, err := t.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	return int32(binary.BigEndian.Uint32(b)), nil
}

func (t *TCPDialer) Close() {
	t.close()
	t.invokePlannedDisconnect()
}

func (t *TCPDialer) unexpectedClose(err error) {
	t.close()
	t.invokeUnexpectedDisconnect(err)
}

func (t *TCPDialer) close() {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.conn != nil {
		_ = t.conn.Close()
		t.conn = nil
	}
}

func (t *TCPDialer) invokePlannedDisconnect() {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if t.plannedDisconnect != nil {
		go t.plannedDisconnect(t)
	}
}

func (t *TCPDialer) invokeUnexpectedDisconnect(err error) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if t.unexpectedDisconnect != nil {
		go t.unexpectedDisconnect(t, err)
	}
}

func (t *TCPDialer) getConn() net.Conn {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.conn
}
