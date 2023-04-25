package main

import (
    "flag"
    "log"
    "os"

    "github.com/r3dsh/text-generation-webui-launcher"
)

func main() {
    // overwrite python config
    launcher.PythonDistURL = "https://www.python.org/ftp/python/%s/python-%s-embed-amd64.zip"

    // split CLI arguments
    launcherArgs, serverArgs := launcher.LauncherArgs()

    os.Args = append([]string{os.Args[0]}, launcherArgs...)
    installDirPtr := flag.String("home", "", "target directory")
    installBranchPtr := flag.String("branch", "main", "git branch to install text-generation-webui from")
    installPythonPtr := flag.String("python", "3.10.11", "python version to use")
    installActionPtr := flag.Bool("install", false, "install text-generation-webui GUI")
    flag.Parse()

    if *installDirPtr == "" {
        log.Fatal("Launcher supports multiple text-generation-webui installations, -home argument is always required!")
    }

    // override python version
    launcher.PythonVersion = *installPythonPtr

    // new launcher instance
    webui := launcher.New(*installDirPtr, *installBranchPtr, append([]string{os.Args[0]}, serverArgs...))

    // install if requested
    if *installActionPtr {
        err := webui.Install()
        if err != nil {
            log.Fatal(err)
        }
    }

    // bring up the UI - see server.py help at https://github.com/oobabooga/text-generation-webui
    err := webui.StartUI()
    if err != nil {
        log.Fatal(err)
    }
}
