"use client";
import { useParams, useSearchParams } from "next/navigation";
import InsideNode from "@/components/InsideNode"; 

export default function NodePage() {
  const searchParams = useSearchParams();
  const nodeData = searchParams.get("data");
  let node = null;
  if (nodeData) {
    try {
      node = JSON.parse(decodeURIComponent(nodeData));
    } catch (err) {
      console.error("Invalid node data:", err);
    }
  }

  return (
    <div style={{ padding: "20px" }}>
      <InsideNode node={node} />
    </div>
  );
}
