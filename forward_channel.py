#!/usr/bin/env python3
"""
forward_channel.py - Tool forward nội dung file và video giữa 2 kênh Telegram
- Sử dụng SESSION_STRING từ file .env
- Không để lại tên người gửi (Forward without author / drop_author=True)
- Sắp xếp theo thứ tự thời gian (từ cũ đến mới: reverse=True)
- Hỗ trợ giữ nguyên album/media group
- Có delay chống spam API Telegram & tự động xử lý FloodWait
- Hỗ trợ lưu tiến trình (Resume) khi bị gián đoạn
"""

import os
import sys
import json
import time
import signal
import asyncio
import argparse
import logging
from datetime import datetime
from typing import List, Optional

from dotenv import load_dotenv
from telethon import TelegramClient, errors, types
from telethon.sessions import StringSession

# Thiết lập logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    datefmt="%H:%M:%S",
)
logger = logging.getLogger("forwarder")

# Tải biến môi trường từ .env
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
load_dotenv(os.path.join(BASE_DIR, ".env"))


def normalize_channel_id(channel_input) -> int:
    """Chuyển đổi channel ID dạng chuỗi hoặc số thành ID kênh hợp lệ (-100...)"""
    s = str(channel_input).strip()
    if s.startswith("-100"):
        return int(s)
    elif s.startswith("-"):
        return int(s)
    elif s.isdigit():
        return int(f"-100{s}")
    else:
        try:
            return int(s)
        except ValueError:
            raise ValueError(f"ID kênh không hợp lệ: {channel_input}")


