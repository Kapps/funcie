package internal

import (
	"os"
	"path"
)

// GetFuncieBaseDir returns the base directory for funcie configuration and storage.
func GetFuncieBaseDir() string {
	return path.Join(os.Getenv("HOME"), ".funcie/")
}

// GetTerraformDir returns the directory where the Terraform repository for funcie is stored.
func GetTerraformDir() string {
	return path.Join(GetFuncieBaseDir(), "terraform-aws-funcie/")
}

// GetTerraformVarsPath returns the path to the Terraform variables file used for init/destroy.
func GetTerraformVarsPath() string {
	return path.Join(GetFuncieBaseDir(), "funcli.tfvars")
}
