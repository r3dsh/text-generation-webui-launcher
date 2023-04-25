package oobabooga

import (
    "errors"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

var (
    PythonDistURL = "https://www.python.org/ftp/python/%s/python-%s-embed-amd64.zip"
    PythonVersion = "3.10.11"
)

type Oobabooga struct {
    HomeDir    string
    Branch     string
    TempDir    string
    WebUIDir   string
    PythonDir  string
    serverArgs []string
}

// @TODO: tbh I forgot why I needed this
func (ooba *Oobabooga) IsInstalled() bool {
    log.Println("@TODO: IsInstalled")
    return true
}

func (ooba *Oobabooga) Install() error {
    err := os.MkdirAll(ooba.TempDir, 0744)
    if err != nil {
        return err
    }

    err = ooba.Python(PythonVersion)
    if err != nil {
        return err
    }

    err = ooba.Configure()
    if err != nil {
        return err
    }

    err = ooba.PipInstall([]string{
        "setuptools",
    })
    if err != nil {
        return err
    }

    // extract final filename for Unzip
    err = Download(ooba.TempDir, fmt.Sprintf("https://github.com/oobabooga/text-generation-webui/archive/refs/heads/%s.zip", ooba.Branch))
    if err != nil {
        return err
    }

    err = Unzip(filepath.Join(ooba.TempDir, fmt.Sprintf("text-generation-webui-%s.zip", ooba.Branch)), ooba.HomeDir)
    if err != nil {
        return err
    }

    err = ooba.InstallRequirements()
    if err != nil {
        return err
    }

    return nil
}

func (ooba *Oobabooga) StartUI() error {
    log.Println("Starting oobabooga web ui from", ooba.WebUIDir)
    forward := ooba.serverArgs[1:]

    args := "server.py"
    if len(forward) > 0 {
        log.Println("passing CLI arguments to server.py", strings.Join(forward, " "))
        args = fmt.Sprintf("server.py %s", strings.Join(forward, " "))
    }

    log.Println("   spawning", args)

    cmd := exec.Command("python.exe", strings.Split(args, " ")...)
    cmd.Dir = ooba.WebUIDir
    cmd.Stdin = os.Stdin

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return err
    }
    stderr, err := cmd.StderrPipe()
    if err != nil {
        return err
    }

    err = cmd.Start()
    if err != nil {
        return err
    }

    // starts go routines for stdout and stderr
    err = scanCmdOutput(stdout, stderr)
    if err != nil {
        return err
    }

    err = cmd.Wait()
    if err != nil {
        return errors.New(fmt.Sprintf("error starting pip install: %v", err))
    }
    log.Println("command done")
    return nil
}

// Python https://docs.python.org/3/using/windows.html#the-embeddable-package
func (ooba *Oobabooga) Python(version string) error {
    downloadURL := fmt.Sprintf(PythonDistURL, version, version)
    zipFile := fmt.Sprintf("python-%s-embed-amd64.zip", version)

    // download Python
    err := Download(ooba.TempDir, downloadURL)
    if err != nil {
        return err
    }

    err = Unzip(filepath.Join(ooba.TempDir, zipFile), ooba.PythonDir)
    if err != nil {
        return err
    }

    // download PIP
    err = Download(ooba.PythonDir, "https://bootstrap.pypa.io/pip/pip.pyz")
    if err != nil {
        return err
    }

    err = Unzip(filepath.Join(ooba.PythonDir, "pip.pyz"), ooba.PythonDir)
    if err != nil {
        return err
    }

    // patch python310._pth file to include Lib\\site-packages and webui folder
    log.Println("writing ", filepath.Join(ooba.PythonDir, "python310._pth"))
    err = os.WriteFile(filepath.Join(ooba.PythonDir)+"/python310._pth", []byte(`Lib\\site-packages
python310.zip
`+ooba.WebUIDir+`
.
`), 0644)
    if err != nil {
        return err
    }

    return nil
}

// Configure - configuration stuff, install real PIP for example.
func (ooba *Oobabooga) Configure() error {
    args := "pip.pyz install -U pip wheel setuptools"
    cmd := exec.Command("python.exe", strings.Split(args, " ")...)
    cmd.Dir = ooba.PythonDir

    err := spawnCommand(cmd)
    if err != nil {
        return err
    }

    err = cmd.Wait()
    if err != nil {
        return err
    }
    return nil
}

func (ooba *Oobabooga) PipInstall(pkgs []string) error {
    args := "install -U " + strings.Join(pkgs, " ")
    cmd := exec.Command("./Scripts/pip.exe", strings.Split(args, " ")...)
    cmd.Dir = ooba.PythonDir

    err := spawnCommand(cmd)
    if err != nil {
        return err
    }

    err = cmd.Wait()
    if err != nil {
        return errors.New(fmt.Sprintf("error starting pip install: %v", err))
    }
    return nil
}

func (ooba *Oobabooga) InstallRequirements() error {
    args := "install -U -r requirements.txt"
    cmd := exec.Command("pip.exe", strings.Split(args, " ")...)
    cmd.Dir = ooba.WebUIDir
    log.Println("install requirements for", ooba.WebUIDir)

    err := spawnCommand(cmd)
    if err != nil {
        return err
    }

    err = cmd.Wait()
    if err != nil {
        return errors.New(fmt.Sprintf("error starting pip install: %v", err))
    }
    return nil
}

func New(home, branch string, serverArgs []string) *Oobabooga {
    InstallBaseDir := home
    InstallTempDir := filepath.Join(InstallBaseDir, "temp")
    InstallOobaboogaDir := filepath.Join(InstallBaseDir, fmt.Sprintf("text-generation-webui-%s", branch))
    // InstallPythonDir := filepath.Join(InstallBaseDir, "python")
    InstallPythonDir := InstallBaseDir

    os.Setenv("PATH", InstallPythonDir+string(filepath.ListSeparator)+os.Getenv("PATH"))
    os.Setenv("PATH", filepath.Join(InstallPythonDir, "scripts")+string(filepath.ListSeparator)+os.Getenv("PATH"))
    os.Setenv("TMPDIR", InstallTempDir)

    os.Setenv("PYTHONHOME", InstallPythonDir)
    os.Setenv("PYTHONPATH", filepath.Join(InstallPythonDir, "Lib/site-packages"))

    return &Oobabooga{
        HomeDir:    home,
        Branch:     branch,
        TempDir:    InstallTempDir,
        WebUIDir:   InstallOobaboogaDir,
        PythonDir:  InstallPythonDir,
        serverArgs: serverArgs,
    }
}
