package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var NameSpace string

var rootCmd = &cobra.Command{
	Use:   "logtail",
	Short: "kubectl tail log for ds/deploy/sts",
	Long: `kubectl tail log for ds/deploy/sts. For example:
kubectl logtail name`,
	Run: runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	fmt.Printf("execute %s with args:%v \n", cmd.Name(), args)
	kubectl, _ := exec.LookPath("kubectl")
	var cmdArgs string
	cmdArgs = kubectl + " get po "
	if NameSpace != "" {
		cmdArgs = cmdArgs + "-n " + NameSpace
	}

	cmd1 := exec.Command("bash", "-c", cmdArgs)
	out, err := cmd1.Output()
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(out), "\n")
	var wg sync.WaitGroup
	for _, v := range lines {
		if strings.HasPrefix(v, args[0]) {
			podName := strings.Split(v, " ")[0]
			cmd2 := kubectl + " logs " + podName + " --tail 10 -f"
			if NameSpace != "" {
				cmd2 = cmd2 + " -n " + NameSpace
			}
			wg.Add(1)
			go Command(cmd.Context(), cmd2, podName, wg.Done)
		}
	}
	wg.Wait()
	<-cmd.Context().Done()
}

func Execute(ctx context.Context) {
	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&NameSpace, "namespace", "n", "", "namespace")
}

func Command(ctx context.Context, cmd string, podName string, done func()) error {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	defer done()
	println("===start to exec:", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			select {
			case <-ctx.Done():
				if ctx.Err() != nil {
					fmt.Printf("程序出现错误: %q\n", ctx.Err())
				} else {
					fmt.Println("程序被终止")
				}
				println("===end to exec:", cmd)
				return
			default:
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					break
				}
				fmt.Print(podName, ": ", readString)
			}
		}
	}()
	return c.Run()
}
