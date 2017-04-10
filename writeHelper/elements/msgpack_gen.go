package elements

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ManageChainMetaDataBytes) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "website":
			z.Website, err = dc.ReadBytes(z.Website)
			if err != nil {
				return
			}
		case "longdesc":
			z.LongDescription, err = dc.ReadBytes(z.LongDescription)
			if err != nil {
				return
			}
		case "shortdesc":
			z.ShortDescription, err = dc.ReadBytes(z.ShortDescription)
			if err != nil {
				return
			}
		case "playlist":
			z.Playlist, err = dc.ReadBytes(z.Playlist)
			if err != nil {
				return
			}
		case "thumbnail":
			z.Thumbnail, err = dc.ReadBytes(z.Thumbnail)
			if err != nil {
				return
			}
		case "banner":
			z.Banner, err = dc.ReadBytes(z.Banner)
			if err != nil {
				return
			}
		case "chantags":
			z.ChannelTags, err = dc.ReadBytes(z.ChannelTags)
			if err != nil {
				return
			}
		case "sugchans":
			z.SuggestedChannels, err = dc.ReadBytes(z.SuggestedChannels)
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
func (z *ManageChainMetaDataBytes) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 8
	// write "website"
	err = en.Append(0x88, 0xa7, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Website)
	if err != nil {
		return
	}
	// write "longdesc"
	err = en.Append(0xa8, 0x6c, 0x6f, 0x6e, 0x67, 0x64, 0x65, 0x73, 0x63)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.LongDescription)
	if err != nil {
		return
	}
	// write "shortdesc"
	err = en.Append(0xa9, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x64, 0x65, 0x73, 0x63)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.ShortDescription)
	if err != nil {
		return
	}
	// write "playlist"
	err = en.Append(0xa8, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Playlist)
	if err != nil {
		return
	}
	// write "thumbnail"
	err = en.Append(0xa9, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Thumbnail)
	if err != nil {
		return
	}
	// write "banner"
	err = en.Append(0xa6, 0x62, 0x61, 0x6e, 0x6e, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.Banner)
	if err != nil {
		return
	}
	// write "chantags"
	err = en.Append(0xa8, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x61, 0x67, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.ChannelTags)
	if err != nil {
		return
	}
	// write "sugchans"
	err = en.Append(0xa8, 0x73, 0x75, 0x67, 0x63, 0x68, 0x61, 0x6e, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.SuggestedChannels)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ManageChainMetaDataBytes) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 8
	// string "website"
	o = append(o, 0x88, 0xa7, 0x77, 0x65, 0x62, 0x73, 0x69, 0x74, 0x65)
	o = msgp.AppendBytes(o, z.Website)
	// string "longdesc"
	o = append(o, 0xa8, 0x6c, 0x6f, 0x6e, 0x67, 0x64, 0x65, 0x73, 0x63)
	o = msgp.AppendBytes(o, z.LongDescription)
	// string "shortdesc"
	o = append(o, 0xa9, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x64, 0x65, 0x73, 0x63)
	o = msgp.AppendBytes(o, z.ShortDescription)
	// string "playlist"
	o = append(o, 0xa8, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74)
	o = msgp.AppendBytes(o, z.Playlist)
	// string "thumbnail"
	o = append(o, 0xa9, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c)
	o = msgp.AppendBytes(o, z.Thumbnail)
	// string "banner"
	o = append(o, 0xa6, 0x62, 0x61, 0x6e, 0x6e, 0x65, 0x72)
	o = msgp.AppendBytes(o, z.Banner)
	// string "chantags"
	o = append(o, 0xa8, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x61, 0x67, 0x73)
	o = msgp.AppendBytes(o, z.ChannelTags)
	// string "sugchans"
	o = append(o, 0xa8, 0x73, 0x75, 0x67, 0x63, 0x68, 0x61, 0x6e, 0x73)
	o = msgp.AppendBytes(o, z.SuggestedChannels)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ManageChainMetaDataBytes) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "website":
			z.Website, bts, err = msgp.ReadBytesBytes(bts, z.Website)
			if err != nil {
				return
			}
		case "longdesc":
			z.LongDescription, bts, err = msgp.ReadBytesBytes(bts, z.LongDescription)
			if err != nil {
				return
			}
		case "shortdesc":
			z.ShortDescription, bts, err = msgp.ReadBytesBytes(bts, z.ShortDescription)
			if err != nil {
				return
			}
		case "playlist":
			z.Playlist, bts, err = msgp.ReadBytesBytes(bts, z.Playlist)
			if err != nil {
				return
			}
		case "thumbnail":
			z.Thumbnail, bts, err = msgp.ReadBytesBytes(bts, z.Thumbnail)
			if err != nil {
				return
			}
		case "banner":
			z.Banner, bts, err = msgp.ReadBytesBytes(bts, z.Banner)
			if err != nil {
				return
			}
		case "chantags":
			z.ChannelTags, bts, err = msgp.ReadBytesBytes(bts, z.ChannelTags)
			if err != nil {
				return
			}
		case "sugchans":
			z.SuggestedChannels, bts, err = msgp.ReadBytesBytes(bts, z.SuggestedChannels)
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

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ManageChainMetaDataBytes) Msgsize() (s int) {
	s = 1 + 8 + msgp.BytesPrefixSize + len(z.Website) + 9 + msgp.BytesPrefixSize + len(z.LongDescription) + 10 + msgp.BytesPrefixSize + len(z.ShortDescription) + 9 + msgp.BytesPrefixSize + len(z.Playlist) + 10 + msgp.BytesPrefixSize + len(z.Thumbnail) + 7 + msgp.BytesPrefixSize + len(z.Banner) + 9 + msgp.BytesPrefixSize + len(z.ChannelTags) + 9 + msgp.BytesPrefixSize + len(z.SuggestedChannels)
	return
}
