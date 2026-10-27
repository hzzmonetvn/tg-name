# Architectural & Engineering Decisions Log

Author: **hzzmonetvn**

---

## 1. Documentation Location & Format

- **Decision**: Author a dedicated `ARCHITECTURE.md` file in the repository root rather than bloating `README.md` or adding heavy docstrings across all files.
- **Considered & Rejected**:
  - *Option A: In-line Docstrings only*: Rejected because docstrings are scattered across multiple files (`profile_engine.py`, `helper_bot.py`, `updater_state.py`, `forward_channel.py`) and do not provide an end-to-end view of the IPC data flow, fallback logic, or architectural diagrams.
  - *Option B: Bloating `README.md`*: Rejected because `README.md` serves as the user-facing handbook and quick-start guide. Adding deep architectural mechanics, MTProto rate limit considerations, and async event loop bindings would overwhelm users looking for quick setup instructions.
- **Why**: `ARCHITECTURE.md` creates a clean separation of concerns: `README.md` for end-user operations, and `ARCHITECTURE.md` for developers and maintainers requiring deep technical insights.

---

## 2. Bilingual Documentation (Tiếng Việt & English)

- **Decision**: Author `ARCHITECTURE.md` with bilingual headings, explanations, and concept mappings.
- **Considered & Rejected**:
  - *Option A: English only*: Rejected because the primary userbase, commit messages, and UI strings originated in Vietnamese.
  - *Option B: Vietnamese only*: Rejected because Telegram MTProto API conventions, Telethon terminology, and system engineering terms are globally understood in English.
- **Why**: Mirrors the bilingual architecture introduced in `i18n.py` (`vi` and `en`), providing consistent accessibility for both local and international developers.

---

## 3. Scope of Subsystems Documented

- **Decision**: Include both the core profile updater engine (`name.py`, `profile_engine.py`, `helper_bot.py`, `updater_state.py`, `i18n.py`, `setup.py`) and auxiliary migration utilities (`forward_channel.py`, `generate_session.py`).
- **Considered & Rejected**:
  - *Option A: Documenting only the Profile Updater*: Rejected because `forward_channel.py` is an active utility in the repository that shares `.env` configuration, uses Telethon MTProto sessions, and handles similar FloodWait constraints.
- **Why**: Provides a single source of truth for all tools existing within the repository.

---

## 4. Concurrency & IPC Communication Architecture

- **Decision**: Document the dynamic `force_update_event` property binding to `asyncio.get_running_loop()`.
- **Considered & Rejected**:
  - *Option A: Static global `asyncio.Event`*: Documenting a static event would misrepresent the real implementation. A static event causes `RuntimeError` if initialized before or across differing asyncio event loops.
- **Why**: Deep technical documentation must reflect the actual solution implemented in `updater_state.py` to prevent future regressions.

---

## 5. Fallback Hierarchy & Rate Limit Mitigation

- **Decision**: Explicitly diagram the two-tier weather resolution (OpenWeatherMap -> wttr.in -> fallback tag) and the Last.fm music vs weather bio prioritization.
- **Considered & Rejected**:
  - *Option A: High-level summary without edge case handling*: Rejected because Telegram's 70-character bio limit and API FloodWait penalties are the primary failure points in profile automation bots.
- **Why**: Gives future contributors immediate context on why string truncation and exception catching are strictly enforced at line boundaries.

---

## 6. PlainUB-Go: Go MTProto Library Selection (`gotd/td`)

- **Decision**: Use `gotd/td` (Pure Go implementation of Telegram MTProto 2.0) over `go-tdlib`.
- **Considered & Rejected**:
  - *Option A: `Arman92/go-tdlib`*: Rejected because it requires CGO, linking against external `libtdjson.so`, which complicates cross-compilation and requires host C libraries.
  - *Option B: Bot API wrapper (`tbot`)*: Rejected because the userbot requires full user MTProto client functionality (`first_name`, direct chat actions, self-account control) which Bot API cannot provide.
- **Why**: `gotd/td` is 100% pure Go, fully concurrent via Goroutines, has built-in session storage, generates clean native binaries without C dependencies, and directly speaks the Telegram MTProto binary protocol.

---

## 7. Migration Scope & Modular Plugin Architecture

- **Decision**: Build a modular dispatcher architecture in `pkg/dispatcher` with dedicated plugin packages for each domain (`admin`, `tgtools`, `ai`, `misc`, `media`, `system`, `fun`, `android`, `files`, `sudo`), consolidating all 80+ commands into clean categorized modules.
- **Considered & Rejected**:
  - *Option A: Monolithic single file router*: Rejected due to maintainability issues and poor separation of concerns.
  - *Option B: Python C-extension bridge*: Rejected because running Python via CGO would defeat the performance and memory benefits of Go.
