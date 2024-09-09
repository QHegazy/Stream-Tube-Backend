import subprocess
import os 
import json
from queue import Queue
from threading import Thread
from dotenv import load_dotenv


load_dotenv()


dev_path: str | None = os.getenv('DEV_PATH')
transcode: str | None = os.getenv('TRANSCODE')

progress_queues = {}

class VideoInfo:
    def __init__(self, bitrate_1080, bitrate_720, bitrate_480, bitrate_360) -> None:
        self.bitrate_1080 = bitrate_1080
        self.bitrate_720 = bitrate_720
        self.bitrate_480 = bitrate_480
        self.bitrate_360 = bitrate_360

def get_video_info(file_path) -> VideoInfo:
    cmd = [
        'ffprobe', '-v', 'quiet', '-print_format', 'json',
        '-show_streams', file_path
    ]
    
    result: subprocess.CompletedProcess[bytes] = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if result.returncode != 0:
        raise Exception("Failed to run ffprobe")
    
    data = json.loads(result.stdout)
    format_info = data.get("streams", {})
    format_info = format_info[0]
    
    bitrate_str = format_info.get("bit_rate", None)
    print(bitrate_str)
    if not bitrate_str:
        raise ValueError("Bitrate information not found")
    
    bitrate = int(bitrate_str) // 1000
    
    # duration_str = format_info.get("duration", None)
    # if not duration_str:
    #     raise ValueError("Duration information not found")
    
    
    return VideoInfo(
        bitrate_1080=int(bitrate * 0.45),
        bitrate_720=int(bitrate * 0.29),
        bitrate_480=int(bitrate * 0.1),
        bitrate_360=int(bitrate * 0.08)
    )

def transcode_video(upload_path, output_dir, video_info, progress_queue, uuid):
    cmd = [
        transcode,
        upload_path,
        output_dir,
        str(video_info.bitrate_360),
        str(video_info.bitrate_480),
        str(video_info.bitrate_720),
        str(video_info.bitrate_1080)
    ]
    
    process = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)
    seen_progress = set()
    for line in process.stdout:
        if line.startswith("Overall Progress:"):
            try:
                progress = int(line.split(":")[1].strip().rstrip('%'))
                if progress not in seen_progress:
                    seen_progress.add(progress)
                    progress_queue.put(progress)
            except ValueError:
                print(f"Failed to parse progress line: {line}")
    
    process.wait()
    
    if process.returncode != 0:
        print(f"Transcoding failed for UUID {uuid}. Error: {process.stderr.read()}")
        progress_queue.put(None)  
    else:
        progress_queue.put(None) 
    return

def encoder(uuid):
    upload_path = f"{dev_path}{uuid}"
    if not os.path.exists(upload_path):
        print(f"File not found: {upload_path}")
        return {"error": "File not found"}, 404
    
    try:
        vid_info = get_video_info(upload_path)
        print(vid_info.__dict__)
    except Exception as e:
        return {"error": f"Failed to get video info: {str(e)}"}
    
    output_dir = f"{dev_path}encoded/{uuid}"
    os.makedirs(output_dir, exist_ok=True)
    
    progress_queue = Queue()
    progress_queues[uuid] = progress_queue

    thread = Thread(target=transcode_video, args=(upload_path, output_dir, vid_info, progress_queue, uuid))
    thread.start()
    progress_queue = progress_queues.get(uuid)

    return "OK"


def status(uuid):
    progress_queue = progress_queues.get(uuid)
    while True:
        progress = progress_queue.get()
        if progress is None:
            break
        yield progress
    del progress_queues[uuid]
    
