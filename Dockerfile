FROM golang:1.17-alpine as build

COPY ./ /go/project
WORKDIR /go/project/http_server/src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/http_server

FROM alpine

COPY --from=build /go/bin/http_server /bin/

RUN apk add --no-cache tini

EXPOSE 80

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["http_server", "-v=5", "-alsologtostderr"]

# docker build -t http_server:v1.0 .
# docker tag http_server:v1.0 chenxinpeint01/http_server:v1.0
# docker login
# docker push chenxinpeint01/http_server:v1.0

# 切换至ubuntu
# docker login
# sudo docker run -p 80:80 --name http_server -d chenxinpeint01/http_server:v1.0
# docker inspect http_server | grep -i pid
# nsenter -t 5577 -n ip addr
# sudo nsenter -t 5577 -n ip a
#1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
#    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
#    inet 127.0.0.1/8 scope host lo
#       valid_lft forever preferred_lft forever
#4: eth0@if5: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
#    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
#    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
#       valid_lft forever preferred_lft forever
# 宿主机
# ip a
#1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
#    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
#    inet 127.0.0.1/8 scope host lo
#       valid_lft forever preferred_lft forever
#    inet6 ::1/128 scope host
#       valid_lft forever preferred_lft forever
#2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
#    link/ether 00:15:5d:02:71:02 brd ff:ff:ff:ff:ff:ff
#    inet 10.168.1.165/24 brd 10.168.1.255 scope global dynamic eth0
#       valid_lft 41141sec preferred_lft 41141sec
#    inet6 fe80::215:5dff:fe02:7102/64 scope link
#       valid_lft forever preferred_lft forever
#3: docker0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default
#    link/ether 02:42:4b:66:fa:42 brd ff:ff:ff:ff:ff:ff
#    inet 172.17.0.1/16 brd 172.17.255.255 scope global docker0
#       valid_lft forever preferred_lft forever
#    inet6 fe80::42:4bff:fe66:fa42/64 scope link
#       valid_lft forever preferred_lft forever
#5: veth34d078d@if4: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue master docker0 state UP group default
#    link/ether 3a:e4:9e:6a:a5:f9 brd ff:ff:ff:ff:ff:ff link-netnsid 0
#    inet6 fe80::38e4:9eff:fe6a:a5f9/64 scope link
#       valid_lft forever preferred_lft forever

