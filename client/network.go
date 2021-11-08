package client

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/pkg/errors"

	"github.com/Mrs4s/MiraiGo/internal/packets"
	"github.com/Mrs4s/MiraiGo/utils"
)

// ConnectionQualityInfo 客户端连接质量测试结果
// 延迟单位为 ms 如为 9999 则测试失败 测试方法为 TCP 连接测试
// 丢包测试方法为 ICMP. 总共发送 10 个包, 记录丢包数
type ConnectionQualityInfo struct {
	// ChatServerLatency 聊天服务器延迟
	ChatServerLatency int64
	// ChatServerPacketLoss 聊天服务器ICMP丢包数
	ChatServerPacketLoss int
	// LongMessageServerLatency 长消息服务器延迟. 涉及长消息以及合并转发消息下载
	LongMessageServerLatency int64
	// LongMessageServerResponseLatency 长消息服务器返回延迟
	LongMessageServerResponseLatency int64
	// SrvServerLatency Highway服务器延迟. 涉及媒体以及群文件上传
	SrvServerLatency int64
	// SrvServerPacketLoss Highway服务器ICMP丢包数.
	SrvServerPacketLoss int
}

var ErrNotConnected = errors.New("no active connection")

func (c *QQClient) ConnectionQualityTest() *ConnectionQualityInfo {
	if !c.Online {
		return nil
	}
	r := &ConnectionQualityInfo{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func(w *sync.WaitGroup) {
		var err error

		if r.ChatServerLatency, err = qualityTest(c.servers[c.currServerIndex].String()); err != nil {
			c.Error("test chat server latency error: %v", err)
			r.ChatServerLatency = 9999
		}

		if addr, err := net.ResolveIPAddr("ip", "ssl.htdata.qq.com"); err == nil {
			if r.LongMessageServerLatency, err = qualityTest((&net.TCPAddr{IP: addr.IP, Port: 443}).String()); err != nil {
				c.Error("test long message server latency error: %v", err)
				r.LongMessageServerLatency = 9999
			}
		} else {
			c.Error("resolve long message server error: %v", err)
			r.LongMessageServerLatency = 9999
		}
		if len(c.srvSsoAddrs) > 0 {
			if r.SrvServerLatency, err = qualityTest(c.srvSsoAddrs[0]); err != nil {
				c.Error("test srv server latency error: %v", err)
				r.SrvServerLatency = 9999
			}
		}

		w.Done()
	}(&wg)
	go func(w *sync.WaitGroup) {
		res := utils.RunICMPPingLoop(&net.IPAddr{IP: c.servers[c.currServerIndex].IP}, 10)
		r.ChatServerPacketLoss = res.PacketsLoss
		if len(c.srvSsoAddrs) > 0 {
			res = utils.RunICMPPingLoop(&net.IPAddr{IP: net.ParseIP(strings.Split(c.srvSsoAddrs[0], ":")[0])}, 10)
			r.SrvServerPacketLoss = res.PacketsLoss
		}
		w.Done()
	}(&wg)
	start := time.Now()
	if _, err := utils.HttpGetBytes("https://ssl.htdata.qq.com", ""); err == nil {
		r.LongMessageServerResponseLatency = time.Now().Sub(start).Milliseconds()
	} else {
		c.Error("test long message server response latency error: %v", err)
		r.LongMessageServerResponseLatency = 9999
	}
	wg.Wait()
	return r
}

func (c *QQClient) getConn() *net.TCPConn {
	return (*net.TCPConn)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&c.conn))))
}

func (c *QQClient) setConn(conn *net.TCPConn) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&c.conn)), unsafe.Pointer(conn))
}

func (c *QQClient) cleanConn() *net.TCPConn {
	return (*net.TCPConn)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.conn)), unsafe.Pointer(nil)))
}

// connectFastest 连接到最快的服务器
// TODO 禁用不可用服务器
func (c *QQClient) connectFastest() error {
	c.Debug("connectFastest")
	// 清理存在的解码句柄
	c.handlers = HandlerMap{}

	ch := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(len(c.servers))
	for _, remote := range c.servers {
		go func(remote *net.TCPAddr) {
			defer wg.Done()
			conn, err := net.DialTCP("tcp", nil, remote)
			if err != nil {
				return
			}
			//addrs = append(addrs, remote)
			if c.getConn() != nil {
				_ = conn.Close()
				return
			}
			c.setConn(conn)
			ch <- nil
		}(remote)
	}
	go func() {
		wg.Wait()
		if c.getConn() == nil {
			ch <- errors.New("All servers are unreachable")
		}
	}()
	err := <-ch
	//c.currServerIndex = 0
	if err != nil {
		return err
	}
	conn := c.getConn()
	c.Info("connected to server: %v [fastest]", conn.RemoteAddr().String())
	go c.netLoop(conn)
	c.retryTimes = 0
	c.ConnectTime = time.Now()
	return nil
}

