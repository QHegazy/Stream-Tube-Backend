from concurrent import futures
import grpc
from grpc_reflection.v1alpha import reflection
import transcoding_pb2
import transcoding_pb2_grpc
import grpc_service

def serve() -> None:
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    # Add servicers to the server
    transcoding_pb2_grpc.add_TranscoderServicer_to_server(grpc_service.TranscoderServicer(), server)
    transcoding_pb2_grpc.add_VideoStatusServiceServicer_to_server(grpc_service.VideoStatusServicer(), server)
    
    # Enable reflection
    SERVICE_NAMES = (
        transcoding_pb2.DESCRIPTOR.services_by_name['Transcoder'].full_name,
        transcoding_pb2.DESCRIPTOR.services_by_name['VideoStatusService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    
    server.add_insecure_port('[::]:50051')
    server.start()
    
    # Add exception handling
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        print("Server interrupted by user.")
    except Exception as e:
        print(f"Server encountered an error: {e}")
    finally:
        server.stop(0)

if __name__ == '__main__':
    serve()
