package client

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/pkg/errors"

	"github.com/Mrs4s/MiraiGo/client/internal/network"
	"github.com/Mrs4s/MiraiGo/client/internal/oicq"
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
	if !c.Online.Load() {
		return nil
	}
	r := &ConnectionQualityInfo{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
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
		if c.highwaySession.AddrLength() > 0 {
			if r.SrvServerLatency, err = qualityTest(c.highwaySession.SsoAddr[0].String()); err != nil {
				c.Error("test srv server latency error: %v", err)
				r.SrvServerLatency = 9999
			}
		}
	}()
	go func() {
		defer wg.Done()
		res := utils.RunICMPPingLoop(&net.IPAddr{IP: c.servers[c.currServerIndex].IP}, 10)
		r.ChatServerPacketLoss = res.PacketsLoss
		if c.highwaySession.AddrLength() > 0 {
			res = utils.RunICMPPingLoop(&net.IPAddr{IP: c.highwaySession.SsoAddr[0].AsNetIP()}, 10)
			r.SrvServerPacketLoss = res.PacketsLoss
		}
	}()
	start := time.Now()
	if _, err := utils.HttpGetBytes("https://ssl.htdata.qq.com", ""); err == nil {
		r.LongMessageServerResponseLatency = time.Since(start).Milliseconds()
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

func (c *QQClient) setConn(conn *net.TCPConn) (swapped bool) {
	return atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.conn)), unsafe.Pointer(nil), unsafe.Pointer(conn))
}

func (c *QQClient) closeConn() *net.TCPConn {
	return (*net.TCPConn)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.conn)), unsafe.Pointer(nil)))
}

