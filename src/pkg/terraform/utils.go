package terraform

import (
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/otterize/otterize-cli/src/pkg/errors"
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
	var err error
	if workingDir == "" {
		workingDir = os.Getenv("PWD")
	}

	terraformPath, err := GetTerraformPath()
	if err != nil {
		return nil, errors.Wrap(err)
	}

	tf, err := tfexec.NewTerraform(workingDir, terraformPath)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return tf, nil
}
