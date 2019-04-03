package main

import (
	"DynamicProxyEip/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)
// check tch port
func checkTcpPort(address string) bool {
	var conn net.Conn
	var err error

	for i:=0; i<3; i++ {
		conn, err = net.DialTimeout("tcp", address, 3*time.Second)
		if err != nil{
			fmt.Println("could not connect to server: ", err)
		}
		time.Sleep(5*time.Second)
	}

	if err != nil{
		return false
	}

	defer conn.Close()
	return true
}
//根据实例id获取Eip
func getEip(err error, client *vpc.Client, instanceId string)  (string,string) {
	request := vpc.CreateDescribeEipAddressesRequest()
	request.Status = "InUse"
	request.AssociatedInstanceId = instanceId
	request.AssociatedInstanceType = "EcsInstance"
	response, err := client.DescribeEipAddresses(request)
	if err != nil {
		fmt.Print(err.Error())
	}

	if len(response.EipAddresses.EipAddress) > 0 {
		eip := response.EipAddresses.EipAddress[0].IpAddress
		allocationId := response.EipAddresses.EipAddress[0].AllocationId
		return eip,allocationId
	}
	return "",""
}
//获取可用Eip
func getAvailableEip(err error, client *vpc.Client)  (string,string) {
	request := vpc.CreateDescribeEipAddressesRequest()
	request.Status = "Available"
	response, err := client.DescribeEipAddresses(request)
	if err != nil {
		fmt.Print(err.Error())
	}

	if len(response.EipAddresses.EipAddress) > 0 {
		eip := response.EipAddresses.EipAddress[0].IpAddress
		allocationId := response.EipAddresses.EipAddress[0].AllocationId
		fmt.Printf("获取可用Eip: %s\n", eip)
		return eip,allocationId
	}
	return "",""
}
//解绑eip
func unassociateEip(allocationId string, instanceId string, client *vpc.Client) error {
	request := vpc.CreateUnassociateEipAddressRequest()
	request.AllocationId = allocationId
	request.InstanceId = instanceId
	i := 0
	for  {
		response, err := client.UnassociateEipAddress(request)

		if err != nil {
			if i < 4 {
				i++
				time.Sleep(1*time.Second)
				continue
			}
			fmt.Printf("解绑eip失败: %s", response)
			return err
		}else {
			fmt.Printf("解绑eip成功，allocationId:  %s\n", allocationId)
			return err
		}

	}

}
//绑定eip
func associateEip(allocationId string, instanceId string, client *vpc.Client) error {
	request := vpc.CreateAssociateEipAddressRequest()

	request.AllocationId = allocationId
	request.InstanceId = instanceId

	i := 0
	for {
		response, err := client.AssociateEipAddress(request)
		if err != nil {
			if i < 4 {
				i++
				time.Sleep(1*time.Second)
				continue
			}
			fmt.Printf("绑定eip失败: %s", response)
			return err
		}else {
			fmt.Printf("绑定eip成功，allocationId:  %s\n", allocationId)
			return err
		}
	}
}
//申请分配Eip
func allocateNewEip(err error, client *vpc.Client) (string,string) {
	request := vpc.CreateAllocateEipAddressRequest()
	request.AutoPay = "true"
	request.Bandwidth = "2"
	request.InternetChargeType = "PayByTraffic"
	response, err := client.AllocateEipAddress(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	eip := response.EipAddress
	allocationId := response.AllocationId
	return eip,allocationId
}
// 释放Eip
func releaseEip(allocationId string, client *vpc.Client) {
	request := vpc.CreateReleaseEipAddressRequest()
	request.AllocationId = allocationId

	for i:=0;i<5;i++ {
		response, err := client.ReleaseEipAddress(request)
		if response.RequestId == ""  {
			i++
			time.Sleep(time.Second)
			fmt.Printf("释放eip异常,需要重试: %s\n", err.Error())
			continue
		}
		fmt.Printf("释放eip成功,allocationId:  %s\n",  allocationId)
		return
	}

	fmt.Printf("释放eip失败,allocationId: %s", allocationId)
}

type Record struct {
	data string
	name string
	ttl int32
	recordType string
}

func getDomainRecord() string {
	apiHost := "api.godaddy.com"
	domain := "servicehub.services"
	name := "devopsproxy"

	req, err := http.NewRequest("GET",
		"https://" + apiHost + "/v1/domains/"+domain+"/records/A/"+name,
		nil)

	if err != nil {
		fmt.Println(err)
		panic(1)
	}


	req.Header.Set("Authorization", os.ExpandEnv("sso-key ${GODADDY_KEY}:${GODADDY_SECRET}"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		// handle err
		fmt.Printf("%s", err)
	}


	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()


	if err2 != nil {
		fmt.Printf("%s", err2)
	}


	var Msg []map[string]string
	json.Unmarshal([]byte(bodyBytes), &Msg)
	fmt.Printf("devopsproxy.servicehub.services 解析信息：%s\n",Msg)

	return Msg[0]["data"]
}

func updateDomainRecord(eip string){

	apiHost := "api.godaddy.com"
	domain := "servicehub.services"
	name := "devopsproxy"

	body := fmt.Sprintf(`[{"data": "%s", "ttl": 600 }]`,eip)
	fmt.Println(body)

	req, err := http.NewRequest("PUT",
		"https://" + apiHost + "/v1/domains/"+domain+"/records/A/"+name,
		bytes.NewBuffer([]byte(body)))


	if err != nil {
		fmt.Println(err.Error())
		panic(1)
	}


	req.Header.Set("Authorization", os.ExpandEnv("sso-key ${GODADDY_KEY}:${GODADDY_SECRET}"))
	req.Header.Set("Content-Type", "application/json")

	resp, err2 := http.DefaultClient.Do(req)

	if err2 != nil {
		fmt.Println(err2.Error())
		panic(1)
	}


	bodyBytes, err3 := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()


	if err3 != nil {
		fmt.Println(err3.Error())
		panic(1)
	}


	fmt.Printf("更新 devopsproxy.servicehub.services %s 解析成功!\n", string(bodyBytes))
	defer resp.Body.Close()

}

func main() {

	regionId := os.Getenv("REGION_ID")
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	instanceId := os.Getenv("INSTANCE_ID")
	check_port := os.Getenv("CHECK_PORT")

	if regionId == "" {
		log.Println("环境变量 REGION_ID 不能为空！")
		panic(1)
		return
	}
	if accessKeyId == "" {
		log.Println("环境变量 ACCESS_KEY_ID 不能为空！")
		panic(1)
		return
	}
	if accessKeySecret == "" {
		log.Println("环境变量 ACCESS_KEY_SECRET 不能为空！")
		panic(1)
		return
	}
	if instanceId == "" {
		log.Println("环境变量 INSTANCE_ID 不能为空！")
		panic(1)
		return
	}
	if check_port == "" {
		log.Println("环境变量 CHECK_PORT 不能为空！")
		panic(1)
		return
	}
	client, err := vpc.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)

	eip,allocationId := getEip(err, client,instanceId)

	// 判断是否绑定了eip
	if eip != "" {
		//绑定了eip 则检查连通性
		fmt.Println("Eip: ", eip)
		fmt.Println("AllocationId: ", allocationId)

		if !checkTcpPort(eip+":"+check_port) {
			log.Println("连接", eip ,check_port,"失败！")
			// 需要开启容器特权，所以注释ping功能
			//不能连接tcp端口，则解绑eip,并释放
			err = unassociateEip(allocationId, instanceId, client)
			if err != nil {
				return
			}
			time.Sleep(time.Second)
			defer releaseEip(allocationId, client)
		} else {
			// 可以连接端口，则返回，什么也不操作
			fmt.Printf("Eip: %s:%s 可以正常连接。\n",eip, check_port)
			// 获取解析信息
			eip2 := getDomainRecord()
			if eip != eip2 {
				fmt.Println("Eip解析不同步，更新解析记录\n")
				updateDomainRecord(eip)
			}
			return
		}
	}

	// 获取可用Eip
	eip,allocationId = getAvailableEip(err,client)

	if eip == "" {
		// 无可用Eip,申请新Eip
		eip, allocationId = allocateNewEip(err, client)

		if eip == "" {
			fmt.Print("申请分配Eip错误。")
			panic(1)
			return
		}
		fmt.Printf("申请分配Eip: %s\n", eip)
	}

	// 绑定Eip
	err = associateEip(allocationId, instanceId, client)
	if err != nil {
		fmt.Printf("绑定Eip错误： %s\n",err.Error())
		panic(1)
		return
	}
	fmt.Printf("绑定 eip %s 到实例 %s\n", eip, instanceId)

	// 更新解析ip
	updateDomainRecord(eip)
	// 获取解析信息
	getDomainRecord()

	smtpTo := os.Getenv("SMTP_TO")
	if smtpTo != "" {
		utils.SendMail("代理EIP替换为: "+eip,"代理EIP替换为: "+eip)
	}
}