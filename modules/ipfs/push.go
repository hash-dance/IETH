package ipfs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	db "github.com/guowenshuai/ieth/db/mongo"
	apicontext "github.com/guowenshuai/ieth/modules/context"
	"github.com/guowenshuai/ieth/types"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	FILE = 2
	DIR  = 1
)

func Push(apiCtx *apicontext.APIContext, path string, recursive bool) {
	add := apiCtx.Config.Ipfs.Address
	logrus.Debugf("ipfs api %s\n", add)
	sh := NewSH(apiCtx.Context, add)
	collection := apiCtx.MongoClient.Collection(db.IpfsCollectionName)

	// NewSH(apiCtx.Context, )
	fileinfo, err := os.Stat(path)
	if err != nil {
		logrus.Errorf("error: %s", err.Error())
		return
	}
	if fileinfo.IsDir() {
		if !recursive { // 上传一个目录
			cid, err := sh.AddDir(path)
			if err != nil {
				logrus.Errorf("error: %s", err)
				return
			}
			inserttime := time.Now()
			_, err = insertOne(collection, types.IpfsData{
				PayloadCid:     cid,
				FileName:       fileinfo.Name(),
				Dir:            path,
				IsDir:          true,
				FileFormat:     filepath.Ext(fileinfo.Name()),
				FileSizeBytes:  fileinfo.Size(),
				Group:          "",
				CuratedDataset: "",
				CreatedTime:    inserttime,
				UpdatedTime:    inserttime,
			})
			write2log(path, cid, 0, DIR)
			// fmt.Printf("%s %s %d %d\n", path, cid, 0, DIR)
			return
		} else { // 递归上传下面的文件
			files := walkDirs(path)
			for _, f := range files {
				pushOneFile(sh, f.path, f.info, collection)
			}
		}
	} else {
		pushOneFile(sh, path, fileinfo, collection)
	}
	return
}

func pushOneFile(sh *shell.Shell, path string, fileinfo os.FileInfo, col *mongo.Collection) {
	// todo 检查是已经传到ipfs, 在数据库中存在
	// for _, eachItem := range existFiles { // 存在，则不继续上传
	// 	if eachItem == path {
	// 		fmt.Printf("path %s already pushed in log\n", path)
	// 		return
	// 	}
	// }
	fh, err := os.Open(path)
	defer fh.Close()
	if err != nil {
		logrus.Errorf("error: %s", err)
		return
	}
	cid, err := sh.Add(fh)
	if err != nil {
		logrus.Errorf("error: %s", err)
		return
	}
	inserttime := time.Now()
	_, err = insertOne(col, types.IpfsData{
		PayloadCid:     cid,
		FileName:       fileinfo.Name(),
		Dir:            path,
		IsDir:          false,
		FileFormat:     filepath.Ext(fileinfo.Name()),
		FileSizeBytes:  fileinfo.Size(),
		Group:          regexp.MustCompile(`^[a-zA-Z]+`).FindString(fileinfo.Name()),
		CuratedDataset: "",
		CreatedTime:    inserttime,
		UpdatedTime:    inserttime,
	})
	if err != nil {
		logrus.Errorf("insert %s err: %s\n", cid, err.Error())
	}

	write2log(path, cid, fileinfo.Size(), FILE)
	// fmt.Printf("%s %s %d %d\n", path, cid, fileinfo.Size(), FILE)
	return
}

type fileInfo struct {
	path string
	info os.FileInfo
}

func walkDirs(dir string) []*fileInfo {
	files := make([]*fileInfo, 0)
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, &fileInfo{
					path: path,
					info: info,
				})
				// fmt.Println(path, info.Size())
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}
	return files
}

func write2log(fullname, hash string, size, ftype int64) {
	logrus.Infof("%s %s %d %d\n", fullname, hash, size, ftype)
}

func insertOne(col *mongo.Collection, d interface{}) (*mongo.InsertOneResult, error) {
	return col.InsertOne(context.Background(), d)
}
