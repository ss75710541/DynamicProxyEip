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

## 运行命令

```
docker run -e REGION_ID=<regionId> -e  ACCESS_KEY_ID=<accessKeyId> -e ACCESS_KEY_SECRET= <accessKeySecret> -e INSTANCE_ID=<instanceId> -e CHECK_PORT=<check_port> dynamic-proxy-eip:latest
```

## 发布到k8s定时任务

```
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: dynamic-proxy-eip
  namespace: dynamic-proxy-eip
spec:
  schedule: '*/2 * * * *'
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: dynamic-proxy-eip
              image: ss75710541/dynamic-proxy-eip:v0.5
              env:
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
          restartPolicy: OnFailure

```
