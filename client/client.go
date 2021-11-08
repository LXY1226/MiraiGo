package client

import (
	"crypto/md5"
	"fmt"
	"math"
	"math/rand"
	"net"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/Mrs4s/MiraiGo/binary"
	"github.com/Mrs4s/MiraiGo/binary/jce"
	"github.com/Mrs4s/MiraiGo/client/pb/msg"
	"github.com/Mrs4s/MiraiGo/internal/crypto"
	"github.com/Mrs4s/MiraiGo/internal/packets"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Mrs4s/MiraiGo/utils"
)

//go:generate go run github.com/a8m/syncmap -o "handler_map_gen.go" -pkg client -name HandlerMap "map[uint16]*handlerInfo"

type QQClient struct {
	Uin         int64
	PasswordMd5 [16]byte

	stat Statistics
	once sync.Once
	conn *net.TCPConn

	// option
	AllowSlider bool

	// account info
	Nickname      string
	Age           uint16
	Gender        uint16
	FriendList    []*FriendInfo
	GroupList     []*GroupInfo
	OnlineClients []*OtherClientInfo
	Online        bool
	QiDian        *QiDianAccountInfo
	GuildService  *GuildService

	// protocol public field
	SequenceId              int32
	OutGoingPacketSessionId []byte
	RandomKey               []byte
	ConnectTime             time.Time

	// internal state
	handlers        HandlerMap
	waiters         sync.Map
	servers         []*net.TCPAddr
	currServerIndex int
	retryTimes      int
	version         *versionInfo
	deviceInfo      *DeviceInfo
	alive           bool
	ecdh            *crypto.EncryptECDH

	// tlv cache
	t104        []byte
	t174        []byte
	g           []byte
	t402        []byte
	t150        []byte
	t149        []byte
	t528        []byte
	t530        []byte
	randSeed    []byte // t403
	rollbackSig []byte

	// sync info
	syncCookie       []byte
	pubAccountCookie []byte
	msgCtrlBuf       []byte
	ksid             []byte

	// session info
	sigInfo        *loginSigInfo
	bigDataSession *bigDataSessionInfo
	dpwd           []byte
	timeDiff       int64
	pwdFlag        bool

	// address
	srvSsoAddrs     []string
	otherSrvAddrs   []string
	fileStorageInfo *jce.FileStoragePushFSSvcList

	// message state
	lastMessageSeq         int32
	msgSvcCache            *utils.Cache
	lastC2CMsgTime         int64
	transCache             *utils.Cache
	lastLostMsg            string
	groupSysMsgCache       *GroupSystemMessages
	groupMsgBuilders       sync.Map
	onlinePushCache        *utils.Cache
	requestPacketRequestID int32
	groupSeq               int32
	friendSeq              int32
	heartbeatEnabled       bool
	groupDataTransSeq      int32
	highwayApplyUpSeq      int32

	// handlers
	groupMessageReceiptHandler sync.Map
	EventHandler
	Logger

	groupListLock sync.Mutex
}

type loginSigInfo struct {
	loginBitmap uint64
	tgt         []byte
	tgtKey      []byte

	srmToken           []byte // study room manager | 0x16a
	t133               []byte
	encryptedA1        []byte
	userStKey          []byte
	userStWebSig       []byte
	sKey               []byte
	sKeyExpiredTime    int64
	d2                 []byte
	d2Key              []byte
	wtSessionTicketKey []byte
	deviceToken        []byte

	psKeyMap    map[string][]byte
	pt4TokenMap map[string][]byte
}

type QiDianAccountInfo struct {
	MasterUin  int64
	ExtName    string
	CreateTime int64

	bigDataReqAddrs   []string
	bigDataReqSession *bigDataSessionInfo
}

type handlerInfo struct {
	fun     func(i interface{}, err error)
	dynamic bool
	params  requestParams
}

