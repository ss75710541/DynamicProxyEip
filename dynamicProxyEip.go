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

var (
	GODADDY_API_HOST string
	GODADDY_DOMAIN string
	GODADDY_DNS_NAME string
)

// check TCP Port
func checkTCPPort(address string) bool {
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
			if i < 10 {
				i++
				time.Sleep(1*time.Second)
				continue
			}
			fmt.Printf("解绑eip失败: %s", response)
			return err
		}

		fmt.Printf("解绑eip成功，allocationId:  %s\n", allocationId)
		return nil

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
			if i < 10 {
				i++
				time.Sleep(1*time.Second)
				continue
			}
			fmt.Printf("绑定eip失败: %s", response)
			return err
		}
		fmt.Printf("绑定eip成功，allocationId:  %s\n", allocationId)
		return nil
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

	var err error
	for i:=0;i<10;i++ {
		var response *vpc.ReleaseEipAddressResponse
		response, err = client.ReleaseEipAddress(request)

		if !response.IsSuccess()  {
			i++
			time.Sleep(2*time.Second)
			fmt.Printf("释放eip异常,需要重试: %s\n", err.Error())
			fmt.Printf(response.String() )
			continue
		}
		fmt.Printf("释放eip成功,allocationId:  %s\n",  allocationId)
		return
	}

	fmt.Printf("释放eip失败,allocationId: %s , error: ", allocationId, err)
}
// Record 解析记录结构体
type Record struct {
	data string
	name string
	ttl int32
	recordType string
}

func getDomainRecord() string {

	req, err := http.NewRequest("GET",
		"https://" + GODADDY_API_HOST + "/v1/domains/"+GODADDY_DOMAIN+"/records/A/"+GODADDY_DNS_NAME,
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
	fmt.Printf("%s.%s 解析信息：%s\n",GODADDY_DNS_NAME, GODADDY_DOMAIN ,Msg)

	return Msg[0]["data"]
}

func updateDomainRecord(eip string){
	body := fmt.Sprintf(`[{"data": "%s", "ttl": 600 }]`,eip)
	fmt.Println(body)

	req, err := http.NewRequest("PUT",
		"https://" + GODADDY_API_HOST + "/v1/domains/"+GODADDY_DOMAIN+"/records/A/"+GODADDY_DNS_NAME,
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


	fmt.Printf("更新 %s.%s %s 解析成功!\n",GODADDY_DNS_NAME, GODADDY_DOMAIN, string(bodyBytes))
	defer resp.Body.Close()

}

func main() {

	regionId := os.Getenv("REGION_ID")
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	instanceId := os.Getenv("INSTANCE_ID")
	checkPort := os.Getenv("CHECK_PORT")

	GODADDY_API_HOST = os.Getenv("GODADDY_API_HOST")
	GODADDY_DOMAIN   = os.Getenv("GODADDY_DOMAIN")
	GODADDY_DNS_NAME = os.Getenv("GODADDY_DNS_NAME")


	if GODADDY_API_HOST == "" {
		log.Fatal("环境变量 GODADDY_API_HOST 不能为空！")
	}
	if GODADDY_DOMAIN == "" {
		log.Fatal("环境变量 GODADDY_DOMAIN 不能为空！")
	}
	if GODADDY_DNS_NAME == "" {
		log.Fatal("环境变量 GODADDY_DNS_NAME 不能为空！")
	}

	if regionId == "" {
		log.Fatal("环境变量 REGION_ID 不能为空！")
	}
	if accessKeyId == "" {
		log.Fatal("环境变量 ACCESS_KEY_ID 不能为空！")
	}
	if accessKeySecret == "" {
		log.Fatal("环境变量 ACCESS_KEY_SECRET 不能为空！")
	}
	if instanceId == "" {
		log.Fatal("环境变量 INSTANCE_ID 不能为空！")
	}
	if checkPort == "" {
		log.Fatal("环境变量 CHECK_PORT 不能为空！")
	}
	client, err := vpc.NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret)

	eip,allocationId := getEip(err, client,instanceId)

	// 判断是否绑定了eip
	if eip != "" {
		//绑定了eip 则检查连通性
		fmt.Println("Eip: ", eip)
		fmt.Println("AllocationId: ", allocationId)

		if !checkTCPPort(eip+":"+checkPort) {
			log.Println("连接", eip ,checkPort,"失败！")
			// 需要开启容器特权，所以注释ping功能
			//不能连接tcp端口，则解绑eip
			err = unassociateEip(allocationId, instanceId, client)
			if err != nil {
				return
			}
			time.Sleep(time.Second)
			// 释放eip
			releaseEip(allocationId, client)
		} else {
			// 可以连接端口，则返回，什么也不操作
			fmt.Printf("Eip: %s:%s 可以正常连接。\n",eip, checkPort)
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