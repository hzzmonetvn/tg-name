# 🏛️ Architecture & Technical Reference / Tài Liệu Kiến Trúc & Kỹ Thuật

Author: **hzzmonetvn**  
Language: **Bilingual (Tiếng Việt & English)**  
Target: **Complete Codebase Deep Architecture**

---

## 📑 Mục lục / Table of Contents

1. [Tổng Quan Hệ Thống / System Architecture Overview](#1-tổng-quan-hệ-thống--system-architecture-overview)
2. [Sơ Đồ Tương Tác / Architectural Interaction Diagram](#2-sơ-đồ-tương-tác--architectural-interaction-diagram)
3. [Phân Tích Chi Tiết Từng Module / Subsystem Deep-Dive](#3-phân-tích-chi-tiết-từng-module--subsystem-deep-dive)
   - [3.1 Entrypoint & Process Lifecycle (`name.py`)](#31-entrypoint--process-lifecycle-namepy)
   - [3.2 Shared State & IPC Singleton (`updater_state.py`)](#32-shared-state--ipc-singleton-updater_statepy)
   - [3.3 Core Updating Engine (`profile_engine.py`)](#33-core-updating-engine-profile_enginepy)
   - [3.4 Interactive Helper Bot (`helper_bot.py`)](#34-interactive-helper-bot-helper_botpy)
   - [3.5 Internationalization Engine (`i18n.py`)](#35-internationalization-engine-i18npy)
   - [3.6 Terminal Configuration Wizard (`setup.py`)](#36-terminal-configuration-wizard-setuppy)
   - [3.7 Channel Media Forwarder (`forward_channel.py`)](#37-channel-media-forwarder-forward_channelpy)
   - [3.8 Session Generator (`generate_session.py`)](#38-session-generator-generate_sessionpy)
4. [Cơ Chế Dự Phòng & Luồng Dữ Liệu / Fallback Hierarchies & Data Flow](#4-cơ-chế-dự-phòng--luồng-dữ-liệu--fallback-hierarchies--data-flow)
   - [4.1 Chuỗi Dự Phòng Thời Tiết / Weather Fallback Chain](#41-chuỗi-dự-phòng-thời-tiết--weather-fallback-chain)
   - [4.2 Ưu Tiên Âm Nhạc vs Thời Tiết / Music vs Weather Priority](#42-ưu-tiên-âm-nhạc-vs-thời-tiết--music-vs-weather-priority)
5. [Ràng Buộc Kỹ Thuật & Giới Hạn Telegram / Technical Constraints & Telegram Limits](#5-ràng-buộc-kỹ-thuật--giới-hạn-telegram--technical-constraints--telegram-limits)
   - [5.1 Giới Hạn Ký Tự Bio (70 chars) / Bio Character Constraint](#51-giới-hạn-ký-tự-bio-70-chars--bio-character-constraint)
   - [5.2 Giới Hạn Tên (64 chars) / First Name Constraint](#52-giới-hạn-tên-64-chars--first-name-constraint)
   - [5.3 Kiểm Soát Tần Suất & FloodWait / Rate Limiting & FloodWait](#53-kiểm-soát-tần-suất--floodwait--rate-limiting--floodwait)
   - [5.4 Đồng Bộ Hóa Event Loop Asyncio / Asyncio Event Loop Binding](#54-đồng-bộ-hóa-event-loop-asyncio--asyncio-event-loop-binding)
6. [Bảng Cấu Hình Biến Môi Trường / Environment Configuration Matrix](#6-bảng-cấu-hình-biến-môi-trường--environment-configuration-matrix)
7. [Vận Hành & Triển Khai / Operations & Deployment](#7-vận-hành--triển-khai--operations--deployment)

---

## 1. Tổng Quan Hệ Thống / System Architecture Overview

### Tiếng Việt
Hệ thống là một giải pháp tự động hóa Telegram toàn diện, kết hợp giữa:
1. **Userbot Engine**: Vận hành dưới danh nghĩa tài khoản cá nhân thông qua giao thức MTProto (Telethon), chịu trách nhiệm cập nhật định kỳ Tên hiển thị (`first_name`) và Tiểu sử (`about` / bio) dựa trên thời gian thực, thời tiết địa phương, và bài hát đang nghe từ Last.fm/Spotify.
2. **Helper Bot**: Vận hành dưới danh nghĩa Bot Telegram chính thức (Bot API Token), đóng vai trò bảng điều khiển quản trị từ xa, hỗ trợ bàn phím inline tương tác, cài đặt đăng nhập OTP/2FA trực tiếp qua chat, xem trước Bio kịch bản và điều khiển vòng đời updater.
3. **Channel Media Forwarder**: Công cụ độc lập xử lý sao chép, chuyển tiếp tập tin, album hình ảnh/video giữa các kênh Telegram ẩn danh, bảo toàn nhóm media và có cơ chế tự động khôi phục tiến trình khi bị gián đoạn.

### English
The system is an end-to-end Telegram profile automation solution comprised of:
1. **Userbot Engine**: Runs as the user's MTProto client via Telethon, orchestrating periodic real-time profile name (`first_name`) and biography (`about`) updates based on accurate local time, dual-source weather conditions, and live Last.fm/Spotify playback scrobbles.
2. **Helper Bot**: Runs as an official Telegram Bot (Bot API Token) providing an interactive inline dashboard, OTP/2FA mobile login wizard over private messages, dynamic bio simulation, and runtime engine management.
3. **Channel Media Forwarder**: An autonomous CLI migration utility designed for forwarding media albums, documents, and videos across Telegram channels while dropping original authors, preserving media groupings, and maintaining crash-resilient checkpoints.

---

## 2. Sơ Đồ Tương Tác / Architectural Interaction Diagram

```
                 +-----------------------------------------------+
                 |                   .env File                   |
                 | (Configuration & Credentials: Hash, Sessions) |
                 +-----------------------+-----------------------+
                                         |
                                         v
                         +-------------------------------+
                         |   updater_state.py (State)    |
                         |  - Shared Global Singleton    |
                         |  - Dynamic asyncio.Event      |
                         |  - Atomic .env update/reload  |
                         +---------------+---------------+
                                         |
                +------------------------+------------------------+
                |                                                 |
                v                                                 v
+-------------------------------+                 +-------------------------------+
|    helper_bot.py (Bot API)    |                 |  profile_engine.py (MTProto)  |
| - TelegramClient (Bot Token)  |                 | - TelegramClient (User Session|
| - Inline Button UI & Callbacks|                 | - Clock & Calendar (pytz)     |
| - Wizard FSM (OTP, 2FA, Setup)|                 | - Weather Engine (OWM + wttr) |
| - Bio Length Simulator        |                 | - Last.fm Scrobbler API       |
| - Admin Security Guard        |                 | - String Truncator (70 chars) |
+---------------+---------------+                 +---------------+---------------+
                |                                                 |
                | trigger_update()                                | UpdateProfileRequest
                +-------------------> [force_update_event] ------>+
                                                                  |
                                                                  v
                                                  +-------------------------------+
                                                  |      Telegram MTProto DC      |
                                                  |   (User Account Name & Bio)   |
                                                  +-------------------------------+
```

---

## 3. Phân Tích Chi Tiết Từng Module / Subsystem Deep-Dive

### 3.1 Entrypoint & Process Lifecycle (`name.py`)

- **Nhiệm vụ / Role**: Điểm khởi động tập trung cho toàn bộ ứng dụng. Tự động kiểm tra cấu hình để phân nhánh kịch bản khởi chạy phù hợp.
- **Cơ chế hoạt động / Execution Flow**:
  1. Nạp và kiểm tra `BOT_TOKEN` từ `state.config`.
  2. **Trường hợp có `BOT_TOKEN`**: Khởi chạy `helper_bot.run_standalone_bot()`. Trong hàm này, Helper Bot sẽ khởi động trước, sau đó tự động kích hoạt `profile_engine.start()` dưới dạng nền (`asyncio.create_task`).
  3. **Trường hợp không có `BOT_TOKEN`**: Khởi chạy chế độ Userbot thuần túy (`engine.start()`). Quá trình này duy trì vòng lặp theo dõi trạng thái `engine.is_running`.

### 3.2 Shared State & IPC Singleton (`updater_state.py`)

- **Nhiệm vụ / Role**: Quản lý trạng thái chia sẻ trong bộ nhớ giữa các tiến trình/task bất đồng bộ và đồng bộ hóa với file `.env`.
- **Cấu trúc & Thuộc tính cốt lõi / Core Properties**:
  - `is_running`, `is_paused`: Trạng thái thực thi của updater engine.
  - `weather_data`, `current_music`, `is_now_playing`: Bộ nhớ cache kết quả truy vấn API mới nhất để tránh spam request.
  - `force_update_event`: Quản lý `asyncio.Event` động.
- **Kỹ thuật đồng bộ Loop động (`force_update_event`)**:
  ```python
  @property
  def force_update_event(self) -> asyncio.Event:
      try:
          current_loop = asyncio.get_running_loop()
      except RuntimeError:
          current_loop = None

      if self._force_update_event is None:
          self._force_update_event = asyncio.Event()
      elif current_loop and getattr(self._force_update_event, "_loop", None) is not current_loop:
          self._force_update_event = asyncio.Event()
      return self._force_update_event
  ```
  *Ý nghĩa*: Tránh ngoại lệ `RuntimeError: Task got Future attached to a different loop` khi một module được import hoặc tái khởi động trên một event loop mới.
- **Ghi cấu hình an toàn (`update_config`)**:
  - Hỗ trợ ký tự đặc biệt (`#`, `$`, quotes, dấu xuống dòng `\n`).
  - Đóng gói giá trị vào dấu nháy kép `""` với cơ chế escape tự động (`\r` loại bỏ, `\` -> `\\`, `"` -> `\"`, `\n` -> `\n`).
  - Duy trì nguyên vẹn các comment `#` sẵn có trong file `.env`.

### 3.3 Core Updating Engine (`profile_engine.py`)

- **Nhiệm vụ / Role**: Lõi xử lý logic sinh chuỗi, truy vấn API bên thứ ba, và giao tiếp MTProto với Telegram servers.
- **Các thành phần xử lý chính**:
  1. `get_weather_icon(desc)`: Ánh xạ mô tả thời tiết (cả tiếng Việt và tiếng Anh) sang Emoji trực quan (ví dụ: `sấm/thunder` -> `⛈️`, `mưa/rain` -> `🌦️`, `mây/cloud` -> `☁️`, `nắng/sun` -> `☀️`).
  2. `get_detailed_weather(city, weather_api_key, timezone, lang)`: Hệ thống truy vấn thời tiết 2 tầng (OpenWeatherMap + fallback wttr.in).
  3. `get_lastfm_track(username, api_key)`: Kết nối `user.getrecenttracks` của Last.fm để lấy bài hát đang nghe (`@attr.nowplaying == "true"`), tổng scrobbles, tên album và nghệ sĩ.
  4. `render_format(template, base_name, weather, city, tz_str, music_info)`: Bộ máy phân giải placeholder & time tags.
  5. `generate_preview()` & `generate_bio_preview()`: Giả lập kết quả hiển thị cho Helper Bot và Userbot command mà không cần đẩy request lên Telegram.
  6. `ProfileEngine.run()`: Vòng lặp chính tính toán khoảng cách đến đầu phút kế tiếp (`sleep_secs = 60 - now.second`), hỗ trợ ngắt tức thì thông qua `state.force_update_event.wait(timeout=sleep_secs)`.

### 3.4 Interactive Helper Bot (`helper_bot.py`)

- **Nhiệm vụ / Role**: Cung cấp giao diện trực quan cho quản trị viên Telegram thông qua Bot API.
- **Tính năng nổi bật**:
  - **Quản trị an toàn (`is_admin`)**: Tự động nhận diện Admin ID từ `.env`. Nếu chưa cấu hình, người đầu tiên gửi tin nhắn tới bot sẽ được tự động gán làm Admin duy nhất.
  - **Wizard Đăng Nhập Điện Thoại**: Quản lý trạng thái đàm thoại (`user_sessions[user_id]`) cho phép gửi số điện thoại, nhận OTP code và mật khẩu 2FA (Two-Factor Authentication / Cloud Password) để sinh `SESSION_STRING` ngay trong chat riêng mà không cần mở terminal.
  - **Inline Dashboard & Menu Phân Cấp**:
    - Menu Chính: Cập nhật ngay, Xem trước Profile, Xem trước Bio, Thời tiết, Start/Stop/Pause/Resume, Cài đặt.
    - Menu Cài đặt: Đổi ngôn ngữ, Đổi tên, Đổi thành phố, Cài đặt Last.fm, Cài đặt Format (Preset mẫu 1 - 8 hoặc nhập tùy biến).
  - **Trình Giả Lập Bio Kịch Bản Kép**: Phân tích độ dài ký tự của mẫu `BIO_FORMAT` khi có nhạc phát và khi không có nhạc phát, cảnh báo trực quan `🟢 Hợp lệ` hoặc `🔴 VƯỢT GIỚI HẠN 70 KÝ TỰ`.
  - **Lắng nghe Userbot Commands**: Cho phép gõ lệnh bắt đầu bằng dấu chấm (`.status`, `.bio`, `.weather`, `.preview`, `.update`, `.name`, v.v.) trực tiếp từ tài khoản user mà không cần tương tác với bot.

### 3.5 Internationalization Engine (`i18n.py`)

- **Nhiệm vụ / Role**: Cung cấp khả năng địa phương hóa song ngữ hoàn chỉnh (Tiếng Việt `vi` và English `en`).
- **Từ điển văn bản (`TEXTS`)**: Bao gồm toàn bộ các nhãn nút bấm, tin nhắn hướng dẫn, tiêu đề thông báo, lỗi và bảng trợ giúp.
- **Xử lý ngày tháng bản địa**:
  - Thứ tiếng Việt (`WEEKDAYS`, `WEEKDAYS_SHORT`): `Thứ Hai` .. `Chủ Nhật`, `T2` .. `CN`.
  - Thứ tiếng Anh: `Monday` .. `Sunday`, `Mon` .. `Sun`.
  - Tháng bản địa: `Tháng 1` .. `Tháng 12` vs `January` .. `December`.
- **Hàm tiện ích**:
  - `get_text(key, lang, **kwargs)`: Trích xuất chuỗi theo ngôn ngữ hiện tại và nội suy chuỗi an toàn (`format(**kwargs)`).
  - `get_weekday(dt, lang, short)`: Trả về thứ tương ứng theo ngôn ngữ.

### 3.6 Terminal Configuration Wizard (`setup.py`)

- **Nhiệm vụ / Role**: Công cụ CLI tương tác đầy màu sắc (ANSI) giúp thiết lập môi trường máy chủ ban đầu.
- **Khả năng cấu hình**:
  - Kiểm tra và tự động cài đặt gói phụ thuộc từ `requirements.txt`.
  - Tạo `SESSION_STRING` qua CLI bằng số điện thoại và mã OTP.
  - Kiểm tra kết nối trực tiếp đến OpenWeatherMap API và Last.fm API trước khi lưu cấu hình.
  - Hỗ trợ tạo file cấu hình tiến trình PM2 hoặc Systemd Service để giữ bot chạy ngầm vĩnh viễn.

### 3.7 Channel Media Forwarder (`forward_channel.py`)

- **Nhiệm vụ / Role**: Script chuyên dụng để di chuyển, sao chép nội dung giữa các Channel Telegram theo thứ tự thời gian.
- **Kiến trúc & Cơ chế kỹ thuật**:
  - **Chuẩn hóa Channel ID (`normalize_channel_id`)**: Tự động chuyển đổi các định dạng ID (`-100...`, `@username`, số nguyên) thành ID chuẩn MTProto.
  - **Forward Ẩn Danh (`drop_author=True`)**: Đảm bảo thông tin tác giả nguồn không bị lộ trên kênh đích.
  - **Bảo Toàn Nhóm Media (`grouped_id`)**:
    - Khi phát hiện một tin nhắn thuộc Album (`msg.grouped_id`), script kích hoạt bộ nhớ đệm gom nhóm toàn bộ các file/ảnh có cùng `grouped_id` và forward dưới dạng media album duy nhất.
  - **Stateful Checkpoint & Resume Engine**:
    - Tự động ghi lại tiến trình quét vào file JSON: `forward_progress_<src>_to_<dst>.json`.
    - Khi bị ngắt bởi `SIGINT` (Ctrl+C) hoặc lỗi mạng, lần chạy tiếp theo sẽ tiếp tục ngay từ `last_scanned_id` mà không cần quét lại từ đầu.
  - **Kiểm Soát Rate Limit Chủ Động**: Cấu hình `delay` tối thiểu giữa các lần forward (mặc định 2.5s) và bẫy lỗi `telethon.errors.FloodWaitError` để tự động `sleep` đúng số giây Telegram yêu cầu.

### 3.8 Session Generator (`generate_session.py`)

- **Nhiệm vụ / Role**: Tiện ích dòng lệnh độc lập, tối giản phục vụ việc đăng nhập MTProto và xuất chuỗi `SESSION_STRING` nạp vào biến môi trường hoặc sao chép lên máy chủ đám mây/container.

---

## 4. Cơ Chế Dự Phòng & Luồng Dữ Liệu / Fallback Hierarchies & Data Flow

### 4.1 Chuỗi Dự Phòng Thời Tiết / Weather Fallback Chain

```
[Trigger Thời Tiết]
       |
       v
+------------------------------------+
| Kiểm tra WEATHER_API & CITY        |
+------------------+-----------------+
                   |
       +-----------+-----------+
       | Có API Key            | Không có API Key / Lỗi mạng
       v                       v
+--------------------+   +---------------------------------------+
|  OpenWeatherMap    |   |                wttr.in                |
|  JSON API (v2.5)   |   |   JSON API (v1 / format=j1)           |
+----------+---------+   | - User-Agent: curl/7.68.0             |
           |             | - Tự động dịch mô tả: vi/en           |
           | Thành công  +-------------------+-------------------+
           v                                 | Thành công
+--------------------+                       v
|  Dữ liệu thời tiết | <---------------------+
|  hoàn chỉnh        |
+----------+---------+
           | Cả 2 đều thất bại
           v
+------------------------------------+
|  Fallback tối giản: "📍 {CITY}"     |
+------------------------------------+
```

### 4.2 Ưu Tiên Âm Nhạc vs Thời Tiết / Music vs Weather Priority

Trong mẫu `BIO_FORMAT`, thẻ `{music_or_weather}` hoạt động theo cơ chế stateful priority:
1. **Trạng thái 1: Đang nghe nhạc (`is_now_playing == True`)**:
   - Dữ liệu Last.fm được ưu tiên tuyệt đối.
   - Định dạng: `🎧 {track} - {artist}`.
   - Thẻ `{weather}` được nhường chỗ hoặc rút gọn để giữ tổng chiều dài Bio <= 70 ký tự.
2. **Trạng thái 2: Không nghe nhạc / Dừng nhạc (`is_now_playing == False`)**:
   - Thẻ `{music_or_weather}` tự động chuyển sang hiển thị chuỗi thời tiết `{weather_medium}` hoặc `{weather_short}`.
   - Không để lại khoảng trống hoặc biểu tượng thừa.

---

## 5. Ràng Buộc Kỹ Thuật & Giới Hạn Telegram / Technical Constraints & Telegram Limits

### 5.1 Giới Hạn Ký Tự Bio (70 chars) / Bio Character Constraint

- **Vấn đề / Problem**: Telegram quy định độ dài tối đa của trường `about` (Bio) tài khoản là **70 ký tự UTF-16 code units**. Vượt quá giới hạn này sẽ gây lỗi `AboutTooLongError` từ MTProto server.
- **Giải pháp xử lý trong code**:
  1. `profile_engine.py`: Áp dụng kỹ thuật slicing an toàn:
     ```python
     p_bio = render_format(bio_format, base_name, weather, city, tz_str, music_info)
     if len(p_bio) > 70:
         p_bio = p_bio[:70]
     ```
  2. `generate_bio_preview()`: Tính toán độ dài trước khi hiển thị cho người dùng trên Helper Bot, xuất cảnh báo `🟢 Hợp lệ (xx/70)` hoặc `🔴 VƯỢT GIỚI HẠN (+xx ký tự)`.

### 5.2 Giới Hạn Tên (64 chars) / First Name Constraint

- Telegram giới hạn trường `first_name` tối đa **64 ký tự**.
- `profile_engine.py` tự động cắt chuỗi nếu độ dài sau khi thay thế template vượt quá 64 ký tự (`p_name = p_name[:64]`).

### 5.3 Kiểm Soát Tần Suất & FloodWait / Rate Limiting & FloodWait

- **Cập nhật Profile định kỳ**: Cố định chu kỳ 1 phút/lần vào đầu mỗi phút (`:00` giây). Telegram cho phép cập nhật profile ở tần suất này mà không kích hoạt FloodWait.
- **Xử lý ngoại lệ**:
  ```python
  except FloodWaitError as e:
      wait_time = e.seconds
      state.last_error = f"FloodWait: phải chờ {wait_time}s"
      await asyncio.sleep(wait_time + 2)
  ```
- **Forwarding Media**: Áp dụng delay cơ sở 2.5s và tự động xử lý khi gặp `FloodWaitError` trong `forward_channel.py`.

### 5.4 Đồng Bộ Hóa Event Loop Asyncio / Asyncio Event Loop Binding

- Khi chạy kết hợp cả Helper Bot và Profile Updater trong cùng một tiến trình Python, việc các task bất đồng bộ tham chiếu đến `asyncio.Event` khởi tạo ở ngoài loop có thể gây lỗi treo tiến trình hoặc `RuntimeError`.
- `updater_state.py` giải quyết triệt để thông qua property `force_update_event` tự động bind lại Event với loop đang chạy.

---

## 6. Bảng Cấu Hình Biến Môi Trường / Environment Configuration Matrix

| Tên biến / Variable | Kiểu dữ liệu | Mặc định / Default | Mô tả chi tiết / Description |
| :--- | :--- | :--- | :--- |
| `API_ID` | `int` | `2040` | Telegram API App ID từ https://my.telegram.org |
| `API_HASH` | `str` | `b18441a1ff...` | Telegram API App Hash tương ứng |
| `SESSION_STRING` | `str` | `""` | Chuỗi Telethon StringSession của tài khoản Userbot |
| `SESSION_NAME` | `str` | `"bot"` | Tên file session sqlite nếu không dùng StringSession |
| `BOT_TOKEN` | `str` | `""` | Telegram Bot Token cấp từ @BotFather cho Helper Bot |
| `ADMIN_ID` | `int` | `""` | ID Telegram của Admin quản lý duy nhất |
| `WEATHER_API` | `str` | `""` | API Key OpenWeatherMap (để trống nếu dùng wttr.in) |
| `CITY` | `str` | `"Bac Ninh"` | Tên thành phố lấy thời tiết mặc định |
| `TIMEZONE` | `str` | `"Asia/Ho_Chi_Minh"` | Múi giờ chuẩn Olson (pytz) |
| `BASE_NAME` | `str` | `"tên"` | Tên hiển thị cố định của bạn |
| `NAME_FORMAT` | `str` | `{base_name} \| HH:mm - DD/MM/YYYY` | Mẫu định dạng First Name Telegram |
| `BIO_FORMAT` | `str` | `{weather} ⏰ HH:mm` | Mẫu định dạng Biography/About Telegram |
| `LASTFM_USERNAME` | `str` | `""` | Tên tài khoản Last.fm để theo dõi nhạc |
| `LASTFM_API_KEY` | `str` | `b25b9595...` | API Key Last.fm tích hợp sẵn |
| `LANGUAGE` | `str` | `"vi"` | Ngôn ngữ hệ thống (`vi` hoặc `en`) |

---

## 7. Vận Hành & Triển Khai / Operations & Deployment

### Khởi chạy trực tiếp (Interactive Execution)
```bash
# Cấu hình qua giao diện Setup Wizard
python3 setup.py

# Khởi chạy hệ thống chính (chạy cả Updater + Helper Bot)
python3 name.py
```

### Vận hành nền qua PM2 (Production Daemon)
```bash
# Khởi chạy và lưu tiến trình
pm2 start name.py --name "tele-updater" --interpreter python3
pm2 save

# Kiểm tra log thời gian thực
pm2 logs tele-updater
```

### Vận hành Channel Forwarder (Auxiliary Tool)
```bash
# Chạy sao chép kênh với tiến trình tự lưu
python3 forward_channel.py --src -1001651870719 --dst -1002638287413 --delay 2.5
```
