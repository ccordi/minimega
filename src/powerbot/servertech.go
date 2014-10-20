package main

import (
	"fmt"
	"telnet"
)

var (
	prompt = "Switched CDU: "
)

// Implementation of a specific Server Tech PDU.
// This should hopefully work for any Server Tech device with an
// up-to-date firmware interface
type ServerTechPDU struct {
	username string
	password string
	c        *telnet.Conn
}

// Create a new instance. The port should point to the device's telnet
// CLI port, which appears to usually be 5214.
func NewServerTechPDU(host, port, username, password string) (PDU, error) {
	var tp ServerTechPDU
	conn, err := telnet.Dial("tcp", host+":"+port)
	if err != nil {
		return tp, err
	}
	tp.c = conn
	tp.username = username
	tp.password = password
	return tp, err
}

// Convenience function, log in.
func (p ServerTechPDU) login() error {
	// wait for login prompt
	_, err := p.c.ReadUntil("Username: ")
	if err != nil {
		return err
	}
	cmd := fmt.Sprintf("%s\r\n", p.username)
	_, err = p.c.Write([]byte(cmd))
	if err != nil {
		return err
	}
	_, err = p.c.ReadUntil("Password: ")
	if err != nil {
		return err
	}
	cmd = fmt.Sprintf("%s\r\n", p.password)
	_, err = p.c.Write([]byte(cmd))
	if err != nil {
		return err
	}
	return nil
}

// Convenience function, log out.
func (p ServerTechPDU) logout() error {
	// send a blank line to make sure we get a prompt
	_, err := p.c.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	_, err = p.c.ReadUntil(prompt)
	if err != nil {
		return err
	}
	_, err = p.c.Write([]byte("exit\r\n"))
	if err != nil {
		return err
	}
	return nil
}

// Log in, say "on <port>" for every port
// the user specified, and log out.
func (p ServerTechPDU) On(ports map[string]string) error {
	p.login()
	for _, port := range ports {
		_, err := p.c.ReadUntil(prompt)
		if err != nil {
			return err
		}
		_, err = p.c.Write([]byte(fmt.Sprintf("on %s\r\n", port)))
		if err != nil {
			return err
		}
	}
	p.logout()
	return nil
}

// Log in, say "off <port>" for every port
// the user specified, and log out.
func (p ServerTechPDU) Off(ports map[string]string) error {
	p.login()
	for _, port := range ports {
		_, err := p.c.ReadUntil(prompt)
		if err != nil {
			return err
		}
		_, err = p.c.Write([]byte(fmt.Sprintf("off %s\r\n", port)))
		if err != nil {
			return err
		}
	}
	p.logout()
	return nil
}

// Log in, say "reboot <port>" for every port
// the user specified, and log out.
func (p ServerTechPDU) Cycle(ports map[string]string) error {
	p.login()
	for _, port := range ports {
		_, err := p.c.ReadUntil(prompt)
		if err != nil {
			return err
		}
		_, err = p.c.Write([]byte(fmt.Sprintf("reboot %s\r\n", port)))
		if err != nil {
			return err
		}
	}
	p.logout()
	return nil
}

// This isn't done yet because I don't want to parse that crap yet.
// Should be pretty easy really, call "loadctl status -o" and then
// look for the lines corresponding to the listed nodes. Print only
// those lines.
func (p ServerTechPDU) Status(ports map[string]string) error {
	fmt.Println("not yet implemented")
	return nil
	// doesn't work right
	/*
		p.login()
		_, err := p.c.ReadUntil("$> ")
		if err != nil {
			return err
		}
		_, err = p.c.Write([]byte("loadctl status -o\r\n"))
		if err != nil {
			return err
		}
		result, err := p.c.ReadUntil("$> ")
		if err != nil {
			return err
		}
		fmt.Println(string(result))
		p.logout()
		return nil
	*/
}