// connect 连接到 QQClient.servers 中的服务器
func (c *QQClient) connect() error {
	return c.connectFastest() // 暂时
	/*c.Info("connect to server: %v", c.servers[c.currServerIndex].String())
	err := c.TCP.Connect(c.servers[c.currServerIndex])
	c.currServerIndex++
	if c.currServerIndex == len(c.servers) {
		c.currServerIndex = 0
	}
	if err != nil {
		c.retryTimes++
		if c.retryTimes > len(c.servers) {
			return errors.New("All servers are unreachable")
		}
		c.Error("connect server error: %v", err)
		return err
	}
	if !c.netAlive {
		c.alive = true
		go c.netLoop()
	}
	c.retryTimes = 0
	c.ConnectTime = time.Now()
	return nil*/
}

// quickReconnect 快速重连
func (c *QQClient) quickReconnect() {
	c.Disconnect()
	time.Sleep(time.Millisecond * 200)
	if err := c.connect(); err != nil {
		c.Error("connect server error: %v", err)
		c.EventHandler.DisconnectHandler(c, &ClientDisconnectedEvent{Message: "quick reconnect failed"})
		return
	}
	if err := c.registerClient(); err != nil {
		c.Error("register client failed: %v", err)
		c.Disconnect()
		c.EventHandler.DisconnectHandler(c, &ClientDisconnectedEvent{Message: "register error"})
		return
	}
}

// Disconnect 中断连接
func (c *QQClient) Disconnect() {
	if conn := c.cleanConn(); conn != nil {
		_ = conn.Close()
	}
	c.Online = false
}

// sendAndWait 向服务器发送一个数据包, 并等待返回
func (c *QQClient) sendAndWait(seq uint16, pkt []byte, params ...requestParams) (interface{}, error) {
	c.Debug("send seq:%v", seq)
	// 整个sendAndWait使用同一个connection防止串线
	conn := c.getConn()

	type T struct {
		Response interface{}
		Error    error
	}
	ch := make(chan T, 1)
	var p requestParams

	if len(params) != 0 {
		p = params[0]
	}

	c.handlers.Store(seq, &handlerInfo{fun: func(i interface{}, err error) {
		ch <- T{
			Response: i,
			Error:    err,
		}
	}, params: p, dynamic: false})

	for retry := 0; retry < 3; retry++ {
		err := c.sendPacketWithConn(conn, pkt)
		if err != nil {
			c.handlers.Delete(seq)
			return nil, err
		}
		select {
		case rsp := <-ch:
			return rsp.Response, rsp.Error
		case <-time.After(time.Second * 5):
			c.handlers.Delete(seq)
		}
	}
	return nil, errors.New("Packet timed out")
}

// sendPacket 向服务器发送一个数据包
func (c *QQClient) sendPacket(pkt []byte) error {
	return c.sendPacketWithConn(c.getConn(), pkt)
}

func (c *QQClient) sendPacketWithConn(conn *net.TCPConn, pkt []byte) error {
	if conn == nil {
		return ErrNotConnected
	}
	n, err := conn.Write(pkt)
	if n != len(pkt) {
		err = io.ErrShortWrite
	}
	if err != nil {
		atomic.AddUint64(&c.stat.PacketLost, 1)
		return errors.Wrap(err, "Packet failed to sendPacket")
	}
	atomic.AddUint64(&c.stat.PacketSent, 1)
	return nil

}

// waitPacket
// 等待一个或多个数据包解析, 优先级低于 sendAndWait
// 返回终止解析函数
func (c *QQClient) waitPacket(cmd string, f func(interface{}, error)) func() {
	c.waiters.Store(cmd, f)
	return func() {
		c.waiters.Delete(cmd)
	}
}

