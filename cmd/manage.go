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
    "fmt"
    // "os"

    "github.com/spf13/cobra"
    "github.com/go-redis/redis"
)

// manageCmd represents the manage command
var manageCmd = &cobra.Command{
    Use:   "manage",
    Short: "Manages number of pod replicas based on queue size",
    Long: `Manages the number of pod replicas based on Redis queue utilisation.
Thresholds are taken from environment variables: 

POD_REPLICA_MIN:         3
POD_REPLICA_MAX:        10
REDIS_TASK_MULTIPLIER:   5    (start a pod after every 5 tasks)
AUTOSCALER_FREQUENCY:   30    (check queue size every 30 seconds) 

Example:
./pod-autoscaler manage 
`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("manage called")
        Example()
    },
}

func init() {
    rootCmd.AddCommand(manageCmd)

    // manageCmd.PersistentFlags().String("foo", "", "A help for foo")
    // manageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Example() {
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       3,
    })

    // pong, err := client.Ping().Result()
    // fmt.Println(pong, err)
    // Output: PONG <nil>

    queues, err := client.Keys("resque:webadmit:queue:*").Result()
    if err != nil {
        panic(err)
    }

    fmt.Println(queues)

    for _, q := range queues {
        fmt.Println(q)
        qs := client.LLen(q)
        fmt.Printf("%v\n", qs)
    }
}
