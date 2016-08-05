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

func getGroupMembers(conf Configuration) []LDAPUserObject {
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
	group_members, err := cli.GetUsersInGroups(conf.LdapServer.SearchGroups)
	ldapUsers := make([]LDAPUserObject, len(group_members))
	for i, g_entry := range group_members {
		ldapUsers[i] = getUserByDn(g_entry, cli)
	}

	cli.Close()
	return ldapUsers

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
