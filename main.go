package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type Task struct {
	Name     string    `toml:"-"`
	Link     string    `toml:"link"`
	Deadline time.Time `toml:"deadline"`
}

type Tasks map[string]Task

func center_pad(str string, padding int) (string, error) {
	l := len(str)
	if l > padding {
		return "", errors.New("string is bigger than padding")
	}
	diff := padding - l
	start := int(diff / 2)
	end := start
	if diff%2 != 0 {
		end++
	}
	result := strings.Repeat(" ", start) + str + strings.Repeat(" ", end)
	return result, nil
}

func window(str string, window_size int) string {
	l := len(str)
	if l < window_size {
		res, err := center_pad(str, window_size)
		if err != nil {
			log.Fatal(err)
		}
		str = res
		l = window_size
	}
	start := int(time.Now().Unix()) % l
	end := window_size + start
	if end > l {
		end = end - l
		return str[start:] + str[:end]
	}
	return str[start:end]
}

func main() {
	home := os.Getenv("HOME")
	path := fmt.Sprintf("%s/todo.toml", home)
	source_code_path := fmt.Sprintf("%s/Code/tomldo/main.go", home)

	block_button := os.Getenv("BLOCK_BUTTON")
	var cmd *exec.Cmd
	switch block_button {
	case "1", "3":
		cmd = exec.Command("alacritty", "-e", "nvim", path)
	case "6":
		cmd = exec.Command("alacritty", "-e", "nvim", source_code_path)
	}

	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		err = cmd.Process.Release()
		if err != nil {
			log.Fatal(err)
		}
	}

	var taskmap = Tasks{}
	if _, err := toml.DecodeFile(path, &taskmap); err != nil {
		log.Fatal(err)
	}
	var tasks []Task
	for name, task := range taskmap {
		task.Name = name
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Name < tasks[j].Name
	})

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Deadline.Before(tasks[j].Deadline)
	})

	index := int(time.Now().Unix()/10) % len(tasks)
	task := tasks[index]

	seconds := int(time.Until(task.Deadline).Seconds())
	days := int(seconds / (60 * 60 * 24))
	hours := int(seconds/(60*60)) % 24
	minutes := int(seconds/60) % 60
	seconds = seconds % 60
	countdown := fmt.Sprintf("%02v:%02v:%02v:%02v", days, hours, minutes, seconds)
	fmt.Printf("😬%s💣%s😬\n", window(task.Name, 20), countdown)
}
