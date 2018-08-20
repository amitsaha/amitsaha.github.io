Title: Docker userns-remap and Container ID cannot be mapped to Host ID
Date: 2018-08-14 20:00
Category: infrastructure
Status: Draft


root@735b8dd78674:~#  find / \( -uid 1004 \)  -ls 2>/dev/null | head -1
   772345      4 drwxr-xr-x   1 1004     sudo         4096 Aug 15 02:29 /usr/share/dotnet
