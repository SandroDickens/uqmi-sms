package sms

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type SMSContext struct {
	Smsc        string `json:"smsc"`
	From        string `json:"sender"`
	Timestamp   string `json:"timestamp"`
	ConcatRef   int    `json:"concat_ref"`
	ConcatPart  int    `json:"concat_part"`
	ConcatParts int    `json:"concat_parts"`
	Text        string `json:"ucs-2"`
}

func hex2byte(s string) []byte {
	var u16bytes []byte
	offset := 0
	for i := 0; i < len(s); i = i + 2 {
		b, _ := strconv.ParseUint(s[i:i+2], 16, 8)
		u16bytes = append(u16bytes, byte(b))
		offset++
	}
	return u16bytes
}

func ReadSMSById(smsId string) {
	cmd := exec.Command("uqmi", "-d", "/dev/cdc-wdm0", "--get-message", smsId, "--storage", "me")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("readSMSById failed with %s:\n%s\n", err, string(out))
		log.Fatalf("readSMSById failed with %s\n", err)
		return
	}
	smsJson := string(out)
	strings.Replace(smsJson, "ucs-2", "text", -1)
	smsCtx := SMSContext{}
	err = json.Unmarshal([]byte(smsJson), &smsCtx)
	if err != nil {
		fmt.Printf("JSON parse failed with %s\n", err)
		log.Fatalf("JSON parse failed with %s\n", err)
		return
	}

	decoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()
	bs := hex2byte(smsCtx.Text)
	text, nil := decoder.Bytes(bs)
	if err != nil {
		fmt.Printf("Can not convert text from UTF-16BE to UTF-8, %s\n", err)
		log.Fatalf("Can not convert text from UTF-16BE to UTF-8, %s\n", err)
		return
	}
	fmt.Printf("%03s, From: %11s, Time: %s\n", smsId, smsCtx.From, smsCtx.Timestamp)
	fmt.Printf("%s\n\n", text)
}

func ReadAllSMS() {
	smsId := listAllSMS()
	if smsId == nil {
		return
	}
	for _, smsId := range smsId {
		ReadSMSById(smsId)
	}
}

func DeleteSMSById(smsId string) {
	cmd := exec.Command("uqmi", "-d", "/dev/cdc-wdm0", "--delete-message", smsId, "--storage", "me")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("delete SMS [%s] failed with %s:\n%s\n", smsId, err, string(out))
		log.Fatalf("deleteSMSById failed with %s\n", err)
		return
	}
	fmt.Printf("SMS [%s] delete success\n", smsId)
}

func DeleteAllSMS() {
	smsId := listAllSMS()
	if smsId == nil {
		return
	}
	for _, smsId := range smsId {
		DeleteSMSById(smsId)
	}
}

func listAllSMS() []string {
	cmd := exec.Command("uqmi", "-d", "/dev/cdc-wdm0", "--list-messages", "--storage", "me")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("listAllSMS failed with %s:\n%s\n", err, string(out))
		log.Fatalf("listAllSMS failed with %s\n", err)
		return nil
	}
	smsIdList := string(out)
	reg, _ := regexp.Compile("[0-9]+")
	smsId := reg.FindAllString(smsIdList, -1)
	return smsId
}

func SendSMS(to string, message string) {
	cmd := exec.Command("uqmi", "-d", "/dev/cdc-wdm0", "--send-message", message, "--send-message-target", to)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("sendSMS failed with %s:\n%s\n", err, string(out))
		log.Fatalf("sendSMS failed with %s\n", err)
		return
	} else {
		fmt.Println("SMS send success")
		log.Println("SMS send success")
	}
}
