package mail

import (
	"strconv"
	"strings"
	"time"
)

func MailTime(fileName string) time.Time {
	// ファイル名の先頭部分がUnix時間
	// 例: 1674617693.M958571P8888.localhost.localdomain,S=545,W=562:2,S
	//     -> 1674617693 がUnix時間
	unixtimePart := strings.Split(fileName, ".")[0]
	unixtime, err := strconv.ParseInt(unixtimePart, 10, 64)
	if err != nil {
		unixtime = 0
	}

	return time.Unix(unixtime, 0)
}
