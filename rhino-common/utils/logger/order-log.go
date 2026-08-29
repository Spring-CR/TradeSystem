package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"rhino-common/utils/byteutils"
	"rhino-core/schema"
	"time"
)

const (
	logTimeFormat = "20060102-15:04:05.000000|"
)

type logData struct {
	order   *schema.TradeOrder
	resp    *schema.TradeActionResp
	logTime time.Time
	data    string
}

type OrderLog struct {
	logDir        string
	logChan       chan *logData
	getKeyFunc    func(tradeOrder *schema.TradeOrder) (key string)
	writerMap     map[string]io.Writer
	flushInterval time.Duration
	running       bool
}

func NewOrderLog(logDir string, logChanLen int, flushInterval time.Duration, getKeyFunc func(tradeOrder *schema.TradeOrder) (key string), running bool) *OrderLog {
	orderLog := &OrderLog{
		logDir:        logDir,
		logChan:       make(chan *logData, logChanLen),
		flushInterval: flushInterval,
		getKeyFunc:    getKeyFunc,
		writerMap:     make(map[string]io.Writer),
		running:       running,
	}
	if running {
		go func() {
			orderLog.keepLogging()
		}()
		time.Sleep(2*time.Second)
	}
	return orderLog
}


func (l *OrderLog) createWriterByKey(key string) io.Writer {
	// 确保日志目录存在
	if err := os.MkdirAll(l.logDir, 0777); err != nil {
		panic(err)
	}
	filePath := filepath.Join(l.logDir, key+".log")
	log.Printf("===>createWriterByKey, key=%s, filePath=%s\n", key, filePath)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	// 64KB 缓冲，兼顾效率与延迟
	return bufio.NewWriterSize(file, 64*1024)
}

// 刷新所有 writer 的缓冲区
func (l *OrderLog) flushAll() {
	for key, w := range l.writerMap {
		if bw, ok := w.(*bufio.Writer); ok {
			if err := bw.Flush(); err != nil {
				// 实际项目中可记录错误日志，这里简单处理
				_ = err
			}
		} else {
			// 如果不是 bufio.Writer（理论上都是），也可尝试断言 syncFlusher
			if flusher, ok := w.(interface{ Flush() error }); ok {
				_ = flusher.Flush()
			}
		}
		_ = key // avoid unused warning
	}
}

func (l *OrderLog) keepLogging() {
	log.Printf("Start keepLogging!")
	ticker := time.NewTicker(l.flushInterval)
	for {
		select {
		case logData := <-l.logChan:
			t := logData.logTime.Format(logTimeFormat)
			key := t[:8] + "_" + l.getKeyFunc(logData.order)
			var logContent string
			if logData.resp == nil {
				logContent = fmt.Sprintf("%s[order=%v|symbol=%v|side=%v]%s\n", t, logData.order.AppOrdID, logData.order.Symbol, logData.order.Side, logData.data)
			} else {
				logContent = fmt.Sprintf("%s[order=%v|symbol=%v|side=%v][resp.ordStatus=%v｜lastShares=%v｜lastPx=%v]%s\n", t, logData.order.AppOrdID, logData.order.Symbol, logData.order.Side, logData.resp.OrdStatus, logData.resp.LastShares, logData.resp.LastPx, logData.data)
			}

			writer, ok := l.writerMap[key]
			if !ok {
				writer = l.createWriterByKey(key)
				l.writerMap[key] = writer
			}
			writer.Write(byteutils.GetZeroCopyBytes(logContent))

		case <-ticker.C:
			// 定期刷新所有缓冲区，确保数据不会长期滞留在内存
			l.flushAll()
		}
	}
}

func (l *OrderLog) Printf(order *schema.TradeOrder, resp *schema.TradeActionResp, content string, args ...interface{}) {
	if !l.running {
		return
	}
	logData := &logData{order: order, resp: resp, logTime: time.Now(), data: fmt.Sprintf(content, args...)}
	l.logChan <- logData
}

// Close 优雅关闭：刷新所有缓冲区，可选的资源释放
func (l *OrderLog) Close() {
	l.flushAll()
	// 如果需要关闭底层文件，可以遍历 writerMap 获取 *os.File，这里省略
}