var decoders = map[string]func(*QQClient, *incomingPacketInfo, []byte) (interface{}, error){
	"wtlogin.login":                                decodeLoginResponse,
	"wtlogin.exchange_emp":                         decodeExchangeEmpResponse,
	"wtlogin.trans_emp":                            decodeTransEmpResponse,
	"StatSvc.register":                             decodeClientRegisterResponse,
	"StatSvc.ReqMSFOffline":                        decodeMSFOfflinePacket,
	"MessageSvc.PushNotify":                        decodeSvcNotify,
	"OnlinePush.ReqPush":                           decodeOnlinePushReqPacket,
	"OnlinePush.PbPushTransMsg":                    decodeOnlinePushTransPacket,
	"ConfigPushSvc.PushReq":                        decodePushReqPacket,
	"MessageSvc.PbGetMsg":                          decodeMessageSvcPacket,
	"MessageSvc.PushForceOffline":                  decodeForceOfflinePacket,
	"PbMessageSvc.PbMsgWithDraw":                   decodeMsgWithDrawResponse,
	"friendlist.getFriendGroupList":                decodeFriendGroupListResponse,
	"friendlist.delFriend":                         decodeFriendDeleteResponse,
	"friendlist.GetTroopListReqV2":                 decodeGroupListResponse,
	"friendlist.GetTroopMemberListReq":             decodeGroupMemberListResponse,
	"group_member_card.get_group_member_card_info": decodeGroupMemberInfoResponse,
	"PttStore.GroupPttUp":                          decodeGroupPttStoreResponse,
	"LongConn.OffPicUp":                            decodeOffPicUpResponse,
	"ProfileService.Pb.ReqSystemMsgNew.Group":      decodeSystemMsgGroupPacket,
	"ProfileService.Pb.ReqSystemMsgNew.Friend":     decodeSystemMsgFriendPacket,
	"OidbSvc.0xd79":                                decodeWordSegmentation,
	"OidbSvc.0x990":                                decodeTranslateResponse,
	"SummaryCard.ReqSummaryCard":                   decodeSummaryCardResponse,
	"LightAppSvc.mini_app_info.GetAppInfoById":     decodeAppInfoResponse,
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

}

// NewClient create new qq client
func NewClient(uin int64, password string) *QQClient {
	return NewClientMd5(uin, md5.Sum([]byte(password)))
}

func NewClientEmpty() *QQClient {
	return NewClient(0, "")
}

func NewClientMd5(uin int64, passwordMd5 [16]byte) *QQClient {
	cli := &QQClient{
		Uin:                     uin,
		PasswordMd5:             passwordMd5,
		SequenceId:              0x3635,
		AllowSlider:             true,
		RandomKey:               make([]byte, 16),
		OutGoingPacketSessionId: []byte{0x02, 0xB0, 0x5B, 0x8B},
		Logger:                  nopLogger{},
		sigInfo:                 &loginSigInfo{},
		requestPacketRequestID:  1921334513,
		groupSeq:                int32(rand.Intn(20000)),
		friendSeq:               22911,
		highwayApplyUpSeq:       77918,
		EventHandler:            nopHandlers,
		msgSvcCache:             utils.NewCache(time.Second * 15),
		transCache:              utils.NewCache(time.Second * 15),
		onlinePushCache:         utils.NewCache(time.Second * 15),
		servers:                 []*net.TCPAddr{},
		ecdh:                    crypto.NewEcdh(),
	}
	cli.GuildService = &GuildService{c: cli}
	cli.ecdh.FetchPubKey(uin)
	cli.UseDevice(SystemDeviceInfo)
	sso, err := getSSOAddress()
	if err == nil && len(sso) > 0 {
		cli.servers = append(sso, cli.servers...)
	}
	adds, err := net.LookupIP("msfwifi.3g.qq.com") // host servers
	if err == nil && len(adds) > 0 {
		var hostAddrs []*net.TCPAddr
		for _, addr := range adds {
			hostAddrs = append(hostAddrs, &net.TCPAddr{
				IP:   addr,
				Port: 8080,
			})
		}
		cli.servers = append(hostAddrs, cli.servers...)
	}
	if len(cli.servers) == 0 {
		cli.servers = []*net.TCPAddr{ // default servers
			{IP: net.IP{42, 81, 172, 22}, Port: 80},
			{IP: net.IP{42, 81, 172, 81}, Port: 80},
			{IP: net.IP{42, 81, 172, 147}, Port: 443},
			{IP: net.IP{114, 221, 144, 215}, Port: 80},
			{IP: net.IP{114, 221, 148, 59}, Port: 14000},
			{IP: net.IP{125, 94, 60, 146}, Port: 80},
		}
	}
	/*	pings := make([]int64, len(cli.servers))
		wg := sync.WaitGroup{}
		wg.Add(len(cli.servers))
		for i := range cli.servers {
			go func(index int) {
				defer wg.Done()
				p, err := qualityTest(cli.servers[index].String())
				if err != nil {
					pings[index] = 9999
					return
				}
				pings[index] = p
			}(i)
		}
		wg.Wait()
		sort.Slice(cli.servers, func(i, j int) bool {
			return pings[i] < pings[j]
		})
		if len(cli.servers) > 3 {
			cli.servers = cli.servers[0 : len(cli.servers)/2] // 保留ping值中位数以上的server
		}*/
	//cli.TCP.PlannedDisconnect(cli.plannedDisconnect)
	//cli.TCP.UnexpectedDisconnect(cli.unexpectedDisconnect)
	rand.Read(cli.RandomKey)
	return cli
}

