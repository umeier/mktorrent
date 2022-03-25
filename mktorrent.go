package mktorrent

import (
	"crypto/sha1"
	"github.com/zeebo/bencode"
	"io"
	"time"
)

const pieceLen = 512000

type InfoDict struct {
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type Torrent struct {
	Info         InfoDict `bencode:"info"`
	AnnounceList []string `bencode:"announce-list,omitempty"`
	Announce     string   `bencode:"announce,omitempty"`
	CreationDate int64    `bencode:"creation date,omitempty"`
	Comment      string   `bencode:"comment,omitempty"`
	CreatedBy    string   `bencode:"created by,omitempty"`
	UrlList      []string `bencode:"url-list,omitempty"`
}

func (t *Torrent) Save(w io.Writer) error {
	enc := bencode.NewEncoder(w)
	return enc.Encode(t)
}

func hashPiece(b []byte) []byte {
	h := sha1.New()
	h.Write(b)
	return h.Sum(nil)
}
func MakeTorrent(r io.Reader, name string, ann []string, url []string) (*Torrent, error) {
	t := &Torrent{
		CreationDate: time.Now().Unix(),
		CreatedBy:    "mktorrent.go",
		Info: InfoDict{
			Name:        name,
			PieceLength: pieceLen,
		},
	}
	if len(ann) == 1 {
		t.Announce = ann[0]
	} else {
		for _, a := range ann {
			t.AnnounceList = append(t.AnnounceList, a)
		}
	}

	if len(url) > 0 {
		for _, u := range url {
			t.UrlList = append(t.UrlList, u)
		}
	}

	b := make([]byte, pieceLen)
	for {
		n, err := io.ReadFull(r, b)
		if err != nil && err != io.ErrUnexpectedEOF {
			return nil, err
		}
		if err == io.ErrUnexpectedEOF {
			b = b[:n]
			t.Info.Pieces += string(hashPiece(b))
			t.Info.Length += n
			break
		} else if n == pieceLen {
			t.Info.Pieces += string(hashPiece(b))
			t.Info.Length += n
		} else {
			panic("short read!")
		}
	}
	return t, nil
}
