email-checker
------------------

Check if email account exists on server through SMTP.

## Usage

```shell
$ email-checker --addrs email@example.com,email2@example.com
smtp: connecting mailsever.example.com.
smtp: connect
smtp: from
email@example.com, false
smtp: connecting mailsever.example.com.
smtp: connect
smtp: from
email2@example.com, true

# You can omit log output with redirect stderr to /dev/null 
$ email-checker --addrs email@example.com,email2@example.com 2>/dev/null
email@example.com, false
email2@example.com, true
```

## Disclaimer

This library is just proof of concept, do not use for spamming, use only for your own SMTP server.