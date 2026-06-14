import { useState, useRef } from "react"

const ACCEPTED = {
  "application/pdf": [".pdf"],
  "image/jpeg": [".jpg", ".jpeg"],
  "image/png": [".png"],
  "image/webp": [".webp"],
  "application/vnd.ms-powerpoint": [".ppt"],
  "application/vnd.openxmlformats-officedocument.presentationml.presentation": [".pptx"],
}

const QUALITY_OPTIONS = {
  pdf: [
    { value: "screen", label: "Screen", desc: "Smallest, 72dpi" },
    { value: "ebook", label: "Ebook", desc: "Balanced, 150dpi" },
    { value: "printer", label: "Print", desc: "High quality, 300dpi" },
  ],
  image: [
    { value: "low", label: "Low", desc: "Smallest size" },
    { value: "medium", label: "Medium", desc: "Balanced" },
    { value: "high", label: "High", desc: "Best quality" },
  ],
  pptx: [
    { value: "ebook", label: "Standard", desc: "Recommended" },
  ],
}

const TARGET_OPTIONS = [
  { value: "0", label: "No limit" },
  { value: "0.5", label: "500 KB" },
  { value: "1", label: "1 MB" },
  { value: "2", label: "2 MB" },
  { value: "5", label: "5 MB" },
]

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8085"

function getFileCategory(file) {
  if (file.type === "application/pdf") return "pdf"
  if (file.type.startsWith("image/")) return "image"
  return "pptx"
}

function formatBytes(bytes) {
  if (bytes < 1024) return bytes + " B"
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB"
  return (bytes / (1024 * 1024)).toFixed(1) + " MB"
}

export default function Uploader({ onResult, onError, onReset }) {
  const [file, setFile] = useState(null)
  const [quality, setQuality] = useState("ebook")
  const [targetSize, setTargetSize] = useState("0")
  const [dragging, setDragging] = useState(false)
  const [loading, setLoading] = useState(false)
  const [progress, setProgress] = useState(0)
  const inputRef = useRef()

  const acceptString = Object.values(ACCEPTED).flat().join(",")

  function handleFile(f) {
    const allExts = Object.values(ACCEPTED).flat()
    const ext = "." + f.name.split(".").pop().toLowerCase()
    if (!allExts.includes(ext)) {
      onError(`Format tidak didukung. Gunakan: PDF, JPG, PNG, WEBP, PPT, PPTX`)
      return
    }
    onReset()
    setFile(f)
    setTargetSize("0")
    const cat = getFileCategory(f)
    setQuality(cat === "image" ? "medium" : "ebook")
  }

  function onDrop(e) {
    e.preventDefault()
    setDragging(false)
    const f = e.dataTransfer.files[0]
    if (f) handleFile(f)
  }

  async function handleCompress() {
    if (!file) return
    setLoading(true)
    setProgress(10)
    onError(null)

    try {
      const formData = new FormData()
      formData.append("file", file)
      formData.append("quality", quality)
      formData.append("target_size", targetSize)

      setProgress(30)

      const res = await fetch(`${API_BASE}/api/v1/compress`, {
        method: "POST",
        body: formData,
      })

      setProgress(80)

      if (!res.ok) {
        const json = await res.json()
        throw new Error(json.error || "Gagal compress file")
      }

      const originalSize = parseInt(res.headers.get("X-Original-Size") || "0")
      const compressedSize = parseInt(res.headers.get("X-Compressed-Size") || "0")
      const ratio = parseFloat(res.headers.get("X-Compression-Ratio") || "0")

      const blob = await res.blob()
      const url = URL.createObjectURL(blob)

      const disposition = res.headers.get("Content-Disposition") || ""
      const match = disposition.match(/filename="?([^"]+)"?/)
      const filename = match ? match[1] : "compressed_" + file.name

      setProgress(100)
      onResult({ url, filename, originalSize, compressedSize, ratio })
    } catch (err) {
      onError(err.message)
    } finally {
      setLoading(false)
      setProgress(0)
    }
  }

  const category = file ? getFileCategory(file) : "pdf"
  const qualityOptions = QUALITY_OPTIONS[category] || QUALITY_OPTIONS.pdf

  return (
    <div className="uploader">
      {/* Drop Zone */}
      <div
        className={`dropzone ${dragging ? "dragging" : ""} ${file ? "has-file" : ""}`}
        onClick={() => inputRef.current.click()}
        onDragOver={(e) => { e.preventDefault(); setDragging(true) }}
        onDragLeave={() => setDragging(false)}
        onDrop={onDrop}
      >
        <input
          ref={inputRef}
          type="file"
          accept={acceptString}
          style={{ display: "none" }}
          onChange={(e) => e.target.files[0] && handleFile(e.target.files[0])}
        />
        {file ? (
          <div className="file-info">
            <span className="file-icon">{category === "pdf" ? "📄" : category === "image" ? "🖼️" : "📊"}</span>
            <div>
              <p className="file-name">{file.name}</p>
              <p className="file-size">{formatBytes(file.size)}</p>
            </div>
            <button className="remove-btn" onClick={(e) => { e.stopPropagation(); setFile(null); onReset() }}>✕</button>
          </div>
        ) : (
          <div className="drop-prompt">
            <span className="drop-icon">↑</span>
            <p className="drop-text">Drop file here or <span className="drop-link">browse</span></p>
            <p className="drop-hint">PDF · JPG · PNG · WEBP · PPT · PPTX — max 150MB</p>
          </div>
        )}
      </div>

      {/* Quality Selector */}
      {file && (
        <div className="quality-selector">
          <p className="quality-label">Quality</p>
          <div className="quality-options">
            {qualityOptions.map((opt) => (
              <button
                key={opt.value}
                className={`quality-btn ${quality === opt.value ? "active" : ""}`}
                onClick={() => setQuality(opt.value)}
              >
                <span className="quality-name">{opt.label}</span>
                <span className="quality-desc">{opt.desc}</span>
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Target Size Selector — hanya tampil untuk PDF */}
      {file && category === "pdf" && (
        <div className="quality-selector">
          <p className="quality-label">Target Size</p>
          <div className="quality-options">
            {TARGET_OPTIONS.map((opt) => (
              <button
                key={opt.value}
                className={`quality-btn ${targetSize === opt.value ? "active" : ""}`}
                onClick={() => setTargetSize(opt.value)}
              >
                <span className="quality-name">{opt.label}</span>
              </button>
            ))}
          </div>
        </div>
      )}

      {/* Compress Button */}
      {file && (
        <button
          className={`compress-btn ${loading ? "loading" : ""}`}
          onClick={handleCompress}
          disabled={loading}
        >
          {loading ? (
            <>
              <span className="spinner" />
              Compressing...
            </>
          ) : "Compress"}
        </button>
      )}

      {/* Progress Bar */}
      {loading && (
        <div className="progress-bar">
          <div className="progress-fill" style={{ width: `${progress}%` }} />
        </div>
      )}
    </div>
  )
}
