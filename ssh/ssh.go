package opssh

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
func PublicKey(keyAsString string) ssh.AuthMethod {
	key, err := ssh.ParsePrivateKey([]byte(keyAsString))
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

type Ssh struct {
	Conn *ssh.Client
}

func NewSession(host, publicKey string) (*Ssh, error) {
	// Read: http://blog.ralch.com/tutorial/golang-ssh-connection/
	s := &Ssh{}
	// log.Infoln("KEYFILE:", publicKey)
	// log.Infoln("host:", host)

	sshConfig := &ssh.ClientConfig{
		User: "rancher",
		Auth: []ssh.AuthMethod{
			// PublicKeyFile(publicKeyFile),
			PublicKey(publicKey),
		},
	}
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", host, 22), sshConfig)
	if err != nil {
		log.Errorln("Failed to dial!", err)
		return nil, fmt.Errorf("Failed to dial: %s", err)
	}
	s.Conn = connection
	return s, nil
}

func (s *Ssh) newSession(w io.Writer) (*ssh.Session, error) {
	session, err := s.Conn.NewSession()
	if err != nil {
		log.Errorln("Failed to create session", err)
		return nil, fmt.Errorf("Failed to create session: %s", err)
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		log.Errorln("request for pseudo terminal failed", err)
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}
	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Errorln("Unable to setup stdout for session", err)
		session.Close()
		return nil, fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(w, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Errorln("Unable to setup stderr for session", err)
		session.Close()
		return nil, fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(w, stderr)
	return session, nil
}

func (s *Ssh) Run(cmd string, w io.Writer) error {
	session, err := s.newSession(w)
	if err != nil {
		log.WithError(err).Errorln("Couldn't create new ssh session")
		return err
	}
	defer session.Close()

	err = session.Run(cmd)
	if serr, ok := err.(*ssh.ExitError); ok {
		// todo: return exit status or something so caller of this method can do something about it.
		log.WithError(serr).Infoln("Exit error, ignoring for now?")
	} else {
		return err
	}
	return nil
}

func (s *Ssh) Close() error {
	return s.Conn.Close()
}
