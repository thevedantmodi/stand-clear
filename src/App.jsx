import React, { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import Platform from "./Platform";

function useIsMobile() {
  const [isMobile, setIsMobile] = useState(window.innerWidth < 768);

  useEffect(() => {
    const handleResize = () => setIsMobile(window.innerWidth < 768);
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  return isMobile;
}

function App() {
  const [params] = useSearchParams();
  const isMobile = useIsMobile();

  const platforms = params.getAll("platform").map((entry) => {
    const [line, stop_id, N] = entry.split(",");
    return { line, stop_id, N };
  });

  const defaultPlatforms = [
    { line: "4", stop_id: "631N", N: "3" },
    { line: "4", stop_id: "631S", N: "3" },
  ];

  // fallback if no URL params
  const list = platforms.length > 0 ? platforms : defaultPlatforms;

  // responsive logic ✅
  const visiblePlatforms = isMobile ? list.slice(0, 1) : list;

  return (
    <div
      style={{
        minHeight: "100vh",
        background: "#e6ecf2",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        gap: "2rem",
        flexWrap: "wrap",
        padding: "1rem",
      }}
    >
      {visiblePlatforms.map((p) => (
        <Platform key={`${p.line}-${p.stop_id}`} {...p} />
      ))}
    </div>
  );
}

export default App;
