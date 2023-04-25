package main

import (
    "flag"
    "log"
    "os"

    "github.com/r3dsh/oobabooga-launcher"
)

func main() {
    // overwrite python config
    oobabooga.PythonDistURL = "https://www.python.org/ftp/python/%s/python-%s-embed-amd64.zip"

    // split CLI arguments
    launcherArgs, serverArgs := oobabooga.LauncherArgs()

    os.Args = append([]string{os.Args[0]}, launcherArgs...)
    installDirPtr := flag.String("home", "", "target directory")
    installBranchPtr := flag.String("branch", "main", "git branch to install oobabooga from")
    installPythonPtr := flag.String("python", "3.10.11", "python version to use")
    installActionPtr := flag.Bool("install", false, "install oobabooga GUI")
    flag.Parse()

    if *installDirPtr == "" {
        log.Fatal("Launcher supports multiple oobabooga installations, -home argument is always required!")
    }

    // override python version
    oobabooga.PythonVersion = *installPythonPtr

    // new oobabooga instance
    ooba := oobabooga.New(*installDirPtr, *installBranchPtr, append([]string{os.Args[0]}, serverArgs...))

    // install if requested
    if *installActionPtr {
        err := ooba.Install()
        if err != nil {
            log.Fatal(err)
        }
    }

    // bring up the UI - see server.py help at https://github.com/oobabooga/text-generation-webui
    err := ooba.StartUI()
    if err != nil {
        log.Fatal(err)
    }
}
