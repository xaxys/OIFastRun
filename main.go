package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gookit/color"
)

var (
	//PATH 目标文件夹
	PATH string
	//TAG 1仅编译，2编译并运行
	TAG int
	//outPutFile 输出文件名
	outPutFile string
	//inPutFile 输入文件名
	inPutFile string
	//O2 是否开启O2优化
	O2 bool
	//COMPAREANS 是否对比答案文件 -1待询问 自然数是编号 -2否
	COMPAREANS = -1
	//EXT 可执行文件后缀名
	EXT string
	//SplitOutput 在测试程序运行结束时进行输出
	SplitOutput bool
)

func init() {
	switch runtime.GOOS {
	case "windows":
		EXT = ".exe"
	case "linux":
		EXT = ""
	}
	PATH, _ = filepath.Abs("")
}

//ParseCmd 解析命令
func ParseCmd() {
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(1)
	} else if os.Args[1] == "build" || os.Args[1] == "b" {
		TAG = 1
	} else if os.Args[1] == "run" || os.Args[1] == "r" {
		TAG = 2
	} else {
		printHelp()
		os.Exit(1)
	}

	flag.StringVar(&outPutFile, "o", "", "Rename Output File")
	flag.StringVar(&inPutFile, "i", "", "Specify Input File")
	flag.BoolVar(&O2, "O2", false, "Rename Output File")
	flag.BoolVar(&SplitOutput, "s", false, "Split Output")

	flag.CommandLine.Parse(os.Args[2:])
}

func main() {
	ParseCmd()
	codeList := SearchCode()
	var list []string
	if TAG == 1 || TAG == 2 {
		list = CompileCode(codeList)
	}
	if TAG == 2 {
		RunCode(list)
	}
}

func printHelp() {
	fmt.Println(
		`[OIFastRun] OIFastRun v1.4.11 2019.11.7
            Author: xaxy
            Description: Fast Compile and Run a CPP Program.
            Usage: oi b[uild] [-i INPUT_FILE] [-o OUTPUT_FILE] [-O2] [-s]
            Usage: oi r[un] [-i INPUT_FILE] [-o OUTPUT_FILE] [-O2] [-s]
            (-O2 : g++ Option "-g" Enabled by Default, It Will be Disabled If You Enabled "-O2")
            (-s  : Print all output at the end of execution)`)
}

//RunCode 运行代码
func RunCode(list []string) {
	for _, file := range list {
		fmt.Println()
		fmt.Println(">>>>>>运行测试", file)
		COMPAREANS = -1

		inputFile, err := filepath.Glob(filepath.Join(filepath.Dir(file), "*.in"))
		if err != nil {
			color.Error.Println("[ERROR]", err.Error())
		}

		tot := 0
		ac := 0
		var statue []string
		if len(inputFile) > 0 {
			fmt.Println()
			fmt.Println("找到", len(inputFile), "份输入数据：")
			for i, v := range inputFile {
				fmt.Printf("> [%d] %s", i, v)
				fmt.Println()
			}
			fmt.Print("是否全部使用？[Y/N] 默认Y 或 输入需要使用的数据编号（多个数据使用','隔开）:")
			var s string
			var testList []string
			fmt.Scanln(&s)
			if s == "Y" || s == "y" || s == "" {
				testList = inputFile
			} else {
				numList := strings.Split(s, ",")
				for _, i := range numList {
					a := getDigit(i)
					if a != -1 && a < len(inputFile) {
						testList = append(testList, inputFile[a])
					}
				}

			}
			if len(testList) > 0 {
				for _, v := range testList {
					res, sta := testCode(file, v)
					if res == 1 {
						ac++
					}
					if res != 0 {
						tot++
					}
					statue = append(statue, sta+" | "+filepath.Base(v))
				}
				if COMPAREANS != -2 {
					fmt.Println()
					fmt.Println(">>>数据统计 AC率:", ac, "/", tot)
					for i, v := range statue {
						fmt.Printf("> [%d] %s", i, v)
						fmt.Println()
					}
				}
				continue
			}
		}
		testCode(file, "")
	}
}

func fileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	return true
}

