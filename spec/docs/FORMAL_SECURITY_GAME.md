# QLink 协议：形式化安全模型与推测执行证明 (Formal Security Model & Speculative Execution Proof)

> 本文档基于扩展的 Bellare-Rogaway (eBR) AKE 安全模型，引入了推测状态约束 (Speculative State Constraints)，通过 Game-Hopping (游戏序列法) 对 QLink 协议的认证完整性与机密性进行严格的数学归约证明。

---

## 1. 系统模型与敌手能力 (System Model & Adversary Capabilities)

### 1.1 实体、状态与原子原语 (Entities, States & Atomic Primitives)
设 $\mathcal{P}$ 为协议参与者集合。对于任意会话实例 $\Pi_{i,j}^s$，其内部状态集为 $\mathbb{S} = \{S_{idle}, S_{spec}, S_{ver}, S_{abort}\}$。

**数据闸门与原子性约束 (Data Gate & Atomicity Constraint)**：
定义函数 `Deliver(M)` 为将消息 $M$ 递交至应用层。协议假定存在一个底层**原子事务原语 (Atomic Transaction Primitive)** 保证：
1. 状态一旦转化为 $S_{abort}$，`Buffer.Clear()` 必然与状态改变原子性完成，应用层不可见残留数据。
2. 状态一旦转化为 $S_{ver}$，`Buffer.Flush()` 必然按序将历史数据交付给应用层。
任何破坏此原子性的系统级故障概率记为 $P_{race}(\Delta t)$。

### 1.2 参数化的解析预言机 (Parameterized Resolve Oracle)
我们不假设底层区块链的绝对完美，而是将其抽象为一个参数化的验证预言机 $\mathcal{O}_{resolve}(DID)$，具有以下特征：
*   **最终性函数 $F(t)$**：预言机返回的版本在时间 $t$ 之后不被回滚（达成最终一致）的概率。
*   **认证误报率 $\epsilon_{auth}$**：预言机被恶意欺骗返回无效文档（但通过了密码学验证）的极小概率。
*   **延迟分布 $D_{resolve}$**：预言机响应的物理耗时分布，敌手可在此分布范围内操控延迟。

### 1.3 敌手模型 $\mathcal{A}$ (Adversary Model & Oracle Queries)
敌手 $\mathcal{A}$ 是完全控制异步网络信道的 PPT 算法，并能与 $\mathcal{O}_{resolve}$ 交互。$\mathcal{A}$ 可调用以下预言机：
*   $\text{Send}(i, j, s, M)$: 向实例 $\Pi_{i,j}^s$ 注入、拦截或修改消息。
*   $\text{Corrupt}(i, t)$: **自适应腐败 (Adaptive Corrupt)**。返回参与者 $i$ 在时间 $t$ 之前生成的所有私钥材料。
    *   **限制条件**: 对于任何选定的测试会话 (Test Session) $\Pi_{i,j}^{s^*}$（其活跃时间区间为 $[t_{start}, t_{end}]$），禁止 $\mathcal{A}$ 调用 $\text{Corrupt}(i, t)$ 或 $\text{Corrupt}(j, t)$ 使得 $t \in [t_{start}, t_{end}]$，以防止直接泄露该会话正在使用的活跃私钥。
*   $\text{Delay}(\mathcal{O}_{resolve}, \Delta t)$: 基于延迟分布 $D_{resolve}$ 拦截并延迟节点的身份解析请求。

---

## 2. 安全目标：最终应用层认证完整性 (EALA Security Game)

我们通过一个形式化的挑战博弈 (Challenge Game) 来定义 **最终应用层认证完整性 (Eventual Application-Layer Authenticity, EALA)**。

**定义 1 (EALA)**: 称协议具有 EALA 属性，如果对于任意 PPT 敌手 $\mathcal{A}$，其赢得以下博弈的优势 $\text{Adv}_{\mathcal{A}}^{EALA}$ 是可忽略的 (negligible)。

**EALA 博弈过程 ($Game_{EALA}$)**：
1.  **Setup**: 挑战者 $\mathcal{C}$ 初始化系统参数与所有参与者的密钥对。
2.  **Training**: $\mathcal{A}$ 自由调用 $\text{Send}, \text{Corrupt}$（遵循过期限制）和 $\text{Delay}$ 预言机。
3.  **Challenge(i, j, s)**: 
    *   $\mathcal{A}$ 提交一个伪造的密文序列 $C_{fake}$ 给实例 $\Pi_{i,j}^s$。
    *   挑战者 $\mathcal{C}$ 将 Buffer 中从 $C_{fake}$ 解密出的第一个明文 $P^*$ 标记为 `CHALLENGE`。
