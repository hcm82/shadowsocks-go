package main

import (
	"encoding/json"
	muconfig "github.com/orvice/shadowsocks-go/mu/config"
	"github.com/orvice/shadowsocks-go/mu/user"
	"gopkg.in/redis.v3"
)

var Redis = new(RedisClient)

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) SetClient(client *redis.Client) {
	r.client = client
}

func (r *RedisClient) GetUserInfo(u user.User) (user.UserInfo, error) {
	var user user.UserInfo
	val, err := r.client.Get(genUserInfoKey(u.GetUserInfo())).Result()
	if err != nil {
		return user, err
	}
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

func (r *RedisClient) StoreUser(user user.UserInfo) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = r.client.Set(genUserInfoKey(user), data, 0).Err()
	return err
}

func (r *RedisClient) Exists(u user.User) (bool, error) {
	return r.client.Exists(genUserInfoKey(u.GetUserInfo())).Result()
}

func (r *RedisClient) Del(u user.User) error {
	return r.client.Del(genUserInfoKey(u.GetUserInfo())).Err()
}

func (r *RedisClient) ClearAll() error {
	return r.client.FlushAll().Err()
}

func InitRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     muconfig.Conf.Redis.Host,
		Password: muconfig.Conf.Redis.Pass, // no password set
		DB:       muconfig.Conf.Redis.Db,   // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return err
	}
	Log.Info(pong)
	Redis.SetClient(client)
	// set storage
	SetStorage(Redis)
	return nil
}
