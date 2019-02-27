# MuddyDNS
Course Project of Computer Networks.(W.I.P)

# Build
Just type
```bash
go get github.com/name1e5s/MuddyDNS
```

and use `go build` to build the project.

# Test whether it works

Run the program and use dig to do the test.

```bash
name1e5s@asgard:~/g/s/g/n/MuddyDNS$ ./MuddyDNS -p 11451 &
[1] 1338
name1e5s@asgard:~/g/s/g/n/MuddyDNS$ 2019/02/27 00:17:14 Listening: [::]:11451

name1e5s@asgard:~/g/s/g/n/MuddyDNS$ dig baidu.com  @127.0.0.1 -p 11451
2019/02/27 00:17:23 Forward baidu.com to 114.114.114.114

; <<>> DiG 9.11.3-1ubuntu1.5-Ubuntu <<>> baidu.com @127.0.0.1 -p 11451
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 62660
;; flags: qr rd ra; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 4096
;; QUESTION SECTION:
;baidu.com.                     IN      A

;; ANSWER SECTION:
baidu.com.              67      IN      A       123.125.115.110
baidu.com.              67      IN      A       220.181.57.216

;; Query time: 3 msec
;; SERVER: 127.0.0.1#11451(127.0.0.1)
;; WHEN: Wed Feb 27 00:17:23 CST 2019
;; MSG SIZE  rcvd: 70

```