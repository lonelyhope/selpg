Linux命令行程序 Selpg  

功能   
从文本输入选入页范围进行处理的实用程序。   

使用方法    
n代表整数  
必选参数（要求有序输入）：  
`-s int`：从第n页开始读取文件.   
`-e int`: 读取文件至第n页.  
例：`./selpg -s 1 -e 2`:从标准输入读取第一、二页至标准输出。  
可选参数：  
`-l int`：设定每页的长度为n行,默认行号为72. 
`-f`:设定按结束符区分页。  
`-d string`: 将读取的数据传至打印机输出。但因并没有打印机测试输出，在本程序中，使用`-dn`参数实现对输入数据添加行号并输出.    
`filename`: 需放在命令末尾，指输入来源。缺失时为标准输入。   
例：  
`./selpg -s 1 -e 2 -l 5 -d 1 selpg.go`：从文件selpg.go读取第1到第10行的数据，为每行添加行号，并将结果输出到屏幕。  

####结果截图
标准输出输出到屏幕：  
![Aaron Swartz](https://github.com/lonelyhope/selpg/blob/master/testPic/stdin_out.png?raw=true) 

将Stdin重定向到文件并输出：   
![Aaron Swartz](https://github.com/lonelyhope/selpg/blob/master/testPic/fileInOut.png?raw=true)

开启子程序并通过管道传递数据，输出子程序处理结果（为每一行添加行号）。使用l 6命令，每一页为6行，输出2、3页，总共输出12行。   
![Aaron Swartz](https://github.com/lonelyhope/selpg/blob/master/testPic/pipe.png?raw=true)




