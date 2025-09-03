import { useEffect,useState } from "react";
import GetGraph from "./GetGraph";
export default function InsideNode({ node }) {
  const [namespaces, setNamespaces] = useState([]);
  const [loading, setLoading] = useState(true);
  const [IsSelectedNamespace, setIsSelectedNamespace] = useState(false)
  const[selectedNs, setSelectedNs] = useState()

  const getGraphForNamespace = (ns)=>{
    setIsSelectedNamespace(true)
    setSelectedNs(ns)
    console.log(selectedNs,setIsSelectedNamespace)
  }
  useEffect(() => {
    async function fetchNamespaces() {
      try {
        const res = await fetch("http://localhost:5000/data");
        const data = await res.json();
        if (res.ok) {
            console.log(data["namespaces"]["namespaces"])
          setNamespaces(data["namespaces"]["namespaces"]);
          console.log(namespaces)
        } else {
          console.error("Error fetching namespaces:", data.error);
        }
      } catch (err) {
        console.error("Failed to fetch namespaces:", err);
      } finally {
        setLoading(false);
      }
    }
    fetchNamespaces();
  }, []);
  if (!node) {
    return <p>No node data provided.</p>;
  }

  return (
    <div
    style={{
        fontFamily: "Arial, sans-serif",
        background: "#f9fafc",
        minHeight: "100vh",   // covers full screen but grows dynamically
        padding: "20px",
        display: "flex",
        flexDirection: "column",
    }}
    >

      {/* Header */}
      <div style={{ textAlign: "center", marginBottom: "20px" }}>
        <h1
          style={{
            fontSize: "28px",
            fontWeight: "bold",
            margin: "0",
            color: "#333",
            fontFamily:"initial"
          }}
        >
          Node Insider
        </h1>
        <hr
          style={{
            marginTop: "10px",
            border: "none",
            borderTop: "2px solid #ddd",
            width: "100%",
          }}
        />
      </div>

      {/* Main Layout */}
      <div style={{ display: "flex", flex: 2, gap: "20px" }}>
        {/* Left Sidebar - Namespaces */}
        <div
          style={{
            width: "220px",
            background: "#fff",
            borderRadius: "12px",
            padding: "15px",
            boxShadow: "0 2px 10px rgba(0,0,0,0.1)",
          }}
        >
          <h3
            style={{
              margin: "0 0 10px 0",
              fontSize: "16px",
              borderBottom: "1px solid #eee",
              paddingBottom: "8px",
              color: "#444",
            }}
          >
            Namespaces 
          </h3>
          {loading ? (
            <p style={{ fontSize: "14px", color: "#777" }}>Loading...</p>
          ) : namespaces.length > 0 ? (
            <ul style={{ listStyle: "none", padding: "0", margin: "0" }}>
                
              {namespaces.map((ns, index) => (
                <li
                  key={index}
                  style={{
                    marginBottom: "8px",
                    cursor: "pointer",
                    color: "#007bff",
                  }}
                  onClick={
                    ()=>{getGraphForNamespace(ns)}
                  }
                >
                  {ns}
                </li>
              ))}
            </ul>
          ) : (
            <p style={{ fontSize: "14px", color: "red" }}>No namespaces found</p>
          )}
        </div>

        {/* Center Content - Empty for graphs */}
        <div
          style={{
            flex: 1,
            background: "#fff",
            borderRadius: "12px",
            boxShadow: "0 2px 10px rgba(0,0,0,0.1)",
            padding: "20px",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            color: "#999",
            fontSize: "18px",
            minWidth: "0", // prevents overflow,
            minHeight: "0", // prevents overflow
          }}
        >
          {IsSelectedNamespace ? (
            <div style={{width:"100%", height:"100%"}}>
              <GetGraph ns={selectedNs}/>
            </div>
          ):(
            <div>
              Graphs will appear here...
            </div>
          )}
        </div>

        {/* Right Sidebar - Search & Dropdown */}
        <div
          style={{
            width: "250px",
            background: "#fff",
            borderRadius: "12px",
            padding: "15px",
            boxShadow: "0 2px 10px rgba(0,0,0,0.1)",
          }}
        >
          <input
            type="text"
            placeholder="Search resources..."
            style={{
              width: "100%",
              padding: "10px",
              marginBottom: "12px",
              borderRadius: "8px",
              border: "1px solid #ccc",
              outline: "none",
              fontSize: "14px",
            }}
          />
          <select
            style={{
              width: "100%",
              padding: "10px",
              borderRadius: "8px",
              border: "1px solid #ccc",
              fontSize: "14px",
              background: "#fff",
            }}
          >
            <option>Pods</option>
            <option>Deployments</option>
            <option>Services</option>
            <option>ConfigMaps</option>
            <option>Secrets</option>
          </select>
        </div>
      </div>

      {/* Node Info - Below */}
      <div
        style={{
          marginTop: "20px",
          background: "#fff",
          borderRadius: "12px",
          padding: "20px",
          boxShadow: "0 2px 10px rgba(0,0,0,0.1)",
        }}
      >
        <h2 style={{ marginBottom: "15px", color: "#333", fontFamily:"initial" }}>Node Details</h2>
        <p><strong>Name:</strong> {node.name}</p>
        <p><strong>Role:</strong> {node.role}</p>
        <p><strong>State:</strong> {node.state}</p>
        <p><strong>CPU:</strong> {node.cpu}</p>
        <p><strong>Memory:</strong> {node.memory}</p>
        <p><strong>CPU %:</strong> {node.cpu_percent}</p>
        <p><strong>Memory %:</strong> {node.mem_percent}</p>
      </div>
    </div>
  );
}
