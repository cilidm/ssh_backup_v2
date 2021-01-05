package client

import (
	"fmt"
	"github.com/cilidm/toolbox/logging"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"ssh_backup/app/config"
	"time"
)

var (
	err       error
	Client    *sftp.Client
	TarClient *sftp.Client
)

func Instance() *sftp.Client {
	if Client != nil {
		return Client
	}
	Client, err = connect(config.Conf.Ssh.SourceUser, config.Conf.Ssh.SourcePwd, config.Conf.Ssh.SourceHost, config.Conf.Ssh.SourcePort)
	if err != nil {
		logging.Error(err)
		return nil
	}
	return Client
}

func TargetInstance() *sftp.Client {
	if TarClient != nil {
		return TarClient
	}
	TarClient, err = connect(config.Conf.Ssh.TargetUser, config.Conf.Ssh.TargetPwd, config.Conf.Ssh.TargetHost, config.Conf.Ssh.TargetPort)
	if err != nil {
		logging.Error(err)
		return nil
	}
	return TarClient
}

func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 15 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}
