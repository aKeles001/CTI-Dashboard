import { useEffect, useState } from "react";
import { Bar, BarChart, XAxis, YAxis, PieChart, Pie }from "recharts";
import { models } from '../../wailsjs/go/models';
import { GetForums, GetChartData } from '../../wailsjs/go/main/App';

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";

const Dashboard: React.FC = () => {
  const [chartData, setChartData] = useState<models.Chart[]>([]);
  const [forums, setForums] = useState<models.Forum[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    GetForums()
      .then((data) => {
        setForums(data || []);
      })
      .catch((err) => {
        console.error("Error fetching forums:", err);
        setError("Failed to load forums.");
      });
  }, []);

  useEffect(() => {
    if (forums.length > 0) {
      const promises = forums.map(forum => GetChartData(forum.forum_id));
      Promise.all(promises)
        .then(results => {
          const flattenedResults = results.flat();
          setChartData(flattenedResults);
        })
        .catch(err => {
          console.error("Error fetching chart data:", err);
          setError("Failed to load chart data.");
        });
    }
  }, [forums]);

  if (error) {
    return <div className="p-4 text-red-500">Error: {error}</div>;
  }

  const chartConfig = {
    unassigned: {
      label: "Unassigned",
      color: "#757575",
    },
    low: {
      label: "Low",
      color: "#105b0b",
    },
    medium: {
      label: "Medium",
      color: "#bc7013",
    },
    high: {
      label: "High",
      color: "#660505",
    },
  } satisfies ChartConfig;

  return (
    <div className="p-4">
    <Card>
      <CardHeader>
        <CardTitle>Forum Post Severity Distribution</CardTitle>
        <CardDescription>A breakdown of post severity levels across different forums.</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <BarChart accessibilityLayer data={chartData}>
            <XAxis
              dataKey="forum_name"
              tickLine={false}
              tickMargin={10}
              axisLine={false}
              angle={-45}
              textAnchor="end"
              height={60}
            />
            <YAxis />
            <ChartTooltip
              content={<ChartTooltipContent />}
              cursor={true}
            />
            <Bar
              dataKey="unassigned"
              stackId="a"
              fill="var(--color-unassigned)"
              radius={[0, 0, 0, 0]}
            />
            <Bar
              dataKey="low"
              stackId="a"
              fill="var(--color-low)"
              radius={[0, 0, 0, 0]}
            />
            <Bar
              dataKey="medium"
              stackId="a"
              fill="var(--color-medium)"
              radius={[0, 0, 0, 0]}
            />
            <Bar
              dataKey="high"
              stackId="a"
              fill="var(--color-high)"
              radius={[4, 4, 0, 0]}
            />
          </BarChart>
        </ChartContainer>
      </CardContent>
      </Card>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-4">
        {chartData.map((chart: any, index) => {
          const pieData = [
            { severity: "high", count: chart.high || 0, fill: "var(--color-high)" },
            { severity: "medium", count: chart.medium || 0, fill: "var(--color-medium)" },
            { severity: "low", count: chart.low || 0, fill: "var(--color-low)" },
            { severity: "unassigned", count: chart.unassigned || 0, fill: "var(--color-unassigned)" },
          ];
          return (
            <Card key={index} className="flex flex-col">
              <CardHeader className="items-center pb-0">
                <CardTitle>{chart.forum_name}</CardTitle>
                <CardDescription>Severity Distribution</CardDescription>
                <p className="text-md font-small mt-2">ID: {chart.forum_id}</p>
                 <div className="w-full mt-4 overflow-hidden rounded-lg border border-border">
                <table className="w-full text-sm text-left">
                  <thead className="bg-muted/50">
                    <tr>
                      <th className="px-3 py-2 font-medium text-muted-foreground">
                        Last Scanned
                      </th>
                      <th className="px-3 py-2 font-medium text-muted-foreground text-center">
                        High
                      </th>
                      <th className="px-3 py-2 font-medium text-muted-foreground text-center">
                        Medium
                      </th>
                      <th className="px-3 py-2 font-medium text-muted-foreground text-center">
                        Low
                      </th>
                      <th className="px-3 py-2 font-medium text-muted-foreground text-center">
                        Overall Posts
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-t">
                      <td className="px-3 py-2">
                        {new Date(chart.last_scaned).toLocaleString()}
                      </td>
                      <td className="px-3 py-2 text-center font-semibold text-red-400">
                        {chart.high}
                      </td>
                      <td className="px-3 py-2 text-center font-semibold text-yellow-400">
                        {chart.medium}
                      </td>
                      <td className="px-3 py-2 text-center font-semibold text-green-400">
                        {chart.low}
                      </td>
                      <td className="px-3 py-2 text-center font-semibold text-gray-400">
                        {chart.count}
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
              </CardHeader>
              <CardContent className="flex-1 pb-0">
                <ChartContainer
                  config={chartConfig}
                  className="mx-auto aspect-square max-h-[250px]"
                >
                  <PieChart>
                    <ChartTooltip
                      cursor={false}
                      content={<ChartTooltipContent hideLabel />}
                    />
                    <Pie data={pieData} dataKey="count" nameKey="severity" innerRadius={60} />
                  </PieChart>
                </ChartContainer>
              </CardContent>
            </Card>
          );
        })}
      </div>
    </div>
    );
}

export default Dashboard;