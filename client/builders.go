package client

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"

	"github.com/Mrs4s/MiraiGo/binary"
	"github.com/Mrs4s/MiraiGo/binary/jce"
	"github.com/Mrs4s/MiraiGo/client/internal/auth"
	"github.com/Mrs4s/MiraiGo/client/internal/network"
	"github.com/Mrs4s/MiraiGo/client/internal/oicq"
	"github.com/Mrs4s/MiraiGo/client/pb"
	"github.com/Mrs4s/MiraiGo/client/pb/cmd0x352"
	"github.com/Mrs4s/MiraiGo/client/pb/msg"
	"github.com/Mrs4s/MiraiGo/client/pb/oidb"
	"github.com/Mrs4s/MiraiGo/client/pb/profilecard"
	"github.com/Mrs4s/MiraiGo/client/pb/structmsg"
	"github.com/Mrs4s/MiraiGo/internal/proto"
	"github.com/Mrs4s/MiraiGo/internal/tlv"
)

var (
	syncConst1 = rand.Int63()
	syncConst2 = rand.Int63()
)

func buildCode2DRequestPacket(seq uint32, j uint64, cmd uint16, bodyFunc func(writer *binary.Writer)) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteByte(2)
		pos := w.FillUInt16()
		w.WriteUInt16(cmd)
		w.Write(make([]byte, 21))
		w.WriteByte(3)
		w.WriteUInt16(0)
		w.WriteUInt16(50) // version
		w.WriteUInt32(seq)
		w.WriteUInt64(j)
		bodyFunc(w)
		w.WriteByte(3)
		w.WriteUInt16At(pos, uint16(w.Len()))
	})
}

func (c *QQClient) buildLoginRequest() *network.Request {
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(c.Uin, 0x0810, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(9)
		if c.AllowSlider {
			w.WriteUInt16(0x17)
		} else {
			w.WriteUInt16(0x16)
		}

		w.Write(tlv.T18(16, uint32(c.Uin)))
		w.Write(tlv.T1(uint32(c.Uin), c.deviceInfo.IpAddress))
		w.Write(tlv.T106(uint32(c.Uin), 0, c.version.AppId, c.version.SSOVersion, c.PasswordMd5, true, c.deviceInfo.Guid, c.deviceInfo.TgtgtKey, 0))
		w.Write(tlv.T116(c.version.MiscBitmap, c.version.SubSigmap))
		w.Write(tlv.T100(c.version.SSOVersion, c.version.SubAppId, c.version.MainSigMap))
		w.Write(tlv.T107(0))
		w.Write(tlv.T142(c.version.ApkId))
		w.Write(tlv.T144(
			[]byte(c.deviceInfo.IMEI),
			c.deviceInfo.GenDeviceInfoData(),
			c.deviceInfo.OSType,
			c.deviceInfo.Version.Release,
			c.deviceInfo.SimInfo,
			c.deviceInfo.APN,
			false, true, false, tlv.GuidFlag(),
			c.deviceInfo.Model,
			c.deviceInfo.Guid,
			c.deviceInfo.Brand,
			c.deviceInfo.TgtgtKey,
		))

		w.Write(tlv.T145(c.deviceInfo.Guid))
		w.Write(tlv.T147(16, []byte(c.version.SortVersionName), c.version.ApkSign))
		/*
			if (miscBitMap & 0x80) != 0{
				w.Write(tlv.T166(1))
			}
		*/
		w.Write(tlv.T154(seq))
		w.Write(tlv.T141(c.deviceInfo.SimInfo, c.deviceInfo.APN))
		w.Write(tlv.T8(2052))
		w.Write(tlv.T511([]string{
			"tenpay.com", "openmobile.qq.com", "docs.qq.com", "connect.qq.com",
			"qzone.qq.com", "vip.qq.com", "gamecenter.qq.com", "qun.qq.com", "game.qq.com",
			"qqweb.qq.com", "office.qq.com", "ti.qq.com", "mail.qq.com", "mma.qq.com",
		}))

		w.Write(tlv.T187(c.deviceInfo.MacAddress))
		w.Write(tlv.T188(c.deviceInfo.AndroidId))
		if len(c.deviceInfo.IMSIMd5) != 0 {
			w.Write(tlv.T194(c.deviceInfo.IMSIMd5))
		}
		if c.AllowSlider {
			w.Write(tlv.T191(0x82))
		}
		if len(c.deviceInfo.WifiBSSID) != 0 && len(c.deviceInfo.WifiSSID) != 0 {
			w.Write(tlv.T202(c.deviceInfo.WifiBSSID, c.deviceInfo.WifiSSID))
		}
		w.Write(tlv.T177(c.version.BuildTime, c.version.SdkVersion))
		w.Write(tlv.T516())
		w.Write(tlv.T521(0))
		w.Write(tlv.T525(tlv.T536([]byte{0x01, 0x00})))
	}))

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         c.Uin,
		CommandName: "wtlogin.login",
		Body:        req,
		Decode:      bindDecoder(c, decodeLoginResponse),
	}
}

func (c *QQClient) buildDeviceLockLoginRequest() *network.Request {
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(c.Uin, 0x0810, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(20)
		w.WriteUInt16(4)

		w.Write(tlv.T8(2052))
		w.Write(tlv.T104(c.sig.T104))
		w.Write(tlv.T116(c.version.MiscBitmap, c.version.SubSigmap))
		w.Write(tlv.T401(c.sig.G))
	}))
	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         c.Uin,
		CommandName: "wtlogin.login",
		Body:        req,
		Decode:      bindDecoder(c, decodeLoginResponse),
	}
}

