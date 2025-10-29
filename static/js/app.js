// static/app.js

document.addEventListener("DOMContentLoaded", () => {
    const uploadForm = document.getElementById("uploadForm");
    const fileInput = document.getElementById("fileInput");
    const uploadStatus = document.getElementById("uploadStatus");
  
    uploadForm.addEventListener("submit", async (e) => {
      e.preventDefault();
  
      const file = fileInput.files[0];
      if (!file) {
        uploadStatus.textContent = "Please select a file first.";
        uploadStatus.style.color = "red";
        return;
      }
  
      const formData = new FormData();
      formData.append("file", file);
  
      try {
        uploadStatus.textContent = "Uploading...";
        const response = await fetch("/upload", {
          method: "POST",
          body: formData,
        });
  
        const result = await response.json();
        if (result.success) {
          uploadStatus.textContent = `✅ Uploaded successfully: ${result.file.original_name}`;
          uploadStatus.style.color = "green";
          fileInput.value = "";
        } else {
          uploadStatus.textContent = "Upload failed.";
          uploadStatus.style.color = "red";
        }
      } catch (err) {
        console.error("Upload error:", err);
        uploadStatus.textContent = "Error uploading file.";
        uploadStatus.style.color = "red";
      }
    });
  });
  