func (c *QQClient) UseDevice(info *DeviceInfo) {
	c.version = genVersionInfo(info.Protocol)
	c.ksid = []byte(fmt.Sprintf("|%s|A8.2.7.27f6ea96", info.IMEI))
	c.deviceInfo = info
}

// Login send login request
func (c *QQClient) Login() (*LoginResponse, error) {
	if c.Online {
		return nil, ErrAlreadyOnline
	}
	err := c.connectFastest()
	if err != nil {
		return nil, err
	}
	rsp, err := c.sendAndWait(c.buildLoginPacket())
	if err != nil {
		c.Disconnect()
		return nil, err
	}
	l := rsp.(LoginResponse)
	if l.Success {
		_ = c.init(false)
	}
	return &l, nil
}

func (c *QQClient) TokenLogin(token []byte) error {
	if c.Online {
		return ErrAlreadyOnline
	}
	err := c.connectFastest()
	if err != nil {
		return err
	}
	{
		r := binary.NewReader(token)
		c.Uin = r.ReadInt64()
		c.sigInfo.d2 = r.ReadBytesShort()
		c.sigInfo.d2Key = r.ReadBytesShort()
		c.sigInfo.tgt = r.ReadBytesShort()
		c.sigInfo.srmToken = r.ReadBytesShort()
		c.sigInfo.t133 = r.ReadBytesShort()
		c.sigInfo.encryptedA1 = r.ReadBytesShort()
		c.sigInfo.wtSessionTicketKey = r.ReadBytesShort()
		c.OutGoingPacketSessionId = r.ReadBytesShort()
		// SystemDeviceInfo.TgtgtKey = r.ReadBytesShort()
		c.deviceInfo.TgtgtKey = r.ReadBytesShort()
	}
	_, err = c.sendAndWait(c.buildRequestChangeSigPacket())
	if err != nil {
		return err
	}
	return c.init(true)
}

func (c *QQClient) ReLogin() error {
	if c.Online {
		return ErrAlreadyOnline
	}
	err := c.connectFastest()
	if err != nil {
		return err
	}
	_, err = c.sendAndWait(c.buildRequestChangeSigPacket())
	if err != nil {
		return err
	}
	err = c.init(true)
	// 登录失败
	if err != nil {
		c.Disconnect()
	}
	return err
}

func (c *QQClient) LoadState(token []byte) {
	r := binary.NewReader(token)
	c.Uin = r.ReadInt64()
	c.sigInfo.d2 = r.ReadBytesShort()
	c.sigInfo.d2Key = r.ReadBytesShort()
	c.sigInfo.tgt = r.ReadBytesShort()
	c.sigInfo.srmToken = r.ReadBytesShort()
	c.sigInfo.t133 = r.ReadBytesShort()
	c.sigInfo.encryptedA1 = r.ReadBytesShort()
	c.sigInfo.wtSessionTicketKey = r.ReadBytesShort()
	c.OutGoingPacketSessionId = r.ReadBytesShort()
	// SystemDeviceInfo.TgtgtKey = r.ReadBytesShort()
	c.deviceInfo.TgtgtKey = r.ReadBytesShort()

	c.deviceInfo.Display = r.ReadBytesShort()
	c.deviceInfo.Product = r.ReadBytesShort()
	c.deviceInfo.Device = r.ReadBytesShort()
	c.deviceInfo.Board = r.ReadBytesShort()
	c.deviceInfo.Brand = r.ReadBytesShort()
	c.deviceInfo.Model = r.ReadBytesShort()
	c.deviceInfo.Bootloader = r.ReadBytesShort()
	c.deviceInfo.FingerPrint = r.ReadBytesShort()
	c.deviceInfo.BootId = r.ReadBytesShort()
	c.deviceInfo.ProcVersion = r.ReadBytesShort()
	c.deviceInfo.BaseBand = r.ReadBytesShort()
	c.deviceInfo.SimInfo = r.ReadBytesShort()
	c.deviceInfo.OSType = r.ReadBytesShort()
	c.deviceInfo.MacAddress = r.ReadBytesShort()
	c.deviceInfo.IpAddress = r.ReadBytesShort()
	c.deviceInfo.WifiBSSID = r.ReadBytesShort()
	c.deviceInfo.WifiSSID = r.ReadBytesShort()
	c.deviceInfo.IMSIMd5 = r.ReadBytesShort()
	c.deviceInfo.IMEI = r.ReadStringShort()
	c.deviceInfo.APN = r.ReadBytesShort()
	c.deviceInfo.VendorName = r.ReadBytesShort()
	c.deviceInfo.VendorOSName = r.ReadBytesShort()
	c.deviceInfo.AndroidId = r.ReadBytesShort()

	SystemDeviceInfo.GenNewGuid()
	SystemDeviceInfo.GenNewTgtgtKey()
}

