/*
Purpose: Simple cat cli command written in golang
Author : NhamLH <lehoainham@gmail.com>
Date   : Jun 12 2016
*/

package main

import (
  "flag"
  "fmt"
  "io/ioutil"
  "strings"
)

func main() {
  nf := flag.Bool("n", false, "Print line number")
  flag.Parse()

  path_list := flag.Args()

  if len(path_list) == 0 {
    fmt.Println("Usage: cat [-n] <file>")
    return
  }

  for _, path := range path_list {
    file, error := ioutil.ReadFile(path)

    if error != nil {
      fmt.Println(error)
      return
    }

    t := string(file)
    // get rid of the last blank line cause by last \n
    t = strings.TrimRight(t, "\n")
    str := strings.Split(t, "\n")

    for i, line := range str {
      if *nf {
        fmt.Println(i+1, line)
      } else {
        fmt.Println(line)
      }
    }
  }
}
