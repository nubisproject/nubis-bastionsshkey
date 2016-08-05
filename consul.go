package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
)

type ConsulClient struct {
	client *consul.KV
}

func NewConsulClient(config Configuration) *ConsulClient {
	conf := consul.DefaultConfig()
	conf.Address = config.Consul.Server
	client, err := consul.NewClient(conf)
	if err != nil {
		log.Fatal("Failed to connect")
	}
	return &ConsulClient{client.KV()}

}
func (c *ConsulClient) Put(obj LDAPUserObject, conf Configuration, user_class string) {
	uidKey := fmt.Sprintf("%s/%s/%s/uid", conf.Consul.Namespace, user_class, obj.Uid)
	uidByteVal := []byte(obj.Uid)
	p := &consul.KVPair{Key: uidKey, Value: uidByteVal}
	_, err := c.client.Put(p, nil)
	if err != nil {
		panic(err)
	}
	for i, sshPublicKey := range obj.SshPublicKey {
		key := fmt.Sprintf("%s/%s/%s/sshPublicKey_Gen%d", conf.Consul.Namespace, user_class, obj.Uid, i)
		byteVal := []byte(sshPublicKey)
		p := &consul.KVPair{Key: key, Value: byteVal}
		_, err := c.client.Put(p, nil)
		if err != nil {
			panic(err)
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