func (c *QQClient) SaveState() []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt64(uint64(c.Uin))
		w.WriteBytesShort(c.sigInfo.d2)
		w.WriteBytesShort(c.sigInfo.d2Key)
		w.WriteBytesShort(c.sigInfo.tgt)
		w.WriteBytesShort(c.sigInfo.srmToken)
		w.WriteBytesShort(c.sigInfo.t133)
		w.WriteBytesShort(c.sigInfo.encryptedA1)
		w.WriteBytesShort(c.sigInfo.wtSessionTicketKey)
		w.WriteBytesShort(c.OutGoingPacketSessionId)
		w.WriteBytesShort(c.deviceInfo.TgtgtKey)

		w.WriteBytesShort(c.deviceInfo.Display)
		w.WriteBytesShort(c.deviceInfo.Product)
		w.WriteBytesShort(c.deviceInfo.Device)
		w.WriteBytesShort(c.deviceInfo.Board)
		w.WriteBytesShort(c.deviceInfo.Brand)
		w.WriteBytesShort(c.deviceInfo.Model)
		w.WriteBytesShort(c.deviceInfo.Bootloader)
		w.WriteBytesShort(c.deviceInfo.FingerPrint)
		w.WriteBytesShort(c.deviceInfo.BootId)
		w.WriteBytesShort(c.deviceInfo.ProcVersion)
		w.WriteBytesShort(c.deviceInfo.BaseBand)
		w.WriteBytesShort(c.deviceInfo.SimInfo)
		w.WriteBytesShort(c.deviceInfo.OSType)
		w.WriteBytesShort(c.deviceInfo.MacAddress)
		w.WriteBytesShort(c.deviceInfo.IpAddress)
		w.WriteBytesShort(c.deviceInfo.WifiBSSID)
		w.WriteBytesShort(c.deviceInfo.WifiSSID)
		w.WriteBytesShort(c.deviceInfo.IMSIMd5)
		w.WriteStringShort(c.deviceInfo.IMEI)
		w.WriteBytesShort(c.deviceInfo.APN)
		w.WriteBytesShort(c.deviceInfo.VendorName)
		w.WriteBytesShort(c.deviceInfo.VendorOSName)
		w.WriteBytesShort(c.deviceInfo.AndroidId)
	})
}

func (c *QQClient) FetchQRCode() (*QRCodeLoginResponse, error) {
	if c.Online {
		return nil, ErrAlreadyOnline
	}
	err := c.connectFastest()
	if err != nil {
		return nil, err
	}
	i, err := c.sendAndWait(c.buildQRCodeFetchRequestPacket())
	if err != nil {
		return nil, errors.Wrap(err, "fetch qrcode error")
	}
	return i.(*QRCodeLoginResponse), nil
}

func (c *QQClient) QueryQRCodeStatus(sig []byte) (*QRCodeLoginResponse, error) {
	i, err := c.sendAndWait(c.buildQRCodeResultQueryRequestPacket(sig))
	if err != nil {
		return nil, errors.Wrap(err, "query result error")
	}
	return i.(*QRCodeLoginResponse), nil
}

func (c *QQClient) QRCodeLogin(info *QRCodeLoginInfo) (*LoginResponse, error) {
	i, err := c.sendAndWait(c.buildQRCodeLoginPacket(info.tmpPwd, info.tmpNoPicSig, info.tgtQR))
	if err != nil {
		return nil, errors.Wrap(err, "qrcode login error")
	}
	rsp := i.(LoginResponse)
	if rsp.Success {
		_ = c.init(false)
	}
	return &rsp, nil
}

// SubmitCaptcha send captcha to server
func (c *QQClient) SubmitCaptcha(result string, sign []byte) (*LoginResponse, error) {
	seq, packet := c.buildCaptchaPacket(result, sign)
	rsp, err := c.sendAndWait(seq, packet)
	if err != nil {
		c.Disconnect()
		return nil, err
	}
	l := rsp.(LoginResponse)
	if l.Success {
		_ = c.init(false)
	}
	return &l, nil
}

func (c *QQClient) SubmitTicket(ticket string) (*LoginResponse, error) {
	seq, packet := c.buildTicketSubmitPacket(ticket)
	rsp, err := c.sendAndWait(seq, packet)
	if err != nil {
		c.Disconnect()
		return nil, err
	}
	l := rsp.(LoginResponse)
	if l.Success {
		_ = c.init(false)
	}
	return &l, nil
}

