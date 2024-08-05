package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"uqmi_sms/sms"
)

func isAllDigit(str string) bool {
	for _, char := range str {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func main() {
	del := flag.Bool("delete", false, "Delete all messages")
	read := flag.Bool("read", false, "Read all messages")
	smsId := flag.Int("id", -1, "message id")
	send := flag.Bool("send", false, "Send messages")
	target := flag.String("target", "", "Target phone number")
	text := flag.String("text", "", "SMS body")
	flag.Parse()
	if *del == *read && *read == *send {
		flag.Usage()
		return
	}
	if *read {
		if *smsId != -1 {
			sms.ReadSMSById(strconv.Itoa(*smsId))
		} else {
			sms.ReadAllSMS()
		}
	} else if *del {
		if *smsId != -1 {
			sms.DeleteSMSById(strconv.Itoa(*smsId))
		} else {
			sms.DeleteAllSMS()
		}
	} else if *send {
		if (len(*target) == 0) || (len(*text) == 0) {
			fmt.Println("Target phone number and SMS body is required")
			log.Fatal("Target phone number and SMS body is required")
			return
		}
		if !isAllDigit(*target) {
			fmt.Println("Malformed target phone number")
			log.Fatal("Malformed target phone number")
			return
		}
		sms.SendSMS(*target, *text)
	}
}
