# Lombard Claimer & Sender

## 📌 Description

This software is designed for **claiming tokens from Lombard** and **sending them to deposit addresses**.
It supports three main modes of operation:

* **Claim** — claim allocation from Lombard.
* **Withdraw** — withdraw tokens from the vault.
* **CheckBalance** — check wallet balances for BARD tokens.

---

## ⚙️ Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/so-many-v2/lombard_claimer.git
   ```

2. **Fill in the input data:**

   * `data/wallets.txt` — private keys of wallets (one per line).
   * `data/deposit_addresses.txt` — deposit addresses for sending tokens.

3. **Configure settings:**
   In the `config.go` file, set the parameter:

   ```go
   SendTokensAfterClaim = true  // if you want to automatically send tokens to deposit addresses after claiming
   SendTokensAfterClaim = false // if you only want to claim without sending
   ```

---

## 🚀 Run

Navigate to the **app** folder:

```bash
cd app
```

Available commands (via `make`):

* **Withdraw tokens (Withdraw):**

  ```bash
  make withdraw
  ```

  > Runs `cmd/withdraw/main.go`

* **Check wallet balance (CheckBalance):**

  ```bash
  make checkBalance
  ```

  > Runs `cmd/checkBalance/main.go`

* **Claim tokens (Claim):**

  ```bash
  make claim
  ```

  > Runs `cmd/claim/main.go`

---

## 📝 Notes

* All private keys and addresses must be correct and correspond to the network specified in `config.go`.
* If `SendTokensAfterClaim = true` is enabled, tokens will be automatically sent to deposit addresses right after a successful claim.

