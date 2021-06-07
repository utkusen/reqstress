# reqstress

reqstress is a benchmarking&stressing tool that can send raw HTTP requests. It's written in Go and uses [fasthttp](https://github.com/valyala/fasthttp) library instead of Go's default http library, because of its lightning-fast performance.

## Why Do We Need Another Benchmarking Tool?

There are really great benchmarking tools out there such as [wrk](https://github.com/wg/wrk), [bombardier](https://github.com/codesenberg/bombardier), [hey](https://github.com/rakyll/hey), [ab](https://httpd.apache.org/docs/2.4/tr/programs/ab.html). Some of them don't support sending custom requests, they are only sending a GET request to a given URL. Some of them support custom requests but it's really hard to craft one by using command line parameters. I wanted to create a tool that can read a raw HTTP request from a text file and replays it. 

So, you can copy your favorite request from Burp Suite, Fiddler etc. and pass it to the reqstresser directly. It would be useful for stressing authenticated endpoints and specific requests that create a huge load.

## reqstress vs. Other Tools

reqstresser is not the fastest benchmarking tool, but it's not bad either. I tested couple of popular tools on a $20 Linode server with same amount of threads. Here is the result:


| Tool         | Num. of Sent Requests | Duration |
|--------------|-----------------------|----------|
| wrk          | ~45000                 | 10s      |
| bombardier   | ~41000                 | 10s      |
| ab           | ~40000                 | 10s      |
| reqstress    | ~39304                 | 10s      |
| hey          | ~35127                 | 10s      |
| goldeneye.py | ~10913                 | 10s      |


# Installation

## From Binary

You can download the pre-built binaries from the [releases](https://github.com/utkusen/reqstress/releases) page and run. For example:

`wget https://github.com/utkusen/reqstress/releases/download/v0.1.3/reqstress_0.1.3_Linux_amd64.tar.gz`

`tar xzvf reqstress_0.1.3_Linux_amd64.tar.gz`

`./reqstress --help`

## From Source

1. Install Go on your system
2. Run: `go get -u github.com/utkusen/reqstress`

# Usage

reqstress requires 6 parameters to run: 

`-r` : Path of the request file. For example: `-r request.txt`. Request file should contain a raw HTTP request. For example:

```http
POST /wp-login.php HTTP/1.1
Host: 1.1.1.1
Content-Length: 107
Cache-Control: max-age=0
Upgrade-Insecure-Requests: 1
Origin: http://1.1.1.1
Content-Type: application/x-www-form-urlencoded
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9
Referer: http://1.1.1.1/wp-login.php?redirect_to=http%3A%2F%2F1.1.1.1%2Fwp-admin%2F&reauth=1
Accept-Encoding: gzip, deflate
Accept-Language: tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7
Cookie: wordpress_test_cookie=WP%20Cookie%20check
Connection: close

log=admin&pwd=asdadsasdads

```

`-w` : The number of workers to run. The default value is 500. You can increase or decrease this by testing out the capability of your system.

`-d` : Duration of the test (in seconds). Default is infinite.

`-https` : Target protocol. Can be `true` or `false`. Default is `true`

`-t` : Request timeout. Default is 5(seconds)