func (c *QQClient) buildQRCodeFetchRequest(size, margin, ecLevel uint32) *network.Request {
	watch := auth.AndroidWatch.Version()
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(0, 0x0812, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteHex(`0001110000001000000072000000`) // trans header
		w.WriteUInt32(uint32(time.Now().Unix()))
		w.Write(buildCode2DRequestPacket(0, 0, 0x31, func(w *binary.Writer) {
			w.WriteUInt16(0)  // const
			w.WriteUInt32(16) // app id
			w.WriteUInt64(0)  // const
			w.WriteByte(8)    // const
			w.WriteBytesShort(EmptyBytes)

			w.WriteUInt16(6)
			w.Write(tlv.T16(watch.SSOVersion, 16, watch.AppId, c.deviceInfo.Guid, []byte(watch.ApkId), []byte(watch.SortVersionName), watch.ApkSign))
			w.Write(tlv.T1B(0, 0, size, margin, 72, ecLevel, 2))
			w.Write(tlv.T1D(watch.MiscBitmap))
			w.Write(tlv.T1F(false, c.deviceInfo.OSType, []byte("7.1.2"), []byte("China Mobile GSM"), c.deviceInfo.APN, 2))
			w.Write(tlv.T33(c.deviceInfo.Guid))
			w.Write(tlv.T35(8))
		}))
	}))

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         0,
		CommandName: "wtlogin.trans_emp",
		Body:        req,
		Decode:      bindDecoder(c, decodeTransEmpResponse),
	}
}

func (c *QQClient) buildQRCodeResultQueryRequest(sig []byte) *network.Request {
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(0, 0x0812, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteHex(`0000620000001000000072000000`) // trans header
		w.WriteUInt32(uint32(time.Now().Unix()))
		w.Write(buildCode2DRequestPacket(1, 0, 0x12, func(w *binary.Writer) {
			w.WriteUInt16(5)  // const
			w.WriteByte(1)    // const
			w.WriteUInt32(8)  // product type
			w.WriteUInt32(16) // app id
			w.WriteBytesShort(sig)
			w.WriteUInt64(0) // const
			w.WriteByte(8)   // const
			w.WriteBytesShort(EmptyBytes)
			w.WriteUInt16(0) // const
		}))
	}))

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         0,
		CommandName: "wtlogin.trans_emp",
		Body:        req,
		Decode:      bindDecoder(c, decodeTransEmpResponse),
	}
}

func (c *QQClient) buildQRCodeLoginRequest(t106, t16a, t318 []byte) *network.Request {
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(c.Uin, 0x0810, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(9)
		w.WriteUInt16(24)

		w.Write(tlv.T18(16, uint32(c.Uin)))
		w.Write(tlv.T1(uint32(c.Uin), c.deviceInfo.IpAddress))

		w.WriteUInt16(0x106)
		w.WriteBytesShort(t106)

		// w.Write(tlv.T106(uint32(c.Uin), 0, c.version.AppId, c.version.SSOVersion, c.PasswordMd5, true, c.deviceInfo.Guid, c.deviceInfo.TgtgtKey, 0))
		w.Write(tlv.T116(c.version.MiscBitmap, c.version.SubSigmap))
		w.Write(tlv.T100(c.version.SSOVersion, c.version.SubAppId, c.version.MainSigMap))
		w.Write(tlv.T107(0))
		w.Write(tlv.T142(c.version.ApkId))
		w.Write(tlv.T144(
			[]byte(c.deviceInfo.IMEI),
			c.deviceInfo.GenDeviceInfoData(),
			c.deviceInfo.OSType,
			c.deviceInfo.Version.Release,
			c.deviceInfo.SimInfo,
			c.deviceInfo.APN,
			false, true, false, tlv.GuidFlag(),
			c.deviceInfo.Model,
			c.deviceInfo.Guid,
			c.deviceInfo.Brand,
			c.deviceInfo.TgtgtKey,
		))

		w.Write(tlv.T145(c.deviceInfo.Guid))
		w.Write(tlv.T147(16, []byte(c.version.SortVersionName), c.version.ApkSign))

		w.WriteUInt16(0x16A)
		w.WriteBytesShort(t16a)

		w.Write(tlv.T154(seq))
		w.Write(tlv.T141(c.deviceInfo.SimInfo, c.deviceInfo.APN))
		w.Write(tlv.T8(2052))
		w.Write(tlv.T511([]string{
			"tenpay.com", "openmobile.qq.com", "docs.qq.com", "connect.qq.com",
			"qzone.qq.com", "vip.qq.com", "gamecenter.qq.com", "qun.qq.com", "game.qq.com",
			"qqweb.qq.com", "office.qq.com", "ti.qq.com", "mail.qq.com", "mma.qq.com",
		}))
		w.Write(tlv.T187(c.deviceInfo.MacAddress))
		w.Write(tlv.T188(c.deviceInfo.AndroidId))
		if len(c.deviceInfo.IMSIMd5) != 0 {
			w.Write(tlv.T194(c.deviceInfo.IMSIMd5))
		}
		w.Write(tlv.T191(0x00))
		if len(c.deviceInfo.WifiBSSID) != 0 && len(c.deviceInfo.WifiSSID) != 0 {
			w.Write(tlv.T202(c.deviceInfo.WifiBSSID, c.deviceInfo.WifiSSID))
		}
		w.Write(tlv.T177(c.version.BuildTime, c.version.SdkVersion))
		w.Write(tlv.T516())
		w.Write(tlv.T521(8))
		// w.Write(tlv.T525(tlv.T536([]byte{0x01, 0x00})))

		w.WriteUInt16(0x318)
		w.WriteBytesShort(t318)
	}))

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         c.Uin,
		CommandName: "wtlogin.login",
		Body:        req,
		Decode:      bindDecoder(c, decodeLoginResponse),
	}
}

func (c *QQClient) buildCaptchaRequest(result string, sign []byte) *network.Request {
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(c.Uin, 0x0810, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(2) // sub command
		w.WriteUInt16(4)

		w.Write(tlv.T2(result, sign))
		w.Write(tlv.T8(2052))
		w.Write(tlv.T104(c.sig.T104))
		w.Write(tlv.T116(c.version.MiscBitmap, c.version.SubSigmap))
	}))

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         c.Uin,
		CommandName: "wtlogin.login",
		Body:        req,
		Decode:      bindDecoder(c, decodeLoginResponse),
	}
}

