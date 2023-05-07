package parse

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type M3u8 struct {
	IsLive  bool // 是否有#EXT-X-ENDLIST
	Version int  // EXT-X-VERSION:version
	TimeMax int

	PlayList   []string
	TimeList   []float32
	TimeLength float32
	Key        Key
}
type Key struct {
	// 'AES-128' or 'NONE'
	//	如果加密方法为 NONE，则 URI 和 IV 属性不得存在
	METHOD string
	URI    string
	IV     string
}

func Parse(reader io.Reader) (*M3u8, error) {
	s := bufio.NewScanner(reader)
	var lines []string
	m3u8 := &M3u8{}
	m3u8.IsLive = true
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if i == 0 {
			if "#EXTM3U" != line {
				return nil, fmt.Errorf("无效的 M3U8，第 1 行缺少#EXTM3U")
			}
			continue
		}
		switch {
		case line == "":
			continue
		case strings.HasPrefix(line, "#EXT-X-VERSION:"):
			m3u8.Version, _ = strconv.Atoi(line[15:])
			continue
		case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
			m3u8.TimeMax, _ = strconv.Atoi(line[22:])
		case line == "#EXT-X-ENDLIST":
			m3u8.IsLive = false
		default:
			continue

		}
	}
	//	fmt.Println(lines)

	return m3u8, nil

}
