/*
Package checks - 3 Docker daemon configuration files
This section covers Docker related files and directory permissions and ownership. Keeping
the files and directories, that may contain sensitive parameters, secure is important for
correct and secure functioning of Docker daemon.
*/
package checks

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/engine-api/client"
)

func CheckServiceOwner(client *client.Client) (res Result) {
	res.Name = "3.1 Verify that docker.service file ownership is set to root:root"
	refUser := "root"
	fileInfo, err := getSystemdFile("docker.service")
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	refUID, refGID := getUserInfo(refUser)
	fileUID, fileGID := getFileOwner(fileInfo)
	if (refUID == fileUID) && (refGID == fileGID) {
		res.Pass()
	} else {
		output := fmt.Sprintf("User/group owner should be : %s", refUser)
		res.Fail(output)
	}

	return
}

func CheckServicePerms(client *client.Client) (res Result) {
	var refPerms uint32
	res.Name = `3.2 Verify that docker.service file permissions are set to
		644 or more restrictive`
	refPerms = 0644
	fileInfo, err := getSystemdFile("docker.service")
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	isLeast, perms := hasLeastPerms(fileInfo, refPerms)
	if isLeast == true {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("File has less restrictive permissions than expected: %v",
			perms)
	}

	return res
}

// func CheckRegistryOwner(client *client.Client) (res Result) {
// 	res.Name = `3.3 Verify that docker-registry.service file ownership is set
// 	to root:root`
// 	refUser := "root"
// 	fileInfo, err := getSystemdFile("docker-registry.service")
// 	if os.IsNotExist(err) {
// 		res.Info("File could not be accessed")
// 		return
// 	}
//
// 	refUid, refGid := getUserInfo(refUser)
// 	fileUid, fileGid := getFileOwner(fileInfo)
// 	if (refUid == fileUid) && (refGid == fileGid) {
// 		res.Status = "PASS"
// 	} else {
// 		res.Status = "WARN"
// 		res.Output = fmt.Sprintf("User/group owner should be : %s", refUser)
// 	}
//
// 	return res
// }

// func CheckRegistryPerms(client *client.Client) (res Result) {
// 	var refPerms uint32
// 	res.Name = `3.4 Verify that docker-registry.service file permissions
// 		are set to 644 or more restrictive`
// 	refPerms = 0644
// 	fileInfo, err := getSystemdFile("docker-registry.service")
// 	if os.IsNotExist(err) {
// 		res.Info("File could not be accessed")
// 		return
// 	}
//
// 	isLeast, perms := hasLeastPerms(fileInfo, refPerms)
// 	if isLeast == true {
// 		res.Status = "PASS"
// 	} else {
// 		res.Status = "WARN"
// 		res.Output = fmt.Sprintf("File has less restrictive permissions than expected: %v",
// 			perms)
// 	}
//
// 	return res
// }

func CheckSocketOwner(client *client.Client) (res Result) {
	res.Name = "3.3 Verify that docker.socket file ownership is set to root:root"
	refUser := "root"
	fileInfo, err := getSystemdFile("docker.socket")
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	refUid, refGid := getUserInfo(refUser)
	fileUid, fileGid := getFileOwner(fileInfo)
	if (refUid == fileUid) && (refGid == fileGid) {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("User/group owner should be : %s", refUser)
	}

	return res
}

func CheckSocketPerms(client *client.Client) (res Result) {
	var refPerms uint32
	res.Name = `3.4 Verify that docker.socket file permissions are set to 644 or more
        restrictive`
	refPerms = 0644
	fileInfo, err := getSystemdFile("docker.socket")
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	isLeast, perms := hasLeastPerms(fileInfo, refPerms)
	if isLeast == true {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("File has less restrictive permissions than expected: %v",
			perms)
	}

	return res
}

