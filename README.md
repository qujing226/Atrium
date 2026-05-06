# Atrium: Post-Quantum DID-Anchored AKE Protocol

Atrium is a research-grade prototype of a decentralized, post-quantum Authenticated Key Exchange (AKE) protocol. It addresses the fundamental **"Security-Latency Paradox"** inherent in Decentralized Identity (DID) networks, achieving 0-RTT session resumption without compromising long-term cryptographic security.

## Core Innovations

1. **S-AKE (Speculative Authenticated Key Exchange)**: Introduces optimistic execution to the AKE lifecycle. By leveraging local caches, S-AKE hides the high latency ($T_{chain}$) of blockchain consensus, enabling instantaneous communication initialization.
2. **Data Isolation Gate (DIG)**: A strict state machine mechanism that permits speculative decryption to maintain cryptographic synchronization (forward secrecy), while strictly isolating plaintexts from the application layer until eventual on-chain consistency is verified.
3. **Hybrid Post-Quantum Security**: Combines Ed25519 for identity binding with Kyber768 (ML-KEM) for quantum-resistant confidentiality, supplemented by a symmetric hash-based ratchet (Q-Ratchet) for efficient packet-level key evolution.

## Academic Documentation

The theoretical foundation, formal security proofs, and protocol specifications are detailed in the following documents:

*   📖 **[Protocol Specification](docs/SPECIFICATION.md)**: Defines the formal state machine, normative rules, and the TCP/Protobuf envelope architecture.
*   🧠 **[Theoretical Model & Formal Security](docs/THEORETICAL_MODEL.md)**: Details the formal definition of the Speculative State ($S_{spec}$), the probabilistic consistency model $P(\Delta t)$, and the overarching security theorem.
*   🛡️ **[Formal Security Game](docs/FORMAL_SECURITY_GAME.md)**: Provides the Game-Hopping proof reducing the protocol's security to EUF-CMA and IND-CCA2 assumptions, alongside realistic network probability bounds.
*   🌍 **[Deployment & Evaluation Plan](docs/DEPLOYMENT_PLAN.md)**: Outlines the WAN topology for empirical evaluation of the 0-RTT TTFB metrics.

## Origin

This project builds upon and heavily modifies architectures from previous works ([gochain](https://github.com/qujing226/go-chain), [easy-im](https://github.com/qujing226/easy-im)), optimizing core lattice-based algorithms and redesigning the DID authentication flow for pure decentralized operation.
 
