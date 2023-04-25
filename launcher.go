package launcher

// Launcher text-generation-webui-launcher
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

type Launcher struct {
    HomeDir    string
    Branch     string
    TempDir    string
    WebUIDir   string
    PythonDir  string
    serverArgs []string
}

// @TODO: tbh I forgot why I needed this
func (t *Launcher) IsInstalled() bool {
    log.Println("@TODO: IsInstalled")
    return true
}

func (t *Launcher) Install() error {
    err := os.MkdirAll(t.TempDir, 0744)
    if err != nil {
        return err
    }

    err = t.Python(PythonVersion)
    if err != nil {
        return err
    }

    err = t.Configure()
    if err != nil {
        return err
    }

    err = t.PipInstall([]string{
        "setuptools",
        "num2words",
        "num2words",
    })
    if err != nil {
        return err
    }

    // extract final filename for Unzip
    err = Download(t.TempDir, fmt.Sprintf("https://github.com/oobabooga/text-generation-webui/archive/refs/heads/%s.zip", t.Branch))
    if err != nil {
        return err
    }

    err = Unzip(filepath.Join(t.TempDir, fmt.Sprintf("text-generation-webui-%s.zip", t.Branch)), t.HomeDir)
    if err != nil {
        return err
    }

    err = t.InstallRequirements()
    if err != nil {
        return err
    }

    return nil
}

func (t *Launcher) StartUI() error {
    log.Println("Starting text-generation-webui web ui from", t.WebUIDir)
    forward := t.serverArgs[1:]

    args := "server.py"
    if len(forward) > 0 {
        log.Println("passing CLI arguments to server.py", strings.Join(forward, " "))
        args = fmt.Sprintf("server.py %s", strings.Join(forward, " "))
    }

    log.Println("   spawning", args)

    cmd := exec.Command("python.exe", strings.Split(args, " ")...)
    cmd.Dir = t.WebUIDir
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

func (t *Launcher) Git() error {
    // @TODO:
    // https://github.com/git-for-windows/git/releases/download/v2.40.1.windows.1/PortableGit-2.40.1-64-bit.7z.exe
    // download
    // extract
    // add to PATH

    return nil
}

// Python https://docs.python.org/3/using/windows.html#the-embeddable-package
func (t *Launcher) Python(version string) error {
    downloadURL := fmt.Sprintf(PythonDistURL, version, version)
    zipFile := fmt.Sprintf("python-%s-embed-amd64.zip", version)

    // download Python
    err := Download(t.TempDir, downloadURL)
    if err != nil {
        return err
    }

    err = Unzip(filepath.Join(t.TempDir, zipFile), t.PythonDir)
    if err != nil {
        return err
    }

    // download PIP
    err = Download(t.PythonDir, "https://bootstrap.pypa.io/pip/pip.pyz")
    if err != nil {
        return err
    }

    err = Unzip(filepath.Join(t.PythonDir, "pip.pyz"), t.PythonDir)
    if err != nil {
        return err
    }

    // patch python310._pth file to include Lib\\site-packages and webui folder
    log.Println("writing ", filepath.Join(t.PythonDir, "python310._pth"))
    err = os.WriteFile(filepath.Join(t.PythonDir)+"/python310._pth", []byte(`Lib\\site-packages
python310.zip
`+t.WebUIDir+`
.
`), 0644)
    if err != nil {
        return err
    }

    return nil
}

// Configure - configuration stuff, install real PIP for example.
func (t *Launcher) Configure() error {
    args := "pip.pyz install -U pip wheel setuptools"
    cmd := exec.Command("python.exe", strings.Split(args, " ")...)
    cmd.Dir = t.PythonDir

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

func (t *Launcher) PipInstall(pkgs []string) error {
    args := "install -U " + strings.Join(pkgs, " ")
    cmd := exec.Command("./Scripts/pip.exe", strings.Split(args, " ")...)
    cmd.Dir = t.PythonDir

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

func (t *Launcher) InstallRequirements() error {
    args := "install -U -r requirements.txt"
    cmd := exec.Command("pip.exe", strings.Split(args, " ")...)
    cmd.Dir = t.WebUIDir
    log.Println("install requirements for", t.WebUIDir)

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

func New(home, branch string, serverArgs []string) *Launcher {
    InstallBaseDir := home
    InstallTempDir := filepath.Join(InstallBaseDir, "temp")
    InstallLauncherDir := filepath.Join(InstallBaseDir, fmt.Sprintf("text-generation-webui-%s", branch))
    // InstallPythonDir := filepath.Join(InstallBaseDir, "python")
    InstallPythonDir := InstallBaseDir

    os.Setenv("PATH", InstallPythonDir+string(filepath.ListSeparator)+os.Getenv("PATH"))
    os.Setenv("PATH", filepath.Join(InstallPythonDir, "scripts")+string(filepath.ListSeparator)+os.Getenv("PATH"))
    os.Setenv("TMPDIR", InstallTempDir)

    os.Setenv("PYTHONHOME", InstallPythonDir)
    os.Setenv("PYTHONPATH", filepath.Join(InstallPythonDir, "Lib/site-packages"))

    return &Launcher{
        HomeDir:    home,
        Branch:     branch,
        TempDir:    InstallTempDir,
        WebUIDir:   InstallLauncherDir,
        PythonDir:  InstallPythonDir,
        serverArgs: serverArgs,
    }
}
