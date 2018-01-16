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

POD_REPLICA_MIN:                             3
POD_REPLICA_MAX:                            10
REDIS_ADDRESS:                  localhost:6379
REDIS_PASS:                             secret
REDIS_DB:                                    3
REDIS_QUEUE:            my-awesome-redis-queue
REDIS_TASK_MULTIPLIER:                       5 (start a pod after every 5 tasks)
AUTOSCALER_FREQUENCY:                       30 (check queue size every 30 seconds) 

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
    POD_REPLICA_MIN, POD_REPLICA_MAX int
    REDIS_ADDRESS, REDIS_PASS, REDIS_QUEUE string
    REDIS_TASK_MULTIPLIER, AUTOSCALER_FREQUENCY, REDIS_DB int
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
        POD_REPLICA_MIN:        toInt(os.Getenv("POD_REPLICA_MIN")),
        POD_REPLICA_MAX:        toInt(os.Getenv("POD_REPLICA_MAX")),
        REDIS_ADDRESS:          os.Getenv("REDIS_ADDRESS"),
        REDIS_PASS:             os.Getenv("REDIS_PASS"),
        REDIS_DB:               toInt(os.Getenv("REDIS_DB")),
        REDIS_QUEUE:            os.Getenv("REDIS_QUEUE"),
        REDIS_TASK_MULTIPLIER:  toInt(os.Getenv("REDIS_TASK_MULTIPLIER")),
        AUTOSCALER_FREQUENCY:   toInt(os.Getenv("AUTOSCALER_FREQUENCY")),
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
        Addr:     envVars.REDIS_ADDRESS,
        Password: envVars.REDIS_PASS,
        DB:       envVars.REDIS_DB,
    }

    redis.isOverwhelmed()
}
