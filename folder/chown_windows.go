//go:build windows

package folder

func ChownInherited(file string) error {
	// Windowsでは何もしない
	return nil
}
