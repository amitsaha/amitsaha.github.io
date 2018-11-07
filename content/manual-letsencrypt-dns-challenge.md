Title: Manual Let's Encrypt, GoDadday DNS and IIS
Date: 2018-11-07
Category: infrastructure

I wanted to create a new SSL certificate for IIS hosted ASP.NET framework application. The key data that may
make this post relevant to you are:

- Let's Encrypt Challenge mode: DNS TXT record
- DNS provider: GoDaddy
- Target web server: IIS
- Target operating system: Windows
- Software used: `certbot` and `openssl`


## Why certbot?

I decided to use [certbot](https://certbot.eff.org/) since it allowed me do create the DNX TXT entries manually. Two other
projects I looked at were [lego](https://github.com/xenolf/lego) and [win-acme](https://github.com/PKISharp/win-acme). While
writing this post, I found out that `lego` has inbuilt support for `godadday` (:facepalm:), so I could have used it to create
the DNS TXT record automatically. Anyway, `win-acme` needed hooks to be provided for the DNS challenge, which seemed like
another thing to do at the moment.

## Generating the certificate

Once you have installed `certbot`:

```
$ certbot certonly --manual --preferred-challenges dns -d <your domain> --config-dir . --logs-dir . --work-dir .
```

The program will pause displaying:

```
Please deploy a DNS TXT record under the name
_acme-challenge.<your domain> with the following value:
random$string
Before continuing, verify the record is deployed.
```

Now, go to your GoDaddy DNS management page, and create the TXT record with the specified string. Once you 
have verified that the domain entry has propagated, press ENTER to continue. To verify, use `nslookup -q=TXT <domain>`
on Windows, or `dig -t` on Linux.


Once the record has propagated, certbot will try to find it, and if successful continue and eventually give an 
output like this:

```
IMPORTANT NOTES:

 - Congratulations! Your certificate and chain have been saved at:
   /home/asaha/letsencrypt/live/<your domain>/fullchain.pem
   Your key file has been saved at:
   /home/asaha/letsencrypt/live/<your domain>/privkey.pem
   ...
```

## Importing into IIS

The directory created will have a bunch of files. We will next create a .pfx file for importing into IIS using `openssl`:

```
$ openssl pkcs12 -export -out certificate.pfx -inkey privkey.pem -in cert.pem -certfile chain.pem
Enter Export Password:
Verifying - Enter Export Password:
```

The resultant file will be certificate.pfx. Now, copy the `certificate.pfx` file to the target IIS box and import
it using this handy [guide](https://www.digicert.com/ssl-support/pfx-import-export-iis-7.htm).

The next step is to attempt to automate at least the certificate generation process using `lego`.
