package network

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"unsafe"
)

type PktHandler func(pkt *Response, netErr error)
type ResponseHandler func(head []byte) (*Response, error)

func (t *Transport) GetConn() *net.TCPConn {
	return (*net.TCPConn)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&t.conn))))
}

func (t *Transport) setConn(conn *net.TCPConn) (swapped bool) {
	return atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&t.conn)), unsafe.Pointer(nil), unsafe.Pointer(conn))
}

func (t *Transport) closeConn() *net.TCPConn {
	return (*net.TCPConn)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&t.conn)), unsafe.Pointer(nil)))
}

// ConnectFastest 连接到最快的服务器
// TODO 禁用不可用服务器
func (t *Transport) ConnectFastest(servers []*net.TCPAddr) (net.Addr, error) {
	ch := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(len(servers))
	for _, remote := range servers {
		go func(remote *net.TCPAddr) {
			defer wg.Done()
			conn, err := net.DialTCP("tcp", nil, remote)
			if err != nil {
				return
			}
			//addrs = append(addrs, remote)
			if !t.setConn(conn) {
				_ = conn.Close()
				return
			}
			ch <- nil
		}(remote)
	}
	go func() {
		wg.Wait()
		if t.GetConn() == nil {
			ch <- errors.New("All servers are unreachable")
		}
	}()
	err := <-ch
	if err != nil {
		return nil, err
	}
	conn := t.GetConn()
	return conn.RemoteAddr(), nil
}

func readPacket(conn *net.TCPConn, minSize, maxSize uint32) ([]byte, error) {
	lBuf := make([]byte, 4)
	_, err := io.ReadFull(conn, lBuf)
	if err != nil {
		return nil, err
	}
	l := binary.BigEndian.Uint32(lBuf)
	if l < minSize || l > maxSize {
		return nil, fmt.Errorf("parse incoming packet error: invalid packet length %v", l)
	}
	data := make([]byte, l-4)
	_, err = io.ReadFull(conn, data)
	return data, err
}

func (t *Transport) NetLoop(pktHandler PktHandler, respHandler ResponseHandler) {
	go t.netLoop(t.GetConn(), pktHandler, respHandler)
}

func (t *Transport) netLoop(conn *net.TCPConn, pktHandler PktHandler, respHandler ResponseHandler) {
	defer func() {
		if r := recover(); r != nil {
			pktHandler(nil, fmt.Errorf("panic: %v", r))
		}
		_ = conn.Close()
	}()
	errCount := 0
	for {
		data, err := readPacket(conn, 4, 10<<20) // max 10MB
		if err != nil {
			// 连接未改变，没有建立新连接
			if t.GetConn() == conn {
				pktHandler(nil, errors.Wrap(ErrConnectionBroken, err.Error()))
			}
			return
		}
		resp, err := respHandler(data)
		if err == nil {
			errCount = 0
			goto ok
		}
		errCount++
		if errCount > 2 {
			err = errors.Wrap(ErrConnectionBroken, err.Error())
		}
	ok:
		go pktHandler(resp, err)
	}
}
