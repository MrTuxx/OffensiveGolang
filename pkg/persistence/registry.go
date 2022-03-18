package persistance

import (
	"golang.org/x/sys/windows/registry"
)

// CreateRegistryKey creates an empty registry key
func CreateRegistryKey() error {
	// inputs: key, path
	path := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`

	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()
	return err
}

// SetRegistryValue
func SetRegistryValue(name string, value string) error {
	// keyName? HKLM
	path := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	err = k.SetStringValue(name, value)
	if err != nil {
		return err
	}
	return err
}

// QueryRegistry queries the specified key and returns its contents
func QueryRegistry(name string) (string, error) {
	path := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue(name)
	if err != nil {
		return "", err
	}
	return s, err
}

func DeleteRegistryKey() error {
	path := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	err := registry.DeleteKey(registry.LOCAL_MACHINE, path)
	if err != nil {
		return err
	}

	return err
}

func DeleteRegistryValue(name string) error {
	path := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()
	err = k.DeleteValue(name)
	if err != nil {
		return err
	}
	return err

}
