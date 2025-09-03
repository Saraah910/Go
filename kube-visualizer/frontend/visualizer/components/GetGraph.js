"use client";
import { useEffect, useState, useRef } from "react";
import dynamic from "next/dynamic";

let ForceGraph3D;
if (typeof window !== "undefined") {
  // Only require on client side
  ForceGraph3D = require("3d-force-graph").default;
}

export default function GetGraph({ ns }) {
  const [graph, setGraph] = useState(null);
  const [loading, setLoading] = useState(true);
  const containerRef = useRef(null);
  const fgInstance = useRef(null);

  useEffect(() => {
    async function fetchGraph() {
      try {
        const res = await fetch(`http://localhost:5000/get-graph-data/${ns}`);
        const data = await res.json();
        setGraph(data["graph"]);
        console.log(data["graph"]);
      } catch (err) {
        console.error("Failed to fetch graph:", err);
      } finally {
        setLoading(false);
      }
    }
    fetchGraph();
  }, [ns]);

  useEffect(() => {
    if (containerRef.current && !fgInstance.current && ForceGraph3D) {
      fgInstance.current = ForceGraph3D()(containerRef.current)
        .width(3000)
        .height(900)
        .nodeLabel((node) => `${node.id}\nType: ${node.type}`)
        .nodeAutoColorBy("type")
        .linkDirectionalParticles(4)
        .linkDirectionalParticleWidth(2)
        .linkDirectionalParticleSpeed((link) =>
          link.rate ? Math.min(link.rate * 0.001, 0.05) : 0.002
        )
        .linkDirectionalParticleColor(() => "orange")
        .graphData(graph);
    } else if (fgInstance.current && graph) {
      fgInstance.current.graphData(graph);
    }
  }, [graph]);

  if (loading) return <div>Loading graph...</div>;
  if (!graph) return <div>No graph data available</div>;

  return (
    <div>
        <h3>Graph Data for namespace: {ns}</h3>
        <div
        ref={containerRef}
        style={{
            width: "1200px",   // must match .width()
            height: "600px",
            boxSizing: "border-box",
            padding: "20px",
            overflow: "hidden", // keeps it tidy
        }}
        />
    </div>
);

}
