---
title: Build and run computer vision applications with Docker
linktitle: Computer vision applications
keywords: Docker Offload, computer vision
summary: |
  Learn how to create computer vision applications and run them locally using Docker Offload.
languages: [python]
tags: [AI]
params:
  time: 15 minutes
---

## Introduction

Computer vision applications are at the heart of many modern AI systems, from
autonomous vehicles to medical imaging to content moderation. These applications
often require significant computational power to process images and video in
real-time, making GPU acceleration essential for production workloads. This
guide demonstrates how to build, containerize, and run a GPU-accelerated
computer vision application using Docker. You'll create a Python application
that performs object detection in videos. While it can be ran on your local
development machine, you'll leverage [Docker Offload](/offload/) to give your
application access to GPU resources in the cloud that let it process video data
faster.

What makes this approach powerful is the seamless workflow: you'll develop
locally using familiar Docker commands, then run it on GPU-accelerated cloud
infrastructure with the exact same commands. No complex cloud configurations, no
infrastructure management, just Docker doing what it does best. By the end of
this guide, you'll be exposed to examples of how to:

- Build computer vision applications that leverage GPU acceleration
- Containerize AI workloads for consistent deployment across environments
- Use Docker Offload to access powerful GPU resources in the cloud

Whether you're building image processing pipelines, training computer vision
models, or deploying AI-powered applications at scale, this guide provides the
foundation for GPU-accelerated containerized workflows.

## Prerequisites

