## JSON tool for golang

This generates a list of Golang structures that can be used for the Marshal/Unmarshal functions, using a given example json.

Reads from stdin, writes to stdout.

##### Example:
```
python main.py < example/twitch.json
```

If you have $GOROOT defined then it will automatically pipe the output through $GOROOT/bin/gofmt. Otherwise you can just pipe the output yourself:

```
python main.py < example/twitch.json | /path/to/gofmt
```
