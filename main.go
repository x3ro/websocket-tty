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


func splitCmd(cmd string) []string {

    x := strings.Split(cmd, "\"")
    cmds := []string{};
    for i:=0; i<len(x); i++ {
        if i%2 == 0 {
            cmds = append(cmds, strings.Split(strings.Trim(x[i], " "), " ")...)
        } else {
            cmds = append(cmds, x[i])
        }
    }

    cmds1 := []string{};
    for i:=0; i<len(cmds); i++ {
        if cmds[i] != "" {
            cmds1 = append(cmds1, cmds[i])
        }
    }

    return cmds1
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
        //cmds := splitCmd(data.Cmd)
        //out, err := exec.Command(cmds[0], cmds[1:]...).CombinedOutput()
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