func (c *QQClient) SubmitSMS(code string) (*LoginResponse, error) {
	rsp, err := c.sendAndWait(c.buildSMSCodeSubmitPacket(code))
	if err != nil {
		c.Disconnect()
		return nil, err
	}
	l := rsp.(LoginResponse)
	if l.Success {
		_ = c.init(false)
	}
	return &l, nil
}

func (c *QQClient) RequestSMS() bool {
	rsp, err := c.sendAndWait(c.buildSMSRequestPacket())
	if err != nil {
		c.Error("request sms error: %v", err)
		return false
	}
	return rsp.(LoginResponse).Error == SMSNeededError
}

func (c *QQClient) init(tokenLogin bool) error {
	if len(c.g) == 0 {
		c.Warning("device lock is disable. http api may fail.")
	}
	if tokenLogin {
		err := func() error {

			notify := make(chan struct{})
			defer c.waitPacket("StatSvc.ReqMSFOffline", func(i interface{}, err error) {
				notify <- struct{}{}
			})()
			defer c.waitPacket("MessageSvc.PushForceOffline", func(i interface{}, err error) {
				notify <- struct{}{}
			})()
			if err := c.registerClient(); err != nil {
				return errors.Wrap(err, "register error")
			}
			select {
			case <-notify:
				return errors.New("token failed")
			case <-time.After(time.Second):
				return nil
			}
		}()
		if err != nil {
			return err
		}
	}
	c.groupSysMsgCache, _ = c.GetGroupSystemMessages()
	if !c.heartbeatEnabled {
		go c.doHeartbeat()
	}
	_ = c.RefreshStatus()
	if c.version.Protocol == QiDian {
		_, _ = c.sendAndWait(c.buildLoginExtraPacket())     // 小登录
		_, _ = c.sendAndWait(c.buildConnKeyRequestPacket()) // big data key 如果等待 config push 的话时间来不及
	}
	seq, pkt := c.buildGetMessageRequestPacket(msg.SyncFlag_START, time.Now().Unix())
	_, _ = c.sendAndWait(seq, pkt, requestParams{"used_reg_proxy": true, "init": true})
	c.syncChannelFirstView()
	return nil
}

func (c *QQClient) GenToken() []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		w.WriteUInt64(uint64(c.Uin))
		w.WriteBytesShort(c.sigInfo.d2)
		w.WriteBytesShort(c.sigInfo.d2Key)
		w.WriteBytesShort(c.sigInfo.tgt)
		w.WriteBytesShort(c.sigInfo.srmToken)
		w.WriteBytesShort(c.sigInfo.t133)
		w.WriteBytesShort(c.sigInfo.encryptedA1)
		w.WriteBytesShort(c.sigInfo.wtSessionTicketKey)
		w.WriteBytesShort(c.OutGoingPacketSessionId)
		w.WriteBytesShort(c.deviceInfo.TgtgtKey)
	})
}

func (c *QQClient) SetOnlineStatus(s UserOnlineStatus) {
	if s < 1000 {
		_, _ = c.sendAndWait(c.buildStatusSetPacket(int32(s), 0))
		return
	}
	_, _ = c.sendAndWait(c.buildStatusSetPacket(11, int32(s)))
}

func (c *QQClient) GetWordSegmentation(text string) ([]string, error) {
	rsp, err := c.sendAndWait(c.buildWordSegmentationPacket([]byte(text)))
	if err != nil {
		return nil, err
	}
	if data, ok := rsp.([][]byte); ok {
		var ret []string
		for _, val := range data {
			ret = append(ret, string(val))
		}
		return ret, nil
	}
	return nil, errors.New("decode error")
}

func (c *QQClient) GetSummaryInfo(target int64) (*SummaryCardInfo, error) {
	rsp, err := c.sendAndWait(c.buildSummaryCardRequestPacket(target))
	if err != nil {
		return nil, err
	}
	return rsp.(*SummaryCardInfo), nil
}

// ReloadFriendList refresh QQClient.FriendList field via GetFriendList()
func (c *QQClient) ReloadFriendList() error {
	rsp, err := c.GetFriendList()
	if err != nil {
		return err
	}
	c.FriendList = rsp.List
	return nil
}

// GetFriendList
// 当使用普通QQ时: 请求好友列表
// 当使用企点QQ时: 请求外部联系人列表
func (c *QQClient) GetFriendList() (*FriendListResponse, error) {
	if c.version.Protocol == QiDian {
		rsp, err := c.getQiDianAddressDetailList()
		if err != nil {
			return nil, err
		}
		return &FriendListResponse{TotalCount: int32(len(rsp)), List: rsp}, nil
	}
	curFriendCount := 0
	r := &FriendListResponse{}
	for {
		rsp, err := c.sendAndWait(c.buildFriendGroupListRequestPacket(int16(curFriendCount), 150, 0, 0))
		if err != nil {
			return nil, err
		}
		list := rsp.(*FriendListResponse)
		r.TotalCount = list.TotalCount
		r.List = append(r.List, list.List...)
		curFriendCount += len(list.List)
		if int32(len(r.List)) >= r.TotalCount {
			break
		}
	}
	return r, nil
}

