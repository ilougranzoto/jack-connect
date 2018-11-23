package sftp

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//MoveFile ...
func MoveFile(localDir string, remoteDir string, host string, port int, user string, key []byte, remoteDirBackup string, fileName string) error {
	var (
		err        error
		sftpClient *sftp.Client
	)
	localDir += fileName

	var remoteFileName = path.Base(fileName)

	sftpClient, err = Connect(user, localDir, remoteDir, port, host, key)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sftpClient.Close()

	buffer, err := ioutil.ReadFile(localDir)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//joga para pasta de backup
	if remoteDirBackup != "" {
		dstFileBackup, err := sftpClient.Create(path.Join(remoteDirBackup, remoteFileName))
		if err != nil {
			log.Println(err)
			return err
		}
		dstFileBackup.Write(buffer)
		defer dstFileBackup.Close()
	}

	//joga para pasta de upload.
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Println(err)
		return err
	}
	dstFile.Write(buffer)

	defer dstFile.Close()
	return err

}

func Connect(userSftp string, localFilePath string, remoteDir string, portSftp int, hostSftp string, key []byte) (*sftp.Client, error) {
	var (
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
		singer       ssh.Signer
	)
	addr = fmt.Sprintf("%s:%d", hostSftp, portSftp)

	singer, err = ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil
	}

	auths := []ssh.AuthMethod{ssh.PublicKeys(singer)}

	clientConfig = &ssh.ClientConfig{
		User:            userSftp,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	clientConfig.SetDefaults()

	// connet to ssh
	sshClient, err = ssh.Dial("tcp", addr, clientConfig)

	if err != nil {
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("Failed to dial: %s", err)
		}
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}