func (c *QQClient) buildSMSRequest() *network.Request {
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(c.Uin, 0x0810, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(8)
		w.WriteUInt16(6)

		w.Write(tlv.T8(2052))
		w.Write(tlv.T104(c.sig.T104))
		w.Write(tlv.T116(c.version.MiscBitmap, c.version.SubSigmap))
		w.Write(tlv.T174(c.sig.T174))
		w.Write(tlv.T17A(9))
		w.Write(tlv.T197())
	}))

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         c.Uin,
		CommandName: "wtlogin.login",
		Body:        req,
		Decode:      bindDecoder(c, decodeLoginResponse),
	}
}

func (c *QQClient) buildSMSCodeSubmitRequest(code string) *network.Request {
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(c.Uin, 0x0810, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(7)
		w.WriteUInt16(7)

		w.Write(tlv.T8(2052))
		w.Write(tlv.T104(c.sig.T104))
		w.Write(tlv.T116(c.version.MiscBitmap, c.version.SubSigmap))
		w.Write(tlv.T174(c.sig.T174))
		w.Write(tlv.T17C(code))
		w.Write(tlv.T401(c.sig.G))
		w.Write(tlv.T198())
	}))

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         c.Uin,
		CommandName: "wtlogin.login",
		Body:        req,
		Decode:      bindDecoder(c, decodeLoginResponse),
	}
}

func (c *QQClient) buildTicketSubmitRequest(ticket string) *network.Request {
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(c.Uin, 0x0810, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(2)
		w.WriteUInt16(4)

		w.Write(tlv.T193(ticket))
		w.Write(tlv.T8(2052))
		w.Write(tlv.T104(c.sig.T104))
		w.Write(tlv.T116(c.version.MiscBitmap, c.version.SubSigmap))
	}))

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         c.Uin,
		CommandName: "wtlogin.login",
		Body:        req,
		Decode:      bindDecoder(c, decodeLoginResponse),
	}
}

func (c *QQClient) buildRequestTgtgtNopicsigRequest() *network.Request {
	seq := c.nextSeq()
	req := binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(15)
		w.WriteUInt16(24)

		w.Write(tlv.T18(16, uint32(c.Uin)))
		w.Write(tlv.T1(uint32(c.Uin), c.deviceInfo.IpAddress))

		w.WriteUInt16(0x106)
		w.WriteBytesShort(c.sig.EncryptedA1)

		w.Write(tlv.T116(c.version.MiscBitmap, c.version.SubSigmap))
		w.Write(tlv.T100(c.version.SSOVersion, 2, c.version.MainSigMap))
		w.Write(tlv.T107(0))
		w.Write(tlv.T108(c.sig.Ksid))
		w.Write(tlv.T144(
			c.deviceInfo.AndroidId,
			c.deviceInfo.GenDeviceInfoData(),
			c.deviceInfo.OSType,
			c.deviceInfo.Version.Release,
			c.deviceInfo.SimInfo,
			c.deviceInfo.APN,
			false, true, false, tlv.GuidFlag(),
			c.deviceInfo.Model,
			c.deviceInfo.Guid,
			c.deviceInfo.Brand,
			c.deviceInfo.TgtgtKey,
		))
		w.Write(tlv.T142(c.version.ApkId))
		w.Write(tlv.T145(c.deviceInfo.Guid))
		w.Write(tlv.T16A(c.sig.SrmToken))
		w.Write(tlv.T154(seq))
		w.Write(tlv.T141(c.deviceInfo.SimInfo, c.deviceInfo.APN))
		w.Write(tlv.T8(2052))
		w.Write(tlv.T511([]string{
			"tenpay.com", "openmobile.qq.com", "docs.qq.com", "connect.qq.com",
			"qzone.qq.com", "vip.qq.com", "qun.qq.com", "game.qq.com", "qqweb.qq.com",
			"office.qq.com", "ti.qq.com", "mail.qq.com", "qzone.com", "mma.qq.com",
		}))
		w.Write(tlv.T147(16, []byte(c.version.SortVersionName), c.version.ApkSign))
		w.Write(tlv.T177(c.version.BuildTime, c.version.SdkVersion))
		w.Write(tlv.T400(c.sig.G, c.Uin, c.deviceInfo.Guid, c.sig.Dpwd, 1, 16, c.sig.RandSeed))
		w.Write(tlv.T187(c.deviceInfo.MacAddress))
		w.Write(tlv.T188(c.deviceInfo.AndroidId))
		w.Write(tlv.T194(c.deviceInfo.IMSIMd5))
		w.Write(tlv.T202(c.deviceInfo.WifiBSSID, c.deviceInfo.WifiSSID))
		w.Write(tlv.T516())
		w.Write(tlv.T521(0))
		w.Write(tlv.T525(tlv.T536([]byte{0x01, 0x00})))
		// w.Write(tlv.545())
		w.Write(tlv.T545([]byte(c.deviceInfo.IMEI)))
	})

	m := oicq.Message{
		Uin:              uint32(c.Uin),
		Command:          0x810,
		EncryptionMethod: oicq.EM_ST,
		Body:             req,
	}

	return &network.Request{
		Type:        network.RequestTypeSimple,
		EncryptType: network.EncryptTypeEmptyKey,
		Uin:         c.Uin,
		SequenceID:  int32(seq),
		CommandName: "wtlogin.exchange_emp",
		Decode:      bindDecoder(c, decodeExchangeEmpResponse),
		Body:        c.oicq.Marshal(&m),
	}
}