func (c *QQClient) GetForwardMessage(resID string) *message.ForwardMessage {
	m := c.DownloadForwardMessage(resID)
	if m == nil {
		return nil
	}
	var (
		item *msg.PbMultiMsgItem
		ret  = &message.ForwardMessage{Nodes: []*message.ForwardNode{}}
	)
	for _, iter := range m.Items {
		if iter.GetFileName() == m.FileName {
			item = iter
		}
	}
	if item == nil {
		return nil
	}
	for _, m := range item.GetBuffer().GetMsg() {
		ret.Nodes = append(ret.Nodes, &message.ForwardNode{
			SenderId: m.Head.GetFromUin(),
			SenderName: func() string {
				if m.Head.GetMsgType() == 82 && m.Head.GroupInfo != nil {
					return m.Head.GroupInfo.GetGroupCard()
				}
				return m.Head.GetFromNick()
			}(),
			Time:    m.Head.GetMsgTime(),
			Message: message.ParseMessageElems(m.Body.RichText.Elems),
		})
	}
	return ret
}

func (c *QQClient) DownloadForwardMessage(resId string) *message.ForwardElement {
	i, err := c.sendAndWait(c.buildMultiApplyDownPacket(resId))
	if err != nil {
		return nil
	}
	multiMsg := i.(*msg.PbMultiMsgTransmit)
	if multiMsg.GetPbItemList() == nil {
		return nil
	}
	var pv string
	for i := 0; i < int(math.Min(4, float64(len(multiMsg.GetMsg())))); i++ {
		m := multiMsg.Msg[i]
		pv += fmt.Sprintf(`<title size="26" color="#777777">%s: %s</title>`,
			func() string {
				if m.Head.GetMsgType() == 82 && m.Head.GroupInfo != nil {
					return m.Head.GroupInfo.GetGroupCard()
				}
				return m.Head.GetFromNick()
			}(),
			message.ToReadableString(
				message.ParseMessageElems(multiMsg.Msg[i].GetBody().GetRichText().Elems),
			),
		)
	}
	return genForwardTemplate(
		resId, pv, "群聊的聊天记录", "[聊天记录]", "聊天记录",
		fmt.Sprintf("查看 %d 条转发消息", len(multiMsg.GetMsg())),
		time.Now().UnixNano(),
		multiMsg.GetPbItemList(),
	)
}

func (c *QQClient) SendGroupPoke(groupCode, target int64) {
	_, _ = c.sendAndWait(c.buildGroupPokePacket(groupCode, target))
}

func (c *QQClient) SendFriendPoke(target int64) {
	_, _ = c.sendAndWait(c.buildFriendPokePacket(target))
}

func (c *QQClient) ReloadGroupList() error {
	c.groupListLock.Lock()
	defer c.groupListLock.Unlock()
	list, err := c.GetGroupList()
	if err != nil {
		return err
	}
	c.GroupList = list
	return nil
}

func (c *QQClient) GetGroupList() ([]*GroupInfo, error) {
	rsp, err := c.sendAndWait(c.buildGroupListRequestPacket(EmptyBytes))
	if err != nil {
		return nil, err
	}
	r := rsp.([]*GroupInfo)
	wg := sync.WaitGroup{}
	batch := 50
	for i := 0; i < len(r); i += batch {
		k := i + batch
		if k > len(r) {
			k = len(r)
		}
		wg.Add(k - i)
		for j := i; j < k; j++ {
			go func(g *GroupInfo, wg *sync.WaitGroup) {
				defer wg.Done()
				m, err := c.GetGroupMembers(g)
				if err != nil {
					return
				}
				g.Members = m
			}(r[j], &wg)
		}
		wg.Wait()
	}
	return r, nil
}

func (c *QQClient) GetGroupMembers(group *GroupInfo) ([]*GroupMemberInfo, error) {
	var nextUin int64
	var list []*GroupMemberInfo
	for {
		data, err := c.sendAndWait(c.buildGroupMemberListRequestPacket(group.Uin, group.Code, nextUin))
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, errors.New("group member list unavailable: rsp is nil")
		}
		rsp := data.(*groupMemberListResponse)
		nextUin = rsp.NextUin
		for _, m := range rsp.list {
			m.Group = group
			if m.Uin == group.OwnerUin {
				m.Permission = Owner
			}
		}
		list = append(list, rsp.list...)
		if nextUin == 0 {
			sort.Slice(list, func(i, j int) bool {
				return list[i].Uin < list[j].Uin
			})
			return list, nil
		}
	}
}

