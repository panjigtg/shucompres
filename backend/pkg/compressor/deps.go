package compressor

import "os/exec"

func GhostscriptAvailable() bool {
	_, err := exec.LookPath(ghostscriptBinary())
	return err == nil
}
