package notify_provider

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"rhino-common/domain_error"
	"rhino-common/utils/http/post_client"
	"rhino-common/utils/tail"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type WechatErrorNotifyProvider struct {
	webhookAddr     string
	logFileWriter   *bufio.Writer
	logFile         *os.File
	logLineChan     chan []byte
	fileTail        *tail.FileTail
	truncateWeekday time.Weekday
	genWebchatData  func(logLine string) (interface{}, error)
	httpClient      *post_client.HTTPClient
}

func NewWechatErrorNotifyProvider(webhookAddr string, logFilePath string, truncateTicker *time.Ticker, flushTicker *time.Ticker, truncateWeekday time.Weekday, genWebchatData func(logLine string) (interface{}, error)) (*WechatErrorNotifyProvider, error) {

	if len(logFilePath) > 0 {
		os.MkdirAll(filepath.Dir(logFilePath), 0755)
	}

	p := &WechatErrorNotifyProvider{webhookAddr: webhookAddr, logLineChan: make(chan []byte, 10000), truncateWeekday: truncateWeekday, genWebchatData: genWebchatData, httpClient:post_client.NewHTTPClient(nil)}

	var err error
	p.logFile, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// 错误处理
		return nil, err
	}
	p.logFileWriter = bufio.NewWriterSize(p.logFile, 4096)
	p.fileTail, err = tail.NewFileTail(logFilePath, false)
	if err != nil {
		// 错误处理
		return nil, err
	}
	p.checkBeijingTimeAndTruncate()

	go func() {
		for {
			select {
			case <-truncateTicker.C:
				p.checkBeijingTimeAndTruncate()
			case <-flushTicker.C:
				p.flush()
			case logLine := <-p.logLineChan:
				logLine = append(logLine, '\n')
				_, err := p.logFileWriter.Write(logLine)
				if err != nil {
					log.Printf("logBuffer.Write error = %v\n", err)
				}
			}
		}
	}()

	go func() {
		for {
			line := <-p.fileTail.Lines
			if p.genWebchatData == nil {
				log.Printf("detect line:%s\n", line)
			} else {
				payload, err := p.genWebchatData(line)
				if err != nil {
					log.Printf("fail to genWebchatData from %s, error=%v\n", line, err)
					continue
				}
				if payload == nil {
					continue
				}
				err = p.httpClient.PostAsync(p.webhookAddr, payload, nil, func (result*post_client.Result)  {
					log.Printf("success send for %s\n", line)
				})
				if err != nil {
					log.Printf("fail to PostAsync to %s, , error=%v\n", p.webhookAddr, err)
					continue
				}
			}
		}
	}()

	return p, nil
}

func (p *WechatErrorNotifyProvider) NotifyError(de *domain_error.Error) {
	if de == nil {
		return
	}
	data, _ := json.Marshal(de)
	select {
	case p.logLineChan <- data:
		//log.Printf("add log: %s\n", data)
	default:
		// 通道满了，可以选择丢弃或阻塞
		log.Printf("写入队列已满! 丢弃, data=%s\n", data)
	}
}

func (p *WechatErrorNotifyProvider) flush() {
	err := p.logFileWriter.Flush()
	if err != nil {
		log.Printf("fail to flush, error = %v\n", err)
	}
}

// truncateLogFile 截断指定的日志文件
func (p *WechatErrorNotifyProvider) TruncateLogFile() error {
	p.flush()
	err := p.logFile.Truncate(0)
	if err != nil {
		log.Printf(">>> fail to Truncate file:%s\n", p.logFile.Name())
	}
	_, err = p.logFile.Seek(0, 0) // 重置文件指针
	if err != nil {
		log.Printf(">>> fail to Seek file:%s\n", p.logFile.Name())
	}
	// err = p.fileTail.Reset()
	// if err != nil {
	// 	log.Printf(">>> fail to Reset fileTail:%s\n", p.logFile.Name())
	// }

	log.Printf("success truncate at %s\n", time.Now())

	return nil
}

// checkBeijingTimeAndTruncate 检查北京时间是否为周日，如果是则截断日志
func (p *WechatErrorNotifyProvider) checkBeijingTimeAndTruncate() {
	// 加载北京时间时区（Asia/Shanghai）
	beijingLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("加载时区失败: %v", err)
		return
	}

	// 获取当前北京时间
	beijingTime := time.Now().In(beijingLocation)

	// 判断是否为周日 (time.Sunday = 0)
	if beijingTime.Weekday() == p.truncateWeekday {
		log.Printf("当前北京时间: %s, 是%s，开始截断日志", beijingTime.Format("2006-01-02 15:04:05 Monday"), p.truncateWeekday)

		// 截断日志文件
		if err := p.TruncateLogFile(); err != nil {
			log.Printf("截断日志文件失败: %v", err)
		}
	} else {
		log.Printf("当前北京时间: %s, 不是%s，跳过截断", beijingTime.Format("2006-01-02 15:04:05 Monday"), p.truncateWeekday)
	}
}
