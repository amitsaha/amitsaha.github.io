Title: Notes on using Cloudflare DNS over HTTPS 
Date: 2017-12-01 16:00
Category: infrastructure
Status: draft


I recently learned about Cloudflare's [1.1.1.1](https://developers.cloudflare.com/1.1.1.1/) DNS service. One of the more
interesting things that caught my attention there was DNS over HTTPS. That is, we can do this:

```
22:27 $ http 'https://cloudflare-dns.com/dns-query?ct=application/dns-json&name=echorand.
me'
HTTP/1.1 200 OK
CF-RAY: 409535ca3b3765bd-SYD
Connection: keep-alive
Content-Length: 281
Content-Type: application/dns-json
Date: Tue, 10 Apr 2018 12:27:53 GMT
Server: cloudflare-nginx
Set-Cookie: __cfduid=dfb12106907c3b55c52b27b8ea99e185a1523363273; expires=Wed, 10-Apr-19 12:27:53 GMT; path=/; domain=.cloudflare-dns.com; HttpOnly; Secure
cache-control: max-age=285

{
    "AD": false,
    "Answer": [
        {
            "TTL": 285,
            "data": "192.30.252.153",
            "name": "echorand.me.",
            "type": 1
        },
        {
            "TTL": 285,
            "data": "192.30.252.154",
            "name": "echorand.me.",
            "type": 1
        }
    ],
    "CD": false,
    "Question": [
        {
            "name": "echorand.me.",
            "type": 1
        }
    ],
    "RA": true,
    "RD": true,
    "Status": 0,
    "TC": false
}


```

Then, I setup a local [DNS over HTTPS](https://developers.cloudflare.com/1.1.1.1/dns-over-https/cloudflared-proxy/) server for all my
DNS resolution by using the `cloudflared` client. This server sets up a local DNS server running on 127.0.0.1 port 53,
listening to your DNS queries and then proxying them over HTTPS to 1.1.1.1:

```

22:31 $ sudo tcpdump host 1.1.1.1
tcpdump: data link type PKTAP
tcpdump: verbose output suppressed, use -v or -vv for full protocol decode
listening on pktap, link-type PKTAP (Apple DLT_PKTAP), capture size 262144 bytes
22:32:09.795071 IP 10.1.1.8.63821 > 1dot1dot1dot1.cloudflare-dns.com.https: Flags [P.], seq 90707220:90707265, ack 994494030, win 8192, length 45
2
...
```

Another interesting bit of information I learned about was Mozilla's plan to integrate DNS over HTTPS in their Firefox
browser. This [article](https://www.ghacks.net/2018/03/20/firefox-dns-over-https-and-a-worrying-shield-study/) has the
details.
