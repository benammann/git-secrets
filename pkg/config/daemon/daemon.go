package daemon

import (
	"fmt"
	cli_config "github.com/benammann/git-secrets/pkg/config/cli"
	config_parser "github.com/benammann/git-secrets/pkg/config/parser"
	"github.com/benammann/git-secrets/pkg/render"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Daemon struct {
	fileWatches map[string]string
}

func NewDaemon() *Daemon {
	return &Daemon{}
}

func (d *Daemon) HandleFileChange(fileName string) {

	parsedRepo, errParse := config_parser.ParseRepository(fileName, make(map[string]string))
	if errParse != nil {
		return
	}

	watchContextName := d.fileWatches[fileName]
	watchContext, watchContextErr := parsedRepo.SetSelectedContext(watchContextName)

	if watchContextErr != nil {
		fmt.Printf("could not resolve context for %s: %s", fileName, watchContextErr.Error())
		return
	}

	fileRenderer := render.NewRenderingEngine(parsedRepo)
	for _, fileToRender := range watchContext.FilesToRender {
		_, errRender := fileRenderer.WriteFile(fileToRender)
		if errRender != nil {
			fmt.Printf("error: could not render %s: %s in %s\n", fileToRender.FileIn, errRender.Error(), fileToRender)
		} else {
			fmt.Println(fileToRender.FileIn, "rendered as", fileToRender.FileOut)
		}
	}
}

func (d *Daemon) Run() error {

	var currentWatches []string

	d.fileWatches = viper.GetStringMapString(cli_config.DaemonWatches)

	if len(d.fileWatches) <= 0 {
		return fmt.Errorf("there are no files to watch")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	addWatch := func(fileName string, watchContext string) {
		errAdd := watcher.Add(fileName)
		if errAdd != nil {
			fmt.Printf("error: could not add watch %s: %s\n", fileName, errAdd.Error())
			return
		}
		currentWatches = append(currentWatches, fileName)
	}

	setWatches := func() {
		for fileToWatch, watchAsContext := range d.fileWatches {
			addWatch(fileToWatch, watchAsContext)
		}
		fmt.Println("Watching", strings.Join(currentWatches, ", "))
	}

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				d.HandleFileChange(event.Name)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}

	}()

	viper.OnConfigChange(func(in fsnotify.Event) {

		d.fileWatches = viper.GetStringMapString(cli_config.DaemonWatches)

		for _, currentWatch := range currentWatches {
			_ = watcher.Remove(currentWatch)
		}

		currentWatches = []string{}

		setWatches()

	})

	setWatches()

	<-done

	return nil
}
