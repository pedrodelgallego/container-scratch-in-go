package main

import (
    "fmt"
    "os"
    "os/exec"
    "syscall"
)

func main() {
    fmt.Printf("Running %v\n", os.Args[2])

    switch os.Args[1] {
    case "run":
        run()
    case "child":
        child()
    default:
        panic("help")
    }
}

func run() {
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
    cmd = setStdInputsAndOutputs(cmd)
    cmd = setNameSpaces(cmd)

    must(cmd.Run())
}

func child() {
    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd = setStdInputsAndOutputs(cmd)

    must(syscall.Sethostname([]byte("container")))
    must(syscall.Chroot("/home/liz/ubuntufs"))
    must(os.Chdir("/"))

    must(cmd.Run())
}

func setNameSpaces(cmd *exec.Cmd) *exec.Cmd {
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
    }
    return cmd
}

func setStdInputsAndOutputs(cmd *exec.Cmd) *exec.Cmd {
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}
