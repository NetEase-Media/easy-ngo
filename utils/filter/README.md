## QA
date: 2023.10.09
Q: 
http router 配置时是否需要有method的概念，因为http path中包含.
A：
tracer 和 metrics 没有method的概念， 但是http模块有method概念，目前结论统一处理，http模块和tracer等均不区分method