package clover

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// EncodeMsg implements msgp.Encodable
func (z *AccessToken) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(5)
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
	return
}

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
func (z *DefaultClient) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(6)
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
	for bzg := range z.GrantType {
		err = en.WriteString(z.GrantType[bzg])
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
	for bai := range z.Scope {
		err = en.WriteString(z.Scope[bai])
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
			for bzg := range z.GrantType {
				z.GrantType[bzg], err = dc.ReadString()
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
			for bai := range z.Scope {
				z.Scope[bai], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "RedirectURI":
			z.RedirectURI, err = dc.ReadString()
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
func (z *RefreshToken) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteMapHeader(5)
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
	err = en.WriteMapHeader(6)
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
	for ajw := range z.Scope {
		err = en.WriteString(z.Scope[ajw])
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
			for ajw := range z.Scope {
				z.Scope[ajw], err = dc.ReadString()
				if err != nil {
					return
				}
			}
		case "r":
			z.RedirectURI, err = dc.ReadString()
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

// MarshalMsg implements msgp.Marshaler
func (z *AccessToken) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 5)
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
	return
}

//UnmarshalMsg implements msgp.Unmarshaler
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

//UnmarshalMsg implements msgp.Unmarshaler
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

// MarshalMsg implements msgp.Marshaler
func (z *DefaultClient) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 6)
	o = msgp.AppendString(o, "ClientID")
	o = msgp.AppendString(o, z.ClientID)
	o = msgp.AppendString(o, "ClientSecret")
	o = msgp.AppendString(o, z.ClientSecret)
	o = msgp.AppendString(o, "GrantType")
	o = msgp.AppendArrayHeader(o, uint32(len(z.GrantType)))
	for bzg := range z.GrantType {
		o = msgp.AppendString(o, z.GrantType[bzg])
	}
	o = msgp.AppendString(o, "UserID")
	o = msgp.AppendString(o, z.UserID)
	o = msgp.AppendString(o, "Scope")
	o = msgp.AppendArrayHeader(o, uint32(len(z.Scope)))
	for bai := range z.Scope {
		o = msgp.AppendString(o, z.Scope[bai])
	}
	o = msgp.AppendString(o, "RedirectURI")
	o = msgp.AppendString(o, z.RedirectURI)
	return
}

//UnmarshalMsg implements msgp.Unmarshaler
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
			for bzg := range z.GrantType {
				z.GrantType[bzg], bts, err = msgp.ReadStringBytes(bts)
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
			for bai := range z.Scope {
				z.Scope[bai], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "RedirectURI":
			z.RedirectURI, bts, err = msgp.ReadStringBytes(bts)
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

func (z *DefaultClient) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 8 + msgp.StringPrefixSize + len(z.ClientID) + msgp.StringPrefixSize + 12 + msgp.StringPrefixSize + len(z.ClientSecret) + msgp.StringPrefixSize + 9 + msgp.ArrayHeaderSize
	for bzg := range z.GrantType {
		s += msgp.StringPrefixSize + len(z.GrantType[bzg])
	}
	s += msgp.StringPrefixSize + 6 + msgp.StringPrefixSize + len(z.UserID) + msgp.StringPrefixSize + 5 + msgp.ArrayHeaderSize
	for bai := range z.Scope {
		s += msgp.StringPrefixSize + len(z.Scope[bai])
	}
	s += msgp.StringPrefixSize + 11 + msgp.StringPrefixSize + len(z.RedirectURI)
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *RefreshToken) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 5)
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
	return
}

//UnmarshalMsg implements msgp.Unmarshaler
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
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *AuthorizeCode) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendMapHeader(o, 6)
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
	for ajw := range z.Scope {
		o = msgp.AppendString(o, z.Scope[ajw])
	}
	o = msgp.AppendString(o, "r")
	o = msgp.AppendString(o, z.RedirectURI)
	return
}

//UnmarshalMsg implements msgp.Unmarshaler
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
			for ajw := range z.Scope {
				z.Scope[ajw], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
			}
		case "r":
			z.RedirectURI, bts, err = msgp.ReadStringBytes(bts)
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

func (z *AuthorizeCode) Msgsize() (s int) {
	s = msgp.MapHeaderSize + msgp.StringPrefixSize + 2 + msgp.StringPrefixSize + len(z.Code) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.ClientID) + msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.UserID) + msgp.StringPrefixSize + 1 + msgp.Int64Size + msgp.StringPrefixSize + 1 + msgp.ArrayHeaderSize
	for ajw := range z.Scope {
		s += msgp.StringPrefixSize + len(z.Scope[ajw])
	}
	s += msgp.StringPrefixSize + 1 + msgp.StringPrefixSize + len(z.RedirectURI)
	return
}