- **Why**: Clean package boundaries make it trivial to add, disable, or test plugins independently while maintaining a shared command registry and safe panic recovery.

---

## 8. Third-Party API Modernization

- **Decision**: Upgrade all outdated scrapers and deprecated Python endpoints to modern REST APIs:
  - Gemini: Modern Google Gen AI v1beta endpoint (`gemini-1.5-flash`, `gemini-1.5-pro`) replacing deprecated Python libraries.
  - Weather: Modern `wttr.in?format=j1` JSON endpoint and OpenWeatherMap integration.
  - Currency: `https://api.frankfurter.app/latest` for real-time European Central Bank rates.
  - Translation: Google Translate API endpoint returning structured translation arrays.
  - QR: Modern vector QR code API service.
- **Considered & Rejected**:
  - *Option A: Keeping legacy scraping scripts*: Rejected because old web scraping endpoints break frequently and return unstructured HTML instead of clean JSON.
- **Why**: Drastically improves reliability, reduces latency, and removes brittle regex parsing.

---

## 9. Universal Custom AI & Image Generation Architecture (`TG-Userboot`)

- **Decision**: Architect a universal, vendor-agnostic AI client supporting custom endpoints (`AI_BASE_URL`, `IMAGE_BASE_URL`) and custom models (`AI_MODEL`, `IMAGE_MODEL`) via standard OpenAI-compatible JSON specifications (`/v1/chat/completions`, `/v1/images/generations`), while preserving native Google Gemini endpoint fallback.
- **Considered & Rejected**:
  - *Option A: Hardcoding specific AI providers (e.g. only OpenAI or only Gemini)*: Rejected because users frequently employ private proxies, local Ollama/vLLM instances, OpenRouter, DeepSeek, Groq, or OneAPI aggregators.
  - *Option B: Separate binary for image generation*: Rejected because consolidating text and image commands (`.ai` and `.img`) into `plugins/ai` keeps the binary unified and simplifies deployment.
- **Why**: Allows users to switch between local models (Ollama, vLLM) and cloud providers (DeepSeek, OpenAI, OpenRouter) merely by changing `.env` variables without recompilation.

---

## 10. Live Endpoint Health Audit & Remediation

- **Decision**: Perform live HTTP probing across all external third-party endpoints and remediate failing/outdated routes:
  1. *Currency API*: `api.frankfurter.app` returned HTTP 301 redirects to `api.frankfurter.dev/v1/latest`. The endpoint was immediately upgraded to the official direct `.dev/v1/latest` URL.
  2. *Translation API*: Google Translate GTX endpoint (`translate.googleapis.com`) returned HTTP 429 (Too Many Requests) due to datacenter IP blocking. Upgraded to a dual-layer architecture: primary translation via MyMemory API (`api.mymemory.translated.net`) with Google fallback.
- **Why**: Proactive endpoint verification ensures zero silent failures or broken commands when deployed to production servers.

---

## 11. MTProto Peer AccessHash Caching, Execution Contexts & Channel-Aware Deletions

- **Decision**:
  1. *Peer AccessHash Caching*: Maintain an in-memory `EntityCache` in `pkg/dispatcher` that dynamically extracts and stores `AccessHash` values for users and channels from incoming MTProto `Updates` / `UpdatesCombined`. When creating `InputPeerUser` or `InputPeerChannel`, access hashes are retrieved to prevent MTProto `CHANNEL_INVALID` or `PEER_ID_INVALID` RPC errors.
  2. *Detached Execution Context*: Run command handler Goroutines with a detached context with a 2-minute deadline (`context.WithTimeout(context.Background(), 2*time.Minute)`) rather than inheriting the short-lived update handler context, ensuring long-running AI queries or external RPCs are not prematurely aborted.
  3. *Channel-Aware Deletion*: Distinguish between basic chats / user dialogs (`MessagesDeleteMessages`) and supergroups / channels (`ChannelsDeleteMessages`) in `types.Bot.Delete()` and `plugins/admin/handleDel`.
  4. *Automatic Replied Message Resolution*: When a command is triggered in reply to a message, automatically fetch and populate `ReplyToUserID` and `ReplyToText` so plugins (`.ai`, `.tr`, `.ban`, `.info`) immediately possess target text and user context.
  5. *Fail-Safe Session Reconnect*: In `profile_engine.py`, explicitly disconnect and reset `self.user_client = None` if authorization verification fails, preventing infinite loops on revoked sessions.
  6. *Reliable Album Resumption*: In `forward_channel.py`, advance `last_scanned_id` checkpoints only upon confirmed delivery of media albums, ensuring failed batches are automatically retried upon resumption.
- **Why**: Eliminates edge-case runtime failures and ensures production-grade reliability across both the Go MTProto userbot and Python profile engine.