func (c *QQClient) buildRequestChangeSigRequest(mainSigMap uint32) *network.Request {
	seq := c.nextSeq()
	req := c.buildOicqRequestPacket(c.Uin, 0x0810, binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt16(11)
		w.WriteUInt16(17)

		w.Write(tlv.T100(c.version.SSOVersion, 100, mainSigMap))
		w.Write(tlv.T10A(c.sig.TGT))
		w.Write(tlv.T116(c.version.MiscBitmap, c.version.SubSigmap))
		w.Write(tlv.T108(c.sig.Ksid))
		h := md5.Sum(c.sig.D2Key)
		w.Write(tlv.T144(
			c.deviceInfo.AndroidId,
			c.deviceInfo.GenDeviceInfoData(),
			c.deviceInfo.OSType,
			c.deviceInfo.Version.Release,
			c.deviceInfo.SimInfo,
			c.deviceInfo.APN,
			false, true, false, tlv.GuidFlag(),
			c.deviceInfo.Model,
			c.deviceInfo.Guid,
			c.deviceInfo.Brand,
			h[:],
		))
		w.Write(tlv.T143(c.sig.D2))
		w.Write(tlv.T142(c.version.ApkId))
		w.Write(tlv.T154(seq))
		w.Write(tlv.T18(16, uint32(c.Uin)))
		w.Write(tlv.T141(c.deviceInfo.SimInfo, c.deviceInfo.APN))
		w.Write(tlv.T8(2052))
		w.Write(tlv.T147(16, []byte(c.version.SortVersionName), c.version.ApkSign))
		w.Write(tlv.T177(c.version.BuildTime, c.version.SdkVersion))
		w.Write(tlv.T187(c.deviceInfo.MacAddress))
		w.Write(tlv.T188(c.deviceInfo.AndroidId))
		w.Write(tlv.T194(c.deviceInfo.IMSIMd5))
		w.Write(tlv.T511([]string{
			"tenpay.com", "openmobile.qq.com", "docs.qq.com", "connect.qq.com",
			"qzone.qq.com", "vip.qq.com", "qun.qq.com", "game.qq.com", "qqweb.qq.com",
			"office.qq.com", "ti.qq.com", "mail.qq.com", "qzone.com", "mma.qq.com",
		}))
		// w.Write(tlv.T202(c.deviceInfo.WifiBSSID, c.deviceInfo.WifiSSID))
	}))

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeEmptyKey,
		SequenceID:  int32(seq),
		Uin:         c.Uin,
		CommandName: "wtlogin.exchange_emp",
		Body:        req,
		Decode:      bindDecoder(c, decodeExchangeEmpResponse),
	}
}

// StatSvc.register
func (c *QQClient) buildClientRegisterPacket() *network.Request {
	seq := c.nextSeq()
	svc := &jce.SvcReqRegister{
		ConnType:     0,
		Uin:          c.Uin,
		Bid:          1 | 2 | 4,
		Status:       11,
		KickPC:       0,
		KickWeak:     0,
		IOSVersion:   int64(c.deviceInfo.Version.SDK),
		NetType:      1,
		RegType:      0,
		Guid:         c.deviceInfo.Guid,
		IsSetStatus:  0,
		LocaleId:     2052,
		DevName:      string(c.deviceInfo.Model),
		DevType:      string(c.deviceInfo.Model),
		OSVer:        string(c.deviceInfo.Version.Release),
		OpenPush:     1,
		LargeSeq:     1551,
		OldSSOIp:     0,
		NewSSOIp:     31806887127679168,
		ChannelNo:    "",
		CPID:         0,
		VendorName:   string(c.deviceInfo.VendorName),
		VendorOSName: string(c.deviceInfo.VendorOSName),
		B769:         []byte{0x0A, 0x04, 0x08, 0x2E, 0x10, 0x00, 0x0A, 0x05, 0x08, 0x9B, 0x02, 0x10, 0x00},
		SetMute:      0,
	}
	b := append([]byte{0x0A}, svc.ToBytes()...)
	b = append(b, 0x0B)
	buf := &jce.RequestDataVersion3{
		Map: map[string][]byte{"SvcReqRegister": b},
	}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		SServantName: "PushService",
		SFuncName:    "SvcReqRegister",
		SBuffer:      buf.ToBytes(),
		Context:      make(map[string]string),
		Status:       make(map[string]string),
	}

	return &network.Request{
		Type:        network.RequestTypeLogin,
		EncryptType: network.EncryptTypeD2Key,
		SequenceID:  int32(seq),
		Uin:         c.Uin,
		CommandName: "StatSvc.register",
		Body:        pkt.ToBytes(),
		Decode:      bindDecoder(c, decodeClientRegisterResponse),
	}
}

func (c *QQClient) buildStatusSetPacket(status, extStatus int32) *network.Request {
	svc := &jce.SvcReqRegister{
		ConnType:        0,
		Uin:             c.Uin,
		Bid:             7,
		Status:          status,
		KickPC:          0,
		KickWeak:        0,
		Timestamp:       time.Now().Unix(),
		IOSVersion:      int64(c.deviceInfo.Version.SDK),
		NetType:         1,
		RegType:         0,
		Guid:            c.deviceInfo.Guid,
		IsSetStatus:     1,
		LocaleId:        2052,
		DevName:         string(c.deviceInfo.Model),
		DevType:         string(c.deviceInfo.Model),
		OSVer:           string(c.deviceInfo.Version.Release),
		OpenPush:        1,
		LargeSeq:        1551,
		ExtOnlineStatus: int64(extStatus),
	}
	buf := &jce.RequestDataVersion3{
		Map: map[string][]byte{"SvcReqRegister": packUniRequestData(svc.ToBytes())},
	}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		SServantName: "PushService",
		SFuncName:    "SvcReqRegister",
		SBuffer:      buf.ToBytes(),
		Context:      make(map[string]string),
		Status:       make(map[string]string),
	}
	return c.uniRequest("StatSvc.SetStatusFromClient", pkt.ToBytes(), nil)
}

// ConfigPushSvc.PushResp
func (c *QQClient) buildConfPushRespPacket(t int32, pktSeq int64, jceBuf []byte) *network.Request {
	req := jce.NewJceWriter()
	req.WriteInt32(t, 1)
	req.WriteInt64(pktSeq, 2)
	req.WriteBytes(jceBuf, 3)
	buf := &jce.RequestDataVersion3{
		Map: map[string][]byte{"PushResp": packUniRequestData(req.Bytes())},
	}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		SServantName: "QQService.ConfigPushSvc.MainServant",
		SFuncName:    "PushResp",
		SBuffer:      buf.ToBytes(),
		Context:      make(map[string]string),
		Status:       make(map[string]string),
	}
	return c.uniRequest("ConfigPushSvc.PushResp", pkt.ToBytes(), nil)
}

