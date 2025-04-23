package models

import "time"

type Proxy struct {
	ID        string    `bson:"_id"`
	ProxyIP   string    `bson:"proxy_ip"`
	Live      int       `bson:"live"`
	Username  string    `bson:"username"`
	CountUsed int       `bson:"count_used"`
	Enable    int       `bson:"enable"`
	Status    int       `bson:"status"`
	LastCall  time.Time `bson:"last_call"`
}
