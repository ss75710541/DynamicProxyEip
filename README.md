# DynamicProxyEip


## 描述

检查阿里云proxy主机Eip端口连通性，如果不能连接则更换Eip

## 运行说明

环境变量|描述
-----|-----
REGION_ID | 机房id
ACCESS_KEY_ID| 阿里云密钥id
ACCESS_KEY_SECRET|阿里云密钥
INSTANCE_ID|虚拟机实例id
CHECK_PORT|检查连通性的TCP端口
SMTP_USER|smtp用户
SMTP_PASSWORD|smtp密码
SMTP_HOST|smtp服务器地址
SMTP_TO|收件箱，多个使用;隔开
GODADDY_KEY| godaddy api key
GODADDY_SECRET| godaddy api 密钥
GODADDY_API_HOST| goadddy api 地址
GODADDY_DOMAIN|域名domain
GODADDY_DNS_NAME|域名解析名称

## 运行命令

```
docker run -e REGION_ID=<regionId> -e  ACCESS_KEY_ID=<accessKeyId> -e ACCESS_KEY_SECRET= <accessKeySecret> -e INSTANCE_ID=<instanceId> -e CHECK_PORT=<check_port> -e GODADDY_KEY=<godaddyKey> -e GODADDY_SECRET=<godaddySecret> ss75710541/dynamic-proxy-eip:latest
```

## 发布到k8s定时任务示例

```
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: dynamic-proxy-eip
  namespace: dynamic-proxy-eip
spec:
  startingDeadlineSeconds: 60
  schedule: '*/5 * * * *'
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: dynamic-proxy-eip
              image: ss75710541/dynamic-proxy-eip:latest
              env:
                - name: GODADDY_API_HOST
                  value: "<GODADDY_API_HOST>"
                - name: GODADDY_DOMAIN
                  value: "<GODADDY_DOMAIN>"
                - name: GODADDY_DNS_NAME
                  value: "<GODADDY_DNS_NAME>"
                - name: REGION_ID
                  value: "<regionId>"
                - name: ACCESS_KEY_ID
                  valueFrom:
                    secretKeyRef:
                      name: access-key
                      key: accessKeyId
                - name: ACCESS_KEY_SECRET
                  valueFrom:
                    secretKeyRef:
                      name: access-key
                      key: accessKeySecret
                - name: INSTANCE_ID
                  value: "<instanceId>"
                - name: CHECK_PORT
                  value: "<checkPort>"
                - name: GODADDY_KEY
                  valueFrom:
                    secretKeyRef:
                      name: access-key
                      key: godaddyKey
                - name: GODADDY_SECRET
                  valueFrom:
                    secretKeyRef:
                      name: access-key
                      key: godaddySecret
          restartPolicy: OnFailure

```
