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
			for cmr := range z.Scope {
				z.Scope[cmr], err = dc.ReadString()
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
	for cmr := range z.Scope {
		err = en.WriteString(z.Scope[cmr])
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
	for cmr := range z.Scope {
		o = msgp.AppendString(o, z.Scope[cmr])
	}
	o = msgp.AppendString(o, "d")
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
			for cmr := range z.Scope {
				z.Scope[cmr], bts, err = msgp.ReadStringBytes(bts)
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

func (z *RefreshToken) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.RefreshToken) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.ClientID) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.UserID) + msgp.StringPrefixSize + 1 + msgp.Int64Size + msgp.StringPrefixSize + 1 + msgp.ArrayHeaderSize
	for cmr := range z.Scope {
		s += msgp.StringPrefixSize + len(z.Scope[cmr])
	}
	s += msgp.StringPrefixSize + 1 + msgp.MapHeaderSize
	if z.Data != nil {
		for ajw, wht := range z.Data {
			_ = wht
			s += msgp.StringPrefixSize + len(ajw) + msgp.GuessSize(wht)
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
			for hct := range z.Scope {
				z.Scope[hct], err = dc.ReadString()
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
				var cua string
				var xhx interface{}
				cua, err = dc.ReadString()
				if err != nil {
					return
				}
				xhx, err = dc.ReadIntf()
				if err != nil {
					return
				}
				z.Data[cua] = xhx
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
	for hct := range z.Scope {
		err = en.WriteString(z.Scope[hct])
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
	for cua, xhx := range z.Data {
		err = en.WriteString(cua)
		if err != nil {
			return
		}
		err = en.WriteIntf(xhx)
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
	for hct := range z.Scope {
		o = msgp.AppendString(o, z.Scope[hct])
	}
	o = msgp.AppendString(o, "r")
	o = msgp.AppendString(o, z.RedirectURI)
	o = msgp.AppendString(o, "d")
	o = msgp.AppendMapHeader(o, uint32(len(z.Data)))
	for cua, xhx := range z.Data {
		o = msgp.AppendString(o, cua)
		o, err = msgp.AppendIntf(o, xhx)
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
			for hct := range z.Scope {
				z.Scope[hct], bts, err = msgp.ReadStringBytes(bts)
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
				var cua string
				var xhx interface{}
				msz--
				cua, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				xhx, bts, err = msgp.ReadIntfBytes(bts)
				if err != nil {
					return
				}
				z.Data[cua] = xhx
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
	for hct := range z.Scope {
		s += msgp.StringPrefixSize + len(z.Scope[hct])
	}
	s += msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.RedirectURI) + msgp.StringPrefixSize + 1 + msgp.MapHeaderSize
	if z.Data != nil {
		for cua, xhx := range z.Data {
			_ = xhx
			s += msgp.StringPrefixSize + len(cua) + msgp.GuessSize(xhx)
		}
	}
	return
}
