package main

import (
	"crypto/tls"
	"crypto/x509"
	"go.mozilla.org/mozldap"
	"log"
	"strings"
)

type LDAPUserObject struct {
	Dn           string
	Uid          string
	SshPublicKey []string
}

func getGroupMembers(conf Configuration) ([]LDAPUserObject, []LDAPUserObject) {
	var (
		cli mozldap.Client
		err error
	)

	tlsconf := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: conf.LdapServer.LDAPInsecure,
		ServerName:         cli.Host,
	}
	if len(conf.LdapServer.TLSCrt) > 0 && len(conf.LdapServer.TLSKey) > 0 {
		cert, err := tls.X509KeyPair([]byte(conf.LdapServer.TLSCrt), []byte(conf.LdapServer.TLSKey))
		if err != nil {
			log.Fatal(err)
		}
		tlsconf.Certificates = []tls.Certificate{cert}
	}

	if len(conf.LdapServer.CACrt) > 0 {
		ca := x509.NewCertPool()
		if ok := ca.AppendCertsFromPEM([]byte(conf.LdapServer.CACrt)); !ok {
			log.Fatal("failed to import CA Certificate")
		}
		tlsconf.RootCAs = ca
	}
	cli, err = mozldap.NewClient(
		conf.LdapServer.LDAPHost,
		conf.LdapServer.LDAPBindUser,
		conf.LdapServer.LDAPBindPassword,
		tlsconf,
		conf.LdapServer.StartTLS,
	)
	if err != nil {
		log.Fatal(err)
	}
	globalAdmins, err := cli.GetUsersInGroups(conf.LdapServer.GlobalAdmins)
	globalAdminsSlice := make([]LDAPUserObject, len(globalAdmins))
	for i, g_entry := range globalAdmins {
		globalAdminsSlice[i] = getUserByDn(g_entry, cli)
	}
	sudoUsers, err := cli.GetUsersInGroups(conf.LdapServer.SudoUsers)
	sudoUsersSlice := make([]LDAPUserObject, len(sudoUsers))
	for i, g_entry := range sudoUsers {
		sudoUsersSlice[i] = getUserByDn(g_entry, cli)
	}

	cli.Close()
	return globalAdminsSlice, sudoUsersSlice

}

func getShortDn(dn string) string {
	short_dn := strings.Split(dn, ",")[0]
	return short_dn
}

func getUserByDn(dn string, cli mozldap.Client) LDAPUserObject {
	tmp := LDAPUserObject{Dn: dn}
	short_dn := getShortDn(dn)
	tmp.Uid, _ = cli.GetUserId(short_dn)
	tmp.SshPublicKey, _ = cli.GetUserSSHPublicKeys(short_dn)
	return tmp
}
