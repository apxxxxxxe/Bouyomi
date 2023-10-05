package data

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/mitchellh/go-ps"
)

// プロセスの存在を返す
func IsProcExist(name string) bool {
	var result bool

	processes, err := ps.Processes()

	if err != nil {
		os.Exit(1)
	}

	result = false
	for _, p := range processes {
		if p.Executable() == name {
			result = true
		}
	}

	return result
}

// 独立したプロセスとして外部コマンドを実行する
func ExecCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{NoInheritHandles: true}
	return cmd.Start()
}
