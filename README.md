# OIFastRun
Fast Compile and Run a CPP Program. ——For OIers

### OIFastRun v1.3.8 2019.8.14

> Author: xaxy
>
> Description: Fast Compile and Run a CPP Program.
>
> Usage: `oi b[uild] [-i INPUT_FILE] [-o OUTPUT_FILE] [-O2]`
>
> Usage: `oi r[un] [-i INPUT_FILE] [-o OUTPUT_FILE] [-O2]`
>
> (g++ Option "-g" Enabled by Default, It Will be Disabled If You Enabled "-O2")

Put binary into `PATH` and Enjoy it XD


##### Example:

```
D:\C++\hello\P3952\P3952>oi r -i complexity.cpp
编译的源代码被重定向到 D:\C++\hello\P3952\P3952\complexity.cpp
找到 1 份源代码：
> D:\C++\hello\P3952\P3952\complexity.cpp
>>>正在编译 D:\C++\hello\P3952\P3952\complexity.cpp -> complexity.exe ... 成功！

>>>>>>运行测试 D:\C++\hello\P3952\P3952\complexity.exe

找到 3 份输入数据：
> [0] D:\C++\hello\P3952\P3952\complexity1.in
> [1] D:\C++\hello\P3952\P3952\complexity2.in
> [2] D:\C++\hello\P3952\P3952\complexity3.in
是否全部使用？[Y/N] 默认Y 或 输入需要使用的数据编号（多个数据使用','隔开）:Y

>>>运行程序 D:\C++\hello\P3952\P3952\complexity.exe 输入重定向至 complexity1.in
Yes
No
No
No

找到 1 份答案数据：
> [0] D:\C++\hello\P3952\P3952\complexity1.ans
是否需要对比答案？ 输入需要使用的数据编号 或 [N]全不使用, [A]全部对比答案 默认使用编号[0]的数据:A
输入需要使用的数据编号 默认使用编号[0]的数据:
对比答案 complexity1.ans
---AC---

>>>运行程序 D:\C++\hello\P3952\P3952\complexity.exe 输入重定向至 complexity2.in
Yes
Yes
对比答案 complexity2.ans
---AC---

>>>运行程序 D:\C++\hello\P3952\P3952\complexity.exe 输入重定向至 complexity3.in
Yes
No
Yes
No
No
对比答案 complexity3.ans
---AC---

>>>数据统计 AC率: 3 / 3
> [ 0 ] AC | complexity1.in
> [ 1 ] AC | complexity2.in
> [ 2 ] AC | complexity3.in
```

