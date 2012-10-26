package main

import (
  "bytes"
  "conf"
  "inotify"
  "log"
  "os"
  "os/exec"
  "strconv"
  "strings"
  "sync/atomic"
  "time"
)

var Config *conf.ConfigFile
var ConfigName string
var Home string
var Lock int64
var Locks []int64
var Paths []string
var RemoteCount int

func handleEventOld() {
  var cmd *exec.Cmd
  //var pid string
  //var pid_out bytes.Buffer
  var cmdStdout bytes.Buffer
  var e error

  if !atomic.CompareAndSwapInt64(&Lock, 0, 1) {
    return
  }

  /*// Get the pid and stop the server.
  log.Println("cat: preparing...")
  cmd = exec.Command("cat", "build/curve_https.pid")
  cmd.Stdout = &pid_out
  e = cmd.Run()
  if e != nil {
    log.Println("cat:", e)
    goto next
  }
  //pid, e := strconv.ParseInt(out.String(), 10, 32)
  //if e != nil { continue }
  pid = pid_out.String()
  log.Println("pid:", pid)
  log.Println("kill: preparing...")
  e = exec.Command("kill", pid).Run()
  if e != nil {
    log.Println("kill:", e)
    goto next
  }*/

  // Wait at least a short while before pulling the repo.
  log.Println("sleeping...")
  time.Sleep(5 * time.Second)

  // Pull from local repository.
  log.Println("git-pull: preparing...")
  e = exec.Command("git", "pull", "origin", "master").Run()
  if e != nil {
    log.Println("git-pull:", e)
    goto next
  }

  // Wait at least a short while before installing.
  log.Println("sleeping...")
  time.Sleep(5 * time.Second)

  // Build and install the code.
  log.Println("make: preparing...")
  cmd = exec.Command("make")
  cmd.Stdout = &cmdStdout
  e = cmd.Run()
  if e != nil {
    log.Println(cmdStdout.String())
    log.Println("make:", e)
  }

  // Start the server.
  log.Println("start: preparing...")
  e = exec.Command("foreman", "start").Start()
  if e != nil {
    log.Println("start:", e)
  }

next:
  atomic.CompareAndSwapInt64(&Lock, 1, 0)
  return
}

func handleEvent(ev *inotify.Event) {
  var cmd *exec.Cmd
  var cmdStdout bytes.Buffer
  var e error

  log.Println(ev)
  name := ev.Name

  log.Println("getting parent path...")
  idx := -1
  path := ""
  for i := 0; i < RemoteCount; i++ {
    if strings.HasPrefix(name, Paths[i]) {
      idx = i
      path = Paths[i]
      break
    }
  }
  if idx == -1 {
    log.Println("couldnt find parent path")
    return
  }

  log.Println("acquiring lock", idx)
  if !atomic.CompareAndSwapInt64(&Locks[idx], 0, 1) {
    log.Println("failed to acquire lock", idx)
    return
  }

  log.Println("sleeping...")
  time.Sleep(5 * time.Second)

  checkoutPath := strings.Join([]string{"/var/eggd", path}, "/")
  e = os.MkdirAll(checkoutPath, 0755)
  if e != nil {
    log.Println("os.MkdirAll:", e)
    goto next
  }
  e = os.Chdir(checkoutPath)
  if e != nil {
    log.Println("os.Chdir:", e)
    goto next
  }
  if os.IsNotExist(os.Chdir(".git")) {
    log.Println("running git-clone...")
    cmd = exec.Command("git", "clone", path, ".")
    cmd.Stdout = &cmdStdout
    e = cmd.Run()
    if e != nil {
      log.Println("git-clone:", e)
      log.Println(cmdStdout.String())
      goto next
    }
  } else {
    os.Chdir("..")
  }

  log.Println("running git-pull...")
  cmd = exec.Command("git", "pull", "origin", "master")
  cmd.Stdout = &cmdStdout
  cmd.Run()
  if e != nil {
    log.Println("git-pull:", e)
    log.Println(cmdStdout.String())
    goto next
  }

  //log.Println("sleeping...")
  //time.Sleep(5 * time.Second)

  log.Println("running make...")
  cmd = exec.Command("make")
  cmd.Stdout = &cmdStdout
  e = cmd.Run()
  if e != nil {
    log.Println("make:", e)
    log.Println(cmdStdout.String())
    goto next
  }

  log.Println("running foreman...")
  cmd = exec.Command("foreman", "start")
  cmd.Stdout = &cmdStdout
  e = cmd.Run()
  if e != nil {
    log.Println("foreman:", e)
    log.Println(cmdStdout.String())
    goto next
  }

next:
  log.Println("releasing lock", idx)
  atomic.CompareAndSwapInt64(&Locks[idx], 1, 0)
  return
}

func initConfig() {
  // Check that the configfile exists.
  Home = os.Getenv("HOME")
  //ConfigName = strings.Join([]string{Home, ".eggconfig"}, "/")
  ConfigName = "/etc/eggconfig"
  fd, e := os.OpenFile(ConfigName, os.O_RDWR | os.O_CREATE, 0644)
  if e != nil { log.Fatal(e) }
  e = fd.Close()
  if e != nil { log.Fatal(e) }

  // Make sure the configfile is set up correctly.
  Config, e := conf.ReadConfigFile(ConfigName)
  if e != nil { log.Fatal(e) }
  Config.AddSection("global")
  if !Config.HasOption("global", "count") {
    Config.AddOption("global", "count", "0")
  }

  // Modify the configfile through cmdline args.
  if len(os.Args) > 2 && os.Args[1] == "add" {
    log.Println("adding", os.Args[2], "to eggd path")
    count, _ := Config.GetInt("global", "count")
    count = count + 1
    Config.AddOption("global", "count", strconv.Itoa(count))
    section := strings.Join([]string{"remote-", strconv.Itoa(count)}, "")
    Config.AddOption(section, "path", os.Args[2])
    Config.WriteConfigFile(ConfigName, 0644, "")
    os.Exit(0)
  }

  RemoteCount, e = Config.GetInt("global", "count")
  if e != nil { log.Fatal(e) }
  Config.WriteConfigFile(ConfigName, 0644, "")
}

func startWatcher() {
  // Start the event loops. Use inotify to watch the configfile and each
  // repository path.
  Config, e := conf.ReadConfigFile(ConfigName)
  if e != nil { log.Fatal(e) }
  watcher, e := inotify.NewWatcher()
  if e != nil { log.Fatal(e) }
  log.Println("watching configfile")
  e = watcher.Watch(ConfigName)
  if e != nil { log.Fatal(e) }
  Locks = make([]int64, RemoteCount)
  Paths = make([]string, RemoteCount)
  for i := 1; i <= RemoteCount; i++ {
    section := strings.Join([]string{"remote-", strconv.Itoa(i)}, "")
    if !Config.HasOption(section, "path") {
      log.Println("doesnt have path")
      continue
    }
    path, e := Config.GetRawString(section, "path")
    if e != nil { log.Fatal(e) }
    log.Println("watching path", path)
    e = watcher.Watch(path)
    if e != nil { log.Fatal(e) }
    Locks[i-1] = 0
    Paths[i-1] = path
  }
  for {
    select {
      case ev := <-watcher.Event:
        go handleEvent(ev)
      case err := <-watcher.Error:
        log.Println("error:", err)
    }
  }
}

func main() {
  initConfig()
  startWatcher()
}
