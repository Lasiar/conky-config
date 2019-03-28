package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"sync"
)

const (
	baseConfig = "base_settings.conf"
)

var (
	configFiles = [6]string{"sys", "info", "battery", "data_storage", "network", "weather"}
	wg          = sync.WaitGroup{}
	rootDir     = path.Join(GetHomeDir(), ".conky")
	logDir      string
	conkyPath   string
)

func init() {
	flag.StringVar(&logDir, "log-dir", path.Join(GetHomeDir(), ".conky/log"), "set folder for log")
	flag.StringVar(&conkyPath, "conky", "/usr/bin/conky", "set conky")
	flag.Parse()
}

func main() {
	run()
}

func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("cant find user dir %v", err)
	}
	return usr.HomeDir
}

func run() {
	setLog()

	baseBuffer := bytes.Buffer{}

	configDir := path.Join(rootDir, "/configs")

	baseFile, err := os.Open(path.Join(configDir, baseConfig))
	if err != nil {
		log.Fatalf("error open basecofnig %v", err)
	}

	_, err = baseBuffer.ReadFrom(baseFile)
	if err != nil {
		log.Fatalf("error read baseconfig %v ", err)
	}

	wg.Add(len(configFiles))
	for _, currentFile := range configFiles {

		go handler(baseBuffer, path.Join(configDir, currentFile+".conf"), currentFile)

	}
	wg.Wait()
}

func handler(buffer bytes.Buffer, filePath string, fileName string) {
	currentConfig, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error open config file %v: %v", filePath, err)
	}

	_, err = buffer.ReadFrom(currentConfig)
	if err != nil {
		log.Fatalf("error read file %v: %v", filePath, err)
	}

	temp, err := ioutil.TempFile("", "conky."+fileName+".")
	if err != nil {
		log.Fatalf("error create tempory file %v", err)
	}

	_, err = io.Copy(temp, &buffer)
	if err != nil {
		log.Fatalf("error write tempory file  %v", err)
	}

	logFile, err := os.OpenFile(path.Join(logDir, fileName+".log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error create log file %v", err)
	}

	cmd := exec.Command("bash", "-c", "cd "+rootDir+" && "+ conkyPath+ " -c "+temp.Name()+"&> "+logFile.Name())
	err = cmd.Run()
	if err != nil {
		log.Fatalf("error run conky %v", err)
	}

	wg.Done()
	return
}

func setLog() {

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatalf("error create dir for log: %v", err)
		}
	}

	f, err := os.OpenFile(path.Join(logDir, "starter.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error open log file: %v", err)
	}
	log.SetOutput(f)
}
