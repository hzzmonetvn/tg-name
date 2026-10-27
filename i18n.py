"""
i18n.py - Hỗ trợ đa ngôn ngữ (Tiếng Việt & English) cho Telegram Profile Updater & Helper Bot.
"""

from typing import Dict, Any

LANGUAGES = {
    "vi": "Tiếng Việt 🇻🇳",
    "en": "English 🇬🇧"
}

WEEKDAYS = {
    "vi": ["Thứ Hai", "Thứ Ba", "Thứ Tư", "Thứ Năm", "Thứ Sáu", "Thứ Bảy", "Chủ Nhật"],
    "en": ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]
}

WEEKDAYS_SHORT = {
    "vi": ["T2", "T3", "T4", "T5", "T6", "T7", "CN"],
    "en": ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]
}

MONTHS = {
    "vi": ["Tháng 1", "Tháng 2", "Tháng 3", "Tháng 4", "Tháng 5", "Tháng 6", "Tháng 7", "Tháng 8", "Tháng 9", "Tháng 10", "Tháng 11", "Tháng 12"],
    "en": ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]
}

TEXTS: Dict[str, Dict[str, str]] = {
    # Nút bấm Menu chính
    "btn_update_now": {
        "vi": "🔄 Cập nhật ngay",
        "en": "🔄 Update Now"
    },
    "btn_preview_profile": {
        "vi": "👁️ Xem trước Profile",
        "en": "👁️ Preview Profile"
    },
    "btn_preview_bio": {
        "vi": "📝 Xem trước Bio",
        "en": "📝 Preview Bio"
    },
    "btn_weather_card": {
        "vi": "🌤️ Thời tiết chi tiết",
        "en": "🌤️ Weather Details"
    },
    "btn_stop_updater": {
        "vi": "⏹️ Dừng Updater",
        "en": "⏹️ Stop Updater"
    },
    "btn_start_updater": {
        "vi": "▶️ Chạy Updater",
        "en": "▶️ Start Updater"
    },
    "btn_pause": {
        "vi": "⏸️ Tạm dừng",
        "en": "⏸️ Pause"
    },
    "btn_resume": {
        "vi": "▶️ Tiếp tục",
        "en": "▶️ Resume"
    },
    "btn_settings": {
        "vi": "⚙️ Cài đặt & Cấu hình",
        "en": "⚙️ Settings & Setup"
    },
    "btn_status": {
        "vi": "📊 Làm mới trạng thái",
        "en": "📊 Refresh Status"
    },
    "btn_help": {
        "vi": "ℹ️ Trợ giúp",
        "en": "ℹ️ Help & Guide"
    },
    "btn_cancel": {
        "vi": "❌ Hủy bỏ thao tác",
        "en": "❌ Cancel Action"
    },
    "btn_back_main": {
        "vi": "⬅️ Quay lại Menu chính",
        "en": "⬅️ Back to Main Menu"
    },
    "btn_back": {
        "vi": "⬅️ Quay lại",
        "en": "⬅️ Back"
    },
    "btn_language": {
        "vi": "🌐 Ngôn ngữ / Language",
        "en": "🌐 Language / Ngôn ngữ"
    },

    # Nút bấm Menu Cài đặt
    "btn_change_account": {
        "vi": "🔑 Đổi tài khoản Telegram",
        "en": "🔑 Change Telegram Account"
    },
    "btn_login_account": {
        "vi": "🔑 Đăng nhập Telegram",
        "en": "🔑 Login Telegram"
    },
    "btn_set_name": {
        "vi": "✏️ Đổi Tên (BASE_NAME)",
        "en": "✏️ Change Name (BASE_NAME)"
    },
    "btn_set_city": {
        "vi": "📍 Đổi Thành phố",
        "en": "📍 Change City"
    },
    "btn_set_weather": {
        "vi": "⛅ Cấu hình Thời tiết",
        "en": "⛅ Weather Settings"
    },
    "btn_set_timezone": {
        "vi": "🕒 Đổi Múi giờ",
        "en": "🕒 Change Timezone"
    },
    "btn_set_api_creds": {
        "vi": "🆔 Cấu hình API ID & Hash",
        "en": "🆔 API ID & Hash Config"
    },
    "btn_set_format": {
        "vi": "🎨 Mẫu hiển thị (Format)",
        "en": "🎨 Display Formats"
    },
    "btn_custom_name_fmt": {
        "vi": "✏️ Tự nhập NAME_FORMAT",
        "en": "✏️ Custom NAME_FORMAT"
    },
    "btn_custom_bio_fmt": {
        "vi": "📝 Tự nhập BIO_FORMAT",
        "en": "📝 Custom BIO_FORMAT"
    },
    "btn_reset_fmt": {
        "vi": "🔄 Khôi phục mặc định",
        "en": "🔄 Reset to Defaults"
    },
    "btn_preview_bio_detail": {
        "vi": "📝 Xem trước Bio chi tiết",
        "en": "📝 Detailed Bio Preview"
    },

    # Trạng thái Dashboard
    "status_stopped": {
        "vi": "🔴 <b>Đang dừng</b>",
        "en": "🔴 <b>Stopped</b>"
    },
    "status_paused": {
        "vi": "⏸️ <b>Đang tạm dừng</b>",
        "en": "⏸️ <b>Paused</b>"
    },
    "status_running": {
        "vi": "🟢 <b>Đang chạy bình thường</b>",
        "en": "🟢 <b>Running Normally</b>"
    },
    "status_disconnected": {
        "vi": "🟡 <b>Mất kết nối Telegram (Đang kết nối lại...)</b>",
        "en": "🟡 <b>Disconnected from Telegram (Reconnecting...)</b>"
    },
    "not_logged_in": {
        "vi": "<i>Chưa đăng nhập</i>",
        "en": "<i>Not logged in</i>"
    },
    "has_session": {
        "vi": "✔ Đã có SESSION_STRING",
        "en": "✔ SESSION_STRING configured"
    },
    "no_session": {
        "vi": "✘ Chưa đăng nhập",
        "en": "✘ Not logged in"
    },
    "not_updated": {
        "vi": "Chưa cập nhật",
        "en": "Not updated yet"
    },
    "no_info": {
        "vi": "Chưa có thông tin",
        "en": "No info available"
    },

    # Tiêu đề & Nội dung Dashboard
    "dashboard_title": {
        "vi": "🤖 <b>BẢNG ĐIỀU KHIỂN TELEGRAM PROFILE UPDATER</b>",
        "en": "🤖 <b>TELEGRAM PROFILE UPDATER DASHBOARD</b>"
    },
    "field_status": {
        "vi": "• <b>Tình trạng:</b>",
        "en": "• <b>Status:</b>"
    },
    "field_account": {
        "vi": "• <b>Tài khoản Userbot:</b>",
        "en": "• <b>Userbot Account:</b>"
    },
    "field_current_name": {
        "vi": "• <b>Tên hiển thị hiện tại:</b>",
        "en": "• <b>Current Display Name:</b>"
    },
    "field_current_bio": {
        "vi": "• <b>Bio hiện tại:</b>",
        "en": "• <b>Current Bio:</b>"
    },
    "field_music": {
        "vi": "• <b>Đang nghe (Last.fm):</b>",
        "en": "• <b>Now Playing (Last.fm):</b>"
    },
    "field_weather": {
        "vi": "• <b>Thời tiết:</b>",
        "en": "• <b>Weather:</b>"
    },
    "field_base_name": {
        "vi": "• <b>Tên cố định (BASE_NAME):</b>",
        "en": "• <b>Fixed Name (BASE_NAME):</b>"
    },
    "field_city": {
        "vi": "• <b>Thành phố:</b>",
        "en": "• <b>City:</b>"
    },
    "field_timezone": {
        "vi": "• <b>Múi giờ:</b>",
        "en": "• <b>Timezone:</b>"
    },
    "field_last_update": {
        "vi": "• <b>Lần cập nhật cuối:</b>",
        "en": "• <b>Last Update:</b>"
    },
    "field_total_updates": {
        "vi": "Tổng",
        "en": "Total"
    },
    "field_last_error": {
        "vi": "• <b>Lỗi gần nhất:</b>",
        "en": "• <b>Last Error:</b>"
    },
    "field_language": {
        "vi": "• <b>Ngôn ngữ:</b>",
        "en": "• <b>Language:</b>"
    },

    # Thẻ thời tiết
    "weather_card_title": {
        "vi": "🌤️ <b>THÔNG TIN THỜI TIẾT CHI TIẾT — {city}</b>",
        "en": "🌤️ <b>DETAILED WEATHER INFORMATION — {city}</b>"
    },
    "weather_condition": {
        "vi": "• <b>Tình trạng:</b> {icon} <b>{desc}</b>",
        "en": "• <b>Condition:</b> {icon} <b>{desc}</b>"
    },
    "weather_temp": {
        "vi": "• <b>Nhiệt độ hiện tại:</b> <code>{temp}</code> (Cảm nhận thực tế: <code>{feels}</code>)",
        "en": "• <b>Current Temp:</b> <code>{temp}</code> (Feels like: <code>{feels}</code>)"
    },
    "weather_humidity": {
        "vi": "• <b>Độ ẩm không khí:</b> 💦 <code>{hum}</code>",
        "en": "• <b>Air Humidity:</b> 💦 <code>{hum}</code>"
    },
    "weather_wind": {
        "vi": "• <b>Tốc độ gió:</b> 💨 <code>{wind}</code>",
        "en": "• <b>Wind Speed:</b> 💨 <code>{wind}</code>"
    },
    "weather_pressure": {
        "vi": "• <b>Áp suất khí quyển:</b> <code>{pressure}</code>",
        "en": "• <b>Pressure:</b> <code>{pressure}</code>"
    },
    "weather_sunrise": {
        "vi": "• <b>Mặt trời mọc (Bình minh):</b> 🌅 <code>{sunrise}</code>",
        "en": "• <b>Sunrise:</b> 🌅 <code>{sunrise}</code>"
    },
    "weather_sunset": {
        "vi": "• <b>Mặt trời lặn (Hoàng hôn):</b> 🌇 <code>{sunset}</code>",
        "en": "• <b>Sunset:</b> 🌇 <code>{sunset}</code>"
    },
    "weather_source": {
        "vi": "• <b>Nguồn cấp dữ liệu:</b> <i>{source}</i>",
        "en": "• <b>Data Source:</b> <i>{source}</i>"
    },
    "weather_tags_title": {
        "vi": "📌 <b>Các từ khóa thời tiết bạn có thể dùng trong Tên & Bio:</b>",
        "en": "📌 <b>Weather placeholders you can use in Name & Bio:</b>"
    },

    # Thẻ xem trước Bio
    "bio_preview_title": {
        "vi": "📝 <b>XEM TRƯỚC TIỂU SỬ TELEGRAM (BIO PREVIEW)</b>",
        "en": "📝 <b>TELEGRAM BIO PREVIEW & ANALYSIS</b>"
    },
    "bio_template_label": {
        "vi": "• <b>Mẫu BIO_FORMAT đang cấu hình:</b>",
        "en": "• <b>Configured BIO_FORMAT template:</b>"
    },
    "bio_actual_label": {
        "vi": "• <b>Bio hiển thị thực tế trên Profile:</b>",
        "en": "• <b>Actual display on Profile:</b>"
    },
    "bio_length_label": {
        "vi": "📊 <b>Độ dài:</b> <code>{len}/70 ký tự</code> ({status})",
        "en": "📊 <b>Length:</b> <code>{len}/70 chars</code> ({status})"
    },
    "bio_valid": {
        "vi": "🟢 Hợp lệ",
        "en": "🟢 Valid"
    },
    "bio_invalid": {
        "vi": "🔴 VƯỢT GIỚI HẠN 70 KÝ TỰ",
        "en": "🔴 EXCEEDS 70 CHARS LIMIT"
    },
    "bio_scenarios_title": {
        "vi": "🎭 <b>CÁC KỊCH BẢN HIỂN THỊ CỦA MẪU NÀY:</b>",
        "en": "🎭 <b>DISPLAY SCENARIOS FOR THIS TEMPLATE:</b>"
    },
    "bio_scenario_music": {
        "vi": "🎧 <b>1. Khi đang phát nhạc (Spotify / Apple Music):</b>",
        "en": "🎧 <b>1. When actively playing music (Spotify / Apple Music):</b>"
    },
    "bio_scenario_no_music": {
        "vi": "🌤️ <b>2. Khi không nghe nhạc (hiển thị thời tiết):</b>",
        "en": "🌤️ <b>2. When no music is playing (fallback to weather):</b>"
    },
    "bio_tip": {
        "vi": "💡 <i>Gợi ý: Dùng lệnh <code>/bioformat &lt;mẫu&gt;</code> để đổi ngay định dạng Bio!</i>",
        "en": "💡 <i>Tip: Use <code>/bioformat &lt;template&gt;</code> to update your Bio format instantly!</i>"
    },

    # Cài đặt
    "settings_title": {
        "vi": "⚙️ <b>CÀI ĐẶT & CẤU HÌNH HỆ THỐNG</b>",
        "en": "⚙️ <b>SYSTEM CONFIGURATION & SETTINGS</b>"
    },
    "settings_prompt": {
        "vi": "<i>Bấm vào một mục bên dưới để cấu hình:</i>",
        "en": "<i>Tap an option below to configure:</i>"
    },
    "lang_switched": {
        "vi": "✔ Đã chuyển ngôn ngữ sang: <b>Tiếng Việt 🇻🇳</b>",
        "en": "✔ Language successfully switched to: <b>English 🇬🇧</b>"
    },
    "choose_language": {
        "vi": "🌐 <b>CHỌN NGÔN NGỮ HIỂN THỊ (SELECT LANGUAGE)</b>\n\nVui lòng chọn ngôn ngữ bạn muốn sử dụng:",
        "en": "🌐 <b>SELECT DISPLAY LANGUAGE (CHỌN NGÔN NGỮ)</b>\n\nPlease select your preferred language:"
    }
}

def get_text(key: str, lang: str = "vi", **kwargs) -> str:
    """Lấy văn bản dịch theo ngôn ngữ đã chọn"""
    l = lang if lang in ("vi", "en") else "vi"
    tmpl = TEXTS.get(key, {}).get(l) or TEXTS.get(key, {}).get("vi") or key
    if kwargs:
        try:
            return tmpl.format(**kwargs)
        except Exception:
            return tmpl
    return tmpl

def get_weekday(idx: int, lang: str = "vi", short: bool = False) -> str:
    """Lấy tên thứ trong tuần (0: Thứ Hai / Monday ... 6: Chủ Nhật / Sunday)"""
    l = lang if lang in ("vi", "en") else "vi"
    table = WEEKDAYS_SHORT[l] if short else WEEKDAYS[l]
    return table[idx % 7]

def get_month_name(month_idx: int, lang: str = "vi") -> str:
    """Lấy tên tháng (1-12)"""
    l = lang if lang in ("vi", "en") else "vi"
    return MONTHS[l][(month_idx - 1) % 12]