func testCode(file string, v string) (int, string) {
	var ac = 0
	var statue = "UK"
	var err error
	var input []byte
	cmd := exec.Command(file)
	fmt.Println()
	if v == "" {
		cmd.Stdin = os.Stdin
		fmt.Println(">>>运行程序", file, "已对接标准输入stdin")
	} else {
		fmt.Println(">>>运行程序", file, "输入重定向至", filepath.Base(v))
		input, err = readFileByte(v)
		if err != nil {
			color.Error.Println("[ERROR]", err.Error())
		}
	}

	stdout, stderr, err := execCommand(cmd, input, !SplitOutput, true, true)
	if err != nil {
		color.Error.Println("[ERROR]", err.Error())
		fmt.Println()
		color.LightMagenta.Println("---RE---")
		ac = -1
		statue = color.FgLightMagenta.Render("RE")
		return ac, statue
	}

	if SplitOutput {
		fmt.Println()
		if len(stdout) > 0 {
			color.Gray.Println("=====[STDOUT]=====")
			for _, v := range stdout {
				fmt.Print(v)
			}
			fmt.Println()
			color.Gray.Println("===[END STDOUT]===")
		}
		if len(stderr) > 0 {
			color.LightRed.Println("=====[STDERR]=====")
			for _, v := range stderr {
				color.Red.Print(v)
			}
			fmt.Println()
			color.LightRed.Println("===[END STDERR]===")
		}
	}

	if COMPAREANS == -2 {
		return ac, statue
	}

	var ansFile []string
	ans := strings.Replace(v, ".in", ".ans", 1)
	out := strings.Replace(v, ".in", ".out", 1)
	if fileExist(ans) {
		ansFile = append(ansFile, ans)
	}
	if fileExist(out) {
		ansFile = append(ansFile, out)
	}

	comp := COMPAREANS
	if len(ansFile) > 0 && COMPAREANS == -1 {
		fmt.Println()
		fmt.Println("找到", len(ansFile), "份答案数据：")
		for i, v := range ansFile {
			fmt.Printf("> [%d] %s", i, v)
			fmt.Println()
		}
		fmt.Print("是否需要对比答案？ 输入需要使用的数据编号 或 [N]全不使用, [A]全部对比答案 默认使用编号[0]的数据:")
		var s string
		fmt.Scanln(&s)
		if s == "A" || s == "a" {
			fmt.Print("输入需要使用的数据编号 默认使用编号[0]的数据:")
			fmt.Scanln(&s)
			a := getDigit(s)
			if a == -1 {
				COMPAREANS = 0
				comp = 0
			} else {
				COMPAREANS = a
				comp = a
			}
		} else if s == "N" || s == "n" {
			COMPAREANS = -2
		} else if a := getDigit(s); a != -1 {
			comp = a
		} else {
			comp = 0
		}

	}

	if len(ansFile)-1 < comp {
		comp = len(ansFile) - 1
	}
	if comp >= 0 {
		fmt.Println("对比答案", filepath.Base(ansFile[comp]))
		f, tip, line := compFile(stdout, ansFile[comp])
		if f {
			fmt.Println()
			color.LightGreen.Println("---Accepted---")
			ac = 1
			statue = color.FgLightGreen.Render("AC")
		} else {
			fmt.Println()
			color.LightRed.Println("---Wrong Answer---")
			fmt.Println(tip, "at line:", line)
			ac = -1
			statue = color.FgLightRed.Render("WA")
		}
	}

	return ac, statue
}

func compFile(ctx []string, file string) (bool, string, int) {
	f, err := os.Open(file)
	if err != nil {
	}
	reader := bufio.NewReader(f)

	line := 0
	for {
		str, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if len(ctx) <= line {
			return false, "Too Few Lines", line
		}

		a := strings.TrimRight(ctx[line], "\n")
		a = strings.TrimRight(a, "\r")
		a = strings.TrimRight(a, " ")
		b := strings.TrimRight(str, "\n")
		b = strings.TrimRight(b, "\r")
		b = strings.TrimRight(b, " ")
		if strings.Compare(a, b) != 0 {
			fmt.Println("[Detail] Ans  :", b)
			fmt.Println("[Detail] Yours:", a)
			return false, "Wrong Answer", line + 1
		}
		line++
	}
	return true, "Accepted", line - 1
}

func getDigit(s string) int {
	x := 0
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			x = -1
			break
		} else {
			x = x*10 + int(s[i]) - int('0')
		}
	}
	return x
}

func readFileByte(file string) ([]byte, error) {
	if file == "" {
		return nil, nil
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	ctx, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		return nil, err2
	}
	return ctx, nil
}

