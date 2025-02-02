import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

interface FileItem {
  name: string;
  uploadTime: string;
}

interface FileListProps {
  files: FileItem[];
  onDownload: (fileName: string) => void;
}

export function FileList({ files, onDownload }: FileListProps) {
  return (
    <Card className="rounded-md shadow-none bg-[#131a1c] text-white">
      <CardHeader className="px-4 py-3 border-b border-gray-400">
        <CardTitle className="text-md font-medium">
          <div className="flex gap-2">
            <img src="./src/assets/folder.svg" /> File Registry
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="p-0">
        {files.map((file, index) => (
          <div key={index}>
            <div className="flex items-center justify-between px-4 py-2">
              <span className="text-sm text-gray-50">
                {file.name} • {file.uploadTime}
              </span>
              <Button
                variant="default"
                size="icon"
                className="text-lg h-8 w-8 rounded-sm text-[#bcfcff] bg-[#131a1c] hover:bg-[#0f1515] hover:text-gray-50"
                onClick={() => onDownload(file.name)}
              >
                <img src="./src/assets/download.svg" />
              </Button>
            </div>
            {index < files.length - 1 && <Separator className="bg-gray-400" />}
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
