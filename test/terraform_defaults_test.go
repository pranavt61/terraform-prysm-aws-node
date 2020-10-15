package test

import (
	"fmt"
	"strings"
	"testing"
  "os"
	"time"
  "io/ioutil"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformSshExample(t *testing.T) {
	t.Parallel()

  // test ID
  uniqueID := random.UniqueId()

  awsRegion := "us-east-1"

  // Create and write SSH keys
  if _, err := os.Stat("./fixtures/ssh_keys"); os.IsNotExist(err) {
    // WARN too open for perm?
    os.MkdirAll("./fixtures/ssh_keys", 0777)
  }

  keyPair := ssh.GenerateRSAKeyPair(t, 4096)
  keyPairName := fmt.Sprintf("terraform_test_%s", uniqueID)

  private_key_bytes := []byte(keyPair.PrivateKey)
  err := ioutil.WriteFile(
    "./fixtures/ssh_keys/" + keyPairName,
    private_key_bytes,
    0644,
  )
  if err != nil {
    t.Fatalf("Failed to write private ssh key")
  }

  public_key_bytes := []byte(keyPair.PublicKey)
  err = ioutil.WriteFile(
    "./fixtures/ssh_keys/" + keyPairName + ".pub",
    public_key_bytes,
    0644,
  )
  if err != nil {
    t.Fatalf("Failed to write public ssh key")
  }

  // Import SSH key
  ec2KeyPair := aws.ImportEC2KeyPair(t, awsRegion, keyPairName, keyPair)

  exampleFolder := "../examples/defaults"
  testFolder, err := os.Getwd()
  if err != nil {
    t.Fatalf("Failed to read working dir")
  }

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: exampleFolder,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
      "private_key_path": testFolder + "/fixtures/ssh_keys/" + keyPairName,
      "public_key_path": testFolder + "/fixtures/ssh_keys/" + keyPairName + ".pub",
		},
	}

  // Clean up
  defer func() {
    terraform.Destroy(t, terraformOptions)

    // Remoce ec2 key
		aws.DeleteEC2KeyPair(t, ec2KeyPair)

    // Remove private SSH key
    err := os.Remove("./fixtures/ssh_keys/" + keyPairName)
    if err != nil {
      t.Fatalf("Failed to remove private ssh key")
    }

    // Remove public SSH key
    err = os.Remove("./fixtures/ssh_keys/" + keyPairName + ".pub")
    if err != nil {
      t.Fatalf("Failed to remove public ssh key")
    }
  }()

  // Create terraform instance
	terraform.InitAndApply(t, terraformOptions)

  //--- Run tests ---//
	testGenericSSHCommand(t, terraformOptions, ec2KeyPair)

  // Docker
	testDockerInstall(t, terraformOptions, ec2KeyPair)
	testDockerComposeInstall(t, terraformOptions, ec2KeyPair)
}

func testGenericSSHCommand(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := fmt.Sprintf("SSH to public host %s", publicInstanceIP)

	// Run a simple echo command on the server
	command := fmt.Sprintf("echo -n '%s'", "Hello, World")
	expectedText := "Hello, World"

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}

func testDockerInstall(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "Check docker version on host"

	// Run a simple echo command on the server
	command := "docker --version"
	expectedText := "Docker version 19.03.13, build 4484c46d9d"

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}

func testDockerComposeInstall(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "Check docker-compose version on host"

	// Run a simple echo command on the server
	command := "docker-compose --version"
	expectedText := "docker-compose version 1.27.4, build 40524192"

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}

func testPrysmDockerComposeFilesFromGit(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "Check for docker-compose files in home dir"

	// Run a simple echo command on the server
	command := "cd ~/prysm-docker-compose && git remote -v | head -n 1"
	expectedText := "origin  https://github.com/pranavt61/prysm-docker-compose.git (fetch)"

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}

func testPrysmKeystoreFileTransfer(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "Check for keystore.json"

	// Run a simple echo command on the server
	command := "ls /home/ubuntu/prysm-docker-compose/launchpad/eth2.0-deposit-cli/validator_keys/keystore.json"
  expectedText := "/home/ubuntu/prysm-docker-compose/launchpad/eth2.0-deposit-cli/validator_keys/keystore.json"

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}

func testPrysmDepositFileTransfer(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "Check for deposit.json"

	// Run a simple echo command on the server
	command := "ls /home/ubuntu/prysm-docker-compose/launchpad/eth2.0-deposit-cli/validator_keys/depost.json"
  expectedText := "/home/ubuntu/prysm-docker-compose/launchpad/eth2.0-deposit-cli/validator_keys/deposit.json"

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}

func testPrysmWalletTransfer(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "Check for deposit.json"

	// Run a simple echo command on the server
	command := "ls /home/ubuntu/prysm-docker-compose/validator/wallets/hash"
  expectedText := "/home/ubuntu/prysm-docker-compose/validator/wallets/hash"

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}

func testPrysmWalletPasswordTransfer(t *testing.T, terraformOptions *terraform.Options, keyPair *aws.Ec2Keypair) {
	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair.KeyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 30
	timeBetweenRetries := 5 * time.Second
	description := "Check for deposit.json"

	// Run a simple echo command on the server
	command := "ls /home/ubuntu/prysm-docker-compose/validator/passwords/wallet-password"
  expectedText := "/home/ubuntu/prysm-docker-compose/validator/passwords/wallet-password"

	// Verify that we can SSH to the Instance and run commands
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		actualText, err := ssh.CheckSshCommandE(t, publicHost, command)

		if err != nil {
			return "", err
		}

		if strings.TrimSpace(actualText) != expectedText {
			return "", fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return "", nil
	})
}
