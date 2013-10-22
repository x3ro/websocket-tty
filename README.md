# socket-tty

This incredible piece of software allows you to execute shell commands on your local machine, from your browser.

Why would anyone do such a thing, you ask? I'm using it to demonstrate shell-based stuff in my browser-based presentations using [Shower](https://github.com/shower/shower/). As soon as I've made the JavaScript side of this process presentable, I will open-source it and link it here.



## Security warning

The server allows to execute **arbitrary commands**, only limited by the permissions of the user it is being run with (read: don't run as root please). The server binds to localhost, and will not accept any remote connections by default.



## "Protocol"

After spending several minutes researching state of the art command execution protocols, I came up with the following JSON messages.

Request:

    { "Cmd": "ls -la" }

Response:

    { "Result": "foo\nbar\nbaz", "Error": "" }

If "Error" is empty in the response, no error occurred.



## Contributing

If that's something you want, I'd be thrilled :) Drop me a pull request.


