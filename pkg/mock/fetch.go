package mock

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binbandit/mocktopus/pkg/utils"
)

func (m *Mock) isRepo() (bool, error) {
	abs, err := filepath.Abs(fmt.Sprintf("%s/.git", m.location))
	if err != nil {
		return false, fmt.Errorf("unable to resolve absolute path to '.git' folder: %v", err)
	}
	_, err = os.Stat(abs)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err

}

func (m *Mock) exists() (bool, error) {
	println(m.location)
	_, err := os.Stat(m.location)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (m *Mock) clone() error {
	repoLocation := fmt.Sprintf("%s/%s", m.config.GetString("host"), m.name)
	_, err := utils.ExecCommand("git", "clone", repoLocation, m.location)
	// We will now check if the folder exists...
	if err != nil {
		return fmt.Errorf("error cloning repo: %v", err)
	}
	if ok, _ := m.exists(); !ok {
		return fmt.Errorf("something went wrong...")
	}
	return nil
}

func (m *Mock) Update() error {
	if ok, _ := m.exists(); !ok {
		m.clone()
	}
	if ok, _ := m.isRepo(); !ok {
		return nil
	}

	_, err := utils.ExecCommand("git", "-C", m.location, "fetch", "--all")
	if err != nil {
		return fmt.Errorf("failed to fetch latest from repo: %v", err)
	}

	// If we are here... the fetch was a success... so we will pull.
	_, err = utils.ExecCommand("git", "-C", m.location, "pull", "--ff-only")
	if err != nil {
		return fmt.Errorf("failed to perform git pull: %v", err)
	}

	return nil
}