func (c *QQClient) GetMemberInfo(groupCode, memberUin int64) (*GroupMemberInfo, error) {
	info, err := c.sendAndWait(c.buildGroupMemberInfoRequestPacket(groupCode, memberUin))
	if err != nil {
		return nil, err
	}
	return info.(*GroupMemberInfo), nil
}

func (c *QQClient) FindFriend(uin int64) *FriendInfo {
	for _, t := range c.FriendList {
		f := t
		if f.Uin == uin {
			return f
		}
	}
	return nil
}

func (c *QQClient) DeleteFriend(uin int64) error {
	if c.FindFriend(uin) == nil {
		return errors.New("friend not found")
	}
	_, err := c.sendAndWait(c.buildFriendDeletePacket(uin))
	return errors.Wrap(err, "delete friend error")
}

func (c *QQClient) FindGroupByUin(uin int64) *GroupInfo {
	for _, g := range c.GroupList {
		f := g
		if f.Uin == uin {
			return f
		}
	}
	return nil
}

func (c *QQClient) FindGroup(code int64) *GroupInfo {
	for _, g := range c.GroupList {
		if g.Code == code {
			return g
		}
	}
	return nil
}

func (c *QQClient) SolveGroupJoinRequest(i interface{}, accept, block bool, reason string) {
	if accept {
		block = false
		reason = ""
	}

	switch req := i.(type) {
	case *UserJoinGroupRequest:
		_, pkt := c.buildSystemMsgGroupActionPacket(req.RequestId, req.RequesterUin, req.GroupCode, func() int32 {
			if req.Suspicious {
				return 2
			} else {
				return 1
			}
		}(), false, accept, block, reason)
		_ = c.sendPacket(pkt)
	case *GroupInvitedRequest:
		_, pkt := c.buildSystemMsgGroupActionPacket(req.RequestId, req.InvitorUin, req.GroupCode, 1, true, accept, block, reason)
		_ = c.sendPacket(pkt)
	}
}

func (c *QQClient) SolveFriendRequest(req *NewFriendRequest, accept bool) {
	_, pkt := c.buildSystemMsgFriendActionPacket(req.RequestId, req.RequesterUin, accept)
	_ = c.sendPacket(pkt)
}

func (c *QQClient) getSKey() string {
	if c.sigInfo.sKeyExpiredTime < time.Now().Unix() && len(c.g) > 0 {
		c.Debug("skey expired. refresh...")
		_, _ = c.sendAndWait(c.buildRequestTgtgtNopicsigPacket())
	}
	return string(c.sigInfo.sKey)
}

func (c *QQClient) getCookies() string {
	return fmt.Sprintf("uin=o%d; skey=%s;", c.Uin, c.getSKey())
}

func (c *QQClient) getCookiesWithDomain(domain string) string {
	cookie := c.getCookies()

	if psKey, ok := c.sigInfo.psKeyMap[domain]; ok {
		return fmt.Sprintf("%s p_uin=o%d; p_skey=%s;", cookie, c.Uin, psKey)
	} else {
		return cookie
	}
}

func (c *QQClient) getCSRFToken() int {
	accu := 5381
	for _, b := range []byte(c.getSKey()) {
		accu = accu + (accu << 5) + int(b)
	}
	return 2147483647 & accu
}

func (c *QQClient) editMemberCard(groupCode, memberUin int64, card string) {
	_, _ = c.sendAndWait(c.buildEditGroupTagPacket(groupCode, memberUin, card))
}

func (c *QQClient) editMemberSpecialTitle(groupCode, memberUin int64, title string) {
	_, _ = c.sendAndWait(c.buildEditSpecialTitlePacket(groupCode, memberUin, title))
}

func (c *QQClient) setGroupAdmin(groupCode, memberUin int64, flag bool) {
	_, _ = c.sendAndWait(c.buildGroupAdminSetPacket(groupCode, memberUin, flag))
}

func (c *QQClient) updateGroupName(groupCode int64, newName string) {
	_, _ = c.sendAndWait(c.buildGroupNameUpdatePacket(groupCode, newName))
}

func (c *QQClient) updateGroupMemo(groupCode int64, newMemo string) {
	_, _ = c.sendAndWait(c.buildGroupMemoUpdatePacket(groupCode, newMemo))
}

func (c *QQClient) groupMuteAll(groupCode int64, mute bool) {
	_, _ = c.sendAndWait(c.buildGroupMuteAllPacket(groupCode, mute))
}

