package crypto

import (
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Mrs4s/MiraiGo/binary"
)

type EncryptECDH struct {
	InitialShareKey []byte
	PublicKey       []byte
	PublicKeyVer    uint16
}

type T133 []byte
type EncryptSession T133

const serverPublicKey = "04EBCA94D733E399B2DB96EACDD3F69A8BB0F74224E2B44E3357812211D2E62EFBC91BB553098E25E33A799ADC7F76FEB208DA7C6522CDB0719A305180CC54A82E"

func (e *EncryptECDH) InitWithDefaultKey() {
	e.generateKey(serverPublicKey)
}

func (e *EncryptECDH) generateKey(sPubKey string) {
	pub, err := hex.DecodeString(sPubKey)
	if err != nil {
		panic(err)
	}
	p256 := elliptic.P256()
	key, sx, sy, err := elliptic.GenerateKey(p256, rand.Reader)
	if err != nil {
		panic(err)
	}
	tx, ty := elliptic.Unmarshal(p256, pub)
	x, _ := p256.ScalarMult(tx, ty, key)
	hash := md5.Sum(x.Bytes()[:16])
	e.InitialShareKey = hash[:]
	e.PublicKey = elliptic.Marshal(p256, sx, sy)
}

func (e *EncryptECDH) DoEncrypt(d, k []byte) []byte {
	w := binary.NewWriter()
	w.WriteByte(0x02)
	w.WriteByte(0x01)
	w.Write(k)
	w.WriteUInt16(0x01_31)
	w.WriteUInt16(e.PublicKeyVer)
	w.WriteUInt16(uint16(len(e.PublicKey)))
	w.Write(e.PublicKey)
	w.EncryptAndWrite(e.InitialShareKey, d)
	return w.Bytes()
}

func (e *EncryptECDH) Id() byte {
	return 0x87
}

func NewEncryptSession(t133 []byte) *EncryptSession {
	return (*EncryptSession)(&t133)
}

func (e *EncryptSession) DoEncrypt(d, k []byte) []byte {
	return binary.NewWriterF(func(w *binary.Writer) {
		encrypt := binary.NewTeaCipher(k).Encrypt(d)
		w.WriteUInt16(uint16(len(*(*[]byte)(e))))
		w.Write(*(*[]byte)(e))
		w.Write(encrypt)
	})
}

func (e *EncryptSession) Id() byte {
	return 69
}

type pubKeyResp struct {
	Meta struct {
		PubKeyVer uint16 `json:"KeyVer"`
		PubKey    string `json:"PubKey"`
	} `json:"PubKeyMeta"`
}

// FetchPubKey 从服务器获取PubKey
func (e *EncryptECDH) FetchPubKey(uin int64) {
	resp, err := http.Get("https://keyrotate.qq.com/rotate_key?cipher_suite_ver=305&uin=" + strconv.FormatInt(uin, 10))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	pubKey := pubKeyResp{}
	err = json.NewDecoder(resp.Body).Decode(&pubKey)
	if err != nil {
		return
	}
	e.PublicKeyVer = pubKey.Meta.PubKeyVer
	e.generateKey(pubKey.Meta.PubKey) // todo check key sign
}
