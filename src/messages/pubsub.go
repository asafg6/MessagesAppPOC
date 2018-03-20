package messages

import (
	"sync"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
)

type ChannelHandle struct {
	channel *BroadcastChannel
	pubsub *redis.PubSub
	subscribers int
}


type PubSubClient struct {
	client *redis.Client
	channels map[string]*ChannelHandle
}

var mutex = &sync.Mutex{}

func (p *PubSubClient) Subscribe(channel string) *BroadcastChannel {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := p.channels[channel]
	if !ok {
		newChannel := MakeNewBroadcastChannel()
		pubsub := p.client.Subscribe(channel)
		channelHandle := ChannelHandle{channel: newChannel, pubsub: pubsub, subscribers: 0}
		go func(){
			log.Printf("Listening to redis channel %v", channel)
			for redisMessage := range pubsub.Channel() {
				log.Printf("Got message from redis: %v", redisMessage.Payload)
				message := Message{}
				if err := json.Unmarshal([]byte(redisMessage.Payload), &message); err != nil {
					log.Println(err)
					continue
				}
				newChannel.Publish(message)
			}
			log.Printf("Stopped listening to redis channel: %v", channel)
		}()
		p.channels[channel] = &channelHandle
	}
	p.channels[channel].subscribers++
	log.Printf("Channel: %v, Subscribers: %v", channel, p.channels[channel].subscribers)
	return p.channels[channel].channel
}


func (p *PubSubClient) UnSubscribe(channel string) {
	mutex.Lock()
	defer mutex.Unlock()
	handle, ok := p.channels[channel]
	if !ok {
		return
	}
	handle.subscribers--
	log.Printf("Channel: %v, Subscribers: %v", channel, handle.subscribers)

	if handle.subscribers <= 0 {
		handle.pubsub.Close()
		delete(p.channels, channel)
		log.Printf("channel %v deleted", channel)
	}
}


func MakeNewClient(Addr string) *PubSubClient {
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	channels := make(map[string]*ChannelHandle)
	return &PubSubClient{client: client, channels: channels}
}

