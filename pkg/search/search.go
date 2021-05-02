package search

import (
	"context"
	"bufio"
	"strings"
	"sync"
	"os"
	//"log"
)
type Result struct {
	Phrase string
	Line string
	LineNum int64
	ColNum int64
} 

func All(ctx context.Context, phrase string, files []string) <-chan []Result{
	ch:=make(chan []Result)
	wg:=sync.WaitGroup{}
	//root:=context.Background()
	ctx,cancel:=context.WithCancel(ctx)
	for i:=0;i<len(files);i++{
		wg.Add(1)
		go func(ctx context.Context, file string,ch chan<-[]Result,i int) {
			defer wg.Done()
			var lines []string
			ress := []Result{}
			f, err := os.Open(file)
			if err != nil {
				return 
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			for j:=0;j<len(lines);j++{
				if strings.Contains(lines[j],phrase){
					res := Result{
						Phrase:  phrase,
						Line:    lines[j],
						LineNum: int64(j + 1),
						ColNum:  int64(strings.Index(lines[j], phrase)) + 1,
					}
					ress = append(ress, res)
				}
			}
			if len(ress) > 0 {
				ch <- ress
			}
		}(ctx,files[i],ch,i)
	}
	go func() {
		defer close(ch)
		wg.Wait()

	}()
	cancel()
	return ch
}