package godis

import (
	"bytes"
	"io/ioutil"
	"os"
	"syscall"
)

// AppendToSCF ...
func AppendToSCF(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|syscall.O_CREAT, 0644)
	if err != nil {
		log.Errorf("scf file open %s failed", fileName+err.Error())
	} else {
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

// ReadSCF ...
func ReadSCF(fileName string) []string {
	f, err := os.Open(fileName)
	if err != nil {
		log.Errorf("open scf file failed")
		return nil
	}
	defer f.Close()

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("read scf file error")
		return nil
	}

	lines := bytes.Split(content, []byte{'*'})
	cmds := make([]string, len(lines)-1)

	for k, v := range lines[1:] {
		v := append(v[:0], append([]byte{'*'}, v[0:]...)...)
		cmds[k] = string(v)
	}
	return cmds
}
