import { Greet2 } from "../wailsjs/go/main/App";
import { useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { FileList } from "./components/file-list";
import { DeviceList } from "./components/device-list";

export default function SocketShare() {
  const [uploadProgress, setUploadProgress] = useState(50);

  const files = [
    { name: "some_video.mp4", uploadTime: "<time of upload>" },
    { name: "some_audio.mp3", uploadTime: "<time of upload>" },
    { name: "some_file.docs", uploadTime: "<time of upload>" },
  ];

  const devices = [
    "192.168.0.102",
    "192.168.0.103",
    "192.168.0.104",
    "192.168.0.105",
    "192.168.0.106",
    "192.168.0.107",
    "192.168.0.108",
    "192.168.0.109",
  ];

  const handleDownload = (fileName: string) => {
    console.log("Downloading:", fileName);
  };

  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      console.log("Selected file:", file.name);
    }
  };

  return (
    <div className="min-h-screen w-full flex flex-col bg-[#0f1515] dark:bg-gray-900 text-white">
      {/* Header */}
      <div className="text-center py-4 shadow-2xl relative">
        <h1 className="text-2xl font-medium">
          <div className="flex items-center gap-2 justify-center">
            <img src="./src/assets/socketshare.svg" alt="logo" />
            Socket Share
          </div>
        </h1>
        <div className="absolute right-4 top-1/2 -translate-y-1/2"></div>
      </div>

      {/* Main content */}
      <div className="flex-1 flex relative">
        {/* Left section - File list and upload */}
        <div className="flex-1 p-6 flex flex-col">
          <div className="flex-1 mb-6">
            <FileList files={files} onDownload={handleDownload} />
          </div>
          <div>
            <Button
              onClick={handleButtonClick}
              variant="outline"
              className="w-fit bg-[#101a1b] hover:text-gray-50 hover:bg-[#1a2c2c] text-gray-50 border-gray-400 text-sm rounded-xl"
            >
              <img src="./src/assets/upload.svg" /> Upload
            </Button>
            <input
              type="file"
              ref={fileInputRef}
              onChange={handleFileChange}
              className="hidden"
            />
          </div>
        </div>

        {/* Right section - Devices */}
        <div className="w-80 p-6">
          <DeviceList devices={devices} />
        </div>
      </div>
    </div>
  );
}
