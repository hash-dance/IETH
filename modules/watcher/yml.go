package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/guowenshuai/ieth/conf"
	apicontext "github.com/guowenshuai/ieth/modules/context"
	"github.com/sirupsen/logrus"
)

func Watch(ctx *apicontext.APIContext) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}
	defer watcher.Close()


	done := make(chan bool)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("catch panic %s", err)
			}
		}()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				logrus.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					logrus.Println("modified file:", event.Name)
					if err := conf.LoadYaml("conf.yaml", ctx.Config); err != nil {
						logrus.Errorf("reload config err: %s\n", err.Error())
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					logrus.Errorf("reload config err: %s\n", err.Error())
				}
				logrus.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("conf.yaml")
	if err != nil {
		logrus.Fatal(err)
	}
	<-done
}