// friendlist.getFriendGroupList
func (c *QQClient) buildFriendGroupListRequest(friendStartIndex, friendListCount, groupStartIndex, groupListCount int16) *network.Request {
	d50, _ := proto.Marshal(&pb.D50ReqBody{
		Appid:                   1002,
		ReqMusicSwitch:          1,
		ReqMutualmarkAlienation: 1,
		ReqKsingSwitch:          1,
		ReqMutualmarkLbsshare:   1,
	})
	req := &jce.FriendListRequest{
		Reqtype: 3,
		IfReflush: func() byte {
			if friendStartIndex <= 0 {
				return 0
			}
			return 1
		}(),
		Uin:         c.Uin,
		StartIndex:  friendStartIndex,
		FriendCount: friendListCount,
		GroupId:     0,
		IfGetGroupInfo: func() byte {
			if groupListCount <= 0 {
				return 0
			}
			return 1
		}(),
		GroupStartIndex: byte(groupStartIndex),
		GroupCount:      byte(groupListCount),
		IfGetMSFGroup:   0,
		IfShowTermType:  1,
		Version:         27,
		UinList:         nil,
		AppType:         0,
		IfGetDOVId:      0,
		IfGetBothFlag:   0,
		D50:             d50,
		D6B:             EmptyBytes,
		SnsTypeList:     []int64{13580, 13581, 13582},
	}
	buf := &jce.RequestDataVersion3{
		Map: map[string][]byte{"FL": packUniRequestData(req.ToBytes())},
	}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		CPacketType:  0x003,
		IRequestId:   1921334514,
		SServantName: "mqq.IMService.FriendListServiceServantObj",
		SFuncName:    "GetFriendListReq",
		SBuffer:      buf.ToBytes(),
		Context:      make(map[string]string),
		Status:       make(map[string]string),
	}
	return c.uniRequest("friendlist.getFriendGroupList", pkt.ToBytes(), decodeFriendGroupListResponse)
}

// SummaryCard.ReqSummaryCard
func (c *QQClient) buildSummaryCardRequest(target int64) *network.Request {
	seq := c.nextSeq()
	packBusinessBuf := func(t int32, buf []byte) []byte {
		return binary.NewWriterF(func(w *binary.Writer) {
			comm, _ := proto.Marshal(&profilecard.BusiComm{
				Ver:      proto.Int32(1),
				Seq:      proto.Int32(int32(seq)),
				Fromuin:  &c.Uin,
				Touin:    &target,
				Service:  &t,
				Platform: proto.Int32(2),
				Qqver:    proto.String("8.4.18.4945"),
				Build:    proto.Int32(4945),
			})
			w.WriteByte(40)
			w.WriteUInt32(uint32(len(comm)))
			w.WriteUInt32(uint32(len(buf)))
			w.Write(comm)
			w.Write(buf)
			w.WriteByte(41)
		})
	}
	gate, _ := proto.Marshal(&profilecard.GateVaProfileGateReq{
		UCmd:           proto.Int32(3),
		StPrivilegeReq: &profilecard.GatePrivilegeBaseInfoReq{UReqUin: &target},
		StGiftReq:      &profilecard.GateGetGiftListReq{Uin: proto.Int32(int32(target))},
		StVipCare:      &profilecard.GateGetVipCareReq{Uin: &target},
		OidbFlag: []*profilecard.GateOidbFlagInfo{
			{
				Fieled: proto.Int32(42334),
			},
			{
				Fieled: proto.Int32(42340),
			},
			{
				Fieled: proto.Int32(42344),
			},
			{
				Fieled: proto.Int32(42354),
			},
		},
	})
	/*
		e5b, _ := proto.Marshal(&oidb.DE5BReqBody{
			Uin:                   proto.Uint64(uint64(target)),
			MaxCount:              proto.Uint32(10),
			ReqAchievementContent: proto.Bool(false),
		})
		ec4, _ := proto.Marshal(&oidb.DEC4ReqBody{
			Uin:       proto.Uint64(uint64(target)),
			QuestNum:  proto.Uint64(10),
			FetchType: proto.Uint32(1),
		})
	*/
	req := &jce.SummaryCardReq{
		Uin:              target,
		ComeFrom:         31,
		GetControl:       69181,
		AddFriendSource:  3001,
		SecureSig:        []byte{0x00},
		ReqMedalWallInfo: 0,
		Req0x5ebFieldId:  []int64{27225, 27224, 42122, 42121, 27236, 27238, 42167, 42172, 40324, 42284, 42326, 42325, 42356, 42363, 42361, 42367, 42377, 42425, 42505, 42488},
		ReqServices:      [][]byte{packBusinessBuf(16, gate)},
		ReqNearbyGodInfo: 1,
		ReqExtendCard:    1,
	}
	head := jce.NewJceWriter()
	head.WriteInt32(2, 0)
	buf := &jce.RequestDataVersion3{Map: map[string][]byte{
		"ReqHead":        packUniRequestData(head.Bytes()),
		"ReqSummaryCard": packUniRequestData(req.ToBytes()),
	}}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		SServantName: "SummaryCardServantObj",
		SFuncName:    "ReqSummaryCard",
		SBuffer:      buf.ToBytes(),
		Context:      make(map[string]string),
		Status:       make(map[string]string),
	}
	return c.uniPacketWithSeq(seq, "SummaryCard.ReqSummaryCard", pkt.ToBytes(), decodeSummaryCardResponse)
}

// friendlist.delFriend
func (c *QQClient) buildFriendDeletePacket(target int64) *network.Request {
	req := &jce.DelFriendReq{
		Uin:     c.Uin,
		DelUin:  target,
		DelType: 2,
		Version: 1,
	}
	buf := &jce.RequestDataVersion3{
		Map: map[string][]byte{"DF": packUniRequestData(req.ToBytes())},
	}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		IRequestId:   c.nextPacketSeq(),
		SServantName: "mqq.IMService.FriendListServiceServantObj",
		SFuncName:    "DelFriendReq",
		SBuffer:      buf.ToBytes(),
		Context:      make(map[string]string),
		Status:       make(map[string]string),
	}
	return c.uniRequest("friendlist.delFriend", pkt.ToBytes(), decodeFriendDeleteResponse)
}

