package main

import (
    "code.google.com/p/go.net/websocket"
    "net/http"
    "os/exec"
    "log"
    "strings"
    "os"
    "path"
    "flag"
)

type Cmd struct {
    Cmd string `json:"cmd"`
}

type Response struct {
    Result string
    Error string
}


func execHandler(ws *websocket.Conn) {
    log.Println("foobar")

    var data Cmd
    err := websocket.JSON.Receive(ws, &data)
    if err != nil {
        log.Fatal(err)
    }

    log.Println(data)

    var out []byte
    if len(data.Cmd) > 2 && data.Cmd[0:3] == "cd " {
        wd, wderr := os.Getwd()
        if wderr != nil {
            log.Fatal(err)
        }

        newwd := path.Clean(wd + "/" + data.Cmd[3:])

        log.Println("changing directory to ", newwd)
        err = os.Chdir(newwd)
        if err != nil {
            log.Println(err)
        } else {
            out = []byte("Changed directory to " + newwd)
        }
    } else {
        out, err = exec.Command("bash", []string{"-c", data.Cmd}...).CombinedOutput()
    }

    var error = ""
    if err != nil {
        log.Println(err)
        error = err.Error()
    }

    err = websocket.JSON.Send(ws, Response{Result:string(out), Error:error})
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    var path string

    flag.StringVar(&path, "path", "", "Working directory, i.e. where commands will be executed.")
    flag.Parse()

    if path == "" {
        log.Fatal("Please specify the 'path' option.")
    }

    err := os.Chdir(path)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Starting with working directory '%s'", path)

    http.Handle("/exec", websocket.Handler(execHandler))
    err = http.ListenAndServe("127.0.0.1:8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: " + err.Error())
    }
}
