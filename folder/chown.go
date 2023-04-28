//go:build !windows

package folder

import (
	"os"
	"path/filepath"
	"syscall"
)

func ChownInherited(file string) error {

	// オーナーを親ディレクトリと同じにする
	stat, err := os.Stat(filepath.Dir(file))
	if err != nil {
		return err
	}
	if sysStat, ok := stat.Sys().(*syscall.Stat_t); ok {

		uid := int(sysStat.Uid)
		gid := int(sysStat.Gid)
		return os.Chown(file, uid, gid)
	}
	return nil
}
