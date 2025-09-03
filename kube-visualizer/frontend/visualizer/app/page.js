"use client";
import { useState, useEffect } from "react";
import ResourceUsage from "../components/resourceUsage";

export default function Home() {
  const [isClient, setIsClient] = useState(false);
  const [kubeconfigUploaded, setKubeconfigUploaded] = useState(false);
  const [file, setFile] = useState(null);

  useEffect(() => {
    setIsClient(true); // mark when client-side is ready
    const uploaded = localStorage.getItem("kubeconfigUploaded") === "true";
    const storedFileName = localStorage.getItem("kubeconfigFileName");
    setKubeconfigUploaded(uploaded);
    if (storedFileName) {
      setFile({ name: storedFileName });
    }
    fetch("http://localhost:5000/set-kubeconfig", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ file_name: storedFileName }),
    })
      .then(res => res.json())
      .then(data => {
        console.log("Backend reinitialized:", data);
      })
      .catch(err => console.error("Failed to reinit kubeconfig:", err));
  }, []);

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    console.log("Selected file:", file);
    setFile(file);
  };
  
  const handleUpload = async () => {
    if (!file) {
      alert("Please select a file first.");
      return;
    }
    const formData = new FormData();
    formData.append("file", file);
    console.log(formData)
    try {
      const res = await fetch("http://localhost:5000/upload", {
        method: "POST",
        body: formData
      });
      if (res.ok) {
        console.log("File uploaded successfully");
        setKubeconfigUploaded(true);
        console.log("Kubeconfig file uploaded");
        localStorage.setItem("kubeconfigUploaded", true);
        localStorage.setItem("kubeconfigFileName", file.name);
      } else {
        console.error("File upload failed");
        alert("File upload failed. Please try again.");
      };
    } catch (error) {
      console.error("Error uploading file:", error);
      alert("Error uploading file. Please try again.");
    }
  };
  if (!isClient) return null;
  return (
    <div style={{ padding: "20px", maxWidth: "800px", margin: "0 auto", textAlign: "center" }}>
      {/* Title */}
      <h1
        style={{
          fontSize: "28px",
          fontWeight: "bold",
          marginBottom: "10px",
        }}
      >
        Kubernetes Visualizer
      </h1>
      <hr style={{ border: "1px solid #ddd", marginBottom: "30px" }} />

      {kubeconfigUploaded ? (
        <div style={{ marginTop: "30px" }}>
          <h2 style={{ color: "#2e7d32", fontWeight: "600" }}>Kubeconfig Uploaded Successfully!</h2>
          <button onClick={() => {
            setKubeconfigUploaded(false);
            setFile(null);
            localStorage.removeItem("kubeconfigUploaded");
            localStorage.removeItem("kubeconfigFileName");
          }} style={{
            marginTop: "15px",
            backgroundColor: "#d32f2f",
            color: "white",
            padding: "6px 12px",
            border: "none",
            borderRadius: "6px",
            cursor: "pointer",
            fontSize: "14px",
            fontWeight: "500",
          }}>
            Go back to upload another file
          </button>
          <p style={{ color: "#555", marginTop: "10px" }}>Visualizing Kubernetes resources...</p>
          <ResourceUsage />
        </div>
      ) : (
        <div style={{ marginTop: "30px" }}>
          <p style={{ fontSize: "16px", color: "#555", marginBottom: "20px" }}>
            Please upload your kubeconfig file
          </p>
          <div style={{ display: "flex", justifyContent: "center", gap: "10px", alignItems: "center" }}>
            {/* Hidden input */}
            <input
              type="file"
              id="fileInput"
              onChange={handleFileChange}
              style={{ display: "none" }}
            />
            <label
              htmlFor="fileInput"
              style={{
                backgroundColor: "#1976d2",
                color: "white",
                padding: "8px 16px",
                borderRadius: "6px",
                cursor: "pointer",
                fontSize: "14px",
                fontWeight: "500",
              }}
            >
              Choose File
            </label>
            <div
              style={{
                minWidth: "220px",
                padding: "8px 12px",
                border: "1px solid #ccc",
                borderRadius: "6px",
                backgroundColor: "#f9f9f9",
                textAlign: "left",
                color: file ? "#222" : "#777",
                fontStyle: file ? "normal" : "italic",
                boxShadow: "inset 0 1px 3px rgba(0,0,0,0.1)",
              }}
            >
              {file ? file.name : "No file chosen"}
            </div>
            <button
              onClick={handleUpload}
              style={{
                backgroundColor: "#4caf50",
                color: "white",
                padding: "8px 16px",
                border: "none",
                borderRadius: "6px",
                cursor: "pointer",
                fontSize: "14px",
                fontWeight: "500",
              }}
            >
              Upload
            </button>
          </div>
        </div>
      )}
    </div>

  );
}
