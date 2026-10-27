"""
updater_state.py - Quản lý trạng thái chung giữa Profile Updater và Helper Bot.
"""

import os
import asyncio
from dotenv import load_dotenv

ENV_FILE = os.path.join(os.path.dirname(os.path.abspath(__file__)), ".env")

class UpdaterState:
    def __init__(self):
        self.is_running = False
        self.is_paused = False
        self.account_name = ""
        self.account_id = ""
        self.last_update_time = "Chưa cập nhật"
        self.update_count = 0
        self.current_name = ""
        self.current_bio = ""
        self.current_weather = ""
        self.weather_data = {}
        self.weather_short = ""
        self.weather_temp = ""
        self.weather_desc = ""
        self.weather_humidity = ""
        self.weather_wind = ""
        self.weather_feels_like = ""
        self.current_music = ""
        self.is_now_playing = False
        self.music_track = ""
        self.music_artist = ""
        self.music_scrobbles = ""
        self.last_error = ""
        self._force_update_event = None
        self.config = {}
        self.reload_config()

    def reload_config(self):
        """Tải lại biến môi trường từ .env"""
        load_dotenv(ENV_FILE, override=True)
        bio_val = os.getenv("BIO_FORMAT", "{weather} ⏰ HH:mm")
        if bio_val:
            bio_val = bio_val.replace("\\n", "\n")
        name_val = os.getenv("NAME_FORMAT", "{base_name} | HH:mm - DD/MM/YYYY")
        if name_val:
            name_val = name_val.replace("\\n", "\n")
        self.config = {
            "API_ID": os.getenv("API_ID", "2040").strip(),
            "API_HASH": os.getenv("API_HASH", "b18441a1ff607e10a989891a5462e627").strip(),
            "SESSION_STRING": os.getenv("SESSION_STRING", "").strip(),
            "SESSION_NAME": os.getenv("SESSION_NAME", "bot").strip(),
            "BOT_TOKEN": os.getenv("BOT_TOKEN", "").strip(),
            "ADMIN_ID": os.getenv("ADMIN_ID", "").strip(),
            "WEATHER_API": os.getenv("WEATHER_API", "").strip(),
            "CITY": os.getenv("CITY", "Bac Ninh").strip(),
            "LASTFM_API_KEY": os.getenv("LASTFM_API_KEY", "b25b959554ed76058ac220b7b2e0a026").strip(),
            "LASTFM_USERNAME": os.getenv("LASTFM_USERNAME", "").strip(),
            "BASE_NAME": os.getenv("BASE_NAME", "tên").strip(),
            "TIMEZONE": os.getenv("TIMEZONE", "Asia/Ho_Chi_Minh").strip(),
            "NAME_FORMAT": name_val.strip(),
            "BIO_FORMAT": bio_val.strip(),
            "LANGUAGE": os.getenv("LANGUAGE", "vi").strip().lower(),
            "AI_API_KEY": os.getenv("AI_API_KEY", "").strip(),
            "AI_BASE_URL": os.getenv("AI_BASE_URL", "https://api.openai.com/v1").strip(),
            "AI_MODEL": os.getenv("AI_MODEL", "gpt-4o-mini").strip(),
            "AI_SYSTEM_PROMPT": os.getenv("AI_SYSTEM_PROMPT", "You are a helpful and concise Telegram AI assistant.").strip(),
            "GEMINI_API_KEY": os.getenv("GEMINI_API_KEY", "").strip(),
            "IMAGE_API_KEY": os.getenv("IMAGE_API_KEY", "").strip(),
            "IMAGE_BASE_URL": os.getenv("IMAGE_BASE_URL", "").strip(),
            "IMAGE_MODEL": os.getenv("IMAGE_MODEL", "dall-e-3").strip(),
        }

    def update_config(self, updates: dict):
        """Cập nhật cấu hình vào file .env và nạp lại, bảo toàn dấu # và ký tự đặc biệt"""
        lines = []
        if os.path.exists(ENV_FILE):
            with open(ENV_FILE, "r", encoding="utf-8") as f:
                lines = f.readlines()

        def format_env_val(v) -> str:
            if v is None:
                return '""'
            s = str(v).replace("\r", "")
            s = s.replace("\\", "\\\\").replace('"', '\\"').replace("\n", "\\n")
            return f'"{s}"'

        handled = set()
        new_lines = []
        for line in lines:
            stripped = line.strip()
            if stripped and not stripped.startswith("#") and "=" in stripped:
                k, _ = stripped.split("=", 1)
                k = k.strip()
                if k in updates:
                    new_lines.append(f"{k}={format_env_val(updates[k])}\n")
                    handled.add(k)
                    continue
            new_lines.append(line)

        for k, v in updates.items():
            if k not in handled:
                new_lines.append(f"{k}={format_env_val(v)}\n")

        with open(ENV_FILE, "w", encoding="utf-8") as f:
            f.writelines(new_lines)

        self.reload_config()

    @property
    def force_update_event(self) -> asyncio.Event:
        """Đảm bảo Event luôn gắn với Event Loop đang hoạt động"""
        try:
            current_loop = asyncio.get_running_loop()
        except RuntimeError:
            current_loop = None

        if self._force_update_event is None:
            self._force_update_event = asyncio.Event()
        elif current_loop and getattr(self._force_update_event, "_loop", None) is not current_loop:
            self._force_update_event = asyncio.Event()

        return self._force_update_event

    def trigger_update(self):
        """Kích hoạt cập nhật profile ngay lập tức"""
        try:
            self.force_update_event.set()
        except Exception as e:
            print(f"[⚠️] Không thể trigger update: {e}")

state = UpdaterState()
