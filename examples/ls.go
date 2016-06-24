/*
Purpose: simple ls cli command written in golang
Author : NhamLH <lehoainham@gmail.com>
Date   : Jun 12 2016
*/

package main

import (
  "fmt"
  "flag"
  "os"
  "io/ioutil"
  "strconv"
  "strings"
  "time"
)

func main() {
  l := flag.Bool("l", false, "Use a long listing format")
  flag.Parse()

  path_list := flag.Args()

  // list current directory if user don't provide one
  if len(path_list) == 0 {
    path_list = append(path_list, ".")
  }

  for _, path := range path_list {
    f, error := os.Stat(path)

    if os.IsNotExist(error) {
      fmt.Println(path, "not found.")
      continue
    }

    if f.IsDir() {
      files, error := ioutil.ReadDir(path)

      fmt.Println(error.Op)

      if os.IsPermission(error) {
        fmt.Println("ls cannot open", path, ": permission denied.")
        continue
      }

      for _,file := range files {
        fmt.Println(compose_string(file, *l))
      }
    } else {
      fmt.Println(compose_string(f, *l))
    }
  }
}

// return a string to print out to console
func compose_string(f os.FileInfo, verbose bool) string {
  if verbose == false {
    return f.Name()
  } else {
    tmp_mtime := f.ModTime()

    mmonth := tmp_mtime.Month().String()[:3]
    mday := strconv.Itoa(tmp_mtime.Day())
    var mtime string

    // only print modify year if it's different from system year
    if cur_year := time.Now().Year(); cur_year != tmp_mtime.Year() {
      myear := strconv.Itoa(tmp_mtime.Year())
      mtime = strings.Join([]string{mmonth, mday, myear}, " ")
    } else { // instead print modify hour and minute
      mhour := strconv.Itoa(tmp_mtime.Hour())
      mminute := strconv.Itoa(tmp_mtime.Minute())
      tmp := strings.Join([]string{mhour, mminute}, ":")
      mtime = strings.Join([]string{mmonth, mday, tmp}, " ")
    }

    result := f.Mode().String() + "\t" + strconv.FormatInt(f.Size(), 10) + "\t" + mtime + "\t" + f.Name()

    return result
  }
}
