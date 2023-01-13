package ssh

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path/filepath"
	"time"
)

type SSHClient struct {
	User string
	Pass string
	Key string			// 私钥文件路径，公钥需写入服务器的authorized_keys
	Host string
	Port int
	Timeout int
	Session *ssh.Session
	Client *ssh.Client
}

// UploadFile 上传文件
func (obj *SSHClient)UploadFile(localFile string, remoteFile string) error {
	// 初始化client、session
	_,err := obj.CreateSession()
	if err != nil {
		return err
	}

	// 创建一个ftp客户端
	ftpClient, err := sftp.NewClient(obj.Client)
	if err != nil {
		return err
	}
	defer ftpClient.Close()

	// 远程文件
	dstFile, err := ftpClient.Create(remoteFile)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 本地文件
	b,err := os.ReadFile(localFile)
	if err != nil {
		return err
	}

	// 传输
	_,err = dstFile.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// RunScriptFile 执行脚本文件
func (obj *SSHClient)RunScriptFile(file string) (string, error) {
	session,err := obj.CreateSession()
	if err != nil {
		return "", err
	}
	// 关闭会话
	defer session.Close()

	b,err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	path := "/tmp/" + filepath.Base(file)
	cmd := fmt.Sprintf("echo '%s'>%s && bash %s", string(b), path, path)
	buf, err := session.CombinedOutput(cmd)
	if err != nil {
		return "",err
	}
	return string(buf), nil
}

// RunOneCmd 执行单个命令，返回结果
func (obj *SSHClient)RunOneCmd(cmd string) (string, error) {
	session,err := obj.CreateSession()
	if err != nil {
		return "", err
	}
	// 关闭会话
	defer session.Close()

	buf, err := session.CombinedOutput(cmd)
	if err != nil {
		return "",err
	}
	return string(buf), nil
}

// CreateSession 创建ssh会话
func (obj *SSHClient)CreateSession() (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	if obj.Key == "" {
		auth = append(auth, ssh.Password(obj.Pass))
	} else {
		pemBytes, err := os.ReadFile(obj.Key)
		if err != nil {
			return nil, err
		}

		var signer ssh.Signer
		if obj.Pass == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(obj.Pass))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}

	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}

	clientConfig = &ssh.ClientConfig{
		User:    obj.User,
		Auth:    auth,
		Timeout: time.Duration(obj.Timeout) * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", obj.Host, obj.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	obj.Client = client

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, err
	}
	obj.Session = session

	return session, nil
}
