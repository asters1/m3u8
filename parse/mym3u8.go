package parse

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type M3u8 struct {
	IsLive  bool // 是否有#EXT-X-ENDLIST
	Version int  // EXT-X-VERSION:version
	TimeMax int

	PlayList   []string
	TimeList   []float64
	TimeLength float64
	Key        Key
}
type Key struct {
	// 'AES-128' or 'NONE'
	//	如果加密方法为 NONE，则 URI 和 IV 属性不得存在
	METHOD string
	URI    string
	IV     string
}

// 用于从一行中提取“键=值”参数的正则表达式模式
var linePattern = regexp.MustCompile(`([a-zA-Z-]+)=("[^"]+"|[^",]+)`)

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
		case strings.HasPrefix(line, "#EXT-X-KEY"):

			params := ParseLineParameters(line)
			if len(params) == 0 {
				return nil, fmt.Errorf("无效的 EXT-X-KEY: %s, 在第 %d 行", line, i+1)
			}
			if params["METHOD"] != "" {
				m3u8.Key.METHOD = params["METHOD"]
				if m3u8.Key.METHOD == "NONE" {
					m3u8.Key.URI = ""
					m3u8.Key.IV = ""
				} else {
					m3u8.Key.URI = params["URI"]
					m3u8.Key.IV = params["IV"]
				}
			} else {
				m3u8.Key.METHOD = "NONE"
				m3u8.Key.URI = ""
				m3u8.Key.IV = ""
			}
			continue
		case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
			m3u8.TimeMax, _ = strconv.Atoi(line[22:])
			continue
		case line == "#EXT-X-ENDLIST":
			m3u8.IsLive = false
		case strings.HasPrefix(line, "#EXTINF:"):
			//			fmt.Println(line[len(line)-1:])
			if line[len(line)-1:] == "," {
				line = line[:len(line)-1]
			}
			//			fmt.Println(line[8:])
			ftime, _ := strconv.ParseFloat(line[8:], 64)
			//			fmt.Println(ftime)
			m3u8.TimeList = append(m3u8.TimeList, ftime)
			continue
		case !strings.HasPrefix(line, "#"):
			m3u8.PlayList = append(m3u8.PlayList, line)
		default:
			continue
		}
	}
	sum := 0.0
	for i := 0; i < len(m3u8.TimeList); i++ {
		sum = sum + m3u8.TimeList[i]
	}
	if m3u8.Key.METHOD == "" {
		m3u8.Key.METHOD = "NONE"
		m3u8.Key.URI = ""
		m3u8.Key.IV = ""
	}
	// fmt.Println(sum)
	m3u8.TimeLength = sum
	//	fmt.Println(lines)
	return m3u8, nil

}
func ParseLineParameters(line string) map[string]string {
	r := linePattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	//	fmt.Println(params)
	return params
}
