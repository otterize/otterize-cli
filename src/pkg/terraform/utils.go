package terraform

import (
	"errors"
	"fmt"
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
		fmt.Println("Error initializing Terraform:", err)
		os.Exit(1)
	}

	return tf, nil
}