// friendlist.GetTroopListReqV2
func (c *QQClient) buildGroupListRequest(vecCookie []byte) *network.Request {
	req := &jce.TroopListRequest{
		Uin:              c.Uin,
		GetMSFMsgFlag:    1,
		Cookies:          vecCookie,
		GroupInfo:        []int64{},
		GroupFlagExt:     1,
		Version:          7,
		CompanyId:        0,
		VersionNum:       1,
		GetLongGroupName: 1,
	}
	buf := &jce.RequestDataVersion3{
		Map: map[string][]byte{"GetTroopListReqV2Simplify": packUniRequestData(req.ToBytes())},
	}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		CPacketType:  0x00,
		IRequestId:   c.nextPacketSeq(),
		SServantName: "mqq.IMService.FriendListServiceServantObj",
		SFuncName:    "GetTroopListReqV2Simplify",
		SBuffer:      buf.ToBytes(),
		Context:      make(map[string]string),
		Status:       make(map[string]string),
	}
	return c.uniRequest("friendlist.GetTroopListReqV2", pkt.ToBytes(), decodeGroupListResponse)
}

// friendlist.GetTroopMemberListReq
func (c *QQClient) buildGroupMemberListRequest(groupUin, groupCode, nextUin int64) *network.Request {
	req := &jce.TroopMemberListRequest{
		Uin:       c.Uin,
		GroupCode: groupCode,
		NextUin:   nextUin,
		GroupUin:  groupUin,
		Version:   2,
	}
	b := append([]byte{0x0A}, req.ToBytes()...)
	b = append(b, 0x0B)
	buf := &jce.RequestDataVersion3{
		Map: map[string][]byte{"GTML": b},
	}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		IRequestId:   c.nextPacketSeq(),
		SServantName: "mqq.IMService.FriendListServiceServantObj",
		SFuncName:    "GetTroopMemberListReq",
		SBuffer:      buf.ToBytes(),
		Context:      make(map[string]string),
		Status:       make(map[string]string),
	}
	return c.uniRequest("friendlist.GetTroopMemberListReq", pkt.ToBytes(), decodeGroupMemberListResponse)
}

// group_member_card.get_group_member_card_info
func (c *QQClient) buildGroupMemberInfoRequest(groupCode, uin int64) *network.Request {
	req := &pb.GroupMemberReqBody{
		GroupCode:       groupCode,
		Uin:             uin,
		NewClient:       true,
		ClientType:      1,
		RichCardNameVer: 1,
	}
	payload, _ := proto.Marshal(req)
	return c.uniRequest("group_member_card.get_group_member_card_info", payload, decodeGroupMemberInfoResponse)
}

// MessageSvc.PbGetMsg
func (c *QQClient) buildGetMessageRequest(flag msg.SyncFlag, msgTime int64) *network.Request {
	cook := c.sig.SyncCookie
	if cook == nil {
		cook, _ = proto.Marshal(&msg.SyncCookie{
			Time:   &msgTime,
			Ran1:   proto.Int64(758330138),
			Ran2:   proto.Int64(2480149246),
			Const1: proto.Int64(1167238020),
			Const2: proto.Int64(3913056418),
			Const3: proto.Int64(0x1D),
		})
	}
	req := &msg.GetMessageRequest{
		SyncFlag:           &flag,
		SyncCookie:         cook,
		LatestRambleNumber: proto.Int32(20),
		OtherRambleNumber:  proto.Int32(3),
		OnlineSyncFlag:     proto.Int32(1),
		ContextFlag:        proto.Int32(1),
		MsgReqType:         proto.Int32(1),
		PubaccountCookie:   EmptyBytes,
		MsgCtrlBuf:         EmptyBytes,
		ServerBuf:          EmptyBytes,
	}
	payload, _ := proto.Marshal(req)
	return c.uniRequest("MessageSvc.PbGetMsg", payload, decodeMessageSvcPacket)
}

// MessageSvc.PbDeleteMsg
func (c *QQClient) buildDeleteMessageRequestPacket(msg []*pb.MessageItem) *network.Request {
	req := &pb.DeleteMessageRequest{Items: msg}
	payload, _ := proto.Marshal(req)
	return c.uniRequest("MessageSvc.PbDeleteMsg", payload, nil)
}

// OnlinePush.RespPush
func (c *QQClient) buildDeleteOnlinePushPacket(uin int64, svrip int32, pushToken []byte, seq uint16, delMsg []jce.PushMessageInfo) []byte {
	req := &jce.SvcRespPushMsg{Uin: uin, Svrip: svrip, PushToken: pushToken, DelInfos: []jce.IJceStruct{}}
	for _, m := range delMsg {
		req.DelInfos = append(req.DelInfos, &jce.DelMsgInfo{
			FromUin:    m.FromUin,
			MsgSeq:     m.MsgSeq,
			MsgCookies: m.MsgCookies,
			MsgTime:    m.MsgTime,
		})
	}
	b := packUniRequestData(req.ToBytes())
	buf := &jce.RequestDataVersion3{
		Map: map[string][]byte{"resp": b},
	}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		IRequestId:   int32(seq),
		SServantName: "OnlinePush",
		SFuncName:    "SvcRespPushMsg",
		SBuffer:      buf.ToBytes(),
		Context:      make(map[string]string),
		Status:       make(map[string]string),
	}
	return c.transport.PackPacket(c.uniPacketWithSeq(seq, "OnlinePush.RespPush", pkt.ToBytes(), nil))
}

