# Building a TCP Server from Scratch 

A minimal, low-level TCP server built in Go that listens on a dedicated port, accepts incoming client connections, and reads raw binary network bytes, translating them back into human-readable text.

---

## Architecture & Core Concepts

Today's task forced us to peel back the layers of network programming. Instead of treating a "TCP Server" as a black box, we broke it down into its physical and digital equivalents:

* **The Mailbox (The Listener):** `net.Listen("tcp", "localhost:6379")` claims ownership over a specific port (6379) at our local address. This is the digital equivalent of reserving a physical mailbox so no other application can take your mail.
* **The Handshake (Accepting):** `listener.Accept()` is a **blocking** operation. It pauses our program completely until a client (like `netcat`) knocks on the door. Once it wakes up, it grants us a dedicated connection tube (`conn`).
* **The Bucket (The Buffer):** Data travels across network wires as raw numbers (bytes). `make([]byte, 1024)` allocates a temporary holding area in memory to catch those bytes as they land.
* **The Translator (Reading):** `conn.Read(buf)` reads the raw binary data into our bucket, returns the exact number of bytes written (`n`), and `string(buf[:n])` translates those ASCII byte codes back into readable characters.

---

