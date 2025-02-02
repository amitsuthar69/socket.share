import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

interface DeviceListProps {
  devices: string[];
}

export function DeviceList({ devices }: DeviceListProps) {
  return (
    <Card className="rounded-md shadow-none border-gray-200 bg-[#131a1c] text-white">
      <CardHeader className="px-4 py-3 border-b border-gray-200">
        <CardTitle className="text-md font-medium">
          <div className="flex gap-2">
            <img src="./src/assets/network.svg" /> Peers Online
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="p-0">
        {devices.map((ip, index) => (
          <div key={index}>
            <div className="px-4 py-2">
              <span className="text-green-300">{ip}</span>
            </div>
            {index < devices.length - 1 && (
              <Separator className="dark:bg-gray-800" />
            )}
          </div>
        ))}
      </CardContent>
    </Card>
  );
}