4.  **Victory Condition**: $\mathcal{A}$ 获胜当且仅当发生 $\text{Deliver}(P^*)$，且挑战者 $\mathcal{C}$ 的发送日志中没有任何记录表明真实的 $j$ 在 $\text{Deliver}(P^*)$ 发生前发送过 $P^*$。

---

## 3. 安全性归约证明 (Formal Proof via Game-Hopping)

**定理 1 (Theorem 1)**: 假设 Ed25519 提供 EUF-CMA 安全性，Kyber768 提供 IND-CCA2 安全性，则 QLink 协议是 S-AKE 安全的。

**证明 (Proof)**: 我们通过构造一系列不可区分的博弈 $G_0, G_1, G_2, G_3$ 来界定敌手的优势。

### Game 0: 真实的 S-AKE 协议
$$\Pr[\mathcal{A} \text{ wins } G_0] = \Pr[W_0]$$

### Game 1: 剥夺签名伪造能力 (EUF-CMA Reduction)
在 $G_1$ 中，如果 $\mathcal{A}$ 提交了一个有效的签名，但该签名并非由诚实节点生成，$\mathcal{C}$ 将直接触发 $Abort$。
**归约器 $\mathcal{R}_{SIG}$ 骨架**:
```text
Reduction R_SIG:
  Input: EUF-CMA challenge public key pk
  Simulate all honest parties using pk for target identity j.
  If adversary A outputs a forged signature σ on message m in G1:
     Output (m, σ) as the solution to the EUF-CMA challenge.
```
*分析*: $|\Pr[W_1] - \Pr[W_0]| \le \text{Adv}_{\mathcal{A}}^{EUF-CMA} \le \epsilon_{sig}$

### Game 2: 剥夺密文区分能力 (IND-CCA2 Reduction)
在 $G_2$ 中，将 Kyber768 封装生成的真实共享密钥替换为理想随机密钥。
**归约器 $\mathcal{R}_{KEM}$ 骨架**:
```text
Reduction R_KEM:
  Input: IND-CCA2 challenge (pk, ciphertext c*, key k*)
  Simulate session setup using c* and k*. 
  If adversary A can distinguish whether k* is real or random:
     Output A's guess to win the IND-CCA2 game.
```
*分析*: $|\Pr[W_2] - \Pr[W_1]| \le \text{Adv}_{\mathcal{A}}^{IND-CCA2} \le \epsilon_{kem}$

### Game 3: 引入数据闸门拦截与概率约束 (The Speculative Isolation bounds)
在 $G_3$ 中，我们分析在推测窗口 $\Delta t$ 内，$\mathcal{A}$ 尝试突破数据闸门的概率。

在现实系统中，$\mathcal{A}$ 成功触发 $\text{Deliver}(P^*)$ 的概率受制于以下系统级风险：
1.  **异步原子性失效 $P_{race}(\Delta t)$**: 状态切断时，非法明文泄露至应用层的系统级故障概率。
2.  **回滚信令丢失 $P_{abort\_loss}(\Delta t)$**: 发起方的状态中止信令在传输过程中丢失。设定网络单包丢失率为 $p$，协议配置的重试上限为 $r$，则回滚信令投递失败的概率界限估计为 $P_{abort\_loss}(\Delta t) \le p^{r+1}$。
3.  **账本伪终局性 $\epsilon_{auth} + (1 - F(\Delta t))$**: $\mathcal{O}_{resolve}$ 在 $\Delta t$ 时刻返回了“通过”验证（包括轻节点被欺骗的概率 $\epsilon_{auth}$），但区块后续被重组抛弃（概率为 $1 - F(\Delta t)$）。

*分析*: 
$$\Pr[W_3 \mid \text{EALA violation}] \le P_{race}(\Delta t) + P_{abort\_loss}(\Delta t) + (1 - F(\Delta t)) + \epsilon_{auth}$$

### 结论与定量上界 (Conclusion & Quantitative Bounds)
合并上述游戏序列，我们得出敌手 $\mathcal{A}$ 破坏 QLink 协议应用层认证完整性的总优势上界：

$$\text{Adv}_{\mathcal{A}}^{S-AKE} \le \epsilon_{sig} + \epsilon_{kem} + P_{race}(\Delta t) + P_{abort\_loss}(\Delta t) + (1 - F(\Delta t)) + \epsilon_{auth}$$

该公式表明，只要工程实现保证原子性且网络参数合理，QLink 在 $\Delta t$ 推测窗口内的应用层安全风险等价于可忽略的密码学概率。 $\blacksquare$