func CheckDockerDirOwner(client *client.Client) (res Result) {
	res.Name = "3.5 Verify that /etc/docker directory ownership is set to root:root "
	refUser := "root"
	fileInfo, err := os.Stat("/etc/docker")
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	refUid, refGid := getUserInfo(refUser)
	fileUid, fileGid := getFileOwner(fileInfo)
	if (refUid == fileUid) && (refGid == fileGid) {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("User/group owner should be : %s", refUser)
	}

	return res
}

func CheckDockerDirPerms(client *client.Client) (res Result) {
	var refPerms uint32
	res.Name = `3.6 Verify that /etc/docker directory permissions
		are set to 755 or more restrictive`
	refPerms = 0755
	fileInfo, err := os.Stat("/etc/docker")
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	isLeast, perms := hasLeastPerms(fileInfo, refPerms)
	if isLeast == true {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("File has less restrictive permissions than expected: %v",
			perms)
	}

	return res
}

func CheckRegistryCertOwner(client *client.Client) (res Result) {
	var badFiles []string
	res.Name = `3.7 Verify that registry certificate file ownership
	 is set to root:root`
	refUser := "root"
	refUid, refGid := getUserInfo(refUser)

	files, err := ioutil.ReadDir("/etc/docker/certs.d/")
	if err != nil {
		res.Status = "INFO"
		res.Output = fmt.Sprintf("Directory is inaccessible")
		return res
	}
	for _, file := range files {
		fmt.Println(file.Name())
		if file.IsDir() {
			certs, err := ioutil.ReadDir(file.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, cert := range certs {
				if err != nil {
					log.Fatal(err)
				}
				fileUid, fileGid := getFileOwner(cert)
				if (refUid != fileUid) || (refGid != fileGid) {
					badFiles = append(badFiles, cert.Name())
				}
			}
		}
	}
	if len(badFiles) == 0 {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("Certificate files do not have %s as owner : %s",
			refUser, badFiles)
	}
	return res
}

func CheckRegistryCertPerms(client *client.Client) (res Result) {
	var badFiles []string
	var refPerms uint32
	res.Name = `3.8 Verify that registry certificate file permissions
		are set to 444 or more restrictive`
	refPerms = 0444
	files, err := ioutil.ReadDir("/etc/docker/certs.d/")
	if err != nil {
		res.Status = "INFO"
		res.Output = fmt.Sprintf("Directory is inaccessible")
		return res
	}
	for _, file := range files {
		fmt.Println(file.Name())
		if file.IsDir() {
			certs, err := ioutil.ReadDir(file.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, cert := range certs {
				if err != nil {
					log.Fatal(err)
				}
				isLeast, _ := hasLeastPerms(cert, refPerms)
				if isLeast == false {
					badFiles = append(badFiles, cert.Name())
				}
			}
		}
	}
	if len(badFiles) == 0 {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("Certificate files do not have required permissions: %s",
			badFiles)
	}
	return res
}

func CheckCACertOwner(client *client.Client) (res Result) {
	res.Name = "3.9 Verify that TLS CA certificate file ownership is set to root:root"
	refUser := "root"
	dockerProc, _ := GetProcCmdline("docker")
	_, certPath := GetCmdOption(dockerProc, "--tlscacert")
	fileInfo, err := os.Stat(certPath)
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	refUid, refGid := getUserInfo(refUser)
	fileUid, fileGid := getFileOwner(fileInfo)
	if (refUid == fileUid) && (refGid == fileGid) {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("User/group owner should be : %s", refUser)
	}

	return res
}

func CheckCACertPerms(client *client.Client) (res Result) {
	var refPerms uint32
	res.Name = `3.10 Verify that TLS CA certificate file permissions
	are set to 444 or more restrictive`
	refPerms = 0444
	dockerProc, _ := GetProcCmdline("docker")
	_, certPath := GetCmdOption(dockerProc, "--tlscacert")
	fileInfo, err := os.Stat(certPath)
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	isLeast, perms := hasLeastPerms(fileInfo, refPerms)
	if isLeast == true {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("File has less restrictive permissions than expected: %v",
			perms)
	}

	return res
}

func CheckServerCertOwner(client *client.Client) (res Result) {
	res.Name = `3.11 Verify that Docker server certificate file ownership is set to
        root:root`
	refUser := "root"
	dockerProc, _ := GetProcCmdline("docker")
	_, certPath := GetCmdOption(dockerProc, "--tlscert")
	fileInfo, err := os.Stat(certPath)
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	refUid, refGid := getUserInfo(refUser)
	fileUid, fileGid := getFileOwner(fileInfo)
	if (refUid == fileUid) && (refGid == fileGid) {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("User/group owner should be : %s", refUser)
	}

	return res
}

func CheckServerCertPerms(client *client.Client) (res Result) {
	var refPerms uint32
	res.Name = `3.12 Verify that Docker server certificate file permissions
		are set to 444 or more restrictive`
	refPerms = 0444
	dockerProc, _ := GetProcCmdline("docker")
	_, certPath := GetCmdOption(dockerProc, "--tlscert")
	fileInfo, err := os.Stat(certPath)
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	isLeast, perms := hasLeastPerms(fileInfo, refPerms)
	if isLeast == true {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("File has less restrictive permissions than expected: %v",
			perms)
	}

	return res
}

func CheckCertKeyOwner(client *client.Client) (res Result) {
	res.Name = `3.13 Verify that Docker server certificate key file ownership is set to
        root:root`
	refUser := "root"
	dockerProc, _ := GetProcCmdline("docker")
	_, certPath := GetCmdOption(dockerProc, "--tlskey")
	fileInfo, err := os.Stat(certPath)
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	refUid, refGid := getUserInfo(refUser)
	fileUid, fileGid := getFileOwner(fileInfo)
	if (refUid == fileUid) && (refGid == fileGid) {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("User/group owner should be : %s", refUser)
	}

	return res
}

func CheckCertKeyPerms(client *client.Client) (res Result) {
	var refPerms uint32
	res.Name = `3.14 Verify that Docker server certificate key file
	permissions are set to 400`
	refPerms = 0400
	dockerProc, _ := GetProcCmdline("docker")
	_, certPath := GetCmdOption(dockerProc, "--tlskey")
	fileInfo, err := os.Stat(certPath)
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	isLeast, perms := hasLeastPerms(fileInfo, refPerms)
	if isLeast == true {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("File has less restrictive permissions than expected: %v",
			perms)
	}

	return res
}

func CheckDockerSockOwner(client *client.Client) (res Result) {
	res.Name = `3.15 Verify that Docker socket file ownership
	is set to root:docker`
	refUser := "root"
	refGroup := "docker"
	fileInfo, err := os.Stat("/var/run/docker.sock")
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	refUid, _ := getUserInfo(refUser)
	refGid := getGroupId(refGroup)
	fileUid, fileGid := getFileOwner(fileInfo)
	if (refUid == fileUid) && (refGid == fileGid) {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("User/group owner should be : %s", refGroup)
	}

	return res
}

func CheckDockerSockPerms(client *client.Client) (res Result) {
	var refPerms uint32
	res.Name = `3.16 Verify that Docker socket file permissions are set to 660`
	refPerms = 0660
	fileInfo, err := os.Stat("/var/run/docker.sock")
	if os.IsNotExist(err) {
		res.Info("File could not be accessed")
		return
	}

	isLeast, perms := hasLeastPerms(fileInfo, refPerms)
	if isLeast == true {
		res.Status = "PASS"
	} else {
		res.Status = "WARN"
		res.Output = fmt.Sprintf("File has less restrictive permissions than expected: %v",
			perms)
	}

	return res
}

func CheckDaemonJSONOwner(client *client.Client) (res Result) {
	return
}

func CheckDaemonJSONPerms(client *client.Client) (res Result) {
	return
}

func CheckDefaultOwner(client *client.Client) (res Result) {
	return
}

func CheckDefaultPerms(client *client.Client) (res Result) {
	return
}
