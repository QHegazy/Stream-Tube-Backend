#!/bin/bash

set -e

input_file=$1
output_dir=$2
bitrate_360=$3
bitrate_480=$4
bitrate_720=$5
bitrate_1080=$6

# Check if input file exists
if [[ ! -f "$input_file" ]]; then
  echo "Error: Input file '$input_file' does not exist."
  exit 1
fi

# Create subdirectories for resolutions
mkdir -p "$output_dir/360p"
mkdir -p "$output_dir/480p"
mkdir -p "$output_dir/720p"
mkdir -p "$output_dir/1080p"

# Common settings
audio_codec="aac"
audio_bitrate="64k"
audio_sample_rate="44100"
video_codec="libx264"
video_profile="main"
hls_time="4"
playlist_type="vod"
keyint_min="120"
g_value="120"
fps="24"
preset="veryslow"
fmt="yuv420p"

# Encode video function
encode_video() {
    local height="$1"
    local bitrate="$2"
    local maxrate=$((bitrate + 500))
    local bufsize=$((bitrate * 2))

    ffmpeg -i "$input_file" \
        -vf "scale=-2:$height,fps=$fps" \
        -c:a "$audio_codec" -ar "$audio_sample_rate" -b:a "$audio_bitrate" \
        -pix_fmt "$fmt" -c:v "$video_codec" \
        -profile:v "$video_profile" -preset "$preset" -tag:v hvc1 \
        -g "$g_value" -keyint_min "$keyint_min"  \
        -hls_time "$hls_time" -hls_playlist_type "$playlist_type" \
        -b:v "${bitrate}k" -maxrate "${maxrate}k" -bufsize "${bufsize}k" \
        -hls_segment_filename "$output_dir/${height}p/${height}p_%03d.ts" \
        -movflags +faststart \
        -progress - \
        "$output_dir/${height}p/${height}p.m3u8"
}

# Encode videos
encode_video 360 "$bitrate_360"
encode_video 480 "$bitrate_480"
encode_video 720 "$bitrate_720"
encode_video 1080 "$bitrate_1080"

# Create master playlist
cat > "$output_dir/master.m3u8" << EOF
#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=$((bitrate_360 * 1000))K,RESOLUTION=640x360
360p/360p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=$((bitrate_480 * 1000))K,RESOLUTION=854x480
480p/480p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=$((bitrate_720 * 1000))K,RESOLUTION=1280x720
720p/720p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=$((bitrate_1080 * 1000))K,RESOLUTION=1920x1080
1080p/1080p.m3u8
EOF

echo "HLS encoding completed. Master playlist created at $output_dir/master.m3u8"
