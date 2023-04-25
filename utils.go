package oobabooga

import (
    "archive/zip"
    "bufio"
    "errors"
    "fmt"
    "io"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"

    "github.com/cavaliergopher/grab/v3"
)

// cleanEnvironment removes all python references
// to disconnect from eventually host installed packages
// paths are additional paths to be pushed to PATH
func cleanEnvironment(currentEnv []string) []string {
    var cleanEnv []string

    for _, e := range currentEnv {
        if strings.HasPrefix(strings.ToLower(e), "path=") {
            // we can get it from parent process
            envPath := strings.Split(os.Getenv("PATH"), string(filepath.ListSeparator))
            var cleanPath []string
            for _, ep := range envPath {
                if ep == "" {
                    continue
                }
                if strings.Contains(strings.ToLower(ep), "python") {
                    continue
                }
                if strings.Contains(strings.ToLower(ep), "conda") {
                    continue
                }
                cleanPath = append(cleanPath, ep)
            }
        } else if strings.Contains(strings.ToLower(e), "python") {
            continue
        } else if strings.Contains(strings.ToLower(e), "conda") {
            continue
        }
        cleanEnv = append(cleanEnv, e)
    }

    return cleanEnv
}

func scanCmdOutput(stdout, stderr io.ReadCloser) error {
    go func() {
        scanner := bufio.NewScanner(stderr)
        scanner.Split(bufio.ScanBytes)
        for scanner.Scan() {
            _, err := os.Stderr.Write(scanner.Bytes())
            if err != nil {
                log.Fatalln("error writing to stderr!")
            }
        }
    }()
    go func() {
        scanner := bufio.NewScanner(stdout)
        scanner.Split(bufio.ScanBytes)
        for scanner.Scan() {
            _, err := os.Stdout.Write(scanner.Bytes())
            if err != nil {
                log.Fatalln("error writing to stdout!")
            }
        }
    }()

    return nil
}

func spawnCommand(cmd *exec.Cmd) error {
    cmd.Env = cleanEnvironment(os.Environ())

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

    return nil
}

func Unzip(zipFile, dst string) error {
    log.Println("unzipping", zipFile, "to", dst)

    archive, err := zip.OpenReader(zipFile)
    if err != nil {
        return err
    }
    defer archive.Close()

    for _, f := range archive.File {
        filePath := filepath.Join(dst, f.Name)
        // fmt.Println("unzipping file ", filePath)

        if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
            // fmt.Println("invalid file path")
            return errors.New("invalid file path")
        }
        if f.FileInfo().IsDir() {
            // fmt.Println("creating directory...")
            os.MkdirAll(filePath, os.ModePerm)
            continue
        }

        if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
            return err
        }

        dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return err
        }

        fileInArchive, err := f.Open()
        if err != nil {
            return err
        }

        if _, err := io.Copy(dstFile, fileInArchive); err != nil {
            return err
        }

        dstFile.Close()
        fileInArchive.Close()
    }

    return nil
}

func Download(dir string, source string) error {
    client := grab.NewClient()
    req, _ := grab.NewRequest(dir, source)

    // start download
    fmt.Printf("Downloading %v...\n", req.URL())
    resp := client.Do(req)
    fmt.Printf("  %v\n", resp.HTTPResponse.Status)

    // start UI loop
    t := time.NewTicker(500 * time.Millisecond)
    defer t.Stop()

Loop:
    for {
        select {
        case <-t.C:
            fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
                resp.BytesComplete(),
                resp.Size,
                100*resp.Progress())

        case <-resp.Done:
            // download is complete
            break Loop
        }
    }

    // check for errors
    if err := resp.Err(); err != nil {
        fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Download saved to ./%v \n", resp.Filename)

    return nil
}

func LauncherArgs() ([]string, []string) {
    launcherArgs := make([]string, 0)
    serverArgs := make([]string, 0)
    skip := false

    for _, arg := range os.Args[1:] {
        if arg == "--" {
            skip = true
            continue
        }

        if skip {
            serverArgs = append(serverArgs, arg)
        } else {
            launcherArgs = append(launcherArgs, arg)
        }
    }

    return launcherArgs, serverArgs
}
