def sse_request(uuid):
    def event_stream(uuid):
        progress_queue = progress_queues.get(uuid)
        if not progress_queue:
            yield f"data: {json.dumps({'error': 'No transcoding job found for this UUID'})}\n\n"
            return

        while True:
            progress = progress_queue.get()
            if progress is None:  # None is our signal that transcoding is complete
                yield f"data: {json.dumps({'progress': 100, 'status': 'complete'})}\n\n"
                break
            yield f"data: {json.dumps({'progress': progress, 'status': 'in_progress'})}\n\n"
        
        # Clean up the queue
        del progress_queues[uuid]

    return Response(event_stream(uuid), mimetype="text/event-stream")