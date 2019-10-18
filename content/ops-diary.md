---
title:  Notes from operations land
date:  2010-10-18 
categories:  infrastructure
---

# PHP dotenv, php-fpm and Bash

Monday was just beginning to roll on as Monday does, I had managed to work out the VPN issues and had
just started to do some planned work. Then, slack tells me that new deployment had just been pushed out successfully, but
the service was actually down. Now, we had HTTP healthchecks which was hitting a specific endpoint but apparently
that was successful, but functionally the service was down. So I check the service logs, which shows something like this:

```
Oct 14 00:03:35 ip-192-168-6-113.eu-central-1.compute.internal docker[20833]: The environment file is invalid!
Oct 14 00:03:35 ip-192-168-6-113.eu-central-1.compute.internal docker[20833]: Failed to parse dotenv file due to an unexpected escape sequence..
```

This was a PHP service which was using [phpdotenv](https://github.com/vlucas/phpdotenv) to load the environment variables
from a file. Okay, so we have found the issue which we can now fix. So, I think why didn't the whole startup script
just abort when it got this error? This is how our startup script looked like at this stage:

```bash
#!/usr/bin/env bash

CHOWN_BIN=/usr/bin/chown
GREP_BIN=/usr/bin/grep


function _check_migrations() {
  php$PHP_VERSION $APPLICATION_ROOT/artisan migrate:status | $GREP_BIN -c 'No'
}

FAILED_MIGRATIONS_COUNT=$(_check_migrations)
if [ $FAILED_MIGRATIONS_COUNT != "0" ]
then
  echo "ERROR: Cannot start while there are $FAILED_MIGRATIONS_COUNT unapplied migrations!!!"
  exit 1
fi


$CHOWN_BIN -R $PHP_FPM_USER:$PHP_FPM_GROUP $APPLICATION_ROOT
$PHP_FPM_BIN -y $PHP_FPM_CONF --nodaemonize &

CHILD_PID=$!
wait $CHILD_PID

```
The error above was happening when the `_check_migrations` function was being called. Since we didn't have a
`-e` in the Bash script, the script continued executing (hence starting the `php-fpm` workers) even though the `php artisan`
command had failed to execute successfully.

So I thought ok, i will just add a `set -e` to the script above:

```diff
diff --git a/scripts/app.sh b/scripts/app.sh
index 6df3af8..3553eff 100755
--- a/scripts/app.sh
+++ b/scripts/app.sh
@@ -1,4 +1,5 @@
 #!/usr/bin/env bash
+set -e
 
 CHOWN_BIN=/usr/bin/chown
 GREP_BIN=/usr/bin/grep
```

So, I tried the above and what happened? The script would error out 
