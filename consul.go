package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"log"
	"strconv"
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
	currentSerial := c.Serial(conf)
	nextSerial, _ := strconv.Atoi(currentSerial)
	nextSerial++
	uidKey := fmt.Sprintf("%s/%s/%d/%s/uid", conf.Consul.Namespace, user_class, nextSerial, obj.Uid)
	uidByteVal := []byte(obj.Uid)
	p := &consul.KVPair{Key: uidKey, Value: uidByteVal}
	_, err := c.client.Put(p, nil)
	if err != nil {
		panic(err)
	}
	for i, sshPublicKey := range obj.SshPublicKey {
		key := fmt.Sprintf("%s/%s/%d/%s/sshPublicKey_Gen%d", conf.Consul.Namespace, user_class, nextSerial, obj.Uid, i)
		byteVal := []byte(sshPublicKey)
		p := &consul.KVPair{Key: key, Value: byteVal}
		_, err := c.client.Put(p, nil)
		if err != nil {
			panic(err)
		}

	}

}

func (c *ConsulClient) CreateSerial(conf Configuration) {
	serialPath := fmt.Sprintf("%s/serial", conf.Consul.Namespace)
	p := &consul.KVPair{Key: serialPath, Value: []byte("0")}
	_, err := c.client.Put(p, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *ConsulClient) IncrementSerial(conf Configuration) int {
	currentSerial := c.Serial(conf)
	newSerial, _ := strconv.Atoi(currentSerial)
	newSerial++
	serialPath := fmt.Sprintf("%s/serial", conf.Consul.Namespace)
	serialByte := []byte(strconv.Itoa(newSerial))
	p := &consul.KVPair{Key: serialPath, Value: serialByte}
	_, err := c.client.Put(p, nil)
	if err != nil {
		log.Fatal(err)
	}

	return newSerial
}

func (c *ConsulClient) Serial(conf Configuration) string {
	currentSerial := []byte("0")
	serialPath := fmt.Sprintf("%s/serial", conf.Consul.Namespace)
	pairs, _, err := c.client.List(serialPath, nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(pairs) == 0 {
		c.CreateSerial(conf)
	} else {
		currentSerial = []byte(pairs[0].Value)
	}

	return string(currentSerial)
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