class ChannelForwarder:
    def __init__(
        self,
        src_channel_id: int,
        dst_channel_id: int,
        delay: float = 2.5,
        mode: str = "files_and_videos",
        silent: bool = True,
        resume: bool = True,
        from_id: int = 0,
        to_id: int = 0,
        limit: int = 0,
        dry_run: bool = False,
    ):
        self.src_channel_id = src_channel_id
        self.dst_channel_id = dst_channel_id
        self.delay = max(0.5, delay)
        self.mode = mode
        self.silent = silent
        self.resume = resume
        self.from_id = from_id
        self.to_id = to_id
        self.limit = limit
        self.dry_run = dry_run

        self.progress_file = os.path.join(
            BASE_DIR,
            f"forward_progress_{abs(self.src_channel_id)}_to_{abs(self.dst_channel_id)}.json",
        )

        self.client: Optional[TelegramClient] = None
        self.src_entity = None
        self.dst_entity = None

        self.last_scanned_id = 0
        self.total_forwarded = 0
        self.total_skipped = 0
        self.is_interrupted = False

    def load_progress(self):
        """Tải tiến trình đã lưu trước đó nếu có"""
        if self.resume and os.path.exists(self.progress_file):
            try:
                with open(self.progress_file, "r", encoding="utf-8") as f:
                    data = json.load(f)
                    self.last_scanned_id = data.get("last_scanned_id", 0)
                    self.total_forwarded = data.get("total_forwarded", 0)
                    logger.info(
                        f"🔄 Đã tìm thấy tiến trình trước đó: Đã quét tới ID {self.last_scanned_id} "
                        f"(Đã forward: {self.total_forwarded} mục). Sẽ tiếp tục từ điểm này."
                    )
            except Exception as e:
                logger.warning(f"Không thể đọc file tiến trình: {e}. Bắt đầu lại từ đầu.")

    def save_progress(self):
        """Lưu tiến trình hiện tại vào file JSON"""
        if self.dry_run:
            return
        data = {
            "src_channel_id": self.src_channel_id,
            "dst_channel_id": self.dst_channel_id,
            "last_scanned_id": self.last_scanned_id,
            "total_forwarded": self.total_forwarded,
            "updated_at": datetime.now().isoformat(),
        }
        try:
            with open(self.progress_file, "w", encoding="utf-8") as f:
                json.dump(data, f, indent=2, ensure_ascii=False)
        except Exception as e:
            logger.error(f"Lỗi khi lưu tiến trình: {e}")

    def is_target_message(self, msg: types.Message) -> bool:
        """Kiểm tra tin nhắn có thuộc loại cần forward hay không"""
        if not msg or not msg.media:
            return False

        has_video = msg.video is not None
        has_document = msg.document is not None
        has_photo = msg.photo is not None

        if self.mode == "strict_files_only":
            return has_video or has_document
        elif self.mode == "all":
            return has_video or has_document or has_photo
        else:
            # files_and_videos (mặc định)
            return has_video or has_document

    async def init_client(self):
        """Khởi tạo Telethon Client từ SESSION_STRING trong .env"""
        api_id_str = os.getenv("API_ID", "2040")
        api_hash = os.getenv("API_HASH", "b18441a1ff607e10a989891a5462e627")
        session_string = os.getenv("SESSION_STRING", "").strip()

        if not session_string:
            raise ValueError(
                "❌ Không tìm thấy SESSION_STRING trong file .env! "
                "Vui lòng chạy 'python3 generate_session.py' để tạo phiên đăng nhập."
            )

        api_id = int(api_id_str)
        self.client = TelegramClient(StringSession(session_string), api_id, api_hash)
        await self.client.connect()

        if not await self.client.is_user_authorized():
            await self.client.disconnect()
            raise ValueError(
                "❌ SESSION_STRING không hợp lệ hoặc đã hết hạn! Vui lòng tạo lại session string."
            )

        me = await self.client.get_me()
        logger.info(f"✅ Đăng nhập thành công với tài khoản: {me.first_name} (ID: {me.id})")

        # Xác thực 2 kênh
        try:
            self.src_entity = await self.client.get_entity(self.src_channel_id)
            logger.info(
                f"📥 Kênh nguồn: {getattr(self.src_entity, 'title', self.src_channel_id)} "
                f"(ID: {self.src_channel_id})"
            )
        except Exception as e:
            raise ValueError(f"Không thể truy cập kênh nguồn {self.src_channel_id}: {e}")

        try:
            self.dst_entity = await self.client.get_entity(self.dst_channel_id)
            logger.info(
                f"📤 Kênh đích:  {getattr(self.dst_entity, 'title', self.dst_channel_id)} "
                f"(ID: {self.dst_channel_id})"
            )
        except Exception as e:
            raise ValueError(f"Không thể truy cập kênh đích {self.dst_channel_id}: {e}")

    async def forward_batch(self, messages: List[types.Message]) -> bool:
        """Forward một danh sách tin nhắn (đơn lẻ hoặc theo album) không để lại tên người gửi"""
        if not messages:
            return False

        msg_ids = [m.id for m in messages]
        info_parts = []
        for m in messages:
            m_type = "video" if m.video else ("file" if m.document else "photo")
            info_parts.append(f"#{m.id} ({m_type})")
        info_str = ", ".join(info_parts)

        if self.dry_run:
            logger.info(f"[DRY-RUN] Sẽ forward mà không lưu tên tác giả: [{info_str}]")
            self.total_forwarded += len(messages)
            return True

        retries = 3
        while retries > 0:
            try:
                # drop_author=True: Không để lại tên người gửi ("forward without name")
                sent = await self.client.forward_messages(
                    entity=self.dst_entity,
                    messages=messages,
                    drop_author=True,
                    silent=self.silent,
                )
                self.total_forwarded += len(messages)
                logger.info(
                    f"✨ [ĐÃ FORWARD] [{info_str}] -> Hoàn thành (Tổng đã chuyển: {self.total_forwarded})"
                )
                return True
            except errors.FloodWaitError as fwe:
                wait_time = fwe.seconds + 2
                logger.warning(
                    f"⚠️ Bị giới hạn tần suất Telegram (FloodWait). Chờ {wait_time}s trước khi thử lại..."
                )
                await asyncio.sleep(wait_time)
                retries -= 1
            except Exception as e:
                logger.error(f"❌ Lỗi khi forward [{info_str}]: {e}")
                retries -= 1
                if retries > 0:
                    await asyncio.sleep(3)
                else:
                    return False
        return False

    async def run(self):
        """Vòng lặp chính xử lý đọc và forward theo thứ tự"""
        self.load_progress()

        min_id = self.from_id if self.from_id > 0 else self.last_scanned_id
        logger.info("=" * 60)
        logger.info("🚀 BẮT ĐẦU FORWARD FILE & VIDEO")
        logger.info(f"   Kênh nguồn:       {self.src_channel_id}")
        logger.info(f"   Kênh đích:        {self.dst_channel_id}")
        logger.info(f"   Chế độ lọc:       {self.mode}")
        logger.info(f"   Bắt đầu từ ID:    > {min_id}")
        logger.info(f"   Kết thúc tại ID:  {self.to_id if self.to_id else 'Đến tin nhắn mới nhất'}")
        logger.info(f"   Delay giữa các lần: {self.delay}s")
        logger.info(f"   Ẩn người gửi:     Có (drop_author=True)")
        logger.info(f"   Chế độ im lặng:   {'Bật' if self.silent else 'Tắt'}")
        if self.limit > 0:
            logger.info(f"   Giới hạn số lượng: {self.limit} mục")
        if self.dry_run:
            logger.info("   ⚠️  CHẾ ĐỘ THỬ NGHIỆM (DRY-RUN) - KHÔNG GỬI THẬT")
        logger.info("=" * 60)

        current_album: List[types.Message] = []
        current_album_id = None
        count_handled = 0

        async def flush_album():
            nonlocal current_album, current_album_id, count_handled
            if not current_album:
                return

            # Kiểm tra xem album này có chứa file/video cần forward không
            has_target_content = any(self.is_target_message(m) for m in current_album)
            if has_target_content:
                if self.mode == "strict_files_only":
                    items_to_send = [m for m in current_album if self.is_target_message(m)]
                else:
                    # Giữ nguyên album bao gồm ảnh bìa và caption nếu album chứa file/video
                    items_to_send = current_album

                if items_to_send:
                    ok = await self.forward_batch(items_to_send)
                    if ok:
                        count_handled += 1
                        await asyncio.sleep(self.delay)
            else:
                self.total_skipped += len(current_album)

            # Cập nhật ID quét lớn nhất trong album
            max_id_in_album = max(m.id for m in current_album)
            if max_id_in_album > self.last_scanned_id:
                self.last_scanned_id = max_id_in_album
                self.save_progress()

            current_album = []
            current_album_id = None

        try:
            # iter_messages với reverse=True: duyệt tin nhắn theo thứ tự từ CŨ đến MỚI
            async for msg in self.client.iter_messages(
                self.src_entity,
                reverse=True,
                min_id=min_id,
            ):
                if self.is_interrupted:
                    break

                if self.to_id and msg.id > self.to_id:
                    logger.info(f"Đã chạm tới ID giới hạn ({self.to_id}). Dừng.")
                    break

                # Xử lý gom nhóm album (nếu có grouped_id)
                if msg.grouped_id is not None:
                    if current_album_id is None:
                        current_album_id = msg.grouped_id
                        current_album.append(msg)
                    elif current_album_id == msg.grouped_id:
                        current_album.append(msg)
                    else:
                        # Kết thúc album cũ, bắt đầu album mới
                        await flush_album()
                        if self.limit > 0 and count_handled >= self.limit:
                            logger.info(f"Đã đạt giới hạn số lượng mục ({self.limit}).")
                            break
                        current_album_id = msg.grouped_id
                        current_album.append(msg)
                else:
                    # Không thuộc album
                    await flush_album()
                    if self.limit > 0 and count_handled >= self.limit:
                        logger.info(f"Đã đạt giới hạn số lượng mục ({self.limit}).")
                        break

                    if self.is_target_message(msg):
                        ok = await self.forward_batch([msg])
                        if ok:
                            count_handled += 1
                            await asyncio.sleep(self.delay)
                    else:
                        self.total_skipped += 1

                    if msg.id > self.last_scanned_id:
                        self.last_scanned_id = msg.id
                        self.save_progress()

                if self.limit > 0 and count_handled >= self.limit:
                    logger.info(f"Đã đạt giới hạn số lượng mục ({self.limit}).")
                    break

            # Flush album cuối cùng nếu còn
            await flush_album()

        except asyncio.CancelledError:
            logger.info("Tiến trình bị hủy.")
        except Exception as e:
            logger.error(f"Lỗi không mong muốn trong quá trình duyệt tin nhắn: {e}", exc_info=True)
        finally:
            self.save_progress()
            logger.info("=" * 60)
            logger.info("🏁 HOÀN THÀNH / TẠM DỪNG TIẾN TRÌNH")
            logger.info(f"   Tổng số tin nhắn/file đã forward: {self.total_forwarded}")
            logger.info(f"   Tin nhắn đã quét tới ID:          {self.last_scanned_id}")
            logger.info(f"   File tiến trình:                 {self.progress_file}")
            logger.info("=" * 60)


