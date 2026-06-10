package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

var conf Config
var auth AuthorizedKey
var ycToken YCToken

func main() {
	conf.New()
	conf.setupLogs()

	log.Debugln(conf.Print())

	if len(conf.vmList) == 0 {
		log.Fatalln("No VM list to checking")
	}
	log.Infof("Start checking %d VMs", len(conf.vmList))

	auth = ParseAuthData()

	ycToken.IAM = getIAMToken()
	vmRunning := true
	for _, vmId := range conf.vmList {
		log.Debugf("VM list: %s", vmId)
		status, err := vmGetStatus(vmId)
		vmRunning = vmIsRunning(status)
		if err != nil {
			log.Errorln(err)
		}
		log.Infof(fmt.Sprintf("%s: %s", vmId, status))
		if vmRunning == false {
			log.Infof("Try to start %s", vmId)
			status, err = vmStart(vmId)
			if err != nil {
				log.Errorln(err)
			}
			for vmRunning == false {
				log.Debugf("Waiting for VM to start...")
				status, err = vmGetStatus(vmId)
				if err != nil {
					log.Errorln(err)
				}
				vmRunning = vmIsRunning(status)
				time.Sleep(time.Duration(conf.delaySec) * time.Second)
			}
		}
	}
	log.Infof("Finish process %d VMs", len(conf.vmList))
}
