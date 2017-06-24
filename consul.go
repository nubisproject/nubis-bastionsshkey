package main

import (
	"bytes"
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
	"os/exec"
	"strings"
)

type ConsulClient struct {
	client *consul.KV
}

type ConsulEntries struct {
	Users []LDAPUserObject
	Group string
}

func (c *ConfigOptions) getConsulACLToken() string {
	var (
		out    bytes.Buffer
		stdErr bytes.Buffer
	)

	unicredsPath := c.UnicredsPath
	cmdArgs := []string{
		"--region", c.Region,
		"get", fmt.Sprintf("%s/%s/consul/acl_token", c.Service, c.Environment),
		"-E", fmt.Sprintf("environment:%s", c.Environment),
		"-E", fmt.Sprintf("service:%s", c.Service),
		"-E", fmt.Sprintf("region:%s", c.Region),
	}
	cmd := exec.Command(unicredsPath, cmdArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &stdErr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("getConsulACLToken unicreds error: %s", stdErr.String())
	}
	output := strings.TrimSpace(out.String())
	return output
}

func NewConsulClient(config Configuration) *ConsulClient {
	conf := consul.DefaultConfig()
	conf.Address = config.Consul.Server

	if useLambda {
		conf.Token = config.Consul.Token
	}

	client, err := consul.NewClient(conf)
	if err != nil {
		log.Fatal("Failed to connect")
	}
	return &ConsulClient{client.KV()}

}
func (c *ConsulClient) Put(obj LDAPUserObject, conf Configuration, user_class string) {
	found := false
	uidKey := fmt.Sprintf("%s/%s/%s/uid", conf.Consul.Namespace, user_class, obj.Uid)
	uidByteVal := []byte(obj.Uid)
	p := &consul.KVPair{Key: uidKey, Value: uidByteVal}
	_, err := c.client.Put(p, nil)
	if err != nil {
		panic(err)
	}
	email := fmt.Sprintf("%s/%s/%s/email", conf.Consul.Namespace, user_class, obj.Uid)
	emailByteVal := []byte(obj.Mail)
	p = &consul.KVPair{Key: email, Value: emailByteVal}
	_, err = c.client.Put(p, nil)
	if err != nil {
		panic(err)
	}
	ldapByteVal := strings.Join(obj.SshPublicKey, "\n")
	key := fmt.Sprintf("%s/%s/%s/sshPublicKey", conf.Consul.Namespace, user_class, obj.Uid)
	consulByteVal, _, consulByteValErr := c.client.Get(key, nil)
	if consulByteValErr == nil {
		found = true
	}

	consulKeyLength := 0
	if consulByteVal != nil {
		consulKeyLength = len(string(consulByteVal.Value))
	}

	if consulKeyLength == 0 {
		if found == true {
			c.client.Delete(key, nil)
		}
	}

	if found == false {
		p = &consul.KVPair{Key: key, Value: []byte(ldapByteVal)}
		_, err = c.client.Put(p, nil)
		if err != nil {
			panic(err)
		}
	}
	if found == true {
		var consulString string
		if consulByteVal != nil {
			consulString = string(consulByteVal.Value[:])
		} else {
			consulString = ""
		}
		if ldapByteVal != consulString {
			log.Printf("Updating Key for: %s", obj.Uid)
			p = &consul.KVPair{Key: key, Value: []byte(ldapByteVal)}
			_, err = c.client.Put(p, nil)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (c *ConsulClient) GetKValues(keys string) {
	pairs, _, err := c.client.List(keys, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v pairs returned\n", len(pairs))

	for i := range pairs {
		pair := pairs[i]
		log.Printf("Key: %s\n", pair.Key)
		log.Printf("Value: %s\n", pair.Value)
	}
	return
}

func GetConsulClient(conf Configuration) *ConsulClient {
	c := NewConsulClient(conf)
	return c
}

func SyncLDAPToConsul(userClass string, usersSet []LDAPUserObject, noop bool, c *ConsulClient, conf Configuration) {
	namespace := fmt.Sprintf("%s/%s/", conf.Consul.Namespace, userClass)
	globalAdminsConsul, _, _ := c.client.Keys(namespace, "/", nil)
	for _, consulAdminUser := range globalAdminsConsul {
		keyPath := consulAdminUser
		consulAdminUser = TrimSuffix(consulAdminUser, "/")
		usernameSplit := strings.Split(consulAdminUser, "/")
		consulAdminUsername := usernameSplit[len(usernameSplit)-1]
		ignoreConsulUser := IgnoreUserLDAPUserObjects(usersSet, consulAdminUsername)
		if consulAdminUsername == userClass {
			continue
		}
		if ignoreConsulUser == false {
			if noop == false {
				log.Printf("Removing %s from %s", consulAdminUsername, userClass)
				log.Printf("KeyPath %s", keyPath)
				c.client.DeleteTree(keyPath, nil)
			} else {
				log.Printf("Should remove %s from %s", consulAdminUsername, userClass)
			}
		}

	}
}
