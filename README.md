# Selfbot

Simple, easy to use (and deploy/hack) bot for matrix.

# Setup

```bash
$ git clone https://git.mrcyjanek.net/mrcyjanek/selfbot.git
$ make run
```
aaand follow on screen instructions

# Usage

To get list of all installed plugins send `!help` to any chat (you can make an empty channel for bot chats btw)

Most of the commands use shell-like parsing of arguments so doing `!wiki en Bitcoin Litecoin` will send you search results for both `Bitcoin` and `Litecoin` while doing `!wiki en Stack\ Overflow` or `!wiki en 'Stack Overflow'` will give you results for `Stack Overflow`