from flask import Flask, request, jsonify
import os
import re
import uuid
import subprocess
import json
from pathlib import Path
from threading import Thread

app = Flask(__name__)
app.config['MAX_CONTENT_LENGTH'] = 5 * 1024 * 1024 * 1024
app.config['UPLOAD_FOLDER'] = 'encode'

class VideoInfo:
    def __init__(self, bitrate_1080, bitrate_720, bitrate_480, bitrate_360, duration):
        self.bitrate_1080 = bitrate_1080
        self.bitrate_720 = bitrate_720
        self.bitrate_480 = bitrate_480
        self.bitrate_360 = bitrate_360
        self.duration = duration

def get_video_info(file_path):
    cmd = [
        'ffprobe', '-v', 'quiet', '-print_format', 'json',
        '-show_format', '-show_streams', file_path
    ]

    result = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if result.returncode != 0:
        raise Exception("Failed to run ffprobe")

    data = json.loads(result.stdout)
    format_info = data.get("format", {})
    
    bitrate_str = format_info.get("bit_rate", None)
    if not bitrate_str:
        raise ValueError("Bitrate information not found")

    bitrate = int(bitrate_str) // 1000

    duration_str = format_info.get("duration", None)
    if not duration_str:
        raise ValueError("Duration information not found")

    duration = float(duration_str)

    return VideoInfo(
        bitrate_1080=int(bitrate * 0.45),
        bitrate_720=int(bitrate * 0.29),
        bitrate_480=int(bitrate * 0.1),
        bitrate_360=int(bitrate * 0.08),
        duration=duration
    )

def transcode_video(upload_path, output_dir, video_info):
    os.makedirs(output_dir, exist_ok=True)
    cmd = [
        './transcode.sh', upload_path, output_dir,
        str(video_info.bitrate_360),
        str(video_info.bitrate_480),
        str(video_info.bitrate_720),
        str(video_info.bitrate_1080)
    ]

    with open(os.path.join(output_dir, "transcode.log"), "w") as stdout_file, \
         open(os.path.join(output_dir, "transcode_err.log"), "w") as stderr_file:
        subprocess.run(cmd, stdout=stdout_file, stderr=stderr_file)

def get_transcoding_progress(log_file_path):
    last_progress = ""
    with open(log_file_path, 'r') as log_file:
        for line in log_file:
            if line.startswith("out_time_ms="):
                last_progress = line.strip()
    return last_progress

def calculate_percentage(progress, total_duration):
    parts = progress.split('=')
    if len(parts) != 2:
        raise ValueError("Invalid progress line format")

    current_progress = int(parts[1])

    if total_duration == 0:
        raise ValueError("Total duration is zero")

    return (current_progress / total_duration) * 100

@app.route('/upload', methods=['POST'])
def upload_video():
    if 'video' not in request.files or len(request.files.getlist('video')) != 1:
        return "Please upload exactly one file.", 400

    file = request.files['video']
    if not re.search(r'\.(mp4|mkv|flv)$', file.filename):
        return "Invalid file type. Only MP4, MKV, and FLV files are allowed.", 400

    file_uuid = str(uuid.uuid4())
    ext = Path(file.filename).suffix
    upload_path = f"./encode/{file_uuid}{ext}"
    
    os.makedirs(os.path.dirname(upload_path), exist_ok=True)
    file.save(upload_path)

    try:
        video_info = get_video_info(upload_path)
    except Exception as e:
        return f"Failed to get video info: {e}", 500

    output_dir = f"./encode/{file_uuid}"
    thread = Thread(target=transcode_video, args=(upload_path, output_dir, video_info))
    thread.start()

    return jsonify({"uuid": file_uuid})

@app.route('/status/<uuid>', methods=['GET'])
def check_status(uuid):
    log_file = f"./encode/{uuid}/transcode.log"
    try:
        progress = get_transcoding_progress(log_file)
        percentage = calculate_percentage(progress, 39)  # Assume total duration is 39 for example
    except Exception as e:
        return f"Failed to get progress: {e}", 500

    return jsonify({
        "progress": progress,
        "percentage": percentage
    })

if __name__ == '__main__':
    app.run(port=3000)
