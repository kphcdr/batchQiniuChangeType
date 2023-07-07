package main

import (
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

const ak = "you ak is here"
const sk = "your sk is here"
const bucket = "your bucket"

// 去掉后面的7个0，就是一个时间戳，只会处理这个时间以前的文件
const endTime = 16794925030000000

// 0 表示标准存储，1 表示低频访问存储，2 表示归档存储，3 表示深度归档存储
const toType = 1

var bucketManage *storage.BucketManager

func main() {
	bucketManage = getBucketManager()
	list()
}
func getBucketManager() *storage.BucketManager {
	mac := qbox.NewMac(ak, sk)
	cfg := storage.Config{
		UseHTTPS: true,
	}
	return storage.NewBucketManager(mac, &cfg)
}
func list() {
	limit := 1000
	prefix := ""
	delimiter := ""
	marker := ""
	page := 1
	for {
		fmt.Println("page ", page)
		entries, _, nextMarker, hasNext, err := bucketManage.ListFiles(bucket, prefix, delimiter, marker, limit)
		if err != nil {
			fmt.Println("list error,", err)
			break
		}
		//print entries
		for _, entry := range entries {
			if entry.PutTime > endTime {
				continue
			}

			if entry.Type == toType {
				continue
			}

			fmt.Println("changeType start: ", entry.Key)
			go changeType(entry)

		}
		if hasNext {
			marker = nextMarker
		} else {
			break
		}
		page++
	}
}

func changeType(entry storage.ListItem) error {
	return bucketManage.ChangeType(bucket, entry.Key, toType)

}
