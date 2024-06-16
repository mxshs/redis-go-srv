package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"mxshs/redis-go/app/utils"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)


func Runner() error {
    port := flag.String("port", "6381", "port of your redis instance")
    timeout := flag.Int("timeout", 0, "command timeout")

    flag.Parse()


    if port == nil || *port == "" {
        return fmt.Errorf("invalid or missing port")
    }


    rdb := redis.NewClient(&redis.Options{
        Addr:             "localhost:" + *port,
        DisableIndentity: true,
        Protocol:         3,
    })

    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {

        cmd := scanner.Text()
        s := strings.Split(cmd, " ")

        ctx, _ := context.WithTimeout(context.Background(), time.Second * time.Duration(*timeout))

        res, err := rdb.Do(ctx, utils.Erase[string, any](s).([]any)...).Text()
        if err != nil {
            io.WriteString(os.Stdout, err.Error() + "\n")
            continue
        }

        io.WriteString(os.Stdout, res + "\n")
    }

    return nil
}

func StressTest() {
    var sm atomic.Int64
    var wg sync.WaitGroup

    for i := 0; i < 10000; i++ {
        wg.Add(1)

        go func(i int) {
            rdb := redis.NewClient(&redis.Options{
                Addr:             "localhost:6379",
                DisableIndentity: true,
                Protocol:         3,
            })

            for j := 0; j < 1000; j++ {
                key := fmt.Sprintf("%d", j * i)
                value := fmt.Sprintf("%d", j - i)

                start := time.Now()

                res := rdb.Set(context.TODO(), key, value, 0)

                msg, err := res.Result()
                if err != nil {
                    fmt.Println(err.Error())
                    return
                }

                if msg != "OK" {
                    fmt.Println(msg)
                    return
                }

                get := rdb.Get(context.TODO(), key)

                _, err = get.Result()
                if err != nil {
                    fmt.Println(err.Error())
                    return
                }

                sm.Add(int64(time.Since(start)))
            }
            wg.Done()
        }(i)
    }

    wg.Wait()
    fmt.Println(float64(sm.Load()) / 1000000 / 1000000)
}

func main() {
    if err := Runner(); err != nil {
        panic(err)
    }
}

