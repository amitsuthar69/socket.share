import { Progress } from "@/components/ui/progress";

interface ProgressBarProps {
  progress: number;
}

export function ProgressBar({ progress }: ProgressBarProps) {
  return (
    <Progress
      value={progress}
      className="h-2 rounded-none bg-gray-200 dark:bg-gray-800"
    />
  );
}