def setup_signal_handlers(forwarder: ChannelForwarder):
    """Bắt phím tắt Ctrl+C để lưu tiến trình an toàn trước khi tắt"""
    def handle_stop(signum, frame):
        logger.info("\n🛑 Đang dừng an toàn... Vui lòng chờ lưu tiến trình.")
        forwarder.is_interrupted = True

    try:
        signal.signal(signal.SIGINT, handle_stop)
        signal.signal(signal.SIGTERM, handle_stop)
    except Exception:
        pass


async def main():
    parser = argparse.ArgumentParser(
        description="Forward nội dung file & video giữa 2 kênh Telegram không để lại tên người gửi."
    )
    parser.add_argument(
        "--src",
        default="1651870719",
        help="ID hoặc username kênh nguồn (mặc định: 1651870719)",
    )
    parser.add_argument(
        "--dst",
        default="2638287413",
        help="ID hoặc username kênh đích (mặc định: 2638287413)",
    )
    parser.add_argument(
        "--delay",
        type=float,
        default=2.5,
        help="Thời gian nghỉ giữa các lần gửi (giây, mặc định: 2.5s)",
    )
    parser.add_argument(
        "--mode",
        choices=["files_and_videos", "strict_files_only", "all"],
        default="files_and_videos",
        help="Chế độ lọc nội dung: 'files_and_videos' (mặc định), 'strict_files_only', hoặc 'all'",
    )
    parser.add_argument(
        "--no-silent",
        action="store_true",
        help="Gửi tin nhắn kèm chuông thông báo (mặc định là gửi im lặng)",
    )
    parser.add_argument(
        "--fresh",
        action="store_true",
        help="Bỏ qua file tiến trình cũ, bắt đầu lại từ đầu",
    )
    parser.add_argument(
        "--from-id",
        type=int,
        default=0,
        help="Bắt đầu forward từ ID tin nhắn này trở đi",
    )
    parser.add_argument(
        "--to-id",
        type=int,
        default=0,
        help="Dừng lại khi chạm tới ID tin nhắn này",
    )
    parser.add_argument(
        "--limit",
        type=int,
        default=0,
        help="Giới hạn số lượng bài viết/album cần forward (0 là không giới hạn)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Chạy thử nghiệm mô phỏng, không gửi thật vào kênh đích",
    )

    args = parser.parse_args()

    src_id = normalize_channel_id(args.src)
    dst_id = normalize_channel_id(args.dst)

    forwarder = ChannelForwarder(
        src_channel_id=src_id,
        dst_channel_id=dst_id,
        delay=args.delay,
        mode=args.mode,
        silent=not args.no_silent,
        resume=not args.fresh,
        from_id=args.from_id,
        to_id=args.to_id,
        limit=args.limit,
        dry_run=args.dry_run,
    )

    setup_signal_handlers(forwarder)

    try:
        await forwarder.init_client()
        await forwarder.run()
    except Exception as e:
        logger.error(f"Lỗi: {e}")
    finally:
        if forwarder.client and forwarder.client.is_connected():
            await forwarder.client.disconnect()


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except (KeyboardInterrupt, SystemExit):
        logger.info("Đã thoát chương trình.")
