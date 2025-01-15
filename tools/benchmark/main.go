package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/GrapefruitCat030/gfc_dcache/tools/benchmark/client"
	"github.com/GrapefruitCat030/gfc_dcache/tools/benchmark/result"
	"github.com/spf13/cobra"
)

type cmdConfig struct {
	serverType     string
	serverAddr     string
	cacheOperation string
	connThreads    int
	totalReq       int
	keySize        int
	valSize        int
	pipelineLen    int
}

var cfg cmdConfig

var rootCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "benchmark is a tool to benchmark cache server",
	Run: func(cmd *cobra.Command, args []string) {
		runBenchmark()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&cfg.serverType, "server-type", "t", "tcp", "server type")
	rootCmd.Flags().StringVarP(&cfg.serverAddr, "server-addr", "a", "localhost:8080", "server address, e.g., localhost:8080")
	rootCmd.Flags().StringVarP(&cfg.cacheOperation, "cache-operation", "o", "get", "cache operation")
	rootCmd.Flags().IntVarP(&cfg.connThreads, "conn-threads", "c", 1, "number of connection threads")
	rootCmd.Flags().IntVarP(&cfg.totalReq, "total-req", "r", 1000, "total number of requests")
	rootCmd.Flags().IntVarP(&cfg.keySize, "key-size", "k", 10, "key size")
	rootCmd.Flags().IntVarP(&cfg.valSize, "val-size", "v", 10, "value size")
	rootCmd.Flags().IntVarP(&cfg.pipelineLen, "pipeline-len", "p", 1, "pipeline length")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func runBenchmark() {
	// 1. Initialize the benchmark result
	var totalRes result.BenchmarkRes
	var resChan = make(chan *result.BenchmarkRes, cfg.connThreads)
	// 2. Start the benchmark workers
	start := time.Now()
	for i := 0; i < cfg.connThreads; i++ {
		go worker(i, cfg.totalReq/cfg.connThreads, resChan)
	}
	// 3. Collect the benchmark result
	for i := 0; i < cfg.connThreads; i++ {
		totalRes.AddResult(<-resChan)
	}
	dur := time.Since(start)
	// 3. Print the benchmark result
	outputResult(&totalRes, dur)
}

func worker(workID, reqNum int, ch chan *result.BenchmarkRes) {
	// 1. Initialize
	cli := client.NewClient(cfg.serverType, cfg.serverAddr)
	pipelineOps := make([]*client.Operation, 0, cfg.pipelineLen)
	res := &result.BenchmarkRes{StatisticBuckets: make([]result.TimeStat, 0)}
	// 2. Start the benchmark
	for i := 0; i < reqNum; i++ {
		// generate the random key and value
		var key, value string
		if cfg.keySize > 0 {
			key = fmt.Sprintf("%d", rand.Intn(cfg.keySize))
		} else {
			key = fmt.Sprintf("%d", workID*reqNum+i)
		}
		value = fmt.Sprintf("%s%s", strings.Repeat("a", cfg.valSize), key)
		// set the request operation
		opType := cfg.cacheOperation
		if opType == client.OperationTypeMix {
			if i%2 == 0 {
				opType = client.OperationTypeSet
			} else {
				opType = client.OperationTypeGet
			}
		}
		// do or pipelineDo
		op := &client.Operation{Name: opType, Key: key, Value: value}
		if cfg.pipelineLen <= 1 {
			justDo(cli, op, res)
		} else {
			pipelineOps = append(pipelineOps, op)
			if len(pipelineOps) == cfg.pipelineLen { // if the pipeline is full, do the pipelineDo
				pipelineDo(cli, pipelineOps, res)
				pipelineOps = pipelineOps[:0] // clear the pipeline
			}
		}
	}
	// 3. Send the result to the channel
	ch <- res
}

func justDo(cli client.Client, op *client.Operation, res *result.BenchmarkRes) {
	originVal := op.Value // if op is GET, the value will be updated
	start := time.Now()
	if err := cli.Do(op); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	dur := time.Since(start)
	if op.Name == client.OperationTypeGet && (op.Value == "" || op.Value != originVal) {
		res.CollectTimeCost(dur, "miss")
		return
	}
	res.CollectTimeCost(dur, op.Name)
}

func pipelineDo(cli client.Client, ops []*client.Operation, res *result.BenchmarkRes) {
	originVals := make([]string, len(ops)) // if op is GET, the value will be updated
	for i, op := range ops {
		originVals[i] = op.Value
	}
	start := time.Now()
	if err := cli.PipelinedDo(ops); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	dur := time.Since(start) / time.Duration(len(ops))
	for i, op := range ops {
		if op.Name == client.OperationTypeGet && (op.Value == "" || op.Value != originVals[i]) {
			res.CollectTimeCost(dur, "miss")
			continue
		}
		res.CollectTimeCost(dur, op.Name)
	}
}

func outputResult(totalRes *result.BenchmarkRes, dur time.Duration) {
	totalCount := totalRes.GetCount + totalRes.MissCount + totalRes.SetCount
	fmt.Printf("%d records get\n", totalRes.GetCount)
	fmt.Printf("%d records miss\n", totalRes.MissCount)
	fmt.Printf("%d records set\n", totalRes.SetCount)
	fmt.Printf("%f seconds total\n", dur.Seconds())
	statCountSum := 0
	statTimeSum := time.Duration(0)
	for b, s := range totalRes.StatisticBuckets {
		if s.Count == 0 {
			continue
		}
		statCountSum += s.Count
		statTimeSum += s.Time
		fmt.Printf("%d%% requests < %d ms\n", statCountSum*100/totalCount, b+1)
	}
	fmt.Printf("%d usec average for each request\n", int64(statTimeSum/time.Microsecond)/int64(statCountSum))
	fmt.Printf("throughput is %f MB/s\n", float64((totalRes.GetCount+totalRes.SetCount)*cfg.valSize)/1e6/dur.Seconds())
	fmt.Printf("rps is %f\n", float64(totalCount)/float64(dur.Seconds()))
}