//SearchCode 获取待编译文件
func SearchCode() []string {
	var cppFile []string
	if inPutFile == "" {
		t, err := filepath.Abs("*.cpp")
		if err != nil {
			color.Error.Println("[ERROR]", err.Error())
		}
		cppFile, err = filepath.Glob(t)
		if err != nil {
			color.Error.Println("[ERROR]", err.Error())
		}
	} else {
		inPutFile, _ = filepath.Abs(inPutFile)
		_, err := os.Stat(inPutFile)
		if err != nil {
			fmt.Println(inPutFile, "文件不存在！")
			os.Exit(1)
		}
		cppFile = []string{inPutFile}
		fmt.Println("编译的源代码被重定向到", inPutFile)
	}
	if len(cppFile) == 0 {
		fmt.Println("文件夹中没有源代码文件！")
		return cppFile
	}
	if len(cppFile) > 1 && outPutFile != "" {
		color.Warn.Tips("文件夹中包含超过一个源代码文件！ '-o' 将不会生效！")
	}

	fmt.Println("找到", len(cppFile), "份源代码：")
	for _, v := range cppFile {
		fmt.Println(">", v)
	}
	return cppFile
}

//CompileCode 编译代码
func CompileCode(cppFile []string) []string {
	var list []string
	for _, v := range cppFile {
		var name string
		var fullName string
		if outPutFile != "" {
			name = outPutFile
			fullName, _ = filepath.Abs(outPutFile)
		} else {
			fullName = strings.Replace(v, ".cpp", EXT, 1)
			name = filepath.Base(fullName)
		}

		if PATH != filepath.Dir(fullName) {
			fmt.Print(">>>正在编译 ", v, " -> ", fullName, " ... ")
		} else {
			fmt.Print(">>>正在编译 ", v, " -> ", name, " ... ")
		}

		var cmd *exec.Cmd
		if O2 {
			cmd = exec.Command("g++", v, "-o", fullName, "-O2")
		} else {
			cmd = exec.Command("g++", v, "-o", fullName, "-g")
		}

		stdout, stderr, err := execCommand(cmd, nil, false, false, true)
		if err != nil {
			color.Red.Println("失败！")
		} else {
			color.Green.Println("成功！")
		}
		if len(stderr) > 0 {
			color.Red.Println("=====[STDERR]=====")
			for _, v := range stderr {
				color.Red.Print(v)
			}
			fmt.Println()
			color.Red.Println("===[END STDERR]===")
		}
		if len(stdout) > 0 {
			for _, v := range stdout {
				fmt.Print(v)
			}
		}
		if err != nil {
			color.Error.Println("[ERROR]", err.Error())
		} else {
			list = append(list, fullName)
		}
	}
	return list
}

func execCommand(cmd *exec.Cmd, input []byte, printStdout bool, printStderr bool, record bool) ([]string, []string, error) {
	var outArray []string
	var errArray []string

	outpipe, err := cmd.StdoutPipe()
	if err != nil {
		color.Error.Println("[ERROR]", err.Error())
		return nil, nil, err
	}
	errpipe, err2 := cmd.StderrPipe()
	if err2 != nil {
		color.Error.Println("[ERROR]", err2.Error())
		return nil, nil, err2
	}

	if input != nil {
		inpipe, err := cmd.StdinPipe()
		if err != nil {
			color.Error.Println("[ERROR]", err.Error())
		}
		go func() {
			_, err = inpipe.Write(input)
			if err != nil {
				color.Error.Println("[ERROR]", err.Error())
				color.Error.Println("若出现 The pipe has been ended 可能是尚未写入所有测试数据而程序过早地结束")
			}
			inpipe.Close()
		}()
	}

	outReader := bufio.NewReader(outpipe)
	errReader := bufio.NewReader(errpipe)
	cmd.Start()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for {
			line, err := outReader.ReadString('\n')
			if line != "" {
				if record {
					outArray = append(outArray, line)
				}
				if printStdout {
					fmt.Print(line)
				}
			}
			if err != nil || io.EOF == err {
				wg.Done()
				break
			}
		}
	}()
	go func() {
		for {
			line, err := errReader.ReadString('\n')
			if line != "" {
				if record {
					errArray = append(errArray, line)
				}
				if printStderr {
					fmt.Fprint(os.Stderr, line)
				}
			}
			if err != nil || io.EOF == err {
				wg.Done()
				break
			}
		}
	}()
	wg.Wait()

	err3 := cmd.Wait()
	return outArray, errArray, err3
}
