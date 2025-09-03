'use client';
import { useState, useEffect } from "react";
import NodeUsage3D from "./NodeUsage3D";

export default function ResourceUsage() {
  const [nodeUsage, setNodeUsage] = useState(null);

  const fetchResourceUsage = async () => {
    const fileName = localStorage.getItem("kubeconfigFileName");
    if (!fileName) return;

    try {
      const res = await fetch("http://localhost:5000/health", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ file_name: fileName }),
      });

      if (res.ok) {
        const data = await res.json();
        console.log("✅ Node Resource Usage:", data["node_resource_usage"]);
        setNodeUsage(data["node_resource_usage"]); 
      } else {
        console.error("❌ Failed to fetch resource usage");
      }
    } catch (err) {
      console.error("⚠️ Error fetching usage:", err);
    }
  };

  useEffect(() => {
    fetchResourceUsage();
  }, []);

  return (
    <div>
      <h3>Namespace wise Resource Usage</h3>
      {nodeUsage ? (
        <NodeUsage3D data={nodeUsage} />
      ) : (
        <p>Loading node usage...</p>
      )}
    </div>
  );
}
