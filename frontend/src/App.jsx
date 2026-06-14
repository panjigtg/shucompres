import { useState } from "react"
import Uploader from "./components/Uploader"
import Result from "./components/Result"
import "./App.css"

const HUGINN_STORY = [
  {
    rune: "H",
    title: "Huginn",
    body: "In the age before memory, Odin sent two ravens across the nine realms each dawn - Huginn, Thought, and Muninn, Memory. Where others carried weight, Huginn carried only what mattered."
  },
  {
    rune: "M",
    title: "The Gift of Lightness",
    body: "The ravens returned each evening to Odin's shoulders, whispering what they had seen. Their power was not in what they carried - but in what they left behind."
  },
  {
    rune: "S",
    title: "Shucompress",
    body: "Named after Huginn's flight. Your files, stripped of excess - swift, lean, and ready. PDF, images, presentations. Nothing leaves your machine."
  }
]

export default function App() {
  const [theme, setTheme] = useState(
    () => localStorage.getItem("theme") || "dark"
  )

  const [phase, setPhase] = useState("story")
  const [storyPage, setStoryPage] = useState(0)
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)

  const isDark = theme === "dark"
  const isLastPage = storyPage === HUGINN_STORY.length - 1
  const current = HUGINN_STORY[storyPage]

  function handleNext() {
    if (isLastPage) {
      setPhase("app")
    } else {
      setStoryPage((p) => p + 1)
    }
  }

  return (
    <div className="app" data-theme={theme}>
      <div className="sky-sheen" aria-hidden="true" />
      <div className="cloud-field" aria-hidden="true">
        <span className="cloud cloud-a" />
        <span className="cloud cloud-b" />
        <span className="cloud cloud-c" />
        <span className="cloud cloud-d" />
      </div>

      <button
        className="theme-toggle"
        onClick={() => {
          const next = isDark ? "light" : "dark"
          setTheme(next)
          localStorage.setItem("theme", next)
        }}
        aria-label="Toggle theme"
      >
        {isDark ? "☾" : "☀︎"}
      </button>

      {phase === "story" ? (
        <div className="story">
          <div className="story-dots">
            {HUGINN_STORY.map((_, i) => (
              <button
                key={i}
                className={`dot ${i === storyPage ? "active" : ""} ${i < storyPage ? "done" : ""}`}
                onClick={() => setStoryPage(i)}
                aria-label={`Go to story page ${i + 1}`}
              />
            ))}
          </div>

          <div className="story-rune" key={storyPage}>
            {current.rune}
          </div>

          <div className="story-content" key={"c" + storyPage}>
            <h1 className="story-title">{current.title}</h1>
            <p className="story-body">{current.body}</p>
          </div>

          <div className="story-nav">
            {storyPage > 0 && (
              <button className="story-btn ghost" onClick={() => setStoryPage((p) => p - 1)}>
                Back
              </button>
            )}
            <button className="story-btn primary" onClick={handleNext}>
              {isLastPage ? "Use Shucompress ->" : "Continue ->"}
            </button>
          </div>

          {!isLastPage && (
            <button className="story-skip" onClick={() => setPhase("app")}>
              skip
            </button>
          )}
        </div>
      ) : (
        <>
          <header className="header">
            <button
              className="logo-btn"
              onClick={() => {
                setPhase("story")
                setStoryPage(0)
                setResult(null)
                setError(null)
              }}
            >
              {/* <span className="logo-rune">H</span> */}
              <span className="logo-shu">shucompress</span>
              {/* <span className="logo-compress">compress</span> */}
            </button>
            <p className="tagline">Swift as Huginn. Nothing left behind.</p>
          </header>

          <main className="main">
            <Uploader
              onResult={setResult}
              onError={setError}
              onReset={() => {
                setResult(null)
                setError(null)
              }}
            />
            {error && (
              <div className="error-box">
                <span className="error-icon">!</span>
                {error}
              </div>
            )}
            {result && <Result result={result} />}
          </main>

          <footer className="footer">
            Named after Huginn, raven of Odin - thought made weightless.
            <br />
            Files processed locally. Never stored. Never sent.
          </footer>
        </>
      )}
    </div>
  )
}
