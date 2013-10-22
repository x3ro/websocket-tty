# websocket-tty

This incredible piece of software allows you to execute shell commands on your local machine, from your browser, via websockets.

Why would anyone do such a thing, you ask? I'm using it to demonstrate shell-based stuff in my browser-based presentations using [Shower](https://github.com/shower/shower/). As soon as I've made the JavaScript side of this process presentable, I will open-source it and link it here.



## Security warning

The server allows to execute **arbitrary commands**, only limited by the permissions of the user it is being run with (read: don't run as root please). The server binds to localhost, and will not accept any remote connections by default.



## How to use?

Clone, run `go get` (because we depend on go.net/websocket) and then `go run main.go -path initial/working/directory`. You can now send websocket messages to `http://localhost:8080/exec`, according to the protocol defined below.



## How are commands executed?

They're executed using `exec.Command`, wrapped into a `bash -c 'your command'` invocation. This allows you to do bash stuff like `touch bar && cat foo > bar`.

An exception is the `cd` command. Since commands executed using `exec.Command` run using the parent processes working directory, running `bash -c 'cd ..'` would be quite useless, as it doesn't change the working directory of the parent. Thus if the command begins with `cd ` (cd followed by a space), the command is not run using `bash -c` and instead changes the working directory of the `websocket-tty` process.



## "Protocol"

After spending several minutes researching state of the art command execution protocols, I came up with the following JSON messages.

Request:

    { "Cmd": "ls -la" }

Response:

    { "Result": "foo\nbar\nbaz", "Error": "" }

If "Error" is empty in the response, no error occurred.



## Contributing

If that's something you want, I'd be thrilled :) Drop me a pull request.


