# 互斥

完全分布式算法的模拟，单个Groutine模拟一个节点进程，当某个进程需要进入临界区的时候，需要向所有集群内的节点发送一条消息，如果收到集群内所有节点的回应，则可以进入临界区。
收到请求的节点根据自己当前的状态决定是否应答消息
* 如果当前正在临界区，则延迟应答
* 如果想要进入临界区，则根据消息的时间戳和它自己请求的时间进行对比，如果自己的时间戳比较小，则延迟应答，否则马上应答
* 如果当前不想进入临界区，则马上应答
