package clover

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *AccessToken) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "a":
			z.AccessToken, err = dc.ReadString()
			if err != nil {
				return
			}
		case "c":
			z.ClientID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "u":
			z.UserID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "e":
			z.Expires, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "s":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for xvk := range z.Scope {
				z.Scope[xvk], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "d":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				msz--
				var bzg string
				var bai interface{}
				bzg, err = dc.ReadString()
				if err != nil {
					return
				}
				bai, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Data[bzg] = bai
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *AccessToken) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(6)
	if err != nil {
		return
	}
	err = en.WriteString("a")
	if err != nil {
		return
	}
	err = en.WriteString(z.AccessToken)
	if err != nil {
		return
	}
	err = en.WriteString("c")
	if err != nil {
		return
	}
	err = en.WriteString(z.ClientID)
	if err != nil {
		return
	}
	err = en.WriteString("u")
	if err != nil {
		return
	}
	err = en.WriteString(z.UserID)
	if err != nil {
		return
	}
	err = en.WriteString("e")
	if err != nil {
		return
	}
	err = en.WriteInt64(z.Expires)
	if err != nil {
		return
	}
	err = en.WriteString("s")
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Scope)))
	if err != nil {
		return
	}
	for xvk := range z.Scope {
		err = en.WriteString(z.Scope[xvk])
		if err != nil {
			return
		}
	}
	err = en.WriteString("d")
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Data)))
	if err != nil {
		return
	}
	for bzg, bai := range z.Data {
		err = en.WriteString(bzg)
		if err != nil {
			return
		}
		err = en.WriteIntf(bai)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *AccessToken) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 6)
	o = msgp.AppendString(o, "a")
	o = msgp.AppendString(o, z.AccessToken)
	o = msgp.AppendString(o, "c")
	o = msgp.AppendString(o, z.ClientID)
	o = msgp.AppendString(o, "u")
	o = msgp.AppendString(o, z.UserID)
	o = msgp.AppendString(o, "e")
	o = msgp.AppendInt64(o, z.Expires)
	o = msgp.AppendString(o, "s")
	o = msgp.AppendArrayHeader(o, uint32(len(z.Scope)))
	for xvk := range z.Scope {
		o = msgp.AppendString(o, z.Scope[xvk])
	}
	o = msgp.AppendString(o, "d")
	o = msgp.AppendMapHeader(o, uint32(len(z.Data)))
	for bzg, bai := range z.Data {
		o = msgp.AppendString(o, bzg)
		o, err = msgp.AppendIntf(o, bai)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AccessToken) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "a":
			z.AccessToken, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "c":
			z.ClientID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "u":
			z.UserID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "e":
			z.Expires, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "s":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for xvk := range z.Scope {
				z.Scope[xvk], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "d":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				var bzg string
				var bai interface{}
				msz--
				bzg, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				bai, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Data[bzg] = bai
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *AccessToken) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.AccessToken) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.ClientID) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.UserID) + msgp.StringPrefixSize + 1 + msgp.Int64Size + msgp.StringPrefixSize + 1 + msgp.ArrayHeaderSize
	for xvk := range z.Scope {
		s += msgp.StringPrefixSize + len(z.Scope[xvk])
	}
	s += msgp.StringPrefixSize + 1 + msgp.MapHeaderSize
	if z.Data != nil {
		for bzg, bai := range z.Data {
			_ = bai
			s += msgp.StringPrefixSize + len(bzg) + msgp.GuessSize(bai)
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PublicKey) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "c":
			z.ClientID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "pu":
			z.PublicKey, err = dc.ReadString()
			if err != nil {
				return
			}
		case "pr":
			z.PrivateKey, err = dc.ReadString()
			if err != nil {
				return
			}
		case "a":
			z.Algorithm, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *PublicKey) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(4)
	if err != nil {
		return
	}
	err = en.WriteString("c")
	if err != nil {
		return
	}
	err = en.WriteString(z.ClientID)
	if err != nil {
		return
	}
	err = en.WriteString("pu")
	if err != nil {
		return
	}
	err = en.WriteString(z.PublicKey)
	if err != nil {
		return
	}
	err = en.WriteString("pr")
	if err != nil {
		return
	}
	err = en.WriteString(z.PrivateKey)
	if err != nil {
		return
	}
	err = en.WriteString("a")
	if err != nil {
		return
	}
	err = en.WriteString(z.Algorithm)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PublicKey) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 4)
	o = msgp.AppendString(o, "c")
	o = msgp.AppendString(o, z.ClientID)
	o = msgp.AppendString(o, "pu")
	o = msgp.AppendString(o, z.PublicKey)
	o = msgp.AppendString(o, "pr")
	o = msgp.AppendString(o, z.PrivateKey)
	o = msgp.AppendString(o, "a")
	o = msgp.AppendString(o, z.Algorithm)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PublicKey) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "c":
			z.ClientID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "pu":
			z.PublicKey, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "pr":
			z.PrivateKey, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "a":
			z.Algorithm, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *PublicKey) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.ClientID) + msgp.StringPrefixSize + 2 + msgp.StringPrefixSize + len(z.PublicKey) + msgp.StringPrefixSize + 2 + msgp.StringPrefixSize + len(z.PrivateKey) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.Algorithm)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DefaultUser) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			z.ID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Username":
			z.Username, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Password":
			z.Password, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Scope":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for cmr := range z.Scope {
				z.Scope[cmr], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "Data":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				msz--
				var ajw string
				var wht interface{}
				ajw, err = dc.ReadString()
				if err != nil {
					return
				}
				wht, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Data[ajw] = wht
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *DefaultUser) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(5)
	if err != nil {
		return
	}
	err = en.WriteString("ID")
	if err != nil {
		return
	}
	err = en.WriteString(z.ID)
	if err != nil {
		return
	}
	err = en.WriteString("Username")
	if err != nil {
		return
	}
	err = en.WriteString(z.Username)
	if err != nil {
		return
	}
	err = en.WriteString("Password")
	if err != nil {
		return
	}
	err = en.WriteString(z.Password)
	if err != nil {
		return
	}
	err = en.WriteString("Scope")
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Scope)))
	if err != nil {
		return
	}
	for cmr := range z.Scope {
		err = en.WriteString(z.Scope[cmr])
		if err != nil {
			return
		}
	}
	err = en.WriteString("Data")
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Data)))
	if err != nil {
		return
	}
	for ajw, wht := range z.Data {
		err = en.WriteString(ajw)
		if err != nil {
			return
		}
		err = en.WriteIntf(wht)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *DefaultUser) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 5)
	o = msgp.AppendString(o, "ID")
	o = msgp.AppendString(o, z.ID)
	o = msgp.AppendString(o, "Username")
	o = msgp.AppendString(o, z.Username)
	o = msgp.AppendString(o, "Password")
	o = msgp.AppendString(o, z.Password)
	o = msgp.AppendString(o, "Scope")
	o = msgp.AppendArrayHeader(o, uint32(len(z.Scope)))
	for cmr := range z.Scope {
		o = msgp.AppendString(o, z.Scope[cmr])
	}
	o = msgp.AppendString(o, "Data")
	o = msgp.AppendMapHeader(o, uint32(len(z.Data)))
	for ajw, wht := range z.Data {
		o = msgp.AppendString(o, ajw)
		o, err = msgp.AppendIntf(o, wht)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DefaultUser) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			z.ID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Username":
			z.Username, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Password":
			z.Password, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Scope":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for cmr := range z.Scope {
				z.Scope[cmr], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "Data":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				var ajw string
				var wht interface{}
				msz--
				ajw, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				wht, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Data[ajw] = wht
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *DefaultUser) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 2 + msgp.StringPrefixSize + len(z.ID) + msgp.StringPrefixSize + 8 + msgp.StringPrefixSize + len(z.Username) + msgp.StringPrefixSize + 8 + msgp.StringPrefixSize + len(z.Password) + msgp.StringPrefixSize + 5 + msgp.ArrayHeaderSize
	for cmr := range z.Scope {
		s += msgp.StringPrefixSize + len(z.Scope[cmr])
	}
	s += msgp.StringPrefixSize + 4 + msgp.MapHeaderSize
	if z.Data != nil {
		for ajw, wht := range z.Data {
			_ = wht
			s += msgp.StringPrefixSize + len(ajw) + msgp.GuessSize(wht)
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *DefaultClient) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ClientID":
			z.ClientID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "ClientSecret":
			z.ClientSecret, err = dc.ReadString()
			if err != nil {
				return
			}
		case "GrantType":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.GrantType) >= int(xsz) {
				z.GrantType = z.GrantType[:xsz]
			} else {
				z.GrantType = make([]string, xsz)
			}
			for hct := range z.GrantType {
				z.GrantType[hct], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "UserID":
			z.UserID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Scope":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for cua := range z.Scope {
				z.Scope[cua], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "RedirectURI":
			z.RedirectURI, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Data":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				msz--
				var xhx string
				var lqf interface{}
				xhx, err = dc.ReadString()
				if err != nil {
					return
				}
				lqf, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Data[xhx] = lqf
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *DefaultClient) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(7)
	if err != nil {
		return
	}
	err = en.WriteString("ClientID")
	if err != nil {
		return
	}
	err = en.WriteString(z.ClientID)
	if err != nil {
		return
	}
	err = en.WriteString("ClientSecret")
	if err != nil {
		return
	}
	err = en.WriteString(z.ClientSecret)
	if err != nil {
		return
	}
	err = en.WriteString("GrantType")
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.GrantType)))
	if err != nil {
		return
	}
	for hct := range z.GrantType {
		err = en.WriteString(z.GrantType[hct])
		if err != nil {
			return
		}
	}
	err = en.WriteString("UserID")
	if err != nil {
		return
	}
	err = en.WriteString(z.UserID)
	if err != nil {
		return
	}
	err = en.WriteString("Scope")
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Scope)))
	if err != nil {
		return
	}
	for cua := range z.Scope {
		err = en.WriteString(z.Scope[cua])
		if err != nil {
			return
		}
	}
	err = en.WriteString("RedirectURI")
	if err != nil {
		return
	}
	err = en.WriteString(z.RedirectURI)
	if err != nil {
		return
	}
	err = en.WriteString("Data")
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Data)))
	if err != nil {
		return
	}
	for xhx, lqf := range z.Data {
		err = en.WriteString(xhx)
		if err != nil {
			return
		}
		err = en.WriteIntf(lqf)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *DefaultClient) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 7)
	o = msgp.AppendString(o, "ClientID")
	o = msgp.AppendString(o, z.ClientID)
	o = msgp.AppendString(o, "ClientSecret")
	o = msgp.AppendString(o, z.ClientSecret)
	o = msgp.AppendString(o, "GrantType")
	o = msgp.AppendArrayHeader(o, uint32(len(z.GrantType)))
	for hct := range z.GrantType {
		o = msgp.AppendString(o, z.GrantType[hct])
	}
	o = msgp.AppendString(o, "UserID")
	o = msgp.AppendString(o, z.UserID)
	o = msgp.AppendString(o, "Scope")
	o = msgp.AppendArrayHeader(o, uint32(len(z.Scope)))
	for cua := range z.Scope {
		o = msgp.AppendString(o, z.Scope[cua])
	}
	o = msgp.AppendString(o, "RedirectURI")
	o = msgp.AppendString(o, z.RedirectURI)
	o = msgp.AppendString(o, "Data")
	o = msgp.AppendMapHeader(o, uint32(len(z.Data)))
	for xhx, lqf := range z.Data {
		o = msgp.AppendString(o, xhx)
		o, err = msgp.AppendIntf(o, lqf)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *DefaultClient) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ClientID":
			z.ClientID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "ClientSecret":
			z.ClientSecret, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "GrantType":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.GrantType) >= int(xsz) {
				z.GrantType = z.GrantType[:xsz]
			} else {
				z.GrantType = make([]string, xsz)
			}
			for hct := range z.GrantType {
				z.GrantType[hct], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "UserID":
			z.UserID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Scope":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for cua := range z.Scope {
				z.Scope[cua], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "RedirectURI":
			z.RedirectURI, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Data":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				var xhx string
				var lqf interface{}
				msz--
				xhx, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				lqf, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Data[xhx] = lqf
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *DefaultClient) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 8 + msgp.StringPrefixSize + len(z.ClientID) + msgp.StringPrefixSize + 12 + msgp.StringPrefixSize + len(z.ClientSecret) + msgp.StringPrefixSize + 9 + msgp.ArrayHeaderSize
	for hct := range z.GrantType {
		s += msgp.StringPrefixSize + len(z.GrantType[hct])
	}
	s += msgp.StringPrefixSize + 6 + msgp.StringPrefixSize + len(z.UserID) + msgp.StringPrefixSize + 5 + msgp.ArrayHeaderSize
	for cua := range z.Scope {
		s += msgp.StringPrefixSize + len(z.Scope[cua])
	}
	s += msgp.StringPrefixSize + 11 + msgp.StringPrefixSize + len(z.RedirectURI) + msgp.StringPrefixSize + 4 + msgp.MapHeaderSize
	if z.Data != nil {
		for xhx, lqf := range z.Data {
			_ = lqf
			s += msgp.StringPrefixSize + len(xhx) + msgp.GuessSize(lqf)
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *RefreshToken) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "r":
			z.RefreshToken, err = dc.ReadString()
			if err != nil {
				return
			}
		case "a":
			z.ClientID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "u":
			z.UserID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "e":
			z.Expires, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "s":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for daf := range z.Scope {
				z.Scope[daf], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "d":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				msz--
				var pks string
				var jfb interface{}
				pks, err = dc.ReadString()
				if err != nil {
					return
				}
				jfb, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Data[pks] = jfb
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *RefreshToken) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(6)
	if err != nil {
		return
	}
	err = en.WriteString("r")
	if err != nil {
		return
	}
	err = en.WriteString(z.RefreshToken)
	if err != nil {
		return
	}
	err = en.WriteString("a")
	if err != nil {
		return
	}
	err = en.WriteString(z.ClientID)
	if err != nil {
		return
	}
	err = en.WriteString("u")
	if err != nil {
		return
	}
	err = en.WriteString(z.UserID)
	if err != nil {
		return
	}
	err = en.WriteString("e")
	if err != nil {
		return
	}
	err = en.WriteInt64(z.Expires)
	if err != nil {
		return
	}
	err = en.WriteString("s")
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Scope)))
	if err != nil {
		return
	}
	for daf := range z.Scope {
		err = en.WriteString(z.Scope[daf])
		if err != nil {
			return
		}
	}
	err = en.WriteString("d")
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Data)))
	if err != nil {
		return
	}
	for pks, jfb := range z.Data {
		err = en.WriteString(pks)
		if err != nil {
			return
		}
		err = en.WriteIntf(jfb)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *RefreshToken) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 6)
	o = msgp.AppendString(o, "r")
	o = msgp.AppendString(o, z.RefreshToken)
	o = msgp.AppendString(o, "a")
	o = msgp.AppendString(o, z.ClientID)
	o = msgp.AppendString(o, "u")
	o = msgp.AppendString(o, z.UserID)
	o = msgp.AppendString(o, "e")
	o = msgp.AppendInt64(o, z.Expires)
	o = msgp.AppendString(o, "s")
	o = msgp.AppendArrayHeader(o, uint32(len(z.Scope)))
	for daf := range z.Scope {
		o = msgp.AppendString(o, z.Scope[daf])
	}
	o = msgp.AppendString(o, "d")
	o = msgp.AppendMapHeader(o, uint32(len(z.Data)))
	for pks, jfb := range z.Data {
		o = msgp.AppendString(o, pks)
		o, err = msgp.AppendIntf(o, jfb)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *RefreshToken) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "r":
			z.RefreshToken, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "a":
			z.ClientID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "u":
			z.UserID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "e":
			z.Expires, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "s":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for daf := range z.Scope {
				z.Scope[daf], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "d":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				var pks string
				var jfb interface{}
				msz--
				pks, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				jfb, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Data[pks] = jfb
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *RefreshToken) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.RefreshToken) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.ClientID) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.UserID) + msgp.StringPrefixSize + 1 + msgp.Int64Size + msgp.StringPrefixSize + 1 + msgp.ArrayHeaderSize
	for daf := range z.Scope {
		s += msgp.StringPrefixSize + len(z.Scope[daf])
	}
	s += msgp.StringPrefixSize + 1 + msgp.MapHeaderSize
	if z.Data != nil {
		for pks, jfb := range z.Data {
			_ = jfb
			s += msgp.StringPrefixSize + len(pks) + msgp.GuessSize(jfb)
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *AuthorizeCode) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "co":
			z.Code, err = dc.ReadString()
			if err != nil {
				return
			}
		case "c":
			z.ClientID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "u":
			z.UserID, err = dc.ReadString()
			if err != nil {
				return
			}
		case "e":
			z.Expires, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "s":
			var xsz uint32
			xsz, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for cxo := range z.Scope {
				z.Scope[cxo], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "r":
			z.RedirectURI, err = dc.ReadString()
			if err != nil {
				return
			}
		case "d":
			var msz uint32
			msz, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				msz--
				var eff string
				var rsw interface{}
				eff, err = dc.ReadString()
				if err != nil {
					return
				}
				rsw, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Data[eff] = rsw
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *AuthorizeCode) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(7)
	if err != nil {
		return
	}
	err = en.WriteString("co")
	if err != nil {
		return
	}
	err = en.WriteString(z.Code)
	if err != nil {
		return
	}
	err = en.WriteString("c")
	if err != nil {
		return
	}
	err = en.WriteString(z.ClientID)
	if err != nil {
		return
	}
	err = en.WriteString("u")
	if err != nil {
		return
	}
	err = en.WriteString(z.UserID)
	if err != nil {
		return
	}
	err = en.WriteString("e")
	if err != nil {
		return
	}
	err = en.WriteInt64(z.Expires)
	if err != nil {
		return
	}
	err = en.WriteString("s")
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Scope)))
	if err != nil {
		return
	}
	for cxo := range z.Scope {
		err = en.WriteString(z.Scope[cxo])
		if err != nil {
			return
		}
	}
	err = en.WriteString("r")
	if err != nil {
		return
	}
	err = en.WriteString(z.RedirectURI)
	if err != nil {
		return
	}
	err = en.WriteString("d")
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Data)))
	if err != nil {
		return
	}
	for eff, rsw := range z.Data {
		err = en.WriteString(eff)
		if err != nil {
			return
		}
		err = en.WriteIntf(rsw)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *AuthorizeCode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 7)
	o = msgp.AppendString(o, "co")
	o = msgp.AppendString(o, z.Code)
	o = msgp.AppendString(o, "c")
	o = msgp.AppendString(o, z.ClientID)
	o = msgp.AppendString(o, "u")
	o = msgp.AppendString(o, z.UserID)
	o = msgp.AppendString(o, "e")
	o = msgp.AppendInt64(o, z.Expires)
	o = msgp.AppendString(o, "s")
	o = msgp.AppendArrayHeader(o, uint32(len(z.Scope)))
	for cxo := range z.Scope {
		o = msgp.AppendString(o, z.Scope[cxo])
	}
	o = msgp.AppendString(o, "r")
	o = msgp.AppendString(o, z.RedirectURI)
	o = msgp.AppendString(o, "d")
	o = msgp.AppendMapHeader(o, uint32(len(z.Data)))
	for eff, rsw := range z.Data {
		o = msgp.AppendString(o, eff)
		o, err = msgp.AppendIntf(o, rsw)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AuthorizeCode) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "co":
			z.Code, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "c":
			z.ClientID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "u":
			z.UserID, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "e":
			z.Expires, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "s":
			var xsz uint32
			xsz, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Scope) >= int(xsz) {
				z.Scope = z.Scope[:xsz]
			} else {
				z.Scope = make([]string, xsz)
			}
			for cxo := range z.Scope {
				z.Scope[cxo], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "r":
			z.RedirectURI, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "d":
			var msz uint32
			msz, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Data == nil && msz > 0 {
				z.Data = make(map[string]interface{}, msz)
			} else if len(z.Data) > 0 {
				for key, _ := range z.Data {
					delete(z.Data, key)
				}
			}
			for msz > 0 {
				var eff string
				var rsw interface{}
				msz--
				eff, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				rsw, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Data[eff] = rsw
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *AuthorizeCode) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 2 + msgp.StringPrefixSize + len(z.Code) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.ClientID) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.UserID) + msgp.StringPrefixSize + 1 + msgp.Int64Size + msgp.StringPrefixSize + 1 + msgp.ArrayHeaderSize
	for cxo := range z.Scope {
		s += msgp.StringPrefixSize + len(z.Scope[cxo])
	}
	s += msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.RedirectURI) + msgp.StringPrefixSize + 1 + msgp.MapHeaderSize
	if z.Data != nil {
		for eff, rsw := range z.Data {
			_ = rsw
			s += msgp.StringPrefixSize + len(eff) + msgp.GuessSize(rsw)
		}
	}
	return
}
