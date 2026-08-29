package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"os/exec"
	"sync"
)

type FileTail struct {
	filePath      string
	tailCmd       *exec.Cmd
	stdout        io.ReadCloser
	Lines         chan string
	tailingDone   sync.WaitGroup
	tailingCtx    context.Context
	tailingCancel context.CancelFunc
	formBegining  bool
}

func NewFileTail(filePath string, formBegining bool) (*FileTail, error) {
	inst := &FileTail{
		filePath:     filePath,
		formBegining: formBegining,
	}

	err := inst.startTailing(formBegining)
	if err != nil {
		return nil, err
	}

	return inst, nil
}

func (t *FileTail) startTailing(formBegining bool) error {

	if formBegining {
		t.tailCmd = exec.Command("tail", "-n", "+1", "-f", t.filePath)
	} else {
		t.tailCmd = exec.Command("tail", "-n", "0",  "-f", t.filePath)
	}

	t.Lines = make(chan string)

	var err error
	t.stdout, err = t.tailCmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = t.tailCmd.Start()
	if err != nil {
		return err
	}
	t.tailingCtx, t.tailingCancel = context.WithCancel(context.Background())
	t.tailingDone.Add(1)
	go func() {
		defer t.tailingDone.Done()
		scanner := bufio.NewScanner(t.stdout)
		for scanner.Scan() {
			select {
			case t.Lines <- scanner.Text():
			case <-t.tailingCtx.Done():
				return
			}
		}
	}()

	return nil
}

func (t *FileTail) Reset() error {
	t.tailingCancel()
	// 发送中断信号（更优雅的停止方式）
	if err := t.tailCmd.Process.Signal(os.Interrupt); err != nil {
		// 如果中断失败，强制杀死
		t.tailCmd.Process.Kill()
	}
	t.stdout.Close()
	t.tailingDone.Wait()
	close(t.Lines)

	return t.startTailing(t.formBegining)
}
