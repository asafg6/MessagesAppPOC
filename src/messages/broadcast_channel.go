package messages

import (
	"sync"
)

type BroadcastChannel struct {
	channelItem *ChannelItem
	lock sync.Mutex
}

func (channel *BroadcastChannel) Listen() *ChannelItem {
	channel.lock.Lock()
	defer channel.lock.Unlock()
	return channel.channelItem
}


func (channel *BroadcastChannel) Publish(data interface{}) {
	channel.lock.Lock()
	defer channel.lock.Unlock()
	newItem := MakeNewChannelItem(data)
	channel.channelItem.setNext(newItem)
	channel.channelItem = newItem
}


func MakeNewBroadcastChannel() *BroadcastChannel {
	m := sync.Mutex{}
	nilItem := MakeNewChannelItem(nil)
	return &BroadcastChannel{channelItem: nilItem, lock: m}
}

type ChannelItem struct {
	next *ChannelItem
	cond *sync.Cond
	data interface{}
}

func (channelItem *ChannelItem) GetNextMessageOrWait() *ChannelItem {
	channelItem.cond.L.Lock()
	defer channelItem.cond.L.Unlock()
	for channelItem.next == nil {
		channelItem.cond.Wait()
	}
	return channelItem.next
}

func (channelItem *ChannelItem) GetNextMessageOrWaitWithClose(close <- chan bool) *ChannelItem {

	waitChan := make( chan struct{})
	go func() {
		channelItem.cond.L.Lock()
		defer channelItem.cond.L.Unlock()
		channelItem.cond.Wait()
		waitChan <- struct {}{}
	} ()
	select {
	case _, _ = <- waitChan:
		return channelItem.next
	case _, ok := <- close:
		if ok {
			return nil
		}

	}
	return channelItem.next
}

func (channelItem *ChannelItem) setNext(next *ChannelItem) {
	channelItem.cond.L.Lock()
	defer channelItem.cond.L.Unlock()
	channelItem.next = next
	channelItem.cond.Broadcast()

}

func (channelItem *ChannelItem) GetData() interface{} {
	return channelItem.data
}

func MakeNewChannelItem(data interface{}) *ChannelItem {
	m := sync.Mutex{}
	cond := sync.NewCond(&m)
	return &ChannelItem{data:data, cond:cond}
}
