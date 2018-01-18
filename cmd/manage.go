// Copyright © 2018 Roland Varga <roland.varga.can@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
    // "fmt"
    "os"
    "strconv"

    "github.com/spf13/cobra"
    "github.com/go-redis/redis"
    log "github.com/sirupsen/logrus"
)

var manageCmd = &cobra.Command{
    Use:   "manage",
    Short: "Manages number of pod replicas based on queue size",
    Long: `Manages the number of pod replicas based on Redis queue utilisation.
Thresholds are taken from environment variables: 

MIN_POD:		                              num
MAX_POD:		                              num
POD_INCREMENT:		                        num
POD_DECREMENT:		                        num
SCALE_UP_POLICY:		                      linear
SCALE_DOWN_POLICY:		                    linear|onlyWhenNoneNeeded
MONITOR_TYPE:		                          redis|haproxy
MONITOR_USERNAME:		                      for haproxy
MONITOR_PASSWORD:		                      for haproxy
MONITOR_HOST:		                          some.address.com
MONITOR_PORT:		                          6379
MONITOR_PATH:		                          haproxy
MONITOR_DB:		                            num (redis)
MONITOR_QUEUE_NAME:		                    resque:app:name
MONITOR_PODS:		
MANAGE_SPEED:		                          rest time in seconds between each control loop


Example:
./pod-autoscaler manage 
`,
    Run: func(cmd *cobra.Command, args []string) {
        log.Info("Executing manage command...")
        Run()
    },
}

type Resource interface {
    isOverwhelmed() bool
}

type RedisQueue struct {
    Addr, Password string
    DB int
}

type EnvVars struct {
    MIN_POD, MAX_POD, POD_INCREMENT  int
    POD_DECREMENT, MONITOR_PORT, MANAGE_SPEED, MONITOR_DB int
    SCALE_UP_POLICY, SCALE_DOWN_POLICY, MONITOR_TYPE, MONITOR_USERNAME string
    MONITOR_PASSWORD, MONITOR_HOST, MONITOR_PATH, MONITOR_QUEUE_NAME, MONITOR_PODS string
}

func (rq RedisQueue) isOverwhelmed(queue string) bool {
    // move connect into other method later
    client := redis.NewClient(&redis.Options{
        Addr:     rq.Addr,
        Password: rq.Password, // no password set
        DB:       rq.DB,
    })


    // gets queue names
    // queues, err := client.Keys("resque:app:queue:*").Result()
    // if err != nil {
    //     panic(err)
    // }

    // iterates those queues
    // for _, q := range queues {
    //     fmt.Println(q)
    //     qs := client.LLen(q)
    //     fmt.Printf("%v\n", qs)
    // }

    q := client.LLen(queue)
    log.Info(q)

    return true // for now FIXME
}


func init() {
    rootCmd.AddCommand(manageCmd)
}

func confLogger() {
    logFmt := new(log.TextFormatter)
    // logFmt.TimestampFormat = "2018-01-15 00:00:00"   // FIXME
    log.SetFormatter(logFmt)
    logFmt.FullTimestamp = true
}

func getVars() EnvVars {
    return EnvVars{
        MIN_POD                 toInt(os.Getenv("MIN_POD")),
        MAX_POD                 toInt(os.Getenv("MAX_POD")),
        POD_INCREMENT           toInt(os.Getenv("POD_INCREMENT")),
        POD_DECREMENT           toInt(os.Getenv("POD_DECREMENT")),
        SCALE_UP_POLICY         os.Getenv("SCALE_UP_POLICY"),
        SCALE_DOWN_POLICY       os.Getenv("SCALE_DOWN_POLICY"),
        MONITOR_TYPE            os.Getenv("MONITOR_TYPE"),
        MONITOR_USERNAME        os.Getenv("MONITOR_USERNAME"),
        MONITOR_PASSWORD        os.Getenv("MONITOR_PASSWORD"),
        MONITOR_HOST            os.Getenv("MONITOR_HOST"),
        MONITOR_PORT            toInt(os.Getenv("MONITOR_PORT")),
        MONITOR_PATH            os.Getenv("MONITOR_PATH"),
        MONITOR_DB              toInt(os.Getenv("MONITOR_PORT")),
        MONITOR_QUEUE_NAME      os.Getenv("MONITOR_QUEUE_NAME"),
        MONITOR_PODS            os.Getenv("MONITOR_PODS"),
        MANAGE_SPEED            toInt(os.Getenv("MANAGE_SPEED")),
    }
}

// helper func for easier parsing
func toInt(s string) (int) {
    i, _ := strconv.Atoi(s)
    return i
}

func Run() {
    confLogger()  // logrus config
    envVars := getVars() // read K8S friendly environment variables

    redis := RedisQueue {
        Addr:     envVars.MONITOR_HOST,
        Password: "",
        DB:       envVars.MONITOR_DB,
    }

    redis.isOverwhelmed()
}
