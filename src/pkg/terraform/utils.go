package terraform

import (
	"errors"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"os/exec"
)

func GetTerraformPath() (string, error) {
	terraformPath, err := exec.LookPath("terraform")
	if err != nil {
		return "", errors.New("terraform binary not found")
	}

	return terraformPath, nil
}

func GetTerraformClient(workingDir string) (*tfexec.Terraform, error) {
	if workingDir == "" {
		workingDir = os.Getenv("PWD")
	}

	terraformPath, err := GetTerraformPath()
	if err != nil {
		return nil, err
	}

	tf, err := tfexec.NewTerraform(workingDir, terraformPath)
	if err != nil {
		return nil, err
	}

	return tf, nil
}
