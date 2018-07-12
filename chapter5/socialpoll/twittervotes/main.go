package main

import (
	"log"

	"github.com/nsqio/go-nsq"

	"gopkg.in/mgo.v2"
)

var db *mgo.Session

type poll struct {
	Options []string
}

func main() {}

func dialdb() error {
	var err error
	log.Println("MongoDBにダイヤル中: localhost")
	db, err = mgo.Dial("localhost")
	return err
}

func closedb() {
	db.Close()
	log.Println("データベース接続が閉じられました")
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}
	iter.Close()
	return options, iter.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote)) // 投票内容をパブリッシュ
		}
		log.Println("Publiser: 停止中です")
		pub.Stop()
		log.Println("Publiser: 停止しました")
		stopchan <- struct{}{}
	}()
	return stopchan
}
