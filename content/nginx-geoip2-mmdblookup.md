Title: Nginx and geoip2 on CentOS 7
Date: 2019-05-24
Category: infrastructure

I wanted to setup Nginx logging so that it would perform GeoIP lookup on the IPv4 address in the `X-Forwarded-For` header.
Here's how I went about doing it on CentOS 7.

This [nginx module](https://github.com/leev/ngx_http_geoip2_module) integrates Maxmind GeoIP2 database with the RPMs
being available by [getpagespeed.com](https://www.getpagespeed.com/server-setup/nginx/upgrade-to-geoip2-with-nginx-on-cens-rhel-7).

Once I had installed the module, the hard part for me was how to get the data I wanted - city, timezone information and others
from nginx and the geoip2 module integration. This is where [mmdblookup](https://maxmind.github.io/libmaxminddb/mmdblookup.html)
helped tremendously.

## mmdblookup

`mmdblookup` can be used to read a MaxMind DB file for an IP address and query various information. We need to give it
a path to the DB file and the IP address and it spits out all that it finds out. For example:

```
$ mmdblookup --file /etc/GeoLite2-City.mmdb --ip 49.255.14.118 

  {
    "city": 
      {
        "geoname_id": 
          2147714 <uint32>
        "names": 
          {
            "de": 
              "Sydney" <utf8_string>
            "en": 
              "Sydney" <utf8_string>
            "es": 
              "Sídney" <utf8_string>
            "fr": 
              "Sydney" <utf8_string>
            "ja": 
              "シドニー" <utf8_string>
            "pt-BR": 
              "Sydney" <utf8_string>
            "ru": 
              "Сидней" <utf8_string>
            "zh-CN": 
              "悉尼" <utf8_string>
          }
      }
    "continent": 
      {
        "code": 
          "OC" <utf8_string>
        "geoname_id": 
          6255151 <uint32>
        "names": 
          {
            "de": 
              "Ozeanien" <utf8_string>
            "en": 
              "Oceania" <utf8_string>
            "es": 
              "Oceanía" <utf8_string>
            "fr": 
              "Océanie" <utf8_string>
            "ja": 
              "オセアニア" <utf8_string>
            "pt-BR": 
              "Oceania" <utf8_string>
            "ru": 
              "Океания" <utf8_string>
            "zh-CN": 
              "大洋洲" <utf8_string>
          }
      }
    "country": 
      {
        "geoname_id": 
          2077456 <uint32>
        "iso_code": 
          "AU" <utf8_string>
        "names": 
          {
            "de": 
              "Australien" <utf8_string>
            "en": 
              "Australia" <utf8_string>
            "es": 
              "Australia" <utf8_string>
            "fr": 
              "Australie" <utf8_string>
            "ja": 
              "オーストラリア" <utf8_string>
            "pt-BR": 
              "Austrália" <utf8_string>
            "ru": 
              "Австралия" <utf8_string>
            "zh-CN": 
              "澳大利亚" <utf8_string>
          }
      }
    "location": 
      {
        "accuracy_radius": 
          5 <uint16>
        "latitude": 
          -33.859100 <double>
        "longitude": 
          151.200200 <double>
        "time_zone": 
          "Australia/Sydney" <utf8_string>
      }
    "postal": 
      {
        "code": 
          "2000" <utf8_string>
      }
    "registered_country": 
      {
        "geoname_id": 
          2077456 <uint32>
        "iso_code": 
          "AU" <utf8_string>
        "names": 
          {
            "de": 
              "Australien" <utf8_string>
            "en": 
              "Australia" <utf8_string>
            "es": 
              "Australia" <utf8_string>
            "fr": 
              "Australie" <utf8_string>
            "ja": 
              "オーストラリア" <utf8_string>
            "pt-BR": 
              "Austrália" <utf8_string>
            "ru": 
              "Австралия" <utf8_string>
            "zh-CN": 
              "澳大利亚" <utf8_string>
          }
      }
    "subdivisions": 
      [
        {
          "geoname_id": 
            2155400 <uint32>
          "iso_code": 
            "NSW" <utf8_string>
          "names": 
            {
              "en": 
                "New South Wales" <utf8_string>
              "fr": 
                "Nouvelle-Galles du Sud" <utf8_string>
              "pt-BR": 
                "Nova Gales do Sul" <utf8_string>
              "ru": 
                "Новый Южный Уэльс" <utf8_string>
            }
        }
      ]
  }

```

Now, let's say we only wanted the name of the city in english, we would do something like this:

```
$ mmdblookup --file /etc/GeoLite2-City.mmdb --ip 49.255.14.118 city names en

"Sydney" <utf8_string>

```

If you look at the first "object" in the output above, you will see that the above three arguments, `city names en` is almost
like accessing a nested key inside a dictionary. I say almost, becomes it's not a JSON format. Anyway, this was the key thing