// LongConn.OffPicUp
func (c *QQClient) buildOffPicUpRequest(target int64, md5 []byte, size int32) *network.Request {
	req := &cmd0x352.ReqBody{
		Subcmd: proto.Uint32(1),
		TryupImgReq: []*cmd0x352.D352TryUpImgReq{
			{
				SrcUin:       proto.Uint64(uint64(c.Uin)),
				DstUin:       proto.Uint64(uint64(target)),
				FileMd5:      md5,
				FileSize:     proto.Uint64(uint64(size)),
				FileName:     []byte(hex.EncodeToString(md5) + ".jpg"),
				SrcTerm:      proto.Uint32(5),
				PlatformType: proto.Uint32(9),
				BuType:       proto.Uint32(1),
				PicOriginal:  proto.Bool(true),
				PicType:      proto.Uint32(1000),
				BuildVer:     []byte("8.2.7.4410"),
				FileIndex:    EmptyBytes,
				SrvUpload:    proto.Uint32(1),
				TransferUrl:  EmptyBytes,
			},
		},
	}
	payload, _ := proto.Marshal(req)
	return c.uniRequest("LongConn.OffPicUp", payload, decodeOffPicUpResponse)
}

// ProfileService.Pb.ReqSystemMsgNew.Friend
func (c *QQClient) buildSystemMsgNewFriendRequest() *network.Request {
	req := &structmsg.ReqSystemMsgNew{
		MsgNum:    20,
		Version:   1000,
		Checktype: 2,
		Flag: &structmsg.FlagInfo{
			FrdMsgDiscuss2ManyChat:       1,
			FrdMsgGetBusiCard:            1,
			FrdMsgNeedWaitingMsg:         1,
			FrdMsgUint32NeedAllUnreadMsg: 1,
			GrpMsgMaskInviteAutoJoin:     1,
		},
		FriendMsgTypeFlag: 1,
	}
	payload, _ := proto.Marshal(req)
	return c.uniRequest("ProfileService.Pb.ReqSystemMsgNew.Friend", payload, decodeSystemMsgFriendPacket)
}

// friendlist.ModifyGroupCardReq
func (c *QQClient) buildEditGroupTagPacket(groupCode, memberUin int64, newTag string) *network.Request {
	req := &jce.ModifyGroupCardRequest{
		GroupCode: groupCode,
		UinInfo: []jce.IJceStruct{
			&jce.UinInfo{
				Uin:  memberUin,
				Flag: 31,
				Name: newTag,
			},
		},
	}
	buf := &jce.RequestDataVersion3{Map: map[string][]byte{"MGCREQ": packUniRequestData(req.ToBytes())}}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		IRequestId:   c.nextPacketSeq(),
		SServantName: "mqq.IMService.FriendListServiceServantObj",
		SFuncName:    "ModifyGroupCardReq",
		SBuffer:      buf.ToBytes(),
		Context:      map[string]string{},
		Status:       map[string]string{},
	}
	return c.uniRequest("friendlist.ModifyGroupCardReq", pkt.ToBytes(), nil)
}

// OidbSvc.0x8fc_2
func (c *QQClient) buildEditSpecialTitlePacket(groupCode, memberUin int64, newTitle string) *network.Request {
	body := &oidb.D8FCReqBody{
		GroupCode: &groupCode,
		MemLevelInfo: []*oidb.D8FCMemberInfo{
			{
				Uin:                    &memberUin,
				UinName:                []byte(newTitle),
				SpecialTitle:           []byte(newTitle),
				SpecialTitleExpireTime: proto.Int32(-1),
			},
		},
	}
	b, _ := proto.Marshal(body)
	payload := c.packOIDBPackage(2300, 2, b)
	return c.uniRequest("OidbSvc.0x8fc_2", payload, nil)
}

// OidbSvc.0x89a_0
func (c *QQClient) buildGroupOperationPacket(body *oidb.D89AReqBody) *network.Request {
	b, _ := proto.Marshal(body)
	payload := c.packOIDBPackage(2202, 0, b)
	return c.uniRequest("OidbSvc.0x89a_0", payload, nil)
}

// OidbSvc.0x89a_0
func (c *QQClient) buildGroupNameUpdateRequest(groupCode int64, newName string) *network.Request {
	body := &oidb.D89AReqBody{
		GroupCode: groupCode,
		StGroupInfo: &oidb.D89AGroupinfo{
			IngGroupName: []byte(newName),
		},
	}
	return c.buildGroupOperationPacket(body)
}

func (c *QQClient) buildGroupMemoUpdatePacket(groupCode int64, newMemo string) *network.Request {
	body := &oidb.D89AReqBody{
		GroupCode: groupCode,
		StGroupInfo: &oidb.D89AGroupinfo{
			IngGroupMemo: []byte(newMemo),
		},
	}
	return c.buildGroupOperationPacket(body)
}

// OidbSvc.0x89a_0
func (c *QQClient) buildGroupMuteAllPacket(groupCode int64, mute bool) *network.Request {
	shutUpTime := int32(0)
	if mute {
		shutUpTime = 268435455
	}
	body := &oidb.D89AReqBody{
		GroupCode: groupCode,
		StGroupInfo: &oidb.D89AGroupinfo{
			ShutupTime: &shutUpTime,
		},
	}
	return c.buildGroupOperationPacket(body)
}

// OidbSvc.0x8a0_0
func (c *QQClient) buildGroupKickPacket(groupCode, memberUin int64, kickMsg string, block bool) *network.Request {
	flagBlock := 0
	if block {
		flagBlock = 1
	}
	body := &oidb.D8A0ReqBody{
		OptUint64GroupCode: groupCode,
		MsgKickList: []*oidb.D8A0KickMemberInfo{
			{
				OptUint32Operate:   5,
				OptUint64MemberUin: memberUin,
				OptUint32Flag:      int32(flagBlock),
			},
		},
		KickMsg: []byte(kickMsg),
	}
	b, _ := proto.Marshal(body)
	payload := c.packOIDBPackage(2208, 0, b)
	return c.uniRequest("OidbSvc.0x8a0_0", payload, nil)
}

// OidbSvc.0x570_8
func (c *QQClient) buildGroupMutePacket(groupCode, memberUin int64, time uint32) *network.Request {
	b, cl := binary.OpenWriterF(func(w *binary.Writer) {
		w.WriteUInt32(uint32(groupCode))
		w.WriteByte(32)
		w.WriteUInt16(1)
		w.WriteUInt32(uint32(memberUin))
		w.WriteUInt32(time)
	})
	payload := c.packOIDBPackage(1392, 8, b)
	cl()
	return c.uniRequest("OidbSvc.0x570_8", payload, nil)
}

