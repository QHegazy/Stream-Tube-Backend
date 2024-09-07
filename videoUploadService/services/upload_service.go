package uploadSerivce

import (
	pbt "VideoUploadService/transcoding"
	pb "VideoUploadService/upload"
	"log"

	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FileServiceServer struct {
	pb.UnimplementedFileServiceServer
}

// UploadVideo receives the video in chunks from the client and writes them to a file.
func (s *FileServiceServer) UploadVideo(stream pb.FileService_UploadVideoServer) error {
	var totalReceived int64
	fileName := uuid.NewString()
	file, err := os.Create("/home/shradwo/backend_stream/videos/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Finished receiving chunks, return the final response
			return stream.SendAndClose(&pb.UploadVideoResponse{
				Status:       200, // Indicate success
				ReceivedSize: totalReceived,
			})
		}
		if err != nil {
			return err
		}

		// Write the chunk to the file
		_, err = file.Write(req.Chunk)
		if err != nil {
			return err
		}

		totalReceived += int64(len(req.Chunk))
		grpc_calls(fileName)
		fmt.Printf("Received chunk of size: %d, total received: %d\n", len(req.Chunk), totalReceived)
	}
}
func grpc_calls(uuid string) {
	// Set a reasonable timeout
	timeout := 30 * time.Second

	// Create a context with timeout for both connection and RPC
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create the gRPC connection
	conn, err := grpc.DialContext(ctx, "localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer conn.Close()

	client := pbt.NewTranscoderClient(conn)

	req := &pbt.UploadCompleteRequest{
		Uuid: uuid,
	}

	// Call the NotifyUploadComplete RPC
	res, err := client.NotifyUploadComplete(ctx, req)
	if err != nil {
		log.Printf("Error calling NotifyUploadComplete: %v", err)
		return
	}

	log.Printf("Response: StatusCode=%d", res.StatusCode)
}