// connectFastest 连接到最快的服务器
// TODO 禁用不可用服务器
func (c *QQClient) connectFastest() error {
	c.Debug("connectFastest")
	// 清理存在的解码句柄
	c.handlers = HandlerMap{}
	c.Disconnect()
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
			if !c.setConn(conn) {
				_ = conn.Close()
				return

			}
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
	c.Debug("connected to server: %v [fastest]", conn.RemoteAddr().String())
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

func (c *QQClient) QuickReconnect() {
	c.quickReconnect("用户请求快速重连")
}

// quickReconnect 快速重连
func (c *QQClient) quickReconnect(message string) {
	c.Disconnect()
	go c.EventHandler.DisconnectHandler(c, &ClientDisconnectedEvent{Message: message})
	time.Sleep(time.Millisecond * 200)
	if err := c.connect(); err != nil {
		c.Error("connect server error: %v", err)
		c.EventHandler.OfflineHandler(c, &ClientOfflineEvent{Message: "快速重连失败"})
		return
	}
	if err := c.registerClient(); err != nil {
		c.Error("register client failed: %v", err)
		c.Disconnect()
		c.EventHandler.OfflineHandler(c, &ClientOfflineEvent{Message: "register error"})
		return
	}
}

// Disconnect 中断连接
func (c *QQClient) Disconnect() {
	c.Online.Store(false)
	if conn := c.closeConn(); conn != nil {
		_ = conn.Close()
	}
}

// sendAndWait 向服务器发送一个数据包, 并等待返回
func (c *QQClient) sendAndWait(seq uint16, pkt []byte, params ...network.RequestParams) (interface{}, error) {
	// 整个sendAndWait使用同一个connection防止串线
	conn := c.getConn()
	type T struct {
		Response interface{}
		Error    error
	}
	ch := make(chan T, 1)
	var p network.RequestParams

	if len(params) != 0 {
		p = params[0]
	}

	c.handlers.Store(seq, &handlerInfo{fun: func(i interface{}, err error) {
		ch <- T{
			Response: i,
			Error:    err,
		}
	}, params: p, dynamic: false})

	for retry := 0; retry < 2; retry++ {
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
	if conn != c.getConn() {
		conn.Close()
		return errors.New("broken connection")
	}
	if conn == nil {
		return ErrNotConnected
	}
	n, err := conn.Write(pkt)
	if n != len(pkt) {
		err = io.ErrShortWrite
	}
	if err != nil {
		c.stat.PacketLost.Add(1)
		return errors.Wrap(err, "Packet failed to sendPacket")
	}
	c.stat.PacketSent.Add(1)
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

// waitPacketTimeoutSyncF
// 等待一个数据包解析, 优先级低于 sendAndWait
func (c *QQClient) waitPacketTimeoutSyncF(cmd string, timeout time.Duration, filter func(interface{}) bool) (r interface{}, e error) {
	notifyChan := make(chan bool)
	defer c.waitPacket(cmd, func(i interface{}, err error) {
		if filter(i) {
			r = i
			e = err
			notifyChan <- true
		}
	})()
	select {
	case <-notifyChan:
		return
	case <-time.After(timeout):
		return nil, errors.New("timeout")
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
//func (c *QQClient) plannedDisconnect(_ *network.TCPListener) {
//	c.Debug("planned disconnect.")
//	c.stat.DisconnectTimes.Add(1)
//	c.Online.Store(false)
//}

// unexpectedDisconnect 非预期断线事件
func (c *QQClient) unexpectedDisconnect(e error) {
	c.Error("unexpected disconnect: %v", e)
	c.stat.DisconnectTimes.Add(1)
	c.Online.Store(false)
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

//func (c *QQClient) procPacket()

// netLoop 通过循环来不停接收数据包
func (c *QQClient) netLoop(conn *net.TCPConn) {
	defer func() {
		if r := recover(); r != nil {
			c.Error("netLoop %+v", r)
		}
		_ = conn.Close()
	}()
	errCount := 0
	for {
		data, err := readPacket(conn, 4, 10<<20) // max 10MB
		if err != nil {
			// 连接未改变，没有建立新连接
			if c.getConn() == conn {
				c.unexpectedDisconnect(err)
			}
			return
		}
		resp, err := c.transport.ReadResponse(data)
		pkt, err := packets.ParseIncomingPacket(data, c.sig.D2Key)
		if err != nil {
			if errors.Is(err, packets.ErrSessionExpired) || errors.Is(err, packets.ErrPacketDropped) {
				c.Disconnect()
				go c.EventHandler.DisconnectHandler(c, &ClientDisconnectedEvent{Message: "session expired"})
				return // 销毁连接
			}
			c.Error("parse incoming packet error: %v", err)
			errCount++
			if errCount > 2 {
				c.quickReconnect("链接错误率过高")
			}
			continue
		}
		if resp.EncryptType == network.EncryptTypeEmptyKey {
			m, err := c.oicq.Unmarshal(resp.Body)
			if err != nil {
				c.Error("decrypt payload error: %v", err)
				if errors.Is(err, packets.ErrUnknownFlag) {
					c.quickReconnect("服务器发送未知响应")
				}
				continue
			}
			resp.Body = m.Body
		}
		errCount = 0
		c.Debug("rev cmd: %v seq: %v", pkt.CommandName, pkt.SequenceId)
		c.stat.PacketReceived.Add(1) // 不再需要atomic
		pkt := &packets.IncomingPacket{
			SequenceId:  uint16(resp.SequenceID),
			CommandName: resp.CommandName,
			Payload:     resp.Body,
		}
		go func(pkt *packets.IncomingPacket) {
			defer func() {
				if pan := recover(); pan != nil {
					c.Error("panic on decoder %v : %v\n%s", pkt.CommandName, pan, debug.Stack())
					c.Dump("packet decode error: %v - %v", pkt.Payload, pkt.CommandName, pan)
				}
			}()

			if decoder, ok := decoders[pkt.CommandName]; ok {
				// found predefined decoder
				info, ok := c.handlers.LoadAndDelete(pkt.SequenceId)
				var decoded interface{}
				decoded = pkt.Payload
				if info == nil || !info.dynamic {
					decoded, err = decoder(c, &network.IncomingPacketInfo{
						SequenceId:  pkt.SequenceId,
						CommandName: pkt.CommandName,
						Params:      info.getParams(),
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
}
