# Brynhild
Brynhild is a parallel execution framework which can be embedded in your blockchain or database. The core of brynhild is an adaptive parallel scheduler following the real workload, you should implement Executor and Storage by your self which can be plugged in brynhild. Benefit from loose couple design of brynhild, you can create customized parallel scheduler easily.

# Reference
- [Aria: A Fast and Practical Deterministic OLTP Database](http://www.vldb.org/pvldb/vol13/p2047-lu.pdf)

- [Scaling Hyperledger Fabric Using Pipelined Execution and Sparse Peers](https://arxiv.org/pdf/2003.05113.pdf)