// OidbSvc.0xed3
func (c *QQClient) buildGroupPokeRequest(groupCode, target int64) *network.Request {
	body := &oidb.DED3ReqBody{
		ToUin:     target,
		GroupCode: groupCode,
	}
	b, _ := proto.Marshal(body)
	payload := c.packOIDBPackage(3795, 1, b)
	return c.uniRequest("OidbSvc.0xed3", payload, nil)
}

// OidbSvc.0xed3
func (c *QQClient) buildFriendPokeRequest(target int64) *network.Request {
	body := &oidb.DED3ReqBody{
		ToUin:  target,
		AioUin: target,
	}
	b, _ := proto.Marshal(body)
	payload := c.packOIDBPackage(3795, 1, b)
	return c.uniRequest("OidbSvc.0xed3", payload, nil)
}

// OidbSvc.0x55c_1
func (c *QQClient) buildGroupAdminSetPacket(groupCode, member int64, flag bool) *network.Request {
	b, cl := binary.OpenWriterF(func(w *binary.Writer) {
		w.WriteUInt32(uint32(groupCode))
		w.WriteUInt32(uint32(member))
		w.WriteBool(flag)
	})
	payload := c.packOIDBPackage(1372, 1, b)
	cl()
	return c.uniRequest("OidbSvc.0x55c_1", payload, nil)
}

// ProfileService.GroupMngReq
func (c *QQClient) buildQuitGroupPacket(groupCode int64) *network.Request {
	jw := jce.NewJceWriter()
	jw.WriteInt32(2, 0)
	jw.WriteInt64(c.Uin, 1)
	b, cl := binary.OpenWriterF(func(w *binary.Writer) {
		w.WriteUInt32(uint32(c.Uin))
		w.WriteUInt32(uint32(groupCode))
	})
	jw.WriteBytes(b, 2)
	cl()
	buf := &jce.RequestDataVersion3{Map: map[string][]byte{"GroupMngReq": packUniRequestData(jw.Bytes())}}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		IRequestId:   c.nextPacketSeq(),
		SServantName: "KQQ.ProfileService.ProfileServantObj",
		SFuncName:    "GroupMngReq",
		SBuffer:      buf.ToBytes(),
		Context:      map[string]string{},
		Status:       map[string]string{},
	}
	return c.uniRequest("ProfileService.GroupMngReq", pkt.ToBytes(), nil)
}

// ProfileService.ReqSetSettings
// groupCode: ~int | string
func (c *QQClient) buildChangeGroupNotifyPacket(groupCode interface{}, level byte) *network.Request {
	payload := make([]byte, 0, 48) // 40 ?
	// jce:List,tag=0;jce:int1,tag=0,value=1
	payload = append(payload, []byte{0x09, 0x00, 0x01}...)
	// jce:StructStart,tag=0
	payload = append(payload, 0x0A)
	// jce:String1,tag=0,Path
	payload = append(payload, 0x06)
	// padLen
	padLen := len(payload)
	payload = append(payload, 0x00)
	payload = append(payload, "message.group.policy."...)
	if code, ok := groupCode.(string); ok {
		payload = append(payload, code...)
	} else if code, ok := groupCode.(int64); ok {
		payload = strconv.AppendInt(payload, code, 10)
	} else {
		panic("invaild parameter")
	}
	// padLen
	payload[padLen] = byte(len(payload) - padLen - 1)
	// jce:String1,tag=1,Value
	payload = append(payload, 0x16)
	// len=1,fixed
	payload = append(payload, 0x01)
	// levelChar
	payload = append(payload, level)
	// jce:StructEnd
	payload = append(payload, 0x0B)
	buf := &jce.RequestDataVersion3{Map: map[string][]byte{
		"ReqSetSettings": packUniRequestData(payload),
		"ReqHeader":      packUniRequestData([]byte{0x0c}),
	}}
	pkt := &jce.RequestPacket{
		IVersion:     3,
		IRequestId:   c.nextPacketSeq(),
		SServantName: "KQQ.ProfileService.ProfileServantObj",
		SFuncName:    "ReqSetSettings",
		SBuffer:      buf.ToBytes(),
		Context:      map[string]string{},
		Status:       map[string]string{},
	}

	return c.uniRequest("ProfileService.ReqSetSettings", pkt.ToBytes(), nil)
	// 我放弃了，我不会用，真的不会
	//payload := jce.NewJceWriter()
	//payload.WriteInt64Slice(make([]int64, 0), 0)
	//payload
	//setting := jce.Setting{
	//	Path:  fmt.Sprintf("message.group.policy.%d", groupCode),
	//	Value: string(level), // 单字节的字符串，没有转换错
	//}
}

/* this function is unused
// LightAppSvc.mini_app_info.GetAppInfoById
func (c *QQClient) buildAppInfoRequestPacket(id string) *network.Request {
	seq := c.nextSeq()
	req := &qweb.GetAppInfoByIdReq{
		AppId:           id,
		NeedVersionInfo: 1,
	}
	b, _ := proto.Marshal(req)
	body := &qweb.QWebReq{
		Seq:        proto.Int64(1),
		Qua:        proto.String("V1_AND_SQ_8.4.8_1492_YYB_D"),
		device: proto.String(c.getWebDeviceInfo()),
		BusiBuff:   b,
		TraceId:    proto.String(fmt.Sprintf("%v_%v_%v", c.Uin, time.Now().Format("0102150405"), rand.Int63())),
	}
	payload, _ := proto.Marshal(body)
	packet := packets.BuildUniPacket(c.Uin, seq, "LightAppSvc.mini_app_info.GetAppInfoById", 1, c.SessionId, EmptyBytes, c.sigInfo.d2Key, payload)
	return seq, packet
}
*/

func (c *QQClient) buildWordSegmentationPacket(data []byte) *network.Request {
	payload := c.packOIDBPackageProto(3449, 1, &oidb.D79ReqBody{
		Uin:     uint64(c.Uin),
		Content: data,
		Qua:     []byte("and_537065262_8.4.5"),
	})
	return c.uniRequest("OidbSvc.0xd79", payload, decodeWordSegmentation)
}
