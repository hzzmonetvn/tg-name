# 🚀 TG-Userboot: Next-Gen High-Performance Telegram Userbot

Author: **hzzmonetvn**  
Language: **Go (Golang 1.23+)**  
Protocol: **MTProto 2.0 via `gotd/td` (Pure Go, Zero CGO)**  
Sources Integrated: [plain-ub](https://github.com/thedragonsinn/plain-ub) & [PlainUB-Extras](https://github.com/R0Xofficial/PlainUB-Extras)

---

## 🌟 Overview / Tổng Quan

**TG-Userboot** is a complete, modern rewrite of the popular Python Pyrogram userbots (`plain-ub` and `PlainUB-Extras`) into pure Go. It consolidates all core userbot mechanics, group administration, multimedia processing, entertainment, system monitoring, dynamic AI text chat, and AI image generation into a single, blazing-fast, static native binary with minimal memory footprint.

### Key Architectural Advantages
- 🏎️ **Pure Go MTProto 2.0**: Built upon `gotd/td` without CGO dependencies; compiles to a single standalone binary.
- ⚡ **Asynchronous Concurrency**: Event-driven update dispatcher leveraging lightweight Goroutines.
- 🤖 **Universal Custom AI Engine**:
  - Support for **Custom AI Base URLs** (`AI_BASE_URL`): OpenAI, OpenRouter, DeepSeek, Groq, Ollama, vLLM, OneAPI, NewAPI, or private proxies.
  - Support for **Custom Chat Models** (`AI_MODEL`): `gpt-4o-mini`, `deepseek-chat`, `deepseek-reasoner`, `gemini-2.0-flash`, `claude-3-5-sonnet`, etc.
  - Built-in **Custom AI Image Generation** (`.img`, `.image`, `.genimg`, `.imagine`): Generates visual art via custom endpoints and models (`dall-e-3`, `flux`, `stable-diffusion-3`).
  - Native Google Gemini fallback support (`GEMINI_API_KEY`).
- 🔄 **Fully Verified & Modernized APIs**:
  - **Live Weather**: Integrated JSON `wttr.in` (`j1`) and OpenWeatherMap APIs.
  - **Real-Time Currency**: High-speed European Central Bank exchange rates via `api.frankfurter.dev/v1/latest` (preventing legacy redirects).
  - **Bulletproof Translation**: Dual-layer architecture with MyMemory Translation API (`api.mymemory.translated.net`) and Google Translate fallback (avoiding datacenter 429 rate limits).
  - **Wikipedia & Search**: Instant summary cards and web queries.
  - **Dynamic QR Engine**: Fast vector QR code generator via `api.qrserver.com`.
- 🛡️ **Fail-Safe Execution**: Built-in panic recovery across all command handlers ensuring 100% userbot uptime.

---

## 📁 Project Structure

```
tg-userboot/
├── cmd/
│   └── userbot/
│       └── main.go          # Application entrypoint & MTProto lifecycle
├── pkg/
│   ├── config/              # Environment & .env parser
│   ├── dispatcher/          # Command router, aliases, and permission engine
│   ├── session/             # Interactive terminal login & session persistence
│   └── types/               # Core domain models, Bot, and Message wrappers
├── plugins/
│   ├── admin/               # Ban, kick, mute, unmute, del, pin, zombies
│   ├── ai/                  # Universal custom AI chat & image generation
│   ├── android/             # Device specifications & root resources
│   ├── files/               # Spoiler and document tools
│   ├── fun/                 # Jokes, cowsay, shrug, decision decider
│   ├── media/               # QR code generation and image utilities
│   ├── misc/                # Weather, currency, translation, wikipedia
│   ├── sudo/                # Dynamic sudoers management and command directory
│   ├── system/              # Neofetch, hardware inspection, runtime specs
│   └── tgtools/             # Ping, alive, user info, chat info, block/unblock
├── bin/
│   └── userboot             # Compiled native ARM64/AMD64 executable
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

---

## 🛠️ Installation & Building / Hướng Dẫn Cài Đặt & Biên Dịch

### Requirements
- Go 1.22+ (Go 1.23.2 verified on Linux ARM64/AMD64)
- Telegram App API ID & Hash from [my.telegram.org](https://my.telegram.org)

### Build
```bash
cd tg-userboot
go mod tidy
go build -o bin/userboot cmd/userbot/main.go
```

### Run Tests
```bash
go test -v ./...
```

### Launch
```bash
# Configure .env
cp .env.example .env
nano .env

# Run the userbot
./bin/userboot
```

---

## ⚙️ Custom AI & Image Configuration (`.env`)

Configure your preferred AI provider in `.env`:

```env
# Custom AI Chat (Universal OpenAI / OneAPI / OpenRouter / DeepSeek / Ollama)
AI_API_KEY=sk-your-custom-key
AI_BASE_URL=https://api.openai.com/v1
AI_MODEL=gpt-4o-mini
AI_SYSTEM_PROMPT=You are a helpful and concise Telegram AI assistant.

# Custom Image Generation
IMAGE_API_KEY=sk-your-custom-key
IMAGE_BASE_URL=https://api.openai.com/v1
IMAGE_MODEL=dall-e-3
```

---

## 📖 Command Directory / Danh Sách Lệnh

| Category | Command | Aliases | Description |
| :--- | :--- | :--- | :--- |
| **AI** | `.ai` | `.ask`, `.aic`, `.chat`, `.gemini` | Query custom AI model via custom endpoint |
| **AI** | `.img` | `.image`, `.genimg`, `.imagine` | Generate images via custom image endpoint |
| **TGTools** | `.ping` | `.p` | Check bot latency and round-trip ping |
| **TGTools** | `.alive` | `.status` | Display bot uptime, Go runtime, and system status |
| **TGTools** | `.id` | `.chatid` | Show current chat, sender, or replied message ID |
| **TGTools** | `.info` | `.whois`, `.user` | Detailed user inspection card |
| **TGTools** | `.leave` | `.kickme` | Leave the current group or channel |
| **TGTools** | `.block` | - | Block specified user in private messages |
| **TGTools** | `.unblock` | - | Unblock specified user |
| **Admin** | `.ban` | `.sban` | Ban user with custom reason |
| **Admin** | `.unban` | - | Unban previously restricted user |
| **Admin** | `.kick` | `.skick`, `.dkick` | Kick user from group |
| **Admin** | `.mute` | `.smute` | Restrict user from sending messages |
| **Admin** | `.unmute` | - | Unmute restricted user |
| **Admin** | `.del` | `.purge` | Delete replied message |
| **Admin** | `.pin` | `.unpin` | Pin or unpin replied message |
| **Misc** | `.weather` | `.wttr` | Live weather forecast for any location |
| **Misc** | `.cash` | `.currency` | Real-time currency exchange conversion (Frankfurter dev API) |
| **Misc** | `.tr` | `.translate` | Multi-language translation (MyMemory + Google fallback) |
| **Misc** | `.wiki` | `.wikipedia` | Search Wikipedia article summaries |
| **Misc** | `.google` | `.g` | Generate Google search link |
| **Media** | `.makeqr` | `.qr`, `.genqr` | Generate QR code from text or link |
| **System** | `.neofetch` | `.sysinfo`, `.sys` | Host hardware, CPU, RAM, and OS specs |
| **Fun** | `.joke` | - | Fetch random humor joke |
| **Fun** | `.cowsay` | - | ASCII cow speech balloon |
| **Fun** | `.shrug` | - | Send `¯\_(ツ)_/¯` |
| **Fun** | `.decide` | - | Random decision maker |
| **Android** | `.magisk` | `.root` | Official Magisk and KernelSU release links |
| **Android** | `.specs` | - | Search GSMArena device specifications |
| **Files** | `.spoiler` | - | Wrap text in Telegram spoiler tag |
| **Sudo** | `.help` | `.cmds` | Interactive commands directory |
| **Sudo** | `.addsudo` | - | Add user ID to active sudo list |
| **Sudo** | `.delsudo` | - | Remove user ID from sudo list |
| **Sudo** | `.sudolist` | - | List all authorized sudo users |