func (c *QQClient) groupMute(groupCode, memberUin int64, time uint32) {
	_, _ = c.sendAndWait(c.buildGroupMutePacket(groupCode, memberUin, time))
}

func (c *QQClient) quitGroup(groupCode int64) {
	_, _ = c.sendAndWait(c.buildQuitGroupPacket(groupCode))
}

func (c *QQClient) kickGroupMember(groupCode, memberUin int64, msg string, block bool) {
	_, _ = c.sendAndWait(c.buildGroupKickPacket(groupCode, memberUin, msg, block))
}

func (g *GroupInfo) removeMember(uin int64) {
	g.Update(func(info *GroupInfo) {
		i := sort.Search(len(info.Members), func(i int) bool {
			return info.Members[i].Uin >= uin
		})
		if i >= len(info.Members) || info.Members[i].Uin != uin { // not found
			return
		}
		info.Members = append(info.Members[:i], info.Members[i+1:]...)
	})
}

func (c *QQClient) SetCustomServer(servers []*net.TCPAddr) {
	c.servers = append(servers, c.servers...)
}

func (c *QQClient) SendGroupGift(groupCode, uin uint64, gift message.GroupGift) {
	_, packet := c.sendGroupGiftPacket(groupCode, uin, gift)
	_ = c.sendPacket(packet)
}

func (c *QQClient) registerClient() error {
	_, err := c.sendAndWait(c.buildClientRegisterPacket())
	if err == nil {
		c.Online = true
	}
	return err
}

func (c *QQClient) nextSeq() uint16 {
	return uint16(atomic.AddInt32(&c.SequenceId, 1) & 0x7FFF)
}

func (c *QQClient) nextPacketSeq() int32 {
	return atomic.AddInt32(&c.requestPacketRequestID, 2)
}

func (c *QQClient) nextGroupSeq() int32 {
	return atomic.AddInt32(&c.groupSeq, 2)
}

func (c *QQClient) nextFriendSeq() int32 {
	return atomic.AddInt32(&c.friendSeq, 1)
}

func (c *QQClient) nextGroupDataTransSeq() int32 {
	return atomic.AddInt32(&c.groupDataTransSeq, 2)
}

func (c *QQClient) nextHighwayApplySeq() int32 {
	return atomic.AddInt32(&c.highwayApplyUpSeq, 2)
}

func (c *QQClient) doHeartbeat() {
	// 不需要atomic/锁
	if c.heartbeatEnabled {
		return
	}
	c.heartbeatEnabled = true
	times := 0
	ticker := time.NewTicker(time.Second * 30)
	for range ticker.C {
		if !c.Online {
			ticker.Stop()
			c.heartbeatEnabled = false
			return
		}
		seq := c.nextSeq()
		sso := packets.BuildSsoPacket(seq, c.version.AppId, c.version.SubAppId, "Heartbeat.Alive", c.deviceInfo.IMEI, []byte{}, c.OutGoingPacketSessionId, []byte{}, c.ksid)
		packet := packets.BuildLoginPacket(c.Uin, 0, []byte{}, sso, []byte{})
		_, err := c.sendAndWait(seq, packet)
		if err != nil {
			continue
		}
		times++
		if times >= 7 {
			_ = c.registerClient()
			times = 0
		}
	}
}

func (c *QQClient) Error(msg string, args ...interface{}) {
	c.Logger.Error(c, msg, args...)
	/*	c.LogEventHandler(&LogEvent{
		Type:    "ERROR",
		Message: fmt.Sprintf(msg, args...),
	})*/
}

func (c *QQClient) Warning(msg string, args ...interface{}) {
	c.Logger.Warning(c, msg, args...)
	/*	c.LogEventHandler(&LogEvent{
		Type:    "WARNING",
		Message: fmt.Sprintf(msg, args...),
	})*/
}

func (c *QQClient) Info(msg string, args ...interface{}) {
	c.Logger.Info(c, msg, args...)
	/*	c.dispatchLogEvent(&LogEvent{
		Type:    "INFO",
		Message: fmt.Sprintf(msg, args...),
	})*/
}

func (c *QQClient) Debug(msg string, args ...interface{}) {
	c.Logger.Debug(c, msg, args...)
	/*	c.dispatchLogEvent(&LogEvent{
		Type:    "DEBUG",
		Message: fmt.Sprintf(msg, args...),
	})*/
}

func (c *QQClient) Trace(msg string, args ...interface{}) {
	c.Logger.Trace(c, msg, args...)
	/*	c.dispatchLogEvent(&LogEvent{
		Type:    "TRACE",
		Message: fmt.Sprintf(msg, args...),
	})*/
}
