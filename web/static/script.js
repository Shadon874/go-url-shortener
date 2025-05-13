document.getElementById("shorten-form").addEventListener("submit", async function (event) {
    event.preventDefault();

    const urlInput = document.getElementById("url-input");
    const url = urlInput.value;

    const response = await fetch("/api/shorten", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ url: url }),
    });
    const data = await response.json();
    const resultDiv = document.getElementById("result");
    if (response.ok) {
        resultDiv.innerHTML = `<p>Shortened URL: <a href="${data.shortened_url}" target="_blank">${data.shortened_url}</a></p>`;
    } else {
        resultDiv.innerHTML = `<p>Error: ${data.error}</p>`;
    }
});