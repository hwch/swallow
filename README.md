# 类型：
## 基础类型： 整形，浮点，字符串，布尔，空值（nil）
## 符合类型： 数组，只读数组（元组），字典<br>
## 支持语句： if, elif, else, for, foreach, break, continue, return
## 支持类型定义： class, func

# 语法格式：
## 赋值语句--->
	```
	a=1
	a,b=f() (函数f返回值为2个)
	a=1,2
	a,b=f(),1
	```

## 判断语句---><br>
	```
	if a=True; a{复合语句}
	if True {复合语句}
	if True {复合语句} else if a=True;a{复合语句}
	if True {复合语句} else {复合语句}
	if True {复合语句} else if True {复合语句} else {复合语句}
	```

## 循环语句---><br>
	```
	for a=0;a<10;a++{复合语句}
	for ;True;a++ {复合语句}
	for ;True; {复合语句}
	for True {复合语句}

	foreach k,v=list {复合语句} (list 取值范围<数组，元组，字典，字符串>)
	```

## 函数定义---><br>
	```
	func 函数名(参数，可为空) {复合语句}
	```

## 类定义---><br>
	```
	class 类名[@父类] {赋值语句，函数定义}
	```


***注意：***
1. 元组里面的值不能被改变
2. break,continue不能再循环外
3. return不能再函数外
4. '_' 放到赋值左边会忽略赋值
5. list(begin[,stop])生成从begin到stop的list
6. 条件表达式不知限于布尔值
7. 两个"之间的字符串会对特殊字符做转义
8. 两个'之间的字符串会原样输出