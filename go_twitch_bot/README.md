## Twitch chat tool

Tool to monitor/spam Twitch chat. Under development.

Currently it only monitors the chat channel and dumps into stdin.

How to use:
* If you haven't already, [create a Twitch account](http://twitch.tv)
* Login then generate an oauth password [here](http://twitchapps.com/tmi)
* Choose an (active) channel to monitor

This will dump the channel's chat into stdin:
```
go run main.go --username <your twitch username> --password <oauth password> --channel <channel name>
```

(All following commands assume the above 3 parameters)

This will only print chat messages that match "SMOrc". This parameter is a golang regexp.
```
go run main.go ... --search "SMOrc"
```

If you have --search, whenever a match is found, you can automatically send an IRC message back to the chat channel via --exec. For example, assuming that you run the program like this:
```
go run main.go ... --search "SMOrc" --exec ./pasta.py
```

Then whenever "SMOrc" is found in a chat message, pasta.py will be run like this:
```
./pasta.py <channel name> <username> <chat message by user>
```
(channel name is prefixed with "#")

For example, pasta.py (with execute permission) could be:
```python
#!/usr/bin/env python
import sys

channel = sys.argv[1]
username = sys.argv[2]
message = sys.argv[3]

print "PRIVMSG %s Stop spamming SMOrc %s" % (channel, username)
```
This will post "Stop spamming SMOrc" whenever "SMOrc" is seen in chat. Beware that every channel has ways to prevent spamming/flooding (r9k, slow mode, etc). You may receive a message from user "jtv" in that case. By default "jtv" messages will always be written to stdin, you can disable it by adding the --ignoreAdmin flag.

If you have a shell-like program (reads newline-separated input from stdin, produces newline-separated output to stdout) you can use --execLoop insead of --exec. Refer to example/replier.go for an example.

The command "PRIVMSG <channel/username name> <message to send>" is an IRC command. For the full list refer to [this](http://tools.ietf.org/html/rfc1459).
