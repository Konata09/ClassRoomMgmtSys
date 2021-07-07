package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	strconv "strconv"
	"strings"
)

var salt = "ustsdnicm3002"

func getPasswordMD5(password string) string {
	data := []byte(password + salt)
	encodePass := md5.Sum(data)
	return hex.EncodeToString(encodePass[:])
}

func checkPassword(passIn string, passMD5 string) bool {
	data := []byte(passIn + salt)
	encodedPass := md5.Sum(data)
	return hex.EncodeToString(encodedPass[:]) == passMD5
}

func trimMACtoShow(in string) string {
	if in == "" {
		return in
	}
	s := strings.ToUpper(strings.Replace(strings.Replace(in, ":", "", -1), "-", "", -1))
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", s[0:2], s[2:4], s[4:6], s[6:8], s[8:10], s[10:12])
}

func trimMACtoStor(in string) string {
	return strings.ToUpper(strings.Replace(strings.Replace(in, ":", "", -1), "-", "", -1))
}

func trimCommandToStor(in string) string {
	return strings.ToUpper(in)
}

func hexStringToByte(hexString string) ([]byte, error) {
	hex, err := hex.DecodeString(hexString)
	if err != nil {
		errMsg := fmt.Sprint(err)
		if strings.Contains(errMsg, "invalid byte") {
			return nil, errors.New("命令不能含有[0-9|A-F]以外的字符")
		} else if strings.Contains(errMsg, "odd length hex string") {
			return nil, errors.New("命令不能为奇数个字符")
		} else {
			return nil, err
		}
	}
	return hex, nil
}

func getSubnetBroadcast(ip string, subnetMask int) string {
	ips := strings.Split(ip, ".")
	var ipbin string
	for _, num := range ips {
		digi, _ := strconv.Atoi(num)
		b := strconv.FormatInt(int64(digi), 2)
		for strings.Count(b, "") <= 8 {
			b = "0" + b
		}
		ipbin += b
	}
	var bcastbin string
	for i := 0; i < 32; i++ {
		if i < subnetMask {
			bcastbin += ipbin[i : i+1]
		} else {
			bcastbin += "1"
		}
	}
	var r [4]int64
	r[0], _ = strconv.ParseInt(bcastbin[0:8], 2, 0)
	r[1], _ = strconv.ParseInt(bcastbin[8:16], 2, 0)
	r[2], _ = strconv.ParseInt(bcastbin[16:24], 2, 0)
	r[3], _ = strconv.ParseInt(bcastbin[24:32], 2, 0)
	bcastip := fmt.Sprintf("%d.%d.%d.%d", r[0], r[1], r[2], r[3])
	return bcastip
}

func getWolPayload(mac string) string {
	if strings.Contains(mac, ":") || strings.Contains(mac, "-") {
		mac = trimMACtoStor(mac)
	}
	var sb strings.Builder
	sb.WriteString("FFFFFFFFFFFF")
	for i := 0; i < 16; i++ {
		sb.WriteString(mac)
	}
	return sb.String()
}

func appendMsg(oldMsg *string, newMsg string) {
	if *oldMsg != "" {
		*oldMsg = fmt.Sprintf("%s|%s", *oldMsg, newMsg)
	} else {
		*oldMsg = newMsg
	}
}
