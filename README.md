###Linux命令行程序 Selpg

####功能 
从文本输入选入页范围进行处理的实用程序。
####使用方法
n代表整数  
必选参数（要求有序输入）：  
`-sn`：从第n页开始读取文件.   
`-en`: 读取文件至第n页.  
例：`./selpg -s1 -e2`:从标准输入读取第一、二页至标准输出。  
可选参数：  
`-ln`：设定每页的长度为n行,默认行号为72. 
`-f`:设定按结束符区分页。  
`-dn`: 将读取的数据传至打印机输出。但因并没有打印机测试输出，在本程序中，使用`-dn`参数实现对输入数据添加行号并输出（n为任意正整数）.  
`filename`: 需放在命令末尾，指输入来源。  
例：  
`./selpg -s1 -e2 -l5 -d1 selpg.go`：从文件selpg.go读取第1到第10行的数据，为每行添加行号，并将结果输出到屏幕。  
注：当可选参数缺失时，默认命令相当于`./selpg -sn -en -l72 os.Stdin`

####结果截图
将文件内容输出到屏幕：  
![Aaron Swartz](https://github.com/lonelyhope/selpg/blob/master/testPic/3.png?raw=true) 

将Stdin重定向到文件并输出：  
![Aaron Swartz](https://github.com/lonelyhope/selpg/blob/master/testPic/5.png?raw=true)

开启子程序并通过管道传递数据，输出子程序处理结果（为每一行添加行号）。使用l4命令，每一页为4行，输出3页，总共输出12行。  
![Aaron Swartz](https://github.com/lonelyhope/selpg/blob/master/testPic/2.png?raw=true)




