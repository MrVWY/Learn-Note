package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type async struct{}

func (a *async) Funcs(funcs ...func() error) func() error {
	return func() error {
		for _, fn := range funcs {
			if err := fn(); err != nil {
				return err
			}
		}
		return nil
	}
}

func (a *async) Run(workName string, initFunc, jobFunc, endFunc func() error) {
	fmt.Println(fmt.Sprintf("%s is Running!", workName))
	if err := initFunc(); err != nil {
		fmt.Println(fmt.Sprintf("%s is init fail", workName))
	}
	fmt.Println(fmt.Sprintf("%s is init successful", workName))

	go func() {
		if err := jobFunc(); err != nil {
			fmt.Println(fmt.Sprintf("%s is job fail", workName))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)
	<-c

	if err := endFunc(); err != nil {
		fmt.Println(fmt.Sprintf("%s is endFunc fail", workName))
	}
	fmt.Println(fmt.Sprintf("%s is endFunc successful", workName))
}

func main() {
	ay := &async{}
	ay.Run("workTest",
		ay.Funcs(
			func() error {
				fmt.Println("init A!")
				return nil
			},
			func() error {
				fmt.Println("init B!")
				return nil
			}),
		ay.Funcs(
			func() error {
				fmt.Println("working C!")
				return nil
			},
			func() error {
				fmt.Println("working D!")
				return nil
			}),
		func() error {
			fmt.Println("End this work!")
			return nil
		})

	fmt.Println("--------------------------------------------")

	ay.Run("workTest",
		ay.Funcs(
			func() error {
				fmt.Println("init A!")
				return nil
			},
			func() error {
				fmt.Println("init B!")
				return nil
			}),
		ay.Funcs(
			func() error {
				fmt.Println("working C!")
				return nil
			},
			func() error {
				return errors.New("Panic!")
			}),
		func() error {
			fmt.Println("End this work!")
			return nil
		})
}
