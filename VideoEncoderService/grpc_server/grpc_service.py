import transcoding_pb2
import transcoding_pb2_grpc
from services.transcode import encoder, status

class TranscoderServicer(transcoding_pb2_grpc.TranscoderServicer):
    def NotifyUploadComplete(self, request, context):
        vid_uuid = request.uuid
        encoder(vid_uuid)
        return transcoding_pb2.TranscodeResponse(status_code=200)

class VideoStatusServicer(transcoding_pb2_grpc.VideoStatusServiceServicer):
    def StatusVideo(self, request, context):
        vid_uuid = request.uuid
        for i in status(vid_uuid):
            yield transcoding_pb2.VideoStatusResponse(status=i)
