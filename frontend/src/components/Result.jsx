function formatBytes(bytes) {
  if (bytes < 1024) return bytes + " B"
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB"
  return (bytes / (1024 * 1024)).toFixed(1) + " MB"
}

export default function Result({ result }) {
  const { url, filename, originalSize, compressedSize, ratio } = result

  function handleDownload() {
    const a = document.createElement("a")
    a.href = url
    a.download = filename
    a.click()
  }

  return (
    <div className="result">
      <div className="result-header">
        <span className="result-check">✓</span>
        <p className="result-title">Compression complete</p>
      </div>

      <div className="result-stats">
        <div className="stat">
          <p className="stat-label">Original</p>
          <p className="stat-value original">{formatBytes(originalSize)}</p>
        </div>

        <div className="stat-arrow">
          <div className="ratio-badge">−{ratio.toFixed(1)}%</div>
          <span>→</span>
        </div>

        <div className="stat">
          <p className="stat-label">Compressed</p>
          <p className="stat-value compressed">{formatBytes(compressedSize)}</p>
        </div>
      </div>

      {/* Visual bar perbandingan */}
      <div className="size-bar-container">
        <div className="size-bar original-bar" style={{ width: "100%" }} />
        <div
          className="size-bar compressed-bar"
          style={{ width: `${(compressedSize / originalSize) * 100}%` }}
        />
      </div>

      <button className="download-btn" onClick={handleDownload}>
        Download {filename}
      </button>
    </div>
  )
}