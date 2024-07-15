function updatePauseText() {
    const end = document.getElementById("toggleBtn").dataset.pauseEnd
    if (end === "") return

    const seconds = (new Date(end) - Date.now()) / 1000
    if (seconds < 1.0) {
        htmx.trigger("#headerForm", "updateApi")
        return
    }

    const m = Math.floor(seconds / 60).toString()
    const s = Math.floor(seconds % 60).toString()
    const text = document.getElementById("toggleBtnText")
    text.innerText = `Paused for ${m}m ${s}s`
}

setInterval(updatePauseText, 500)