// sendAndWaitDynamic
// 发送数据包并返回需要解析的 response
func (c *QQClient) sendAndWaitDynamic(seq uint16, pkt []byte) ([]byte, error) {
	ch := make(chan []byte, 1)
	c.handlers.Store(seq, &handlerInfo{fun: func(i interface{}, err error) { ch <- i.([]byte) }, dynamic: true})
	err := c.sendPacket(pkt)
	if err != nil {
		c.handlers.Delete(seq)
		return nil, err
	}
	select {
	case rsp := <-ch:
		return rsp, nil
	case <-time.After(time.Second * 15):
		c.handlers.Delete(seq)
		return nil, errors.New("Packet timed out")
	}
}

// plannedDisconnect 计划中断线事件
//func (c *QQClient) plannedDisconnect(_ *utils.TCPDialer) {
//	c.Debug("planned disconnect.")
//	atomic.AddUint32(&c.stat.DisconnectTimes, 1)
//	c.Online = false
//}

// unexpectedDisconnect 非预期断线事件
func (c *QQClient) unexpectedDisconnect(e error) {
	c.Error("unexpected disconnect: %v", e)
	atomic.AddUint32(&c.stat.DisconnectTimes, 1)
	c.Online = false
	if err := c.connect(); err != nil {
		c.Error("connect server error: %v", err)
		c.EventHandler.DisconnectHandler(c, &ClientDisconnectedEvent{Message: "connection dropped by server."})
		return
	}
	if err := c.registerClient(); err != nil {
		c.Error("register client failed: %v", err)
		c.Disconnect()
		c.EventHandler.DisconnectHandler(c, &ClientDisconnectedEvent{Message: "register error"})
		return
	}
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
	buf := make([]byte, l-4)
	_, err = io.ReadFull(conn, buf)
	return buf, err
}

// netLoop 通过循环来不停接收数据包
func (c *QQClient) netLoop(conn *net.TCPConn) {
	errCount := 0
	for {
		data, err := readPacket(conn, 4, 10<<20) // max 10MB
		if err != nil {
			// 连接未改变，没有建立新连接
			if c.getConn() == conn {
				c.unexpectedDisconnect(err)
			}
			break
		}
		pkt, err := packets.ParseIncomingPacket(data, c.sigInfo.d2Key)
		if err != nil {
			c.Error("parse incoming packet error: %v", err)
			if errors.Is(err, packets.ErrSessionExpired) || errors.Is(err, packets.ErrPacketDropped) {
				c.Disconnect()
				go c.EventHandler.DisconnectHandler(c, &ClientDisconnectedEvent{Message: "session expired"})
				break
			}
			errCount++
			if errCount > 2 {
				go c.quickReconnect()
			}
			break
		}
		if pkt.Flag2 == 2 {
			pkt.Payload, err = pkt.DecryptPayload(c.ecdh.InitialShareKey, c.RandomKey, c.sigInfo.wtSessionTicketKey)
			if err != nil {
				c.Error("decrypt payload error: %v", err)
				continue
			}
		}
		errCount = 0
		c.Debug("rev cmd: %v seq: %v", pkt.CommandName, pkt.SequenceId)
		atomic.AddUint64(&c.stat.PacketReceived, 1)
		go func(pkt *packets.IncomingPacket) {
			defer func() {
				if pan := recover(); pan != nil {
					c.Error("panic on decoder %v : %v\n%s", pkt.CommandName, pan, debug.Stack())
				}
			}()

			if decoder, ok := decoders[pkt.CommandName]; ok {
				// found predefined decoder
				info, ok := c.handlers.LoadAndDelete(pkt.SequenceId)
				var decoded interface{}
				decoded = pkt.Payload
				if info == nil || !info.dynamic {
					decoded, err = decoder(c, &incomingPacketInfo{
						SequenceId:  pkt.SequenceId,
						CommandName: pkt.CommandName,
						Params: func() requestParams {
							if !ok {
								return nil
							}
							return info.params
						}(),
					}, pkt.Payload)
					if err != nil {
						c.Debug("decode pkt %v error: %+v", pkt.CommandName, err)
					}
				}
				if ok {
					info.fun(decoded, err)
				} else if f, ok := c.waiters.Load(pkt.CommandName); ok { // 在不存在handler的情况下触发wait
					f.(func(interface{}, error))(decoded, err)
				}
			} else if f, ok := c.handlers.LoadAndDelete(pkt.SequenceId); ok {
				// does not need decoder
				f.fun(pkt.Payload, nil)
			} else {
				c.Debug("Unhandled Command: %s Seq: %d This message can be ignored.", pkt.CommandName, pkt.SequenceId)
			}
		}(pkt)
	}
	_ = conn.Close()
}
