"use client";
import React, { useState,useEffect, useRef } from "react";
import InsideNode from "./InsideNode";
import * as echarts from "echarts";
import { useRouter } from "next/navigation";
import Link from "next/link";

function NodeUsageBarChart({ nodes, title }) {
  const chartRef = useRef(null);
  const router = useRouter();
  useEffect(() => {
    if (!nodes || nodes.length === 0) return;

    const chart = echarts.init(chartRef.current);

    const option = {
      title: { text: title, left: "center" },
      tooltip: {
        trigger: "axis",
        axisPointer: { type: "shadow" },
        formatter: (params) => {
          const node = nodes[params[0].dataIndex];
          return `
            <b>${node.name}</b><br/>
            Role: ${node.role}<br/>
            State: ${node.state}<br/>
            CPU: ${node.cpu}m (${node.cpu_percent.toFixed(2)}%)<br/>
            Memory: ${node.memory}Mi (${node.mem_percent.toFixed(2)}%)
          `;
        }
      },
      legend: { data: ["CPU (m)", "Memory (Mi)"], top: 30, right: 10 },
      grid: { left: "5%", right: "5%", bottom: "10%", containLabel: true },
      xAxis: {
        type: "category",
        data: nodes.map((n) => n.name),
        axisLabel: { rotate: 30, fontSize: 12 }
      },
      yAxis: { type: "value", name: "Usage" },
      series: [
        {
          name: "CPU (m)",
          type: "bar",
          data: nodes.map((n) => ({
            value: n.cpu,
            itemStyle: {
              color: n.state === "Ready" ? "#42a5f5" : "red",
              shadowColor: n.state === "Ready" ? "transparent" : "rgba(255,0,0,0.7)",
              shadowBlur: n.state === "Ready" ? 0 : 15
            }
          }))
        },
        {
          name: "Memory (Mi)",
          type: "bar",
          data: nodes.map((n) => ({
            value: n.memory,
            itemStyle: {
              color: n.state === "Ready" ? "#66bb6a" : "red",
              shadowColor: n.state === "Ready" ? "transparent" : "rgba(255,0,0,0.7)",
              shadowBlur: n.state === "Ready" ? 0 : 15
            }
          }))
        }
      ]
    };

    chart.setOption(option);

    // Make bars clickable
    chart.on("click", (params) => {
      const node = nodes[params.dataIndex];
      const encoded = encodeURIComponent(JSON.stringify(node))
      router.push(`/node/${node.name}?data=${encoded}` );
      alert(`Clicked on node: ${node.name}\nRole: ${node.role}\nState: ${node.state}`);      
    });

    return () => chart.dispose();
  }, [nodes, title, router]);

  return (
    <div>
      <div ref={chartRef} style={{ width: "100%", height: "400px" }} />
    </div>
    
  );
}


export default function NodeUsage({ data }) {
  if (!data || !data.nodes) return null;

  const controlPlaneNodes = data.nodes.filter((n) => n.role === "control-plane");
  const workerNodes = data.nodes.filter((n) => n.role === "worker");

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "60px" }}>
      <NodeUsageBarChart nodes={controlPlaneNodes} title="Control Plane Nodes Usage" />
      <NodeUsageBarChart nodes={workerNodes} title="Worker Nodes Usage" />
    </div>
  );
}
