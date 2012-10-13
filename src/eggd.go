package main

import (
  "bytes"
  "inotify"
  "log"
  "os"
  "os/exec"
  "sync/atomic"
  "time"
)

var Lock int64

func handle() {
  var cmd *exec.Cmd
  var pid string
  var pid_out bytes.Buffer
  var make_out bytes.Buffer
  var e error

  // Atomic Lock.
  if !atomic.CompareAndSwapInt64(&Lock, 0, 1) {
    return
  }

  // Get the pid and stop the server.
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
  }

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
  cmd.Stdout = &make_out
  e = cmd.Run()
  if e != nil {
    log.Println(make_out.String())
    log.Println("make:", e)
  }

  // Start the server.
  log.Println("start: preparing...")
  e = exec.Command("build/curve_https").Start()
  if e != nil {
    log.Println("start:", e)
  }

next:
  // Release the event Lock after all events have been parsed.
  atomic.CompareAndSwapInt64(&Lock, 1, 0)
  return
}

func main() {
  watcher, e := inotify.NewWatcher()
  if e != nil { log.Fatal(e) }
  e = watcher.Watch("/home/ubuntu/deploy/server.git")
  if e != nil { log.Fatal(e) }
  e = os.Chdir("/home/ubuntu/deploy_local/server")
  if e != nil { log.Fatal(e) }

  for {
    // Event loop. The first event Locks it.
    select {
    case ev := <-watcher.Event:
      log.Println("event:", ev)
      go handle()
    case err := <-watcher.Error:
      log.Println("error:", err)
    }
  }
}