To follow this guide, you need to:

 - [Install the latest version of Docker Desktop](../manuals/desktop/release-notes.md)
 - [Set up Docker Offload billing](/offload/usage/#docker-offload-billing) by purchasing credits or enabling on-demand usage
 - Have basic knowledge of Python and Docker. To learn Docker basics, see [Get
   started with Docker](/get-started/).

## Step 1: Create the computer vision application

First, create a simple Python app that performs object detection in videos.

1. Create the project structure. In your terminal and a directory of your
   choice, run the following commands to create the project directory and
   sub-directories:

   ```console
   $ mkdir cv-docker-app
   $ cd cv-docker-app
   $ mkdir input output
   ```

2. Create `app.py` in `cv-docker-app` using a text or code editor of your
   choice, and copy-paste the following code.

   The sample application uses the [Ultralytics
   YOLO](https://github.com/ultralytics/yolov5) model for object detection, the
   `cv2` library for video processing, and the `torch` library for running the
   model.

   ```python{title="app.py", collapse=true}
   import cv2
   import numpy as np
   import os
   import time
   import torch
   from ultralytics import YOLO
   
   def check_gpu():
       """Check GPU availability for video processing"""
       print("=== GPU Information ===")
       
       if torch.cuda.is_available():
           gpu_name = torch.cuda.get_device_name(0)
           total_memory = torch.cuda.get_device_properties(0).total_memory / 1024**3
           
           print(f"✅ GPU Available: {gpu_name}")
           print(f"GPU Memory: {total_memory:.1f} GB")
           
           if total_memory >= 20.0:
               print("🚀 EXCELLENT: L4-class GPU detected - optimal for video processing")
           elif total_memory >= 12.0:
               print("✅ GOOD: High-end GPU - suitable for video processing")
           elif total_memory >= 6.0:
               print("⚠️  MARGINAL: May struggle with high-resolution video processing")
           else:
               print("❌ INSUFFICIENT: Video processing requires 8GB+ GPU memory")
               print("   Recommendation: Use Docker Offload with L4 GPU")
       else:
           print("❌ No GPU detected - video processing will be extremely slow")
           
   def load_detection_model():
       """Load YOLO model for object detection"""
       print("Loading YOLO model for video processing...")
       
       model = YOLO('yolov8n.pt')  # Use nano for faster video processing
       
       if torch.cuda.is_available():
           model.to('cuda')
           print("✅ Model loaded on GPU")
       else:
           print("⚠️  Model running on CPU - video processing will be slow")
       
       return model
   
   def process_video(video_path, output_path, model):
       """Process video with object detection frame by frame"""
       print(f"Processing video: {video_path}")
       
       # Open video
       cap = cv2.VideoCapture(video_path)
       if not cap.isOpened():
           print(f"❌ Error opening video: {video_path}")
           return
       
       # Get video properties
       fps = int(cap.get(cv2.CAP_PROP_FPS))
       width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
       height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
       total_frames = int(cap.get(cv2.CAP_PROP_FRAME_COUNT))
       duration = total_frames / fps
       
       print(f"Video specs: {width}x{height}, {fps}fps, {duration:.1f}s ({total_frames} frames)")
       
       # Check if this is demanding for GPU
       pixel_throughput = width * height * fps
       if pixel_throughput > 200_000_000:  # >200M pixels/sec
           print("🔥 HIGH-RESOLUTION VIDEO: This will stress test your GPU")
           if not torch.cuda.is_available() or torch.cuda.get_device_properties(0).total_memory < 8 * 1024**3:
               print("⚠️  WARNING: High-resolution video processing needs 8GB+ GPU")
               print("   Consider using Docker Offload with L4 GPU for smooth processing")
       
       # Setup video writer
       fourcc = cv2.VideoWriter_fourcc(*'mp4v')
       out = cv2.VideoWriter(output_path, fourcc, fps, (width, height))
       
       frame_count = 0
       detection_count = 0
       start_time = time.time()
       
       print("Processing frames...")
       
       while True:
           ret, frame = cap.read()
           if not ret:
               break
           
           # Run object detection on frame
           results = model(frame, verbose=False)
           
           # Draw detections
           frame_detections = 0
           if results[0].boxes is not None:
               boxes = results[0].boxes
               for box in boxes:
                   conf = box.conf[0]
                   if conf > 0.5:
                       # Get coordinates and class
                       x1, y1, x2, y2 = box.xyxy[0].cpu().numpy().astype(int)
                       class_id = int(box.cls[0])
                       class_name = model.names[class_id]
                       
                       # Draw bounding box and label
                       cv2.rectangle(frame, (x1, y1), (x2, y2), (0, 255, 0), 2)
                       label = f"{class_name}: {conf:.2f}"
                       cv2.putText(frame, label, (x1, y1-10), 
                                  cv2.FONT_HERSHEY_SIMPLEX, 0.5, (0, 255, 0), 2)
                       frame_detections += 1
           
           detection_count += frame_detections
           
           # Add frame info
           frame_info = f"Frame: {frame_count+1}/{total_frames} | Objects: {frame_detections}"
           cv2.putText(frame, frame_info, (10, 30), 
                      cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 255, 0), 2)
           
           # Write processed frame
           out.write(frame)
           frame_count += 1
           
           # Progress update
           if frame_count % 30 == 0:  # Every 30 frames
               elapsed = time.time() - start_time
               fps_actual = frame_count / elapsed if elapsed > 0 else 0
               progress = (frame_count / total_frames) * 100
               print(f"  Progress: {progress:.1f}% | Processing: {fps_actual:.1f} fps")
       
       # Cleanup
       cap.release()
       out.release()
       
       total_time = time.time() - start_time
       avg_fps = frame_count / total_time if total_time > 0 else 0
       
       print(f"✅ Video processing complete!")
       print(f"   Processed {frame_count} frames in {total_time:.1f}s")
       print(f"   Average processing speed: {avg_fps:.1f} fps")
       print(f"   Total objects detected: {detection_count}")
       print(f"   Output saved to: {output_path}")
       
       # Performance assessment
       if avg_fps >= fps * 0.8:  # Processing at 80%+ of original speed
           print("🚀 EXCELLENT: Real-time processing achieved")
       elif avg_fps >= 10:
           print("✅ GOOD: Fast processing speed")
       elif avg_fps >= 2:
           print("⚠️  SLOW: Consider using Docker Offload for faster processing")
       else:
           print("❌ VERY SLOW: Consider using Docker Offload for faster processing")
   
   def create_sample_video():
       """Create a sample video with moving objects"""
       print("Creating sample video for testing...")
       
       # Video properties
       width, height = 1280, 720
       fps = 30
       duration = 10  # seconds
       total_frames = fps * duration
       
       # Create video writer
       fourcc = cv2.VideoWriter_fourcc(*'mp4v')
       out = cv2.VideoWriter('/app/input/sample_video.mp4', fourcc, fps, (width, height))
       
       for frame_num in range(total_frames):
           # Create frame with moving objects
           frame = np.zeros((height, width, 3), dtype=np.uint8)
           
           # Moving rectangle
           rect_x = int((frame_num * 5) % (width + 200)) - 100
           cv2.rectangle(frame, (rect_x, 300), (rect_x + 150, 380), (0, 0, 255), -1)
           
           # Moving circle
           circle_x = int(width - (frame_num * 3) % (width + 100))
           cv2.circle(frame, (circle_x, 200), 40, (255, 0, 0), -1)
           
           # Static objects
           cv2.rectangle(frame, (100, 500), (200, 600), (0, 255, 0), -1)
           cv2.rectangle(frame, (width-200, 500), (width-100, 600), (0, 255, 255), -1)
           
           # Add frame number
           cv2.putText(frame, f"Frame: {frame_num+1}", (50, 50), 
                      cv2.FONT_HERSHEY_SIMPLEX, 1, (255, 255, 255), 2)
           
           out.write(frame)
       
       out.release()
       print("✅ Sample video created: /app/input/sample_video.mp4")
   
   def main():
       print("🎥 Video Object Detection Processor")
       check_gpu()
       
       model = load_detection_model()
       
       input_dir = "/app/input"
       output_dir = "/app/output"
       
       # Look for existing video files first
       video_files = [f for f in os.listdir(input_dir) 
                     if f.lower().endswith(('.mp4', '.avi', '.mov', '.mkv'))]
       
       if video_files:
           print(f"Found video file: {video_files[0]}")
           video_path = os.path.join(input_dir, video_files[0])
       else:
           print("No video files found. Creating sample video...")
           create_sample_video()
           video_path = "/app/input/sample_video.mp4"
       
       # Process the video
       output_video = "/app/output/processed_video.mp4"
       process_video(video_path, output_video, model)
       
       print("\n✨ Processing complete!")
       print("💡 TIP: Place your own .mp4 files in the input/ directory")
       print("     for processing real videos with object detection")
       print("Check /app/output/ for the processed video")
   
   if __name__ == "__main__":
       main()
   ```

3. Create a `requirements.txt` file in the `cv-docker-app` directory with the
   following contents:

   ```text{title="requirements.txt"}
   opencv-python==4.8.1.78
   ultralytics==8.0.196
   numpy==1.24.3
   pillow>=8.3.2
   ```

At this point, you have a simple Python application that can perform object
detection on video files. The application uses the Ultralytics YOLO model for
real-time object detection. You can run this application locally using Python,
but by using Docker, you can easily containerize it so you can run it on powerful
GPU-accelerated cloud infrastructure.

You should have the following directory structure:

```text
cv-docker-app/
├── app.py
├── input/
├── output/
└── requirements.txt
```

## Step 2: Containerize the application

Containerizing this application is a simple as creating a [Dockerfile](/reference/dockerfile/) that
specifies the base image, installs the required dependencies, copies the
application code into the container, and then executes it.

In the `cv-docker-app` directory, create a file named `Dockerfile` with the
following contents:

```dockerfile{title="Dockerfile", collapse=true}
# Use PyTorch official image from Docker Hub
FROM pytorch/pytorch:latest

# Set environment variables
ENV PYTHONUNBUFFERED=1

# Install system dependencies
RUN apt-get update && apt-get install -y \
    libgl1-mesa-glx \
    libglib2.0-0 \
    libsm6 \
    libxext6 \
    libxrender-dev \
    libgomp1 \
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy requirements file
COPY requirements.txt .

# Install Python packages
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY app.py .

# Create input and output directories
RUN mkdir -p /app/input /app/output

# Default command
CMD ["python", "app.py"]
```

You should have the following directory structure:

```text
cv-docker-app/
├── app.py
├── Dockerfile
├── input/
├── output/
└── requirements.txt
```

## Step 3: Optional. Start Docker Offload

You can run this application locally if your GPU is powerful enough. The
recommended option is to use [Docker Offload](/offload/).

Docker Offload is a feature that lets you build and run containers in a powerful
GPU-accelerated cloud environment. This ensures that the application can perform
object detection at real-time speeds. Before you can use it, you must either
purchase credits or enable on-demand usage in [Docker Offload
billing](/offload/usage/#docker-offload-billing).

To start Docker Offload, in a terminal, run:

```console
$ docker offload start
```

When prompted, select your account, and then select **Yes** to enable GPU
support.

You can stop Docker Offload at any time by running `docker offload stop`.

You can also check the status of Docker Offload by running:

```console
$ docker offload status
```

## Step 4: Build and run the container

1. In the terminal, navigate to the `cv-docker-app` directory and build the
   Docker image:

   ```console
   $ docker build -t cv-gpu-app .
   ```

2. Run the container with GPU support. Also, bind mount the `input` and `output`
   directories to share data between the host and the container. In a terminal,
   run:

   ```console
   $ docker run --rm \
     --gpus all \
     --mount type=bind,source=$(pwd)/input,target=/app/input \
     --mount type=bind,source=$(pwd)/output,target=/app/output \
     cv-gpu-app
    ```

> [!NOTE]
>
>  The application will create its own video if none are present. You can place
>  your own video files in the `input/` directory before running the container.


You'll see output similar to the following if using Docker Offload:

```console{title="Output", collapse=true}
🎥 Video Object Detection Processor
=== GPU Information ===
✅ GPU Available: NVIDIA L4
GPU Memory: 22.0 GB
🚀 EXCELLENT: L4-class GPU detected - optimal for video processing
Loading YOLO model for video processing...
Downloading https://github.com/ultralytics/assets/releases/download/v0.0.0/yolov8n.pt to 'yolov8n.pt'...
100%|██████████| 6.23M/6.23M [00:00<00:00, 372MB/s]
✅ Model loaded on GPU
No video files found. Creating sample video...
Creating sample video for testing...
✅ Sample video created: /app/input/sample_video.mp4
Processing video: /app/input/sample_video.mp4
Video specs: 1280x720, 30fps, 10.0s (300 frames)
Processing frames...
  Progress: 10.0% | Processing: 23.1 fps
  Progress: 20.0% | Processing: 37.7 fps
  Progress: 30.0% | Processing: 47.6 fps
  Progress: 40.0% | Processing: 54.6 fps
  Progress: 50.0% | Processing: 60.1 fps
  Progress: 60.0% | Processing: 64.7 fps
  Progress: 70.0% | Processing: 67.9 fps
  Progress: 80.0% | Processing: 70.7 fps
  Progress: 90.0% | Processing: 72.9 fps
  Progress: 100.0% | Processing: 74.8 fps
✅ Video processing complete!
   Processed 300 frames in 4.0s
   Average processing speed: 74.7 fps
   Total objects detected: 111
   Output saved to: /app/output/processed_video.mp4
🚀 EXCELLENT: Real-time processing achieved

✨ Processing complete!
💡 TIP: Place your own .mp4 files in the input/ directory
     for processing real videos with object detection
Check /app/output/ for the processed video
```

You can now find the processed video in the `output/` directory. The output
video will have bounding boxes around detected objects.

## Summary

Computer vision applications represent a critical category of AI workloads that
demand significant computational resources. In this guide, you built a
GPU-accelerated object detection application and learned how to containerize and
run it using Docker's AI-focused tools.

The power of this approach lies in its simplicity: you developed a computer
vision application locally, containerized it with standard Docker commands, and
then seamlessly ran it on powerful GPUs in the cloud using Docker Offload. No
cloud configuration complexity, no infrastructure management overhead, just
familiar Docker workflows that scale from development to production. You
explored how Docker simplifies GPU-accelerated AI development through:

- Docker Offload: Access to powerful GPU resources in the cloud using the same
  commands you use locally, eliminating the need for complex cloud
  infrastructure setup
- Containerization: Using the official PyTorch image from Docker Hub
  that includes pre-configured CUDA support and optimized Python environments
- Volume mounting: Seamless file sharing between your local development
  environment and GPU-accelerated containers
- GPU resource management: Automatic GPU detection and utilization through
  Docker's `--gpus` flag

This foundation extends beyond object detection to any GPU-intensive computer
vision workload, from real-time video processing to large-scale image analysis
pipelines. The same containerized workflow supports model training, inference
serving, and batch processing tasks.