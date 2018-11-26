package model

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//MoveFile realiza a conexão e envia o arquivo via sftp para um servidor.
func MoveFile(localDir, remoteDir, host, port, user, key, remoteDirBackup, fileName, fileNameBkp, flgRemoverArq string) error {
	var (
		err        error
		sftpClient *sftp.Client
	)

	localDir += fileName

	var remoteFileName = path.Base(fileName)

	sftpClient, err = Connect(localDir, remoteDir, host, port, user, key)
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

		var remoteFileNameBkp = path.Base(fileNameBkp)
		dstFileBackup, err := sftpClient.Create(path.Join(remoteDirBackup, remoteFileNameBkp))
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

	if flgRemoverArq == "Y" {
		os.Remove(localDir + fileName)
		os.Remove(localDir + fileNameBkp)
	}

	return err

}

//Connect realiza conexão via sftp em algum servidor.
func Connect(localDir string, remoteDir string, host string, port string, user string, key string) (*sftp.Client, error) {
	var (
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
		singer       ssh.Signer
	)
	portSftp, err := strconv.Atoi(port)
	bytes := []byte(key)

	addr = fmt.Sprintf("%s:%d", host, portSftp)

	singer, err = ssh.ParsePrivateKey(bytes)
	if err != nil {
		return nil, nil
	}

	auths := []ssh.AuthMethod{ssh.PublicKeys(singer)}

	clientConfig = &ssh.ClientConfig{
		User:            user,
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
