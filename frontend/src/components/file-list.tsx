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
    <Card className="rounded-md shadow-none border-gray-200 dark:border-gray-800">
      <CardHeader className="px-4 py-3 border-b border-gray-200 dark:border-gray-800">
        <CardTitle className="text-sm font-medium">file index</CardTitle>
      </CardHeader>
      <CardContent className="p-0">
        {files.map((file, index) => (
          <div key={index}>
            <div className="flex items-center justify-between px-4 py-2">
              <span className="text-sm text-gray-600 dark:text-gray-400 font-mono">
                {file.name} • {file.uploadTime}
              </span>
              <Button
                variant="outline"
                size="icon"
                className="h-6 w-6 rounded-sm"
                onClick={() => onDownload(file.name)}
              >
                ↓
              </Button>
            </div>
            {index < files.length - 1 && (
              <Separator className="dark:bg-gray-800" />
            )}
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
