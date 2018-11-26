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

//MoveFile realize the connection and send the files through sftp server.
func MoveFile(localDir, remoteDir, host, port, user, key, remoteDirBackup, fileName, fileNameBackup, flagRemoveFiles string) error {
	var (
		err        error
		sftpClient *sftp.Client
	)

	//open the connection with the server
	sftpClient, err = Connect(host, port, user, key)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer sftpClient.Close()

	//create the buffer of the file
	buffer, err := ioutil.ReadFile(localDir + fileName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	dstFile, err := sftpClient.Create(path.Join(remoteDir, fileName))
	if err != nil {
		log.Println(err)
		return err
	}
	dstFile.Write(buffer)

	defer dstFile.Close()

	if flagRemoveFiles == "Y" {
		os.Remove(localDir + fileName)
	}

	//Move the backup file to the specified folder
	if remoteDirBackup != "" {

		//create the buffer of the backup file
		bufferBkp, err := ioutil.ReadFile(localDir + fileNameBackup)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		//realize the upload of the backup file
		dstFileBackup, err := sftpClient.Create(path.Join(remoteDirBackup, fileNameBackup))
		if err != nil {
			log.Println(err)
			return err
		}
		dstFileBackup.Write(bufferBkp)
		defer dstFileBackup.Close()

		if flagRemoveFiles == "Y" {
			os.Remove(localDir + fileNameBackup)
		}
	}

	return err

}

//Connect realize the connection through sftp server.
func Connect(host, port, user, key string) (*sftp.Client, error) {
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

	// connect to ssh